package routes

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func StaticFileRouter(r chi.Router, path string) {

	v1 := fmt.Sprintf("/%s/*", path)
	pathPrefix := fmt.Sprintf("/%s/", path)

	r.Handle(v1, http.StripPrefix(pathPrefix, http.FileServer(http.Dir(path))))
}
