package main

import (
	"fmt"
	"github.com/go-pkgz/routegroup"
	"io"
	"log"
	"net/http"
	"net/url"
)

type handlers struct {
	cfg Config
}

func buildUrl(baseUrl string, path string, query *url.Values) string {
	u, _ := url.Parse(baseUrl)
	u = u.JoinPath(path)
	if query != nil {
		u.RawQuery = query.Encode()
	}

	return u.String()
}

func (c *handlers) getResponse(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("everything")
	query := r.URL.Query()
	query.Add("api_key", c.cfg.TmdbApiKey)

	u := buildUrl(c.cfg.TmdbApiUrl, path, &query)

	resp, err := http.Get(u)
	if err != nil {
		http.Error(w, "Failed to get data from remote server.", http.StatusBadGateway)
		return
	}

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("%d: Failed to get data from remote server.", resp.StatusCode), http.StatusBadGateway)
		return
	}

	w.Header().Add("Content-Type", resp.Header.Get("Content-Type"))
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}

func main() {
	config := New()

	mux := http.NewServeMux()
	router := routegroup.New(mux)
	proxyGroup := router.Mount("/proxy")

	h := handlers{cfg: *config}

	proxyGroup.HandleFunc("GET /content/{everything...}", h.getResponse)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.AppPort), router))
}
