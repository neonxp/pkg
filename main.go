package main

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	templates = template.Must(template.ParseGlob("./tpl/*.gohtml"))
	packages  = Packages{}
)

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	cfg := &Config{}

	fp, err := os.Open("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewDecoder(fp).Decode(cfg); err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		templates.ExecuteTemplate(w, "index.gohtml", indexRenderContext{
			Title:    cfg.Title,
			Packages: cfg.Packages,
			Host:     cfg.Host,
			Doc:      "//pkg.go.dev/",
		})
	})

	r.Get("/{pkg}", func(w http.ResponseWriter, r *http.Request) {
		pkgID := chi.URLParam(r, "pkg")
		pkg, ok := (*cfg.Packages)[pkgID]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		templates.ExecuteTemplate(w, "package.gohtml", pageRenderContext{
			Title:   cfg.Title,
			Package: &pkg,
			Host:    cfg.Host,
			Doc:     "//pkg.go.dev/",
		})
	})

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "static"))
	FileServer(r, "/static", filesDir)

	srv := &http.Server{
		Handler: r,
		Addr:    addr,
	}
	go func() {
		<-ctx.Done()
		srv.Shutdown(context.Background())
	}()
	log.Println("Start at", addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
