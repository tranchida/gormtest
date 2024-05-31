package controlers

import "github.com/gin-gonic/gin"

func SetupRouter() {

	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.GET("/products", AllProducs)

	r.Run()
}
