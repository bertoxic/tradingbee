package config

import (
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

type AppConfig struct {
	InProduction bool
	Client       *mongo.Client
	Config  	Config
}

type Config struct {
	DBHost     string
    DBPort     string
    DBNAME     string
    DBUser     string
    DBPassword string
    MONGODB_URI  string
    JWTSecret  string
}

func NewAppConfig() (*AppConfig, error){
	config, err :=LoadConfig() 
	if err != nil {
		return nil, err
	}
	return &AppConfig{Config: *config}, nil
}


func LoadConfig() (*Config, error) {
    return &Config{
        DBHost:     os.Getenv("DB_HOST"),
        DBNAME:     os.Getenv("DB_NAME"),
        DBPort:     os.Getenv("DB_PORT"),
        DBUser:     os.Getenv("DB_USER"),
        DBPassword: os.Getenv("DB_PASSWORD"),
        MONGODB_URI:  os.Getenv("MONGODB_URI"),
        JWTSecret:  os.Getenv("JWT_SECRET"),
    }, nil
}