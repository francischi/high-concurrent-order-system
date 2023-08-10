package implement

import (
	"log"
	"sync"
	"time"
	"errors"
	"fmt"
	"encoding/json"
	"golang/pkg/base"
	"golang/pkg/helpers"
	"golang/pkg/repos/models"
	"github.com/go-redis/redis"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

type ProductRepo struct {
	base.Repository
	DBconn *gorm.DB
	RDconn *redis.Client
	ConnPool *helpers.ConnPool
}

func NewProductRepo(db *gorm.DB , rd *redis.Client ,ConnPool *helpers.ConnPool) *ProductRepo{
	return &ProductRepo{
		DBconn:db,
		RDconn:rd,
		ConnPool: ConnPool,
	}
}

type productsContent struct {
	ProductIds  map[string]int
}

func (p *ProductRepo) ReduceProducts(productIds map[string]int)(err error){

	locks := p.acquireLocks(productIds , 5)
	defer p.releaseLock(locks)

	if len(locks) != len(productIds){
		return p.InvalidArgument("create order error")
	}
	
    err = p.reduceRedisProducts(productIds)
	if err!=nil{
		return p.InvalidArgument(err.Error())
	}

	if err := p.pushIntoQue(productIds);err!=nil{
		return err
	}

	return nil
}

func (p *ProductRepo)reduceRedisProducts(productIds map[string]int)error{
	if err := p.checkProducts(productIds) ; err!=nil{
		return err
	}

	errChan := make(chan string, len(productIds))
	var wg sync.WaitGroup
	wg.Add(len(productIds))

	for productId ,quantity := range productIds{
		go func(productId string, quantity int){
			defer wg.Done()
			err := p.reduceRedisProduct(productId , quantity)
			if err!=nil{
				errChan<-err.Error()
				return
			}

		}(productId , quantity)
	}

	wg.Wait()
	close(errChan)
	for err := range errChan{
        return errors.New(err)
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
	var model models.Product
	var available int
	err := p.DBconn.Table(model.TableName()).Where("product_uuid = ?", productId).Select("available").Scan(&available).Error
	if err != nil {
		if err==gorm.ErrRecordNotFound{
			return available , err
		}
		return available , err
	}

	return available , nil
}

func (p *ProductRepo)reduceRedisProduct(productId string,quantity int)error{
	if err := p.RDconn.DecrBy(productId , int64(quantity)).Err() ;err!=nil{
		return err
	}
	if err := p.RDconn.Expire(productId, 600*time.Second).Err();err!=nil{
		return err
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

func(p *ProductRepo) pushIntoQue(productIds map[string]int )(err error){
	conn ,err := p.ConnPool.GetConn()

	ch , queue ,err := p.prepareQueChannel(conn , helpers.GetEnvStr("amqp.productchannel"))
	if err!=nil{
		return p.SystemError("product_repo_error :"+err.Error())
	}

	products := productsContent{
		ProductIds: productIds,
	}

	encodedproducts, err := json.Marshal(products)

	err = ch.Publish(
		"",     // exchange
		queue.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
		  ContentType: "application/json",
		  Body:        encodedproducts,
	})
	if err!=nil{
		return p.SystemError("product_repo_error :"+err.Error())
	}

	ch.Close() 
	if err = p.ConnPool.ReturnConn(conn);err!=nil{
		return p.SystemError("product_repo_error :"+err.Error())
	}

	return nil
}

func (p *ProductRepo)prepareQueChannel(conn *amqp.Connection , queueName string)( *amqp.Channel ,amqp.Queue ,  error) {
	var channel *amqp.Channel
	var queue amqp.Queue
	
	channel, err := conn.Channel()
	if err!=nil{
		return channel, queue ,err
	}

	queue , err = channel.QueueDeclare(
		queueName, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err!=nil{
		return channel , queue , err
	}
	return channel , queue , nil
}