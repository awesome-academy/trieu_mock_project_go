package config

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	mysqlmigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// BuildDSN builds the MySQL DSN from database config
func BuildDSN() string {
	dbConfig := LoadConfig().Database
	return dbConfig.User + ":" + dbConfig.Password + "@tcp(" + dbConfig.Host + ":" + dbConfig.Port + ")/" + dbConfig.Database + "?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true"
}

// ConnectToMySQL establishes a connection to the MySQL database
func ConnectToMySQL(dsn string) error {
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	dbConfig := LoadConfig().Database
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)

	log.Println("✅ Database connected successfully")
	return nil
}

// RunMigrations applies database migrations using golang-migrate
func RunMigrations() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("Failed to get database instance: %w", err)
	}

	driver, err := mysqlmigrate.WithInstance(sqlDB, &mysqlmigrate.Config{})
	if err != nil {
		return fmt.Errorf("Failed to create MySQL driver: %w", err)
	}

	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		return fmt.Errorf("Failed to get migrations path: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"mysql",
		driver,
	)
	if err != nil {
		return fmt.Errorf("Failed to initialize migrations: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("Failed to apply migrations: %w", err)
	}

	log.Println("Migrations applied successfully")
	return nil
}

// InitDB initializes database connection and performs migration
func InitDB() error {
	dsn := BuildDSN()
	if err := ConnectToMySQL(dsn); err != nil {
		return err
	}

	if err := RunMigrations(); err != nil {
		return err
	}

	return nil
}
