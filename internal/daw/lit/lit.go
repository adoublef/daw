//go:generate npm install
//go:generate npm run build
package lit

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/benbjohnson/hashfs"
)

const (
	prefix = "/daw/assets"
)

//go:embed all:dist
var fsys embed.FS
var sub, _ = fs.Sub(fsys, "dist")
var hashFS = hashfs.NewFS(sub)

var FuncMap = template.FuncMap{
	"daw": func(s string) string {
		return filepath.Join(prefix, hashFS.HashName(s))
	},
}

var Handler = http.StripPrefix(prefix, hashfs.FileServer(hashFS))
