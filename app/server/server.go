package server

import (
	"errors"
	"net/http"
	"orange-backstage-api/app/router"
	"orange-backstage-api/infra/config"

	"github.com/gin-gonic/gin"
)

type Server struct {
	http   *http.Server
	router *router.Router

	cfg config.Server
}

func New(r *router.Router, cfg config.Server) *Server {
	return &Server{
		http: &http.Server{
			Addr:              ":" + cfg.Port,
			ReadTimeout:       cfg.ReadTimeout,
			WriteTimeout:      cfg.WriteTimeout,
			MaxHeaderBytes:    1 << 20,
			ReadHeaderTimeout: cfg.ReadTimeout,
		},
		router: r,
		cfg:    cfg,
	}
}

func (s *Server) Serve() error {
	gin.SetMode(s.cfg.RunMode)

	r := gin.New()

	r.MaxMultipartMemory = 10 << 20 // 10 MiB

	s.router.Register(r)
	s.http.Handler = r

	var err error
	if s.IsSSL() {
		err = s.http.ListenAndServeTLS(
			s.cfg.CertFilePath,
			s.cfg.KeyFilePath,
		)
	} else {
		err = s.http.ListenAndServe()
	}
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s Server) IsSSL() bool {
	return s.cfg.CertFilePath != "" && s.cfg.KeyFilePath != ""
}

func (s Server) Addr() string {
	return s.http.Addr
}
