package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Set Gin to production mode
	gin.SetMode(gin.ReleaseMode)

	// Create a new Gin router
	router := gin.New()

	// Add middleware to log requests and recover from panics
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Get the port number from the PORT environment variable, or use 8080 as a default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Get the database connection information from environment variables
	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")
	dbuser := os.Getenv("DB_USER")
	dbpassword := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Construct the database connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbhost, dbport, dbuser, dbpassword, dbname)

	// Open a connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}
	defer db.Close()

	// Test the database connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error testing database connection: %v", err)
	}

	// Set a timeout for all requests
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Add a health check endpoint
	router.GET("/healthz", func(c *gin.Context) {
		// Test the database connection
		err = db.Ping()
		if err != nil {
			c.String(http.StatusInternalServerError, "Error pinging database")
			return
		} else {
			log.Println("Database connection OK")
		}

		// Return a success message
		c.String(http.StatusOK, "OK")
	})

	// Start the server
	log.Printf("Starting server on :%s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %v", err)
	}
}
