package web

import (
	"embed"
	log "github.com/sirupsen/logrus"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

//go:embed *
var content embed.FS

func StaticContentHandler(w http.ResponseWriter, r *http.Request) {
	ext := filepath.Ext(r.URL.Path)
	file, err := content.ReadFile(strings.TrimPrefix(r.URL.Path, "/"))
	if err != nil {
		log.Errorf("Unable to send file to client: %s", err)
		w.WriteHeader(404)
	} else {
		w.Header().Set("Content-Type", mime.TypeByExtension(ext))
		w.Write(file)
	}
	log.Infof("serving content: %s", r.URL)
}
