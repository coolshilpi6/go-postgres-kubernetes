package main

import (
  "database/sql"
  "fmt"

  _ "github.com/lib/pq"
)

const (
  host     = "192.168.216.217"
  port     = 31372
  user     = "postgresadmin"
  password = "admin123"
  dbname   = "postgresdb"
)

func main() {
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
  defer rows.Close()
  for rows.Next() {
    var id int
    var name string
    err = rows.Scan(&id, &name)
    if err != nil {
      // handle this error
      panic(err)
    }
    fmt.Println(id, name)
  }
  // get any error encountered during iteration
  err = rows.Err()
  if err != nil {
    panic(err)
  }
  defer db.Close()
}