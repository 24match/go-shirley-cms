package main

import (
	"fmt"

	"medical-device-cms/backend/config"
	"medical-device-cms/backend/routes"
)

func main() {
	if err := config.InitDB(); err != nil {
		panic(err)
	}

	r := config.SetupRouter()
	routes.RegisterRoutes(r)

	fmt.Println("Server starting on :8080")
	r.Run(":8080")
}