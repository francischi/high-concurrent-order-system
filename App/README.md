# **App**

## **Introduction**
此為主要後端APP，有會員、訂單、商品相關功能，詳細請參考 api-document.md
此後端購物下單功能透過緩存 (redis) 與queue (rabbitMQ) 之設計，能夠容納更大量訂單同時發生。