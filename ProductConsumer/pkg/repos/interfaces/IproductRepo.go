package interfaces

type IproductRepo interface {

	ReduceProducts(productIds map[string]int)(err error)

}