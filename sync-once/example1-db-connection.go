package main

import (
   "fmt"
   "sync"
)

type dbConnection struct {}

var (
   dbConnOnce sync.Once
   conn *dbConnection
)

func GetConnection() *dbConnection {
   dbConnOnce.Do(func() {
      conn = &DbConnection{}
   })
   return conn
}
