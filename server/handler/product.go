package handler

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	broadcast "github.com/lailacha/go-MarketAPI_esgi/server/broadcaster"
	"github.com/lailacha/go-MarketAPI_esgi/server/payement"
	"github.com/lailacha/go-MarketAPI_esgi/server/product"
)



type ProductHandler interface {
	CreateProduct(c *gin.Context)
	GetProduct(c *gin.Context)
	UpdateProduct(c *gin.Context)
	DeleteProduct(c *gin.Context)
	GetProducts(c *gin.Context)
}

type productHandler struct {
	broadcaster broadcast.Broadcaster
	payementService payement.Service
	productService product.Service
}


func NewProductHandler(broadcaster broadcast.Broadcaster, productService product.Service, payementService payement.Service) *productHandler {
	return &productHandler{
		broadcaster: broadcaster,
		payementService: payementService,
		productService: productService,
	}
}


func (adapter *productHandler) CreateProduct (c *gin.Context) {
	

	var inputProduct product.InputProduct


	err := c.ShouldBindJSON(&inputProduct)


	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data: inputProduct,
		})
		return
	}

	b := adapter.broadcaster


	// save the payement
	product, err := adapter.productService.Create(inputProduct);

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data: inputProduct,
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


func (adapter *productHandler) UpdateProduct (c *gin.Context) {
	

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "invalid id",
			Data: err.Error(),
		})
		return
	}

	var inputProduct product.InputProduct

	err = c.ShouldBindJSON(&inputProduct)

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Status:  http.StatusBadRequest,
			Message: "invalid product",
			Data: err.Error(),
		})
		return
	}

	updatedProduct, err := adapter.productService.Update(id, inputProduct)

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
		Text: updatedProduct.Name + " is updated to price " + strconv.Itoa(updatedProduct.Price),
	})


	c.JSON(http.StatusOK, &Response{
		Status:  http.StatusOK,
		Message: "product updated",
		Data: updatedProduct,
	})

}


func (adapter *productHandler) DeleteProduct (c *gin.Context) {
	

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


func (adapter *productHandler) GetProduct (c *gin.Context) {

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


func (adapter *productHandler) GetProducts (c *gin.Context) {

	products, err := adapter.productService.GetAll()

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