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
	var manifestPath string
	
	if _, err := os.Stat("/app/static/build/manifest.json"); err == nil {
		manifestPath = "/app/static/build/manifest.json"
	} else {
		manifestPath = filepath.Join(rootDir, "public", "build", "manifest.json")
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

	// Try to move Vite manifest if main manifest doesn't exist
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		viteManifestPath := filepath.Join(filepath.Dir(manifestPath), ".vite", "manifest.json")
		if err := os.Rename(viteManifestPath, manifestPath); err != nil {
			panic(fmt.Sprintf("manifest file not found at %s and failed to move from %s: %v", manifestPath, viteManifestPath, err))
		}
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
	f, err := os.Open(manifestPath)
	if err != nil {
		panic(fmt.Errorf("cannot open vite manifest file at %s: %w", manifestPath, err))
	}
	defer f.Close()

	var viteAssets map[string]struct {
		File string   `json:"file"`
		Css  []string `json:"css,omitempty"`
	}
	
	if err := json.NewDecoder(f).Decode(&viteAssets); err != nil {
		panic(fmt.Errorf("cannot decode vite manifest: %w", err))
	}

	return func(p string) (template.HTML, error) {
		asset, ok := viteAssets[p]
		if !ok {
			return "", fmt.Errorf("asset %q not found in vite manifest", p)
		}
		
		var tags strings.Builder
		
		// Add the main JS file
		tags.WriteString(fmt.Sprintf(`<script type="module" crossorigin src="%s"></script>`, 
			path.Join(buildDir, asset.File)))
		
		// Add CSS files if any
		for _, cssFile := range asset.Css {
			tags.WriteString(fmt.Sprintf(`<link rel="stylesheet" href="%s">`, 
				path.Join(buildDir, cssFile)))
		}
		
		return template.HTML(tags.String()), nil
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
	// Return the full URL as-is for hot reloading
	return url, nil
}

// viteReactRefresh Generate React refresh runtime script
func viteReactRefresh(url string) func() (template.HTML, error) {
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
