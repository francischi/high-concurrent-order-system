package implement

import (
	"log"
	"sync"
	"time"
	"errors"
	"fmt"
	// "strconv"
	"golang/pkg/base"
	"golang/pkg/repos/models"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type ProductRepo struct {
	base.Repository
	DBconn *gorm.DB
	RDconn *redis.Client
}

func NewProductRepo(db *gorm.DB , rd *redis.Client) *ProductRepo{
	return &ProductRepo{
		DBconn:db,
		RDconn:rd,
	}
}

func (p *ProductRepo) ReduceProducts(productIds map[string]int)(err error){

	locks := p.acquireLocks(productIds , 3)
	defer p.releaseLock(locks)

	if len(locks) != len(productIds){
		return p.InvalidArgument("create order error")
	}
	
    err = p.reduceProductsImpl(productIds)
	if err!=nil{
		return p.InvalidArgument(err.Error())
	}

	return nil
}

func (p *ProductRepo)reduceProductsImpl(productIds map[string]int)error{
	if err := p.checkProducts(productIds) ; err!=nil{
		return err
	}
	return nil
}

func (p *ProductRepo)checkProducts(productIds map[string]int)error{

	errChan := make(chan string, len(productIds))
	var wg sync.WaitGroup
	wg.Add(len(productIds))

	for productId , quantity := range productIds{
		go func(productId string, quantity int ){
			defer wg.Done()

			CurQuantity,err := p.RDconn.Get(productId).Int()
			if err!=nil{
				if err == redis.Nil{
					dbQuantity , err := p.getDbQuantity(productId)
					if err!=nil{
						errChan <- "product repo error " + err.Error()
						return 
					}
					err = p.RDconn.Set(productId , dbQuantity , 600*time.Second).Err()
					if err!=nil{
						errChan <- "product repo error " + err.Error()
						return 
					}
					CurQuantity = dbQuantity
					if dbQuantity < quantity{
						errChan <- "product shortage"
						return 
					}
				}else{
					errChan<- "product repo error " + err.Error()
					return 
				}
			}
			if CurQuantity < quantity {
				errChan <- "product shortage"
				return 
			}
			return
		}(productId , quantity)
	}
	wg.Wait()
	close(errChan)
	for err := range errChan{
        return errors.New(err)
    }
	return nil
}

func (p *ProductRepo)getDbQuantity(productId string) (num int,error error) {
	var model models.ProductModel
	var available int
	err := p.DBconn.Table(model.TableName()).Where("product_id = ?", productId).Select("available").Scan(&available).Error
	if err != nil {
		if err==gorm.ErrRecordNotFound{
			return available , err
		}
		return available , err
	}
	fmt.Println(available)

	return available , nil
}

func (p *ProductRepo)reduceProductImpl(client redis.Pipeliner,productId string,quantity int)error{
	prodQuantity,err := client.Get(productId).Int()
	if err!=nil{
		return p.SystemError("productrepo error")
	}

	newProdQuantity := prodQuantity - quantity
	if newProdQuantity < 0{
		return p.InvalidArgument("product shortage")
	}

	if err := client.Set(productId , newProdQuantity , 600*time.Second).Err();err!=nil{
		return p.SystemError("productrepo set new quantity error")
	}

	return nil
}

func (p *ProductRepo)acquireLocks(productIds map[string]int , maxIterNum int) (locks chan string) {
	// 拿到鎖並返回

	var wg sync.WaitGroup
	locks = make(chan string, len(productIds))
	for productId := range productIds{
		wg.Add(1)
		go func(productId string){
			defer wg.Done()
			lock := p.acquireLock(productId , maxIterNum)
			if lock{
				locks <- productId
			}
		}(productId)
	}
	wg.Wait()
	return locks
}

func (p *ProductRepo)acquireLock(productId string , maxIterNum int) (bool) {
	for i := 0; i < maxIterNum; i++ {
		isLockAcquired:= p.RDconn.SetNX(fmt.Sprintf("%s-lock",productId) , 1 , 3*time.Second )
		if isLockAcquired.Val() {
			return true
		}
		// 等待一段時間後重試
		time.Sleep(500 * time.Millisecond)
	}
	return false
}

func (p *ProductRepo)releaseLock(locks chan string){
	keyNum := len(locks)
	for i:=0 ; i<keyNum ; i++{
		productId := <- locks
		go func(productId string){
			_ ,err := p.RDconn.Del(fmt.Sprintf("%s-lock",productId)).Result()
			if err!=nil{
				log.Println(err.Error())
			}
		}(productId)
	}
}