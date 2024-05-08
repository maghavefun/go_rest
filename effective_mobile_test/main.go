package main

import (
	"os"

	_ "github.com/maghavefun/effective_mobile_test/docs"
	"github.com/maghavefun/effective_mobile_test/initializer"
	"github.com/maghavefun/effective_mobile_test/router"
)

func init() {
	initializer.LoadEnvVariables()
	initializer.ConnectToDB()
	if os.Getenv("GO_ENV") == "DEV" {
		initializer.Migrate()
	}
}

// @title           E_M swagger API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:4000
// @BasePath  /v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

// @schemes http
func main() {
	router.NewRouter().Serve()
}
