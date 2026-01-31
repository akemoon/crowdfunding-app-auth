package api

import (
	"net/http"

	"github.com/akemoon/golib/myhttp/middleware"
)

type Router struct {
	mux *http.ServeMux
	mws []middleware.Midddleware
}

func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
		mws: []middleware.Midddleware{},
	}
}

func (r *Router) Use(mws ...middleware.Midddleware) *Router {
	r.mws = append(r.mws, mws...)
	return r
}

func (r *Router) Handle(pattern string, h http.Handler) {
	r.mux.Handle(pattern, r.wrap(h))
}

func (r *Router) HandleFunc(pattern string, hf http.HandlerFunc) {
	r.mux.HandleFunc(pattern, func(w http.ResponseWriter, req *http.Request) {
		r.wrap(http.HandlerFunc(hf)).ServeHTTP(w, req)
	})
}

func (r *Router) Handler() http.Handler {
	return r.mux
}

func (r *Router) wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		wrapped := h
		for i := len(r.mws) - 1; i >= 0; i-- {
			wrapped = r.mws[i](wrapped)
		}
		wrapped.ServeHTTP(w, req)
	})
}
