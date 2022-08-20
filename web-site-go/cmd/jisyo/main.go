package main

import (
	"web-site-go/db"
	"web-site-go/server"
)

func main() {
	db.Init()
	server.Init()
	db.Close()
}
