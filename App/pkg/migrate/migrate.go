package main

import (
	"golang/pkg/helpers"
	"golang/pkg/migrate/migrations"
)

func main(){
	helpers.InitMySql()
	migrations.Member(helpers.SqlSession)
	migrations.Order(helpers.SqlSession)
	migrations.Product(helpers.SqlSession)
}