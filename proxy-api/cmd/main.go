package main

import (
	"fmt"
	"github.com/go-pkgz/routegroup"
	"log"
	"net/http"
	"proxy-api/internal/config"
	proxydelivery "proxy-api/internal/delivery/http/routes/proxy"
	"proxy-api/internal/service/proxy"
)

func main() {
	cfg := config.New()

	mux := http.NewServeMux()
	router := routegroup.New(mux)

	proxyService := proxy.New(proxy.Opts{Cfg: *cfg})
	proxydelivery.New(proxydelivery.Opts{Router: router, ProxyService: proxyService})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.AppPort), router))
}
