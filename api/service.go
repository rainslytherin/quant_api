package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

type Service struct {
	Logger *slog.Logger
	cfg    *Config
	isInit bool

	*http.Server

	*gin.Engine
}

func CreateService(cfg *Config) *Service {
	return &Service{
		isInit: false,
		cfg:    cfg,
		Logger: slog.Default(),
		Engine: gin.New(),
	}
}

func (s *Service) Gin() *gin.Engine {
	return s.Engine
}

func (s *Service) Init() {
	s.Server = &http.Server{}
	s.InitGin()
	s.isInit = true
	s.Logger.Info("Init HTTP Service.")
}

func (s *Service) InitGin() {
	s.Engine.Use(gin.Recovery())
	s.Engine.Use(cors.Default())
	s.InitHandlers()
	s.Server.Handler = s.Engine
}

func (s *Service) WithLogger(log *slog.Logger) {
	s.Logger = log.With("service", "http")
}

func (s *Service) Start() error {
	s.Logger.Info("Start HTTP Service.")
	if !s.isInit {
		return fmt.Errorf("HTTP Service is not init.")
	}
	addr := s.ServiceAddress()

	s.Server.Addr = addr

	go func() {
		if err := s.ListenAndServe(); err != nil {
			s.Logger.Error("HTTP Server has error.", "error", err)
			os.Exit(1)
		}
		s.Logger.Info("HTTP Server is exited.")
	}()

	s.Logger.Info("HTTP Server is listen.", slog.String("Listen", s.Addr))

	s.Logger.Info("HTTP Server is started.")

	return nil
}

func (s *Service) ServiceAddress() string {
	return fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
}

// TODO: finish stop func
func (s *Service) Close() error {
	s.Logger.Info("Stop HTTP Service.")
	return nil
}
