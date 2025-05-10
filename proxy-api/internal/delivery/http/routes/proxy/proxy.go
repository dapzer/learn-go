package proxy

import (
	"fmt"
	"github.com/go-pkgz/routegroup"
	"io"
	"net/http"
	"proxy-api/internal/service/proxy"
)

type Opts struct {
	Router       *routegroup.Bundle
	ProxyService *proxy.Service
}

func New(opts Opts) {
	proxyGroup := opts.Router.Mount("/proxy")
	h := handlers{service: opts.ProxyService}
	proxyGroup.HandleFunc("GET /content/{everything...}", h.getContent)
	proxyGroup.HandleFunc("GET /image/{everything...}", h.getImage)
}

type handlers struct {
	service *proxy.Service
}

func (s *handlers) getContent(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("everything")
	query := r.URL.Query()

	resp, err := s.service.GetResponse(path, query)
	if err != nil {
		http.Error(w, "Failed to get data from remote server.", http.StatusBadGateway)
		return
	}

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("%d: Failed to get data from remote server.", resp.StatusCode), http.StatusBadGateway)
		return
	}

	w.Header().Add("Content-Type", resp.Header.Get("Content-Type"))
	if _, err = io.Copy(w, resp.Body); err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}

func (s *handlers) getImage(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("everything")

	resp, err := s.service.GetImage(path)
	if err != nil {
		http.Error(w, "Failed to get data from remote server.", http.StatusBadGateway)
		return
	}

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("%d: Failed to get data from remote server.", resp.StatusCode), http.StatusBadGateway)
		return
	}

	if _, err = io.Copy(w, resp.Body); err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}
