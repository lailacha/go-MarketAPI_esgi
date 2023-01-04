package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	broadcast "github.com/lailacha/go-MarketAPI_esgi/server/broadcaster"
	"github.com/lailacha/go-MarketAPI_esgi/server/payement"
	"github.com/lailacha/go-MarketAPI_esgi/server/product"
)


type PayementHandler interface {
	CreatePayement(c *gin.Context)
	GetPayement(c *gin.Context)
	UpdatePayement(c *gin.Context)
	DeletePayement(c *gin.Context)
	GetPayements(c *gin.Context)
}

type payementHandler struct {
	broadcaster broadcast.Broadcaster
	payementService payement.Service
	productService product.Service
}


func NewPayementHandler(broadcaster broadcast.Broadcaster, productService product.Service, payementService payement.Service) *payementHandler {
	return &payementHandler{
		broadcaster: broadcaster,
		payementService: payementService,
		productService: productService,
	}
}


func (adapter *payementHandler) CreatePayement (c *gin.Context) {
	
	
	productId, err := strconv.Atoi(c.PostForm("productId"))

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "invalid product id",
		})
		return
	}

	// get the product 

	product, err := adapter.productService.Get(productId)

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})

		fmt.Println(err.Error())
		return
	}

	payement, err := adapter.payementService.Create(product)

	b := adapter.broadcaster

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})

		return
	}

	b.Submit(Message{
		UserId: "1",
		Text: "Payement is created",
	})

	 c.JSON(http.StatusOK, &Response{
		Status:  http.StatusOK,
		Message: "Payement is created",
		Data: payement,
	})

}

func (adapter *payementHandler) GetPayement (c *gin.Context) {
	
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "invalid id",
		})
		return
	}

	payement, err := adapter.payementService.Get(id)

	b := adapter.broadcaster

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "something went wrong",
		})
		return
	}

	b.Submit(Message{
		UserId: "1",
		Text: "Payement price is " + strconv.Itoa(payement.PricePaid),
	})

	 c.JSON(http.StatusOK, &Response{
		Status:  http.StatusOK,
		Message: "Payement price is " + strconv.Itoa(payement.PricePaid),
		Data: payement,
	})

}

func (adapter *payementHandler) UpdatePayement(c *gin.Context) {


	id, err := strconv.Atoi(c.Param("id"))


	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "invalid id",
		})
		return
		}

	var payement payement.InputPayement

	err = c.ShouldBindJSON(&payement)


	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "invalid payement",
		})
		return
	}

	newPayement, err := adapter.payementService.Update(id, payement)

	b := adapter.broadcaster

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "something went wrong",
		})
		return
	}

	b.Submit(Message{
		UserId: "1",
		Text: "Payement is updated",
	})
	
	 c.JSON(http.StatusOK, &Response{
		Status:  http.StatusOK,
		Message: "Payement is updated",
		Data: newPayement,
	})


	
}

func (adapter *payementHandler) DeletePayement(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "invalid id",
		})
		return
	}

	err = adapter.payementService.Delete(id)

	b := adapter.broadcaster

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "something went wrong",
		})
		return
	}

	b.Submit(Message{
		UserId: "1",
		Text: "Payement is deleted",
	})

	 c.JSON(http.StatusOK, &Response{
		Status:  http.StatusOK,
		Message: "Payement is deleted",
	})

}


func (adapter *payementHandler) GetPayements (c *gin.Context) {

	payements, err := adapter.payementService.GetAll()

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "something went wrong",
		})
		return
	}

	b := adapter.broadcaster

	b.Submit(Message{
		UserId: "1",
		Text: "Payements are fetched",
	})

	 c.JSON(http.StatusOK, &Response{
		Status:  http.StatusOK,
		Message: "Payements are fetched",
		Data: payements,
	})

}