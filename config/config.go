package config

import (
	"fmt"
	models "java-gem/graph/model"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Repositpry struct {
	DB *gorm.DB
}

func Configure() *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), //
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  true,
		},
	)

	if os.Getenv("RAILWAY_ENVIRONMENT") == "" {
		// Load .env file only in development
		if err := godotenv.Load(); err != nil {
			errorMessage := fmt.Sprintf("Failed to load the .env file: %v", err)
			log.Fatal(errorMessage)
		}
	}

	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbName := os.Getenv("DATABASE_NAME")
	dbPassword := os.Getenv("DATABASE_PASSWORD")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Africa/Kigali", dbHost, dbUser, dbPassword, dbName, dbPort)

	fmt.Println("CONNECTION STRING___", dbHost, dbPort, dbUser, dbName, dbPassword)

	DB, err := gorm.Open(
		postgres.New(
			postgres.Config{
				DSN:                  dsn,
				PreferSimpleProtocol: true,
			},
		),
		&gorm.Config{
			Logger: newLogger,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: false,
			},
		},
	)

	if err != nil {
		log.Fatal("Error connecting to database")
	}

	DB.AutoMigrate(&models.User{}, &models.Coffee{})

	return DB
}

func CloseConnection(db *gorm.DB) {
	dbSQL, err := db.DB()

	if err != nil {
		panic("Failed to close the connection with Database")
	}

	dbSQL.Close()
}

var RequiredOperations map[string]bool = map[string]bool{
	"login":  false,
	"signUp": false,
}
