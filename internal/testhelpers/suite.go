package testhelpers

import (
	"os"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"remy/internal/config"
)

type TestSuite struct {
	suite.Suite
	DB *gorm.DB
}

func (s *TestSuite) SetupSuite() {
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

	clearDatabase(db)

	s.DB = db
}

func (s *TestSuite) TearDownTest() {
	clearDatabase(s.DB)
}

func (s *TestSuite) TearDownSuite() {
	sqlDB, err := s.DB.DB()
	if err != nil {
		s.T().Fatalf("Failed to get database connection: %v", err)
	}
	sqlDB.Close()
}

func clearDatabase(db *gorm.DB) {
	var tables []string
	db.Raw("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'").Scan(&tables)

	for _, table := range tables {
		db.Exec("DELETE FROM " + table)
	}

	for _, table := range tables {
		db.Exec("DELETE FROM sqlite_sequence WHERE name='" + table + "'")
	}
}
