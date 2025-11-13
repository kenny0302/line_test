package main

import "main/app"

func main() {
	router := app.SetupRouter()
	router.Run(":8080")
}
