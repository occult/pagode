# Deployment Guide

This document describes how to deploy the Pagode application using nixpacks-compatible platforms like Railway, Render, or Heroku.

## Overview

Pagode uses a multi-language build process that combines a React/TypeScript frontend (built with Vite) with a Go backend. The deployment strategy embeds the built frontend assets directly into the Go binary using Go's `embed.FS`, creating a single self-contained executable.

## Build Process

The deployment follows these phases:

1. **Setup**: Install Node.js and Go
2. **Install**: Install npm dependencies 
3. **Frontend Build**: Build React/Vite assets to `dist/` directory
4. **Go Build**: Compile Go binary with embedded frontend assets
5. **Start**: Run the single Go binary

## Nixpacks Configuration

The build is configured via `nixpacks.toml`:

```toml
# Multi-language provider configuration for Go + Vite application  
providers = ["node"]

[variables]
NODE_ENV = "production"

[phases.setup]
nixPkgs = ["nodejs", "go"]

[phases.install]  
cmds = [
    "npm ci"
]

[phases.frontend-build]
cmds = ["npm run build"]
dependsOn = ["install"]

[phases.build]
cmds = ["go build -o app ./cmd/web"]
dependsOn = ["frontend-build"]

[start]
cmd = "./app"
```

### Key Configuration Details

- **Providers**: Uses only `"node"` provider to avoid conflicts with Go versions
- **nixPkgs**: Manually installs both `nodejs` and `go` packages
- **Dependencies**: Frontend build must complete before Go build
- **Output**: Single `app` binary containing embedded assets

## Frontend Configuration

### Vite Configuration (`vite.config.mts`)

```typescript
export default defineConfig({
  plugins: [
    laravel({
      input: ["resources/js/app.jsx", "resources/css/app.css"],
      publicDirectory: "dist",        // Changed from "public"
      buildDirectory: ".",           // Changed from "build"
      refresh: true,
    }),
    // ... other plugins
  ],
  build: {
    manifest: true,
    outDir: "dist",                  // Changed from "public/build"
    // ... other build options
  },
});
```

### Package.json Scripts

```json
{
  "scripts": {
    "build": "vite build",
    "dev": "vite"
  }
}
```

## Backend Configuration

### Asset Embedding (`assets/assets.go`)

```go
package assets

import "embed"

//go:embed dist
var Assets embed.FS
```

### Asset Serving (`pkg/handlers/build.go`)

```go
func (h *Build) Routes(g *echo.Group) {
    // Serve the embedded build directory
    distFS, err := fs.Sub(assets.Assets, "dist")
    if err != nil {
        panic(err)
    }
    fs := http.StripPrefix("/build/", http.FileServer(http.FS(distFS)))
    g.GET("/build/*", echo.WrapHandler(fs))
}
```

### Router Integration (`cmd/web/main.go`)

```go
import (
    "github.com/occult/pagode/assets"
    "github.com/occult/pagode/pkg/handlers"
)

func main() {
    c := services.NewContainer()
    
    // Build router with embedded assets
    if err := handlers.BuildRouter(c, assets.Assets); err != nil {
        fatal("failed to build the router", err)
    }
    
    // Start server...
}
```

## Directory Structure

```
pagode/
├── nixpacks.toml           # Deployment configuration
├── vite.config.mts         # Frontend build configuration
├── package.json            # Node.js dependencies and scripts
├── assets/
│   └── assets.go          # Go embed.FS for built assets
├── dist/                  # Vite build output (created during build)
│   ├── assets/
│   │   ├── app.*.css
│   │   └── app.*.js
│   └── .vite/
│       └── manifest.json
├── cmd/web/               # Go application entry point
├── pkg/handlers/          # HTTP handlers including asset serving
└── resources/js/          # React/TypeScript source code
```

## Deployment Platforms

### Railway

1. Connect your GitHub repository to Railway
2. Railway will automatically detect `nixpacks.toml`
3. Set environment variables as needed
4. Deploy

### Render

1. Create a new Web Service
2. Connect your repository
3. Render will auto-detect nixpacks configuration
4. Configure environment variables
5. Deploy

### Vercel/Netlify

These platforms may require additional configuration as they're optimized for frontend deployments. Consider using Railway or Render for full-stack Go applications.

## Environment Variables

Set these environment variables in your deployment platform:

```bash
# Required
PAGODA_APP_ENCRYPTION_KEY=your-32-char-encryption-key

# Database
PAGODA_DATABASE_CONNECTION=path-to-your-database

# Optional
PAGODA_HTTP_PORT=8000
PAGODA_APP_ENVIRONMENT=prod
```

## Local Testing

Test the nixpacks build locally:

```bash
# Install nixpacks
curl -sSL https://nixpacks.com/install.sh | bash

# Build the image
nixpacks build . --name pagode-test

# Run the container
docker run -p 8000:8000 pagode-test
```

## Build Assets

The build process creates these key outputs:

- **Frontend Assets**: `dist/assets/app.*.js` and `dist/assets/app.*.css`
- **Asset Manifest**: `dist/.vite/manifest.json` for cache-busting
- **Go Binary**: `app` executable with embedded assets
- **Docker Image**: Self-contained container ready for deployment

## Troubleshooting

### Build Failures

1. **"vite: not found"**: Ensure `npm ci` completes successfully before frontend build
2. **"no Go files in /app"**: Check that Go build command specifies `./cmd/web`
3. **Asset serving 404s**: Verify `assets/dist` directory exists and contains built files

### Runtime Issues

1. **Database connection errors**: Set `PAGODA_DATABASE_CONNECTION` environment variable
2. **Static file 404s**: Check that assets are properly embedded and routes are configured
3. **Port binding**: Ensure `PAGODA_HTTP_PORT` matches your platform's expected port

## Performance Considerations

- **Single Binary**: Eliminates need for separate static file serving
- **Embedded Assets**: Faster startup and simpler deployment
- **Build Caching**: nixpacks caches dependencies for faster subsequent builds
- **Asset Compression**: Vite automatically compresses assets during build

## Security

- Assets are embedded in the binary, reducing attack surface
- No separate static file server needed
- Environment-based configuration prevents secrets in code
- CSRF protection enabled by default for non-GET requests