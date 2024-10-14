package server

import (
	"net/http"

	"RescueSupport.sv/handlers"
	"RescueSupport.sv/model"
	"github.com/gin-gonic/gin"
)

type Company struct {
	company handlers.Company
}

func NewCompany(c handlers.Company) Company {
	return Company{
		company: c,
	}
}

func (c Company) CompanySignUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var company model.SignUp

		//Serialize the request body
		if err := ctx.ShouldBindJSON(&company); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"data": handleServerResponse(400, "bad request", "", err.Error(), nil)})
			return

		}

		//Call a handler to process this request
		if err := c.company.SignUp(company); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"data": handleServerResponse(500, "internal server error", "", err.Error(), nil)})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": handleServerResponse(200, "success", "", nil, nil)})
	}
}
