package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"fmt"
	"github.com/gorilla/mux"
	"database/sql"
	"strconv"
	"github.com/joho/godotenv"
  _ "github.com/lib/pq"
)

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	if name == "" {
		name = "Guest"
	}
	now := time.Now() // find the time right now
	log.Printf("Received request for %s\n", name)
	w.Write([]byte(fmt.Sprintf("Hello, %s\n", name)))
	w.Write([]byte(fmt.Sprintf("%s\n", now.Format("02-01-2006"))))
	w.Write([]byte(fmt.Sprintf("%s\n", now.Format("15:04:05"))))
	
	//Getting the database from .env file
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	portstr := os.Getenv("DB_PORT")
	host := os.Getenv("DB_HOST")
	dbname := os.Getenv("DB_NAME")
	port, err := strconv.Atoi(portstr)
	
	 psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }


  err = db.Ping()
  if err != nil {
    panic(err)
  }

  fmt.Println("Successfully connected!")
  
  rows, err := db.Query("SELECT id, name from COMPANY")
  if err != nil {
    // handle this error better than this
    panic(err)
  }
  w.Write([]byte(fmt.Sprintf("Company Table Details\n")))
  w.Write([]byte(fmt.Sprintf("ID\tName\n")))
  defer rows.Close()
  for rows.Next() {
    var id int
    var name string
    err = rows.Scan(&id, &name)
    if err != nil {
      // handle this error
      panic(err)
    }
	w.Write([]byte(fmt.Sprintf("%d", id)))
	w.Write([]byte(fmt.Sprintf("\t%s", name)))
	w.Write([]byte(fmt.Sprintf("\n")))
  }
  // get any error encountered during iteration
  err = rows.Err()
  if err != nil {
    panic(err)
  }
  defer db.Close()
    
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	// Create Server and Route Handlers
	r := mux.NewRouter()

	r.HandleFunc("/", handler)
	r.HandleFunc("/health", healthHandler)
	r.HandleFunc("/readiness", readinessHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8000",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start Server
	go func() {
		log.Println("Starting Server")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful Shutdown
	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}