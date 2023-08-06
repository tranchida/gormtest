package controlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"tranchida.github.com/gormtest/services"
)

func AllProducs(ctx *gin.Context) {

	products, err := services.FindAll()

	if err != nil {
		log.Fatal(err)
	}

	ctx.JSON(http.StatusOK, gin.H{"data": products})

}
