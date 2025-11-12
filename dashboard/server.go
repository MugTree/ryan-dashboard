package dashboard

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/MugTree/ryan_dashboard/dashboard/public"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	env    *EnvVars
	db     *sqlx.DB
	mux    chi.Router
	server *http.Server
}

type EnvVars struct {
	IsProd        bool
	LogLocation   string
	SensorAddress string
}

var ServerEnv EnvVars

func NewServer(db *sqlx.DB, address string, env *EnvVars) *Server {

	mux := chi.NewMux()

	return &Server{
		env: env,
		db:  db,
		mux: mux,
		server: &http.Server{
			Addr:              address,
			Handler:           mux,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       5 * time.Second,
		},
	}
}

func (s *Server) Start() error {

	slog.Info("Starting webserver", "address", s.server.Addr)

	s.mux.Group(func(router chi.Router) {

		// images, css, js
		router.Group(func(r chi.Router) {
			r.Use(versionedAssetsMiddleware)
			setStaticAssests(r, s.env.IsProd)
		})

		// api calls
		// router.Group(func(r chi.Router) {
		// 	r.Use(headerAuthMiddleware(s.env.JsonApiKey))
		// 	apiRoutes(r, s.db)
		// })

		// admin, dasboard - login required
		// router.Group(func(r chi.Router) {
		// 	r.Use(cookieAuthMiddleware(true))
		// 	adminRoutes(r, s.db, s.env)
		// })

		// front end admin tasks where we require a logged in user
		// router.Group(func(r chi.Router) {
		// 	r.Use(cookieAuthMiddleware(true))
		// 	webAdminRoutes(r, s.db)
		// })

		// front end no logging required
		router.Group(func(r chi.Router) {
			webRoutes(r, s.db, s.env)
		})

		router.NotFound(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			w.Write([]byte("route does not exist"))
			slog.Info("route does not exist" + r.URL.Path)
		})

	})

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Stop() error {
	slog.Info("Stopping the http server")
	// ensure shutdown doesnt hang
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	slog.Info("Stopped the http server")
	return nil
}

func setStaticAssests(r chi.Router, isProd bool) {

	staticFolderProduction := func(n func(h http.Handler) http.Handler) http.Handler {
		return n(http.FileServer(http.FS(public.AssetsFS)))
	}

	staticFolderDevelopment := func(n func(h http.Handler) http.Handler) http.Handler {
		return n(http.FileServer(http.Dir(AssetsPathDev)))
	}

	stopDirectoryListing := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "/") {
				http.NotFound(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	var staticHandler http.Handler
	if isProd {
		staticHandler = staticFolderProduction(stopDirectoryListing)
	} else {
		staticHandler = staticFolderDevelopment(stopDirectoryListing)
	}
	r.Get(`/{:[^.]+\.[^.]+}`, staticHandler.ServeHTTP)
	r.Get(`/{:img|js|css}/*`, staticHandler.ServeHTTP)
}

// versionedAssetMatcher matches versioned assets like "app.abc123.js".
// See https://regex101.com/r/bGfflm/latest
var versionedAssetMatcher = regexp.MustCompile(`^(?P<name>[^.]+)\.[a-z0-9]+(?P<extension>\.[a-z0-9]+)$`)

// versionedAssetsMiddleware is Middleware to help serve versioned assets without the version.
// It basically strips the version from the asset path and forwards the request, probably to a static file handler.
func versionedAssetsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if versionedAssetMatcher.MatchString(r.URL.Path) {
			r.URL.Path = versionedAssetMatcher.ReplaceAllString(r.URL.Path, `$1$2`)
		}

		next.ServeHTTP(w, r)
	})
}

// ugly but useful
// func headerAuthMiddleware(apiKey string) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			key := r.Header.Get("X-API-KEY")

// 			if key != apiKey {
// 				slog.Error("api login error", "error", fmt.Sprintf("incorrect api key val: %v", key))
// 				http.Error(w, "unauthorized", http.StatusUnauthorized)
// 				return
// 			}

// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }

// func cookieAuthMiddleware(strict bool) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 			auth, err := authenticateUser(r)

// 			if err != nil {
// 				if err.Error() == "session expired" {
// 					http.Redirect(w, r, "/", http.StatusUnauthorized)
// 					return
// 				} else {
// 					slog.Error("App error", "error", err)
// 					http.Error(w, http.StatusText(http.StatusInternalServerError), 500)
// 					return
// 				}
// 			}

// 			const url = "/login"

// 			if strict && !auth.Check() && r.URL.Path != url {
// 				http.Redirect(w, r, url, http.StatusSeeOther)
// 				return
// 			}

// 			ctx := setUser(r.Context(), &auth) //context.WithValue(r.Context(), AuthKey{}, auth)

// 			next.ServeHTTP(w, r.WithContext(ctx))
// 		})
// 	}
// }

func logAndError(w http.ResponseWriter, err error) {
	slog.Error("App error", "error", err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), 500)
}

// func writeJSON(w http.ResponseWriter, status int, v any) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(status)
// 	err := json.NewEncoder(w).Encode(v)
// 	if err != nil {
// 		slog.Error("json encode error", "error", err, "item", v)
// 		http.Error(w, http.StatusText(http.StatusInternalServerError), 500)
// 	}
// }

var hashOnce sync.Once
var appCSSPath, bulmaCSSPath, dataStarPath string

func getHashedPath(path string, assetsPath string, isLive bool) string {

	//   fmt.Printf("path: %s\n", path)
	//   fmt.Printf("assetsPath: %s\n", assetsPath)
	//   fmt.Printf("isLive: %v\n", isLive)

	externalPath := strings.TrimPrefix(path, assetsPath)

	// fmt.Printf("externalPath: %s\n", externalPath)
	// fmt.Println("-------------------------------")

	ext := filepath.Ext(path)

	if ext == "" {
		panic("no extension found")
	}

	var err error

	if isLive {
		_, err = public.AssetsFS.ReadFile(path)

	} else {
		_, err = os.ReadFile(path)
	}

	if err != nil {
		fmt.Printf("error getting hold of assets: %v", err)
	}

	return fmt.Sprintf("/%v.%x%v", strings.TrimSuffix(externalPath, ext), sha256.Sum256([]byte(path)), ext)
}

func setAssetPaths(isProd bool) {
	hashOnce.Do(func() {

		assetsPath := ""

		if isProd {
			assetsPath = AssetsPathProd
		} else {
			assetsPath = AssetsPathDev
		}

		appCSSPath = getHashedPath(assetsPath+"css/main.css", assetsPath, isProd)
		bulmaCSSPath = getHashedPath(assetsPath+"css/bulma.css", assetsPath, isProd)
		dataStarPath = getHashedPath(assetsPath+"js/datastar.js", assetsPath, isProd)
	})
}

func formatError(prefix string, r *http.Request, err error) error {
	return fmt.Errorf(prefix+" - url:%v error:%v", chi.RouteContext(r.Context()).RoutePattern(), err)
}
