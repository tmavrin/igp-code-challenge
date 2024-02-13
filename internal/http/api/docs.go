package api

import (
	"net/http"
	"path"
	"strings"
)

//go:generate go run github.com/swaggo/swag/cmd/swag fmt --dir ../../../.
//go:generate go run github.com/swaggo/swag/cmd/swag init --parseDependency --dir ../../../cmd/signup,signup  --output ../../../docs --outputTypes json,yaml

func DocsHandler(base string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file := ""
		if file == "" || strings.HasSuffix(file, "/") {
			file += "index.html"
		} else if !strings.HasSuffix(file, "yaml") &&
			!strings.HasSuffix(file, "html") &&
			!strings.HasSuffix(file, "js") &&
			!strings.HasSuffix(file, "css") {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.ServeFile(w, r, path.Join(base, file))
	}
}
