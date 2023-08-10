package productModule

import (
	"productConsumer/pkg/base"
	"productConsumer/pkg/repos/interfaces"
)

type ProductService struct{
	base.Service
	ProductRepo interfaces.IproductRepo
}

type productsContent struct {
	ProductIds  map[string]int
}

func NewProductService(productRepo interfaces.IproductRepo) *ProductService{
	var ProductService ProductService
	ProductService.ProductRepo = productRepo
	return &ProductService
}

func(s *ProductService) Reduce(dto *ReduceDto)(err error){
	if err := dto.Check();err!=nil{
		return s.InvalidArgument(err.Error())
	}

	if err:=s.ProductRepo.ReduceProducts(dto.ProductIds);err!=nil{
		return err
	}

	return nil
}