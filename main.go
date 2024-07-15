// main.go

package main

import (
	"fmt"
	"parking_lot_service/internal/di" // Import your container package
)

func main() {
	container := di.NewContainer()
	if container == nil {
		return
	}
	e := container.GetEchoInstance()
	router := container.GetRouter()
	router.MapRoutes(e)
	port := ":8080"
	fmt.Printf("Server started on port %s\n", port)
	e.Logger.Fatal(e.Start(port))
}
