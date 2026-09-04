#!/usr/bin/env node
// Playwright runner for verify-pagode scenarios.
//
//   node scripts/pw.mjs scenarios/<name>.mjs [--name label] [--headed] [-- scenario args]
//   node scripts/pw.mjs --doctor        # resolve playwright, launch chromium, exit
//
// Evidence for a run lands in tmp/verify-pagode/evidence/<timestamp>-<label>/:
//   NN-<step>.png   screenshot after each step (NN-<step>-FAILED.png on failure)
//   trace.zip       Playwright trace (npx playwright show-trace trace.zip)
//   run.log         step log + browser console errors + HTTP >= 400 responses
//   summary.json    steps, status, base URL, timings
//
// The scenario module default-exports async ({ page, context, step, db, expect, unique, log, args, baseURL }).
import { createRequire } from "node:module";
import { pathToFileURL, fileURLToPath } from "node:url";
import { execFileSync } from "node:child_process";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";

const here = path.dirname(fileURLToPath(import.meta.url));
const skillDir = path.resolve(here, "..");
const repo = path.resolve(skillDir, "../../..");
const port = process.env.VERIFY_PORT || "18000";
const baseURL = `http://localhost:${port}`;
const runDir = path.join(repo, "tmp/verify-pagode/run", port);
const dbFile = path.join(runDir, "verify.db");
const evidenceRoot = path.join(repo, "tmp/verify-pagode/evidence");

function findPlaywright() {
  const found = [];
  try {
    found.push(createRequire(path.join(repo, "package.json")).resolve("playwright"));
  } catch {}
  const npxCache = path.join(os.homedir(), ".npm/_npx");
  if (fs.existsSync(npxCache)) {
    for (const d of fs.readdirSync(npxCache)) {
      const p = path.join(npxCache, d, "node_modules/playwright/index.js");
      if (fs.existsSync(p)) found.push(p);
    }
  }
  try {
    const g = execFileSync("npm", ["root", "-g"], { encoding: "utf8" }).trim();
    const p = path.join(g, "playwright/index.js");
    if (fs.existsSync(p)) found.push(p);
  } catch {}
  return found[0] ?? null;
}

async function loadPlaywright() {
  const entry = findPlaywright();
  if (!entry) {
    throw new Error(
      "playwright not found in <repo>/node_modules, ~/.npm/_npx, or the global npm root.\n" +
        "Install without touching package.json:  npm install --no-save playwright && npx playwright install chromium",
    );
  }
  const mod = await import(pathToFileURL(entry).href);
  const pw = mod.default ?? mod;
  const version = JSON.parse(fs.readFileSync(path.join(path.dirname(entry), "package.json"), "utf8")).version;
  return { chromium: pw.chromium, entry, version };
}

function parseArgs(argv) {
  const out = { scenario: null, name: null, headed: false, doctor: false, args: [] };
  let i = 0;
  while (i < argv.length) {
    const a = argv[i];
    if (a === "--doctor") out.doctor = true;
    else if (a === "--headed") out.headed = true;
    else if (a === "--name") out.name = argv[++i];
    else if (a === "--") { out.args = argv.slice(i + 1); break; }
    else if (!out.scenario) out.scenario = a;
    else out.args.push(a);
    i++;
  }
  return out;
}

const slug = (s) => s.toLowerCase().replace(/[^a-z0-9]+/g, "-").replace(/^-|-$/g, "");
const stamp = () => new Date().toISOString().replace(/[-:]/g, "").replace(/\..+/, "").replace("T", "-");

async function doctor() {
  const { chromium, entry, version } = await loadPlaywright();
  const browser = await chromium.launch({ headless: true });
  console.log(`playwright ${version} at ${entry}; chromium ${browser.version()}`);
  await browser.close();
}

async function runScenario(opts) {
  const scenarioPath = path.resolve(opts.scenario);
  const label = opts.name ?? path.basename(scenarioPath, ".mjs");
  const evidence = path.join(evidenceRoot, `${stamp()}-${slug(label)}`);
  fs.mkdirSync(evidence, { recursive: true });
  const logStream = fs.createWriteStream(path.join(evidence, "run.log"));
  const log = (line) => {
    const msg = `[${new Date().toISOString()}] ${line}`;
    logStream.write(msg + "\n");
    console.log(msg);
  };

  if (!fs.existsSync(runDir)) {
    throw new Error(`no instance state at ${runDir}; run scripts/launch.sh first`);
  }

  const { chromium, version } = await loadPlaywright();
  const browser = await chromium.launch({ headless: !opts.headed });
  const context = await browser.newContext({ baseURL, viewport: { width: 1280, height: 800 } });
  await context.tracing.start({ screenshots: true, snapshots: true });
  const page = await context.newPage();
  page.on("console", (m) => {
    if (m.type() === "error" || m.type() === "warning") log(`browser console.${m.type()}: ${m.text()}`);
  });
  page.on("pageerror", (e) => log(`browser pageerror: ${e.message}`));
  page.on("requestfailed", (r) => log(`request failed: ${r.method()} ${r.url()} ${r.failure()?.errorText ?? ""}`));
  page.on("response", (r) => {
    if (r.status() >= 400) log(`http ${r.status()} ${r.request().method()} ${r.url()}`);
  });

  const steps = [];
  let n = 0;
  const step = async (name, fn) => {
    n++;
    const id = `${String(n).padStart(2, "0")}-${slug(name)}`;
    log(`STEP ${id}`);
    const started = Date.now();
    try {
      const result = await fn();
      await page.screenshot({ path: path.join(evidence, `${id}.png`) });
      steps.push({ id, name, status: "ok", ms: Date.now() - started });
      return result;
    } catch (err) {
      await page.screenshot({ path: path.join(evidence, `${id}-FAILED.png`) }).catch(() => {});
      steps.push({ id, name, status: "failed", ms: Date.now() - started, error: String(err) });
      throw err;
    }
  };
  const db = (sql) => {
    const out = execFileSync("sqlite3", ["-json", dbFile, sql], { encoding: "utf8" }).trim();
    log(`db: ${sql} -> ${out || "[]"}`);
    return out ? JSON.parse(out) : [];
  };
  const expect = (cond, message) => {
    if (!cond) throw new Error(`expectation failed: ${message}`);
  };
  const unique = (prefix) => `${prefix}-${Date.now().toString(36)}${Math.random().toString(36).slice(2, 6)}`;

  const summary = { label, scenario: scenarioPath, baseURL, playwright: version, startedAt: new Date().toISOString(), status: "failed", steps, evidence };
  let exitCode = 1;
  try {
    const mod = await import(pathToFileURL(scenarioPath).href);
    await mod.default({ page, context, step, db, expect, unique, log, args: opts.args, baseURL });
    summary.status = "passed";
    exitCode = 0;
    log("RESULT passed");
  } catch (err) {
    summary.error = String(err?.stack ?? err);
    log(`RESULT failed: ${err?.stack ?? err}`);
  } finally {
    summary.finishedAt = new Date().toISOString();
    await context.tracing.stop({ path: path.join(evidence, "trace.zip") }).catch(() => {});
    await browser.close().catch(() => {});
    fs.writeFileSync(path.join(evidence, "summary.json"), JSON.stringify(summary, null, 2));
    log(`EVIDENCE ${evidence}`);
    logStream.end();
  }
  process.exit(exitCode);
}

const opts = parseArgs(process.argv.slice(2));
if (opts.doctor) {
  doctor().catch((e) => { console.error(String(e.message ?? e)); process.exit(1); });
} else if (!opts.scenario) {
  console.error("usage: node scripts/pw.mjs <scenario.mjs> [--name label] [--headed] [-- args]  |  --doctor");
  process.exit(2);
} else {
  runScenario(opts).catch((e) => { console.error(String(e.stack ?? e)); process.exit(1); });
}
