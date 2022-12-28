package model

type User struct {
	Id       int
	Username string
	Email    string
	Password string
	Wallet   string
}

type Transaction struct {
	Id      int
	Reciver string
	Sender  string
	Value   string
	Message string
}
