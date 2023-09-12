// web is used by http/pages packages.
package web

import (
	"context"
	"io/fs"
	"mime"
	"path/filepath"

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

type Refresher struct {
}

func NewRefresher() Refresher {
	return Refresher{}
}

func (r Refresher) String() string {
	return "web.Refresher"
}

func (r Refresher) Serve(ctx context.Context) error {
	reloadVite()
	return suture.ErrDoNotRestart
}
