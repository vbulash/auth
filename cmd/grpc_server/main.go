package main

import (
	"github.com/vbulash/auth/internal/api"
	"log"
)

func main() {
	err := api.AppRun()
	if err != nil {
		log.Fatal(err)
	}
}
