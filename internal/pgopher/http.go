package pgopher

import (
	"net/http"
	"path/filepath"
	"strings"
)

func readinessProbe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func livenessProbe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	switch s.cfg.Sink.Type {
	case "file":
		profile := strings.TrimPrefix(r.RequestURI, "/api/v1/profile/")
		http.ServeFile(w, r, filepath.Join(s.cfg.Sink.FileSinkOptions.Folder, profile))
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
