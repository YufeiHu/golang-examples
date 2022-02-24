package main

import (
    "sync"
)

type dbConnection struct {}

var (
   dbConnOnce sync.Once
   conn *dbConnection
)

func GetConnection() *dbConnection {
   dbConnOnce.Do(func() {
      conn = &dbConnection{}
   })
   return conn
}
