# File upload

A signed-in user uploads a file from `/files` and sees it listed with size and modification time; the file is written to the instance's uploads directory.

## Sub-features

- `files-list` shows previously uploaded files under `Uploaded Files`.
- `files-upload` accepts one file and confirms with a toast naming it.
- `files-required` rejects an empty submission.

## How to get to it (user POV)

- Choose `Upload Files` in the sidebar, or open `/files`.

## Driving it with Playwright

Preconditions:

- Instance healthy; a signed-in user.
- A small local file to upload, created by the scenario in the evidence directory (for example `hello.txt`).

- **Open.** `page.goto("/files")`. `Uploaded Files` heading is visible; the list is empty on a fresh instance.
- **Choose a file.** `page.setInputFiles("#file", "<path>")`.
- **Submit.** Click the upload button (the only submit button on the page). Toast `<filename> was uploaded successfully.` appears and a row with the file name, size, and time is listed.
- **On disk.** `fs.existsSync("tmp/verify-pagode/run/18000/uploads/<filename>")` is true and the size matches. This is the `PAGODE_FILES_DIRECTORY` the instance was launched with.
- **Empty submit.** Submitting with no file shows `A file is required.`.
- **Proof.** Screenshot of the listed row plus the `ls -l` of the uploads directory in the report.

## Gotchas

- Files are stored by their original name; uploading the same name twice overwrites silently. Use `unique()` in the filename.
- The uploads directory is per instance and removed by cleanup; copy the file into the evidence directory if the bytes are part of the proof.
- There is no download route; the list is the only user-visible view.
