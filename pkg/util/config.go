// pkg/util/config.go
package util

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

// Declare global variables to hold the database connection details
var (
	host     string
	port     string
	user     string
	password string
	dbname   string
	pcName   string
	dirname  string
	job      string
	month    string
	joboder  string
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	PCName     string
	DirName    string
}

var (
	config          *Config
	fsnotifyStarted bool
	contentCond     *sync.Cond
	initialStartup  bool
	contentMutex    sync.Mutex
	logPathMutex    sync.Mutex
	logPath         string
	contentReady    bool
	content         string
	contentCh       = make(chan string)
)

func GetConfig() *Config {
	return config
}

func SetConfig(dbHost, dbPort, dbUser, dbPassword, dbName, PCname, dirName string) {
	host = dbHost
	port = dbPort
	user = dbUser
	password = dbPassword
	dbname = dbName
	pcName = PCname
	dirname = dirName

	// Save the configuration to the environment file
	err := godotenv.Write(map[string]string{
		"DB_HOST":     host,
		"DB_PORT":     port,
		"DB_USER":     user,
		"DB_PASSWORD": password,
		"DB_NAME":     dbname,
		"PC_NAME":     pcName,
		"DIR_NAME":    dirname,
	}, ".env")
	if err != nil {
		log.Println("error saving configuration:", err)
	}
}

func InitConfig() {
	// Load the configuration from the environment file
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("error loading configuration from .env file:", err)
	}

	// Initialize the config instance
	config = &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		PCName:     os.Getenv("PC_NAME"),
		DirName:    os.Getenv("DIR_NAME"),
	}
}

func SetJob(jobOrder string, jobMonth string) string {
	return jobMonth + jobOrder
}

func LoadConfigGo() string {
	// Read the configuration from environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	pcName := os.Getenv("PC_NAME")
	dirname := os.Getenv("DIR_NAME")

	// Create a JSON object containing the configuration values
	configJSON := fmt.Sprintf(`{
		"dbHost": "%s",
		"dbPort": "%s",
		"dbUser": "%s",
		"dbPassword": "%s",
		"dbName": "%s",
		"PCname": "%s",
		"dirName": "%s"
	}`, host, port, user, password, dbname, pcName, dirname)

	return configJSON
}
