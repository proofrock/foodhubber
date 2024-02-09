package params

// Called in a mutex-protected environment, so safe

var orders, stock int

func TouchOrdersGen() {
	orders++
}

func TouchStockGen() {
	stock++
}

func GetOrderGen() int {
	return orders
}

func GetStockGen() int {
	return stock
}
