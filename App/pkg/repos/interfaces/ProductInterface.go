package interfaces

type ProductRepo interface {

	ReduceProducts(productIds map[string]int)(err error)

}