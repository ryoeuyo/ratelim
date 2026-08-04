package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ProxyHandler struct{}

func (h *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	target, err := url.Parse("http://localhost:8080")
	if err != nil {
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	r.Host = target.Host
	proxy.ServeHTTP(w, r)
}
