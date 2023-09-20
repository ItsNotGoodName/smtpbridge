// web is used by http/pages packages.
package web

import (
	"context"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/thejerf/suture/v4"
)

//go:generate pnpm install
//go:generate pnpm run build

func init() {
	// Fix mimetype on Windows when using http.FileServer
	mime.AddExtensionType(".js", "application/javascript")
}

// https://github.com/labstack/echo/blob/deb17d2388a74cd4133f46c2dedfb7601da1db0a/echo_fs.go#LL144C2-L144C2
func mustSubFS(currentFs fs.FS, fsRoot string) fs.FS {
	fsRoot = filepath.ToSlash(filepath.Clean(fsRoot)) // note: fs.FS operates only with slashes. `ToSlash` is necessary for Windows
	subFs, err := fs.Sub(currentFs, fsRoot)
	if err != nil {
		panic(err)
	}
	return subFs
}

type ViteHotReload struct {
}

func NewViteHotReload() ViteHotReload {
	return ViteHotReload{}
}

func (r ViteHotReload) String() string {
	return "web.ViteHotReload"
}

func (r ViteHotReload) Serve(ctx context.Context) error {
	reloadVite()
	return suture.ErrDoNotRestart
}

func CacheControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, "/assets") {
			// Files in assets have a hash associated with them
			w.Header().Set("Cache-Control", "max-age=31536000,immutable")
		} else {
			w.Header().Set("Cache-Control", "max-age=3600")
		}

		h.ServeHTTP(w, r)
	})
}
