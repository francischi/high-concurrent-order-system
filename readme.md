# **hight preformence order system**
## **Introduction**
使用訂單系統為範例，以符合clean architecture 之方式撰寫，示範如何使用redis對資料進行緩存，
以及使用rabbitmq進行非同步任務之處理，提高系統乘載極限以即可擴展性。

## **技術點**
* Clean Architecture
* Message Broker ( RabbitMQ )
* Cache ( Redis )
* DB ( MySQL )
* Message Broker Connection Pool
* Docker

## 系統架構
![image](https://github.com/francischi/high-concurrent-order-system/images/system_design.jpg)

## 使用方式
* rabbitmq
```
> docker pull rabbitmq:management
> docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:management

透過localhost:15672查看rabbitmq是否啟動
```

* app
```
> cd ./App
> go run main.go
```

* consumer
```
> cd ../MailConsumer
> go run main.go
```

<br>

```
PS 此專案架構較適合中大型專案開發使用
```