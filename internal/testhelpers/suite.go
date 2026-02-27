package testhelpers

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"remy/internal/config"
	"remy/internal/infrastructure"
)

type IntegrationSuite struct {
	suite.Suite
	DB     *gorm.DB
	Engine *gin.Engine
}

func (s *IntegrationSuite) SetupSuite() {
	os.Setenv("APP_ENV", "test")

	config.LoadEnv()
	cfg, err := config.LoadConfig()
	if err != nil {
		s.T().Fatalf("Failed to load config: %v", err)
	}

	if cfg.DatabaseURL == "" {
		s.T().Fatal("Database URL is not set in config")
	}

	db, err := config.SetupDatabase(cfg)
	if err != nil {
		s.T().Fatalf("Failed to connect to database: %v", err)
	}

	router := infrastructure.NewRouter()
	infrastructure.SetupRoutes(router, db)

	clearDatabase(db)

	gin.SetMode(gin.TestMode)

	s.DB = db
	s.Engine = router
}

func (s *IntegrationSuite) TearDownTest() {
	clearDatabase(s.DB)
}

func (s *IntegrationSuite) TearDownSuite() {
	sqlDB, err := s.DB.DB()
	if err != nil {
		s.T().Fatalf("Failed to get database connection: %v", err)
	}
	sqlDB.Close()
}

func clearDatabase(db *gorm.DB) {
	var tables []string
	if err := db.Raw("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'").Scan(&tables).Error; err != nil {
		panic(fmt.Sprintf("Failed to get table names: %v", err))
	}

	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("DELETE FROM %s", table)).Error; err != nil {
			panic(fmt.Sprintf("Failed to clear table %s: %v", table, err))
		}
	}

	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("DELETE FROM sqlite_sequence WHERE name='%s'", table)).Error; err != nil {
			panic(fmt.Sprintf("Failed to reset auto-increment for table %s: %v", table, err))
		}
	}
}
