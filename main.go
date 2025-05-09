package main

import (
	"app/db"
	"app/routes"
)

func main() {
	db.Conn()
	routes.Run()
}
