package routes

import (
	"github.com/gin-gonic/gin"
)

func serveSwaggerYAML(c *gin.Context) {
    c.File("./docs/swagger.yaml")
}

// Unused since no swagger.json needed currently
// func serveSwaggerJSON(c *gin.Context) {
//     c.File("./docs/swagger.json")
// }



