package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
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
	service := &Service{
		isInit: false,
		cfg:    cfg,
		Logger: slog.Default(),
		Engine: gin.New(),
	}

	service.Init()
	return service
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
	s.Engine.Use(sloggin.New(s.Logger))
	s.InitCors()
	s.InitHandlers()
	s.Server.Handler = s.Engine
}

func (s *Service) InitCors() {
	s.Engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://121.37.182.188:8080", "http://121.37.182.188"},
		AllowMethods:     []string{"GET", "HEAD", "DELETE", "OPTIONS", "POST", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}

func (s *Service) WithLogger(log *slog.Logger) {
	s.Logger = log.With("service", "http")
}

func (s *Service) Start() {
	s.Logger.Info("Start HTTP Service.")
	if !s.isInit {
		panic("HTTP Service is not init.")
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
}

func (s *Service) ServiceAddress() string {
	return fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
}

// TODO: finish stop func
func (s *Service) Close() {
	s.Logger.Info("Stop HTTP Service.")
	return
}
