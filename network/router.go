package network

import "github.com/bonzonkim/url-shortener/handlers"

type Router struct {
	network *Network
	handlers *handlers.Handlers
}

func NewRouter(n *Network, h *handlers.Handlers) *Router {
	r := &Router{
		network:	n,
		handlers:	h,
	}

	n.engine.GET("/metrics", r.handlers.GetTopDomains)
	n.engine.POST("/shorten", r.handlers.ShortenURL)
	n.engine.GET("/:shortenURL", r.handlers.RedirectURL)

	return r
}
