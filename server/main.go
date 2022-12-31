package main;

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lailacha/go-MarketAPI_esgi/server/broadcaster"
	"github.com/lailacha/go-MarketAPI_esgi/server/payement"
	"github.com/lailacha/go-MarketAPI_esgi/server/product"
	adapter "github.com/lailacha/go-MarketAPI_esgi/server/adapter"
	"gorm.io/gorm"
	"log"
	"gorm.io/driver/mysql"
)

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

	// create the repository and the services
	//productRepository := product.NewProductRepository(db)
	//productService := product.NewService(productRepository)


	payementRepository := payement.NewPayementRepository(db)
	payementService := payement.NewService(payementRepository)


	productRepository := product.NewProductRepository(db)
	productService := product.NewService(productRepository)


	// get the broadcaster
	b := broadcast.NewBroadcaster(20)
	
	ginAdapter := adapter.NewGinAdapter(b, productService, payementService)

	router.GET("/stream", ginAdapter.Stream)
	router.POST("/createPayement", ginAdapter.CreatePayement)

	router.POST("/createProduct", ginAdapter.CreateProduct)
	router.PUT("/updateProduct/:id", ginAdapter.UpdateProduct)
	router.DELETE("/deleteProduct/:id", ginAdapter.DeleteProduct)
	router.GET("/getProduct/:id", ginAdapter.GetProduct)

	router.Run(fmt.Sprintf(":%v", 8084))

	// run the broadcaster

}