package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	broadcast "github.com/lailacha/go-MarketAPI_esgi/server/broadcaster"
	"github.com/lailacha/go-MarketAPI_esgi/server/handler"
	"github.com/lailacha/go-MarketAPI_esgi/server/payement"
	"github.com/lailacha/go-MarketAPI_esgi/server/product"
	
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func handleCreateToken(c *gin.Context) {

	type Request struct {
		UserId   string `json:"userId"`
		Password string `json:"password"`
	}

	type Response struct {
		Token string `json:"token"`
	}

	type ResponseError struct {
		Status int `json:"status"`
	}

	//parsing login body
	var request Request
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(400, &ResponseError{Status: 400})
		return
	}

	fmt.Println(request.UserId)

	if request.UserId == "admin" && request.Password == "admin" {

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userId": request.UserId,
			"exp":    time.Now().Add(time.Hour * 24).Unix(),
			"iat":    time.Now().Unix(),
		})

		tokenString, err := token.SignedString([]byte("secret"))

		if err != nil {
			c.JSON(400, &ResponseError{Status: 400})
			return
		}

		c.JSON(
			200,
			&Response{
				Token: tokenString,
			},
		)

		return
	}

	c.JSON(400, &ResponseError{Status: 400})

}


func main() {

	router := gin.Default()

	//create the connection to the db
	dbURL := "user:password@tcp(127.0.0.1:3306)/go-exam?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	// run a migration to create the tables
	db.AutoMigrate(&payement.Payement{})
	db.AutoMigrate(&product.Product{})

	payementRepository := payement.NewPayementRepository(db)
	payementService := payement.NewService(payementRepository)

	productRepository := product.NewProductRepository(db)
	productService := product.NewService(productRepository)

	// get the broadcaster
	b := broadcast.NewBroadcaster(200)

	payementHandler := handler.NewPayementHandler(b, productService, payementService)

	productHandler := handler.NewProductHandler(b, productService, payementService)

	router.POST("/createToken", handleCreateToken)

	router.GET("/stream", func(c *gin.Context) {

		//create a new channel to handle the stream
		listener := make(chan interface{})

		// get the broadcaster

		b.Register(listener)

		//close the channel when error message or client is gone
		defer b.Unregister(listener)

		clientGone := c.Request.Context().Done()

		c.Stream(func(w io.Writer) bool {
			select {
			case <-clientGone:
				fmt.Println("client gone")
				return false
			case message := <-listener:
				serviceMsg, ok := message.(handler.Message)
				if !ok {
					fmt.Println("not a message")
					c.SSEvent("message", message)
					return false
				}
				c.SSEvent("message", " "+serviceMsg.UserId+" ??? "+serviceMsg.Text)
				return true
			}
		})

	})

	router.POST("/payement", payementHandler.CreatePayement)
	router.GET("/payement/:id", payementHandler.GetPayement)
	router.PUT("/payement/:id", payementHandler.UpdatePayement)
	router.DELETE("/payement/:id", payementHandler.DeletePayement)
	router.GET("/payement", payementHandler.GetPayements)

	router.POST("/product", productHandler.CreateProduct)
	router.PUT("/product/:id", productHandler.UpdateProduct)
	router.DELETE("/product/:id", productHandler.DeleteProduct)
	router.GET("/product/:id", productHandler.GetProduct)
	router.GET("/product", productHandler.GetProducts)

	router.Run(fmt.Sprintf(":%v", 8084))

	// run the broadcaster

}
