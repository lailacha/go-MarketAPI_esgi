package adapter

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lailacha/go-MarketAPI_esgi/server/broadcaster"
	"github.com/lailacha/go-MarketAPI_esgi/server/payement"
	"github.com/lailacha/go-MarketAPI_esgi/server/product"
)

type GinAdapter interface {
	Stream(c *gin.Context)

	CreatePayement(c *gin.Context)

	UpdateProduct(c *gin.Context)
	CreateProduct(c *gin.Context)
	DeleteProduct(c *gin.Context)
	GetProduct(c *gin.Context)
}

type ginAdapter struct {
	broadcaster broadcast.Broadcaster
	productService product.Service
	payementService payement.Service
}

type Message struct
{
	UserId string
	Text string
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data	interface{} `json:"data"`
}


func NewGinAdapter(broadcaster broadcast.Broadcaster, productService product.Service, payementService payement.Service) *ginAdapter {
	return &ginAdapter{
		broadcaster: broadcaster,
		payementService: payementService,
		productService: productService,
	}
}

// Stream is the handler for the stream endpoint
func (adapter *ginAdapter) Stream(c *gin.Context) {
	

	//create a new channel to handle the stream
	listener := make(chan interface{})

	// get the broadcaster

	adapter.broadcaster.Register(listener)

	//close the channel when error message or client is gone
	defer adapter.broadcaster.Unregister(listener)

	clientGone := c.Request.Context().Done()

	c.Stream(func(w io.Writer) bool {
		select {
		case <-clientGone:
			return false
		case message := <-listener:
			serviceMsg, ok := message.(Message)
			if !ok {
				fmt.Println("not a message")
				c.SSEvent("message", message)
				return false
			}
			c.SSEvent("message", " "+serviceMsg.UserId+" → "+serviceMsg.Text)
			return true
		}
	})



	fmt.Println("stream is OK")
}

func (adapter *ginAdapter) CreatePayement (c *gin.Context) {
	
	//get POST data


	fmt.Println("create payement", c.PostForm("id"))

	id, err := strconv.Atoi(c.PostForm("id"))


	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		fmt.Println(err)
		return
	}
		price := c.PostForm("price")


		// get the broadcaster
		b := adapter.broadcaster

		// save the payement
		adapter.payementService.Create(id, price);


		b.Submit(Message{
			UserId: "1",
			Text: "Payement is created",
		})

		response := &Response{
			Status:  http.StatusOK,
			Message: "Payement is created",
			Data: nil,
		}

		c.JSON(http.StatusOK, response)
	

}

func (adapter *ginAdapter) CreateProduct (c *gin.Context) {
	

	name := c.PostForm("name")


	if name == "" { 
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		fmt.Println(name)
		return
	}

	price := c.PostForm("price")

	if price == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid price"})
		return
	}


	fmt.Println("create product", c.PostForm("name"))
	// get the broadcaster
	b := adapter.broadcaster


	// save the payement
	product, err := adapter.productService.Create(name, price);

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "invalid product",
			Data: err.Error(),
		})
		return
	}


	b.Submit(Message{
		UserId: "1",
		Text: "Product is created",
	})

	response := &Response{
		Status:  http.StatusOK,
		Message: "Product is created",
		Data: product,
	}

	c.JSON(http.StatusOK, response)

}


func (adapter *ginAdapter) UpdateProduct (c *gin.Context) {
	

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "invalid id",
			Data: err.Error(),
		})
		return
	}

	var product product.Product

	err = c.ShouldBindJSON(&product)

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "invalid product",
			Data: err.Error(),
		})
		return
	}

	updatedProduct, err := adapter.productService.Update(id, product)


	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data: nil,
		})
		return
	}

	c.JSON(http.StatusOK, &Response{
		Status:  http.StatusOK,
		Message: "product updated",
		Data: updatedProduct,
	})

}


func (adapter *ginAdapter) DeleteProduct (c *gin.Context) {
	

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "invalid id",
			Data: err.Error(),
		})
		return
	}

	err = adapter.productService.Delete(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "product not found",
			Data: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &Response{
		Status:  http.StatusOK,
		Message: "product deleted",
		Data: nil,
	})

}


func (adapter *ginAdapter) GetProduct (c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "invalid id",
			Data: err.Error(),
		})
		return
	}

	product, err := adapter.productService.Get(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "product not found",
			Data: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &Response{
		Status:  http.StatusOK,
		Message: "product found",
		Data: product,
	})

}

