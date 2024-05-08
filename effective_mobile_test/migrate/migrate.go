package main

import (
	"github.com/maghavefun/effective_mobile_test/initializer"
)

func init() {
	initializer.LoadEnvVariables()
	initializer.ConnectToDB()
}

func main() {
	// автомиграции запускаются в main.go, если окружение "DEV"
	// initializer.Migrate()
}
