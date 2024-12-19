package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mdedealf/go-api/internal/delivery/http"
	"github.com/mdedealf/go-api/internal/delivery/http/middleware"
	"github.com/mdedealf/go-api/internal/repository"
	"github.com/mdedealf/go-api/internal/usecase"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type AppConfig struct {
	DB     *gorm.DB
	App    *fiber.App
	Log    *logrus.Logger
	Config *viper.Viper
}

func (cfg *AppConfig) Run() {
	// setup repositories
	productRepository := repository.NewProductRepository(cfg.Log)
	// setup use cases
	productUseCase := usecase.NewProductUsecase(productRepository, cfg.Log, cfg.DB)
	// setup controller
	productController := http.NewProductController(&productUseCase, cfg.Log)
	// setup middleware
	authMiddleware := middleware.NewAuth()
	routeConfig := http.Router{
		App:               cfg.App,
		ProductController: productController,
		AuthMiddleware:    authMiddleware,
	}
	routeConfig.Setup()
}
