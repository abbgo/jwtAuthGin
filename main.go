package main

import (
	"jwtAuthGin/routes"
	"log"
	"os"
)

func main() {
	r := routes.Routes()

	// static file
	os.Mkdir("./uploads", os.ModePerm)
	r.Static("/image", "./uploads")

	// run routes
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
