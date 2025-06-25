package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mikestefanello/backlite"
	"github.com/occult/pagode/config"
	"github.com/occult/pagode/ent"
	"github.com/occult/pagode/pkg/log"
	inertia "github.com/romsar/gonertia/v2"
	"github.com/spf13/afero"

	// Required by ent.
	_ "github.com/occult/pagode/ent/runtime"
)

// Container contains all services used by the application and provides an easy way to handle dependency
// injection including within tests.
type Container struct {
	// Validator stores a validator
	Validator *Validator

	// Web stores the web framework.
	Web *echo.Echo

	// Config stores the application configuration.
	Config *config.Config

	// Cache contains the cache client.
	Cache *CacheClient

	// Database stores the connection to the database.
	Database *sql.DB

	// Files stores the file system.
	Files afero.Fs

	// ORM stores a client to the ORM.
	ORM *ent.Client

	// Graph is the entity graph defined by your Ent schema.
	Graph *gen.Graph

	// Mail stores an email sending client.
	Mail *MailClient

	// Auth stores an authentication client.
	Auth *AuthClient

	// Tasks stores the task client.
	Tasks *backlite.Client

	// Inertia for React
	Inertia *inertia.Inertia
}

// NewContainer creates and initializes a new Container.
func NewContainer() *Container {
	c := new(Container)
	c.initConfig()
	c.initValidator()
	c.initWeb()
	c.initCache()
	c.initDatabase()
	c.initFiles()
	c.initORM()
	c.initAuth()
	c.initMail()
	c.initTasks()
	c.initInertia()
	return c
}

// Shutdown gracefully shuts the Container down and disconnects all connections.
func (c *Container) Shutdown() error {
	// Shutdown the web server.
	webCtx, webCancel := context.WithTimeout(context.Background(), c.Config.HTTP.ShutdownTimeout)
	defer webCancel()
	if err := c.Web.Shutdown(webCtx); err != nil {
		return err
	}

	// Shutdown the task runner.
	taskCtx, taskCancel := context.WithTimeout(context.Background(), c.Config.Tasks.ShutdownTimeout)
	defer taskCancel()
	c.Tasks.Stop(taskCtx)

	// Shutdown the ORM.
	if err := c.ORM.Close(); err != nil {
		return err
	}

	// Shutdown the database.
	if err := c.Database.Close(); err != nil {
		return err
	}

	// Shutdown the cache.
	c.Cache.Close()

	return nil
}

// initConfig initializes configuration.
func (c *Container) initConfig() {
	cfg, err := config.GetConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	c.Config = &cfg

	// Configure logging.
	switch cfg.App.Environment {
	case config.EnvProduction:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	default:
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
}

// initValidator initializes the validator.
func (c *Container) initValidator() {
	c.Validator = NewValidator()
}

// initWeb initializes the web framework.
func (c *Container) initWeb() {
	c.Web = echo.New()
	c.Web.HideBanner = true
	c.Web.Validator = c.Validator
}

// initCache initializes the cache.
func (c *Container) initCache() {
	store, err := newInMemoryCache(c.Config.Cache.Capacity)
	if err != nil {
		panic(err)
	}

	c.Cache = NewCacheClient(store)
}

// initDatabase initializes the database.
func (c *Container) initDatabase() {
	var err error
	var connection string

	switch c.Config.App.Environment {
	case config.EnvTest:
		// TODO: Drop/recreate the DB, if this isn't in memory?
		connection = c.Config.Database.TestConnection
	default:
		connection = c.Config.Database.Connection
	}

	c.Database, err = openDB(c.Config.Database.Driver, connection)
	if err != nil {
		panic(err)
	}
}

// initFiles initializes the file system.
func (c *Container) initFiles() {
	// Use in-memory storage for tests.
	if c.Config.App.Environment == config.EnvTest {
		c.Files = afero.NewMemMapFs()
		return
	}

	fs := afero.NewOsFs()
	if err := fs.MkdirAll(c.Config.Files.Directory, 0755); err != nil {
		panic(err)
	}
	c.Files = afero.NewBasePathFs(fs, c.Config.Files.Directory)
}

// initORM initializes the ORM.
func (c *Container) initORM() {
	drv := entsql.OpenDB(c.Config.Database.Driver, c.Database)
	c.ORM = ent.NewClient(ent.Driver(drv))

	// Run the auto migration tool.
	if err := c.ORM.Schema.Create(context.Background()); err != nil {
		panic(err)
	}

	// Load the graph.
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	p := filepath.Join(filepath.Dir(d), "../ent/schema")
	g, err := entc.LoadGraph(p, &gen.Config{})
	if err != nil {
		panic(err)
	}
	c.Graph = g
}

// initAuth initializes the authentication client.
func (c *Container) initAuth() {
	c.Auth = NewAuthClient(c.Config, c.ORM)
}

// initMail initialize the mail client.
func (c *Container) initMail() {
	var err error
	c.Mail, err = NewMailClient(c.Config)
	if err != nil {
		panic(fmt.Sprintf("failed to create mail client: %v", err))
	}
}

// initTasks initializes the task client.
func (c *Container) initTasks() {
	var err error
	// You could use a separate database for tasks, if you'd like. but using one
	// makes transaction support easier.
	c.Tasks, err = backlite.NewClient(backlite.ClientConfig{
		DB:              c.Database,
		Logger:          log.Default(),
		NumWorkers:      c.Config.Tasks.Goroutines,
		ReleaseAfter:    c.Config.Tasks.ReleaseAfter,
		CleanupInterval: c.Config.Tasks.CleanupInterval,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to create task client: %v", err))
	}

	if err = c.Tasks.Install(); err != nil {
		panic(fmt.Sprintf("failed to install task schema: %v", err))
	}
}

func ProjectRoot() string {
	currentDir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		_, err := os.ReadFile(filepath.Join(currentDir, "go.mod"))
		if os.IsNotExist(err) {
			if currentDir == filepath.Dir(currentDir) {
				return ""
			}
			currentDir = filepath.Dir(currentDir)
			continue
		} else if err != nil {
			return ""
		}
		break
	}
	return currentDir
}

func (c *Container) getInertia() *inertia.Inertia {
	rootDir := ProjectRoot()
	viteHotFile := filepath.Join(rootDir, "public", "hot")
	rootViewFile := filepath.Join(rootDir, "resources", "views", "root.html")
	
	// Use different manifest paths based on environment
	var manifestPath, viteManifestPath string
	
	// Check if we're in Docker environment
	if _, err := os.Stat("/app/static/build/manifest.json"); err == nil {
		log.Default().Info("Using Docker container paths for manifest")
		manifestPath = "/app/static/build/manifest.json"
		viteManifestPath = "/app/static/build/.vite/manifest.json"
	} else {
		manifestPath = filepath.Join(rootDir, "public", "build", "manifest.json")
		viteManifestPath = filepath.Join(rootDir, "public", "build", ".vite", "manifest.json")
	}

	// check if laravel-vite-plugin is running in dev mode (it puts a "hot" file in the public folder)
	url, err := viteHotFileUrl(viteHotFile)
	if err != nil {
		panic(err)
	}
	if url != "" {
		i, err := inertia.NewFromFile(
			rootViewFile,
		)
		if err != nil {
			panic(err)
		}

		i.ShareTemplateFunc("getEnvironment", func() string {
			return string(c.Config.App.Environment)
		})

		i.ShareTemplateFunc("vite", func(entry string) (template.HTML, error) {
			if entry != "" && !strings.HasPrefix(entry, "/") {
				entry = "/" + entry
			}
			htmlTag := fmt.Sprintf(`<script type="module" src="%s%s"></script>`, url, entry)
			return template.HTML(htmlTag), nil
		})
		// Always define viteReactRefresh, but return empty content in production
		if c.Config.App.Environment == "local" {
			i.ShareTemplateFunc("viteReactRefresh", viteReactRefresh(url))
		} else {
			i.ShareTemplateFunc("viteReactRefresh", func() (template.HTML, error) {
				return template.HTML(""), nil
			})
		}

		return i
	}

	// laravel-vite-plugin not running in dev mode, use build manifest file
	// Add debug logging
	log.Default().Info("Production asset loading", "environment", c.Config.App.Environment, "manifestPath", manifestPath, "viteManifestPath", viteManifestPath)
	
	// check if the manifest file exists, if not, rename it
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		log.Default().Info("Main manifest not found, checking alternate location", "manifestPath", manifestPath)
		
		// Check if vite manifest exists
		if _, viteErr := os.Stat(viteManifestPath); viteErr != nil {
			log.Default().Error("Vite manifest also not found", "viteManifestPath", viteManifestPath, "error", viteErr)
		}
		
		// move the manifest from ./public/build/.vite/manifest.json to ./public/build/manifest.json
		// so that the vite function can find it
		err := os.Rename(viteManifestPath, manifestPath)
		if err != nil {
			log.Default().Error("Failed to move manifest", "error", err)
			return nil
		}
		log.Default().Info("Successfully moved manifest", "from", viteManifestPath, "to", manifestPath)
	} else {
		log.Default().Info("Manifest found at expected location", "manifestPath", manifestPath)
	}

	i, err := inertia.NewFromFile(
		rootViewFile,
		inertia.WithVersionFromFile(manifestPath),
	)
	if err != nil {
		panic(err)
	}

	// Share environment with the template
	i.ShareTemplateFunc("getEnvironment", func() string {
		return string(c.Config.App.Environment)
	})

	// Always use the standard path `/files/` as defined in config.StaticPrefix
	// This works in both local and Docker environments since the static files server
	// is configured to serve the `static` directory under the `/files` URL prefix
	i.ShareTemplateFunc("vite", vite(manifestPath, "/files/"))
	// Always define viteReactRefresh, but return empty content in production
	if c.Config.App.Environment == "local" {
		i.ShareTemplateFunc("viteReactRefresh", viteReactRefresh(url))
	} else {
		i.ShareTemplateFunc("viteReactRefresh", func() (template.HTML, error) {
			return template.HTML(""), nil
		})
	}

	return i
}

func (c *Container) initInertia() {
	c.Inertia = c.getInertia()
}

func vite(manifestPath, buildDir string) func(path string) (template.HTML, error) {
	// Add more detailed logging for debugging paths
	log.Default().Info("Vite manifest configuration", "manifestPath", manifestPath, "buildDir", buildDir)
	
	// Check if the manifest path exists
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		log.Default().Error("Manifest file does not exist", "manifestPath", manifestPath)
	} else {
		log.Default().Info("Manifest file exists", "manifestPath", manifestPath)
	}
	
	// Check related directories to understand the filesystem layout
	staticDir := filepath.Join(ProjectRoot(), "static")
	log.Default().Info("Static directory details", "staticDir", staticDir)
	
	// List static directory
	files, err := os.ReadDir(staticDir)
	if err != nil {
		log.Default().Error("Cannot read static directory", "error", err)
	} else {
		for _, file := range files {
			log.Default().Info("Static dir file", "name", file.Name(), "isDir", file.IsDir())
		}
	}
	
	// Check the public/build directory too
	publicBuildDir := filepath.Join(ProjectRoot(), "static", "public", "build")
	log.Default().Info("Public build directory", "publicBuildDir", publicBuildDir)
	
	// Try to read the manifest
	f, err := os.Open(manifestPath)
	if err != nil {
		log.Default().Error("cannot open provided vite manifest file", "error", err)
		panic(err)
	}
	defer f.Close()

	viteAssets := make(map[string]*struct {
		File   string `json:"file"`
		Source string `json:"src"`
		Css    []string `json:"css,omitempty"`
	})
	err = json.NewDecoder(f).Decode(&viteAssets)
	
	// Debug - print content of viteAssets
	log.Default().Info("Available assets in manifest:", "count", len(viteAssets))
	for k, v := range viteAssets {
		log.Default().Info("Asset entry", "path", k, "file", v.File, "cssCount", len(v.Css))
	}

	if err != nil {
		log.Default().Error("cannot unmarshal vite manifest file to json", "error", err)
		panic(err)
	}

	return func(p string) (template.HTML, error) {
		log.Default().Info("Vite template function called", "path", p, "buildDir", buildDir)
		
		if val, ok := viteAssets[p]; ok {
			log.Default().Info("Found asset in manifest", "path", p, "file", val.File)
			
			// Build HTML tags based on file type
			var tags strings.Builder
			
			// Add the main JS file
			jsPath := path.Join(buildDir, val.File)
			tags.WriteString(fmt.Sprintf("<script type=\"module\" crossorigin src=\"%s\"></script>", jsPath))
			log.Default().Info("Added JS script tag", "src", jsPath)
			
			// Add CSS files if any
			for _, cssFile := range val.Css {
				cssSrc := path.Join(buildDir, cssFile)
				tags.WriteString(fmt.Sprintf("\n<link rel=\"stylesheet\" href=\"%s\">", cssSrc))
				log.Default().Info("Added CSS link tag", "href", cssSrc)
			}
			
			// Get the final HTML output
			tagsHTML := tags.String()
			log.Default().Info("Generated HTML tags", "html", tagsHTML)
			return template.HTML(tagsHTML), nil
		}
		
		// Asset not found in manifest
		log.Default().Error("Asset not found in manifest", "path", p)
		return "", fmt.Errorf("asset %q not found in vite manifest", p)
	}
}

// openDB opens a database connection.
func openDB(driver, connection string) (*sql.DB, error) {
	if driver == "sqlite3" {
		// Helper to automatically create the directories that the specified sqlite file
		// should reside in, if one.
		d := strings.Split(connection, "/")
		if len(d) > 1 {
			dirpath := strings.Join(d[:len(d)-1], "/")

			if err := os.MkdirAll(dirpath, 0755); err != nil {
				return nil, err
			}
		}

		// Replace any random value placeholder, which is often used for in-memory test databases.
		connection = strings.Replace(connection, "$RAND", fmt.Sprint(rand.Int()), 1)
	}

	return sql.Open(driver, connection)
}

// viteHotFileUrl Get the vite hot file url
func viteHotFileUrl(viteHotFile string) (string, error) {
	_, err := os.Stat(viteHotFile)
	if err != nil {
		return "", nil
	}
	content, err := os.ReadFile(viteHotFile)
	if err != nil {
		return "", err
	}
	url := strings.TrimSpace(string(content))
	// Instead of conditionals, just use unconditional string replacements
	// to fix the linting issue (this is local development only)
	url = strings.Replace(url, "http://localhost", "", 1)
	url = strings.Replace(url, "http://127.0.0.1", "", 1)
	if url == "" {
		url = "//localhost:1323"
	}
	return url, nil
}

// viteReactRefresh Generate React refresh runtime script
func viteReactRefresh(url string) func() (template.HTML, error) {
	// Use unconditional strings.Replace to fix linting warning
	url = strings.Replace(url, "http:", "", 1)
	url = strings.Replace(url, "https:", "", 1)
	return func() (template.HTML, error) {
		if url == "" {
			return "", nil
		}
		script := fmt.Sprintf(`
<script type="module">
    import RefreshRuntime from '%s/@react-refresh'
    RefreshRuntime.injectIntoGlobalHook(window)
    window.$RefreshReg$ = () => {}
    window.$RefreshSig$ = () => (type) => type
    window.__vite_plugin_react_preamble_installed__ = true
</script>`, url)

		return template.HTML(script), nil
	}
}
