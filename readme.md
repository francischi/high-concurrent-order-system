# **hight preformence order system**
## **Introduction**
使用訂單系統為範例，以符合clean architecture 之方式撰寫，示範如何使用redis對資料進行緩存，
以及使用rabbitmq進行非同步任務之處理，提高系統乘載極限以及可擴展性。

## 場景分析
* 大流量
* 需要即時得到結果
* 注重資料正確性( 避免race condition )
* 搶多個鎖機制

## **技術點**
* Clean Architecture
* Message Broker ( RabbitMQ )
* Cache ( Redis )
* Lock ( mutex )
* DB ( MySQL )
* Message Broker Connection Pool
* Docker

## 系統架構
![image](https://github.com/francischi/high-concurrent-order-system/blob/main/images/system_design.jpg)

為了確保扣庫存時不會發生race condition，在redis操作訂單前加入了鎖的機制，
基於訂單會有多筆商品需要扣庫存，因此加入搶多個互斥鎖以及自旋鎖機制，
並將扣庫存之動作移進redis中處理，並將扣完庫存之訂單存入queue，
讓訂單可以被consumer消耗，若訂單量大也可啟動多個consumer達到水平擴展。

## 使用方式
* start rabbitmq
    ```
    > docker pull rabbitmq:management
    > docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:management

    透過localhost:15672查看rabbitmq是否啟動
    ```

* start redis
    ```
    > docker pull redis
    > docker run --name redis -p 6379:6379 -d redis
    ```
* create .env
    ```
    copy .env.example in each folder(App、mailConsumer、prductConsumer)
    and change to your own setting
    ```
* migrate 
    ```
    create a DB name it as same as .env
    > cd App
    > go run ./pkg/migrate/migrate.go

    will create members、orders、products three tables and seed product data
    ```
* app
    ```
    > cd App
    > go run main.go
    ```

* mailConsumer
    ```
    > cd MailConsumer
    > go run main.go
    ```

* productConsumer
    ```
    > cd ProductConsumer
    > go run main.go
    ```

<br>

```
PS 此專案架構較適合中大型專案開發使用
```