package config

import (
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
		log.Fatalf("Failed to connect to database: %v", err)
		return err
	}

	// Configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
		return err
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	log.Println("✅ Database connected successfully")
	return nil
}

// RunMigrations applies database migrations using golang-migrate
func RunMigrations() error {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
		return err
	}

	driver, err := mysqlmigrate.WithInstance(sqlDB, &mysqlmigrate.Config{})
	if err != nil {
		log.Fatalf("Failed to create MySQL driver: %v", err)
		return err
	}

	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		log.Fatalf("Failed to get migrations path: %v", err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to initialize migrations: %v", err)
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
		return err
	}

	log.Println("✅ Migrations applied successfully")
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
