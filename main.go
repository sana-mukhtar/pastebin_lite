package main

import (
	"log"
	"pastebin/internal"
)

func main() {

	internal.InitDB()
	log.Println("DB verification successful")
}
