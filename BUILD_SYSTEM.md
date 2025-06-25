# Build System and Asset Management

This document explains how Pagode's build system works, including asset compilation, manifest generation, and static file serving across different environments.

## Overview

Pagode uses a sophisticated build system that handles frontend assets differently based on the environment:
- **Development**: Assets served directly from Vite dev server with hot reloading
- **Production**: Pre-built assets served from static directories with cache busting

## Build Process (nixpacks.toml)

### 1. Frontend Asset Building
```bash
npx vite build
```
- Compiles TypeScript/React code
- Generates hashed filenames for cache busting (e.g., `app-abc123.js`)
- Creates build output in `public/build/`
- Generates Vite manifest at `public/build/.vite/manifest.json`

### 2. File Organization in Container
```bash
# Copy all public assets to container static directory
cp -r public/* /app/static/

# Copy assets to files directory (matches /files/ URL prefix)
mkdir -p /app/files
cp -r public/build/assets /app/files/
```

### 3. Binary Compilation
```bash
CGO_ENABLED=1 go build -o pagode ./cmd/web
```

## Asset Resolution System

### Development Mode (`container.go:289-323`)

When `public/hot` file exists (Vite dev server running):
- Assets served directly from Vite dev server (typically `localhost:5173`)
- Hot module replacement enabled
- React refresh runtime injected
- No manifest file needed

### Production Mode (`container.go:325-377`)

When no hot file detected:
1. **Manifest Loading**: Attempts to load manifest from:
   - `/app/static/build/manifest.json` (Docker container)
   - `public/build/manifest.json` (local development)

2. **Manifest Normalization**: If main manifest missing, moves Vite manifest:
   ```go
   // Move from .vite/manifest.json to manifest.json
   os.Rename(viteManifestPath, manifestPath)
   ```

3. **Asset URL Generation**: Creates `vite()` template function that:
   - Maps entry points to actual file paths
   - Generates HTML `<script>` and `<link>` tags
   - Uses `/files/` URL prefix consistently

### Vite Template Function (`container.go:383-469`)

The `vite()` function handles asset resolution:

```go
// Input: "resources/js/app.tsx"
// Manifest lookup: finds "assets/app-abc123.js"
// Output: <script type="module" src="/files/assets/app-abc123.js"></script>
```

Key features:
- Reads manifest file to map entry points to built assets
- Handles both JavaScript and CSS files
- Generates proper HTML tags with integrity and crossorigin attributes
- Logs extensively for debugging asset loading issues

## Static File Serving (`router.go:23-41`)

### Route Configuration by Environment

**Local Development**:
```go
// Serve public/build at /files/ URL
c.Web.Static("/files", filepath.Join(services.ProjectRoot(), "public/build"))
```

**Production**:
```go
// Serve /app/files at /files/ URL  
c.Web.Static("/files", "/app/files")
```

**All Environments**:
```go
// General static files from static/ directory
c.Web.Static(config.StaticPrefix, config.StaticDir)
```

### Cache Control

All static routes include cache control middleware:
```go
middleware.CacheControl(c.Config.Cache.Expiration.StaticFile)
```

## File Structure

```
pagode/
├── public/build/           # Vite build output (local)
│   ├── assets/            # Compiled JS/CSS with hashes
│   ├── manifest.json      # Asset mapping (created from .vite/)
│   └── .vite/
│       └── manifest.json  # Original Vite manifest
├── static/                # General static files
└── /app/                  # Container paths
    ├── static/            # Copy of public/ in container
    └── files/             # Assets accessible at /files/ URL
```

## URL Mapping

| Environment | Entry Point | Manifest Lookup | Final URL |
|-------------|-------------|-----------------|-----------|
| Development (hot) | `resources/js/app.tsx` | Vite dev server | `http://localhost:5173/resources/js/app.tsx` |
| Development (build) | `resources/js/app.tsx` | `public/build/manifest.json` | `/files/assets/app-abc123.js` |
| Production | `resources/js/app.tsx` | `/app/static/build/manifest.json` | `/files/assets/app-abc123.js` |

## Key Design Decisions

1. **Consistent URL Prefix**: `/files/` prefix works across all environments
2. **Manifest Normalization**: Automatically handles Vite's nested manifest structure
3. **Environment Detection**: Uses file existence checks rather than explicit config
4. **Graceful Fallbacks**: Extensive logging and error handling for asset loading
5. **Cache Busting**: Leverages Vite's automatic hash-based filenames

## Debugging Asset Issues

The system includes extensive logging. Check logs for:
- Manifest file locations and existence
- Asset mapping from entry points to files  
- Generated HTML tag content
- Static directory listings
- File copy operations during build

Common issues:
- Missing manifest files (check build completed successfully)
- Wrong file paths (verify container file structure)
- Cache issues (check cache control headers)
- Development vs production asset serving mismatches