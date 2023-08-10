package interfaces

type OrderRepo interface{

	AddOrder(memberId string , orderProducta map[string]int)(err error)

}