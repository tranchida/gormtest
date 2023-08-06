package controlers

import "github.com/gin-gonic/gin"

func SetupRouter() {

	r := gin.Default()

	r.GET("/products", AllProducs)

	r.Run()
}
