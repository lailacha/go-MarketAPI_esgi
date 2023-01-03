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
	GetPayement(c *gin.Context)
	UpdatePayement(c *gin.Context)
	DeletePayement(c *gin.Context)
	GetPayements(c *gin.Context)


	UpdateProduct(c *gin.Context)
	CreateProduct(c *gin.Context)
	DeleteProduct(c *gin.Context)
	GetProduct(c *gin.Context)
	GetProducts(c *gin.Context)
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
			Message: "something went wrong",
		})
		return
	}

	payement, err := adapter.payementService.Create(product)

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
		Text: "Payement is created",
	})

	 c.JSON(http.StatusOK, &Response{
		Status:  http.StatusOK,
		Message: "Payement is created",
		Data: payement,
	})

}

func (adapter *ginAdapter) GetPayement (c *gin.Context) {
	
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
		Text: "Payement price is " + payement.PricePaid,
	})

	 c.JSON(http.StatusOK, &Response{
		Status:  http.StatusOK,
		Message: "Payement price is " + payement.PricePaid,
		Data: payement,
	})

}

func (adapter *ginAdapter) UpdatePayement(c *gin.Context) {


	id, err := strconv.Atoi(c.Param("id"))


	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "invalid id",
		})
		return
		}

	var payement payement.Payement

	err = c.ShouldBindJSON(&payement)


	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "invalid payement",
		})
		return
	}

	payement, err = adapter.payementService.Update(id, payement)

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
		Data: payement,
	})


	
}

func (adapter *ginAdapter) DeletePayement(c *gin.Context) {

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

func (adapter *ginAdapter) GetPayements (c *gin.Context) {

	payements, err := adapter.payementService.FindAll()

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

func (adapter *ginAdapter) GetProducts (c *gin.Context) {

	products, err := adapter.productService.FindAll()

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
		Text: "Products are fetched",
	})

	 c.JSON(http.StatusOK, &Response{
		Status:  http.StatusOK,
		Message: "Products are fetched",
		Data: products,
	})

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
		Text: product.Name + " is created",
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

	b := adapter.broadcaster


	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data: nil,
		})
		return
	}

	b.Submit(Message{
		UserId: "1",
		Text: updatedProduct.Name + " is updated to price " + updatedProduct.Price,
	})


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

	b := adapter.broadcaster

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "product not found",
			Data: err.Error(),
		})
		return
	}

	b.Submit(Message{
		UserId: "1",
		Text: "Product is deleted",
	})

	c.JSON(http.StatusOK, &Response{
		Status:  http.StatusOK,
		Message: "product deleted",
		Data: "product deleted",
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

	b := adapter.broadcaster


	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "product not found",
			Data: err.Error(),
		})
		return
	}

	b.Submit(Message{
		UserId: "1",
		Text: product.Name + " is found",
	})

	c.JSON(http.StatusOK, &Response{
		Status:  http.StatusOK,
		Message: "product found",
		Data: product,
	})

}

