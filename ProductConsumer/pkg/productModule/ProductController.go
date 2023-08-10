package productModule

import (
	// "fmt"
	"log"
	"encoding/json"
)

type ProductController struct{
	ProductService ProductService
}

func NewProductController(ProductService *ProductService) *ProductController{
	return &ProductController{
		ProductService : *ProductService,
	}
}

func failOnError(err error) {
	if err != nil {
	  log.Fatalf("%s",err)
	}
}

func (c *ProductController)Reduce(body []byte){

	var dto ReduceDto
	
	err := json.Unmarshal(body, &dto)
	failOnError(err)

	err = c.ProductService.Reduce(&dto)
	failOnError(err)
}