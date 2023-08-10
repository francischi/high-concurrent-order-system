//+build wireinject

package productModule

import (
	"github.com/google/wire"
	impl "productConsumer/pkg/repos/implement"
	"productConsumer/pkg/repos/interfaces"
	"productConsumer/pkg/helpers"
)

var ProductRepo = wire.NewSet(impl.NewProductRepo,wire.Bind(new(interfaces.IproductRepo), new(*impl.ProductRepo)))

func InitialProductController() *ProductController{
	wire.Build(
		helpers.NewSqlSession ,
		ProductRepo,
		NewProductService,
		NewProductController,
	)
	return &ProductController{}
}