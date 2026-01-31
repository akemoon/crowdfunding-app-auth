package api

import (
	"net/http"

	"github.com/akemoon/crowdfunding-app-auth/api/handler"
	"github.com/akemoon/crowdfunding-app-auth/service/auth"
	"github.com/akemoon/crowdfunding-app-auth/service/token"
	"github.com/akemoon/golib/myhttp/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	r *Router
}

func NewServer() *Server {
	return &Server{
		r: NewRouter().Use(
			middleware.BaseMetrics(),
		),
	}
}

func (s *Server) AddAuthHandlers(svc *auth.Service) {
	s.r.HandleFunc("POST /signup", handler.SignUp(svc))
	s.r.HandleFunc("POST /signin", handler.SignIn(svc))
}

func (s *Server) AddTokenHandlers(svc *token.Service) {
	s.r.HandleFunc("GET /check", handler.CheckAccessToken(svc))
}

func (s *Server) AddSwaggerUI() {
	s.r.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
}

func (s *Server) AddMetrics() {
	s.r.Handle("/metrics", promhttp.Handler())
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.r.Handler())
}
