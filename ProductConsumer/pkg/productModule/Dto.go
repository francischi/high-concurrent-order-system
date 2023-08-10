package productModule

import "errors"

type ReduceDto struct {
	ProductIds map[string]int
}

func (dto *ReduceDto) Check() error {
	for productIds, quantity := range dto.ProductIds {
		if productIds == "" {
			return errors.New("product ids required")
		}
		if quantity <=0 {
			return errors.New("invalid quantity")
		}
	}
	return nil
}
