package model

type User struct {
	Id       int
	Username string
	Email    string
	Password string
	Wallet   string
}

type Transaction struct {
	Id        int
	Reciver   string
	Sender    string
	S_Reciver string
	S_Sender  string
	Value     string
	Message   string
}

type Block struct {
	Height      uint32
	CurrHash    string
	CurrHashSrt string
	Miner       string
	MinerSrt    string
	TimeStamp   string
	Output      uint64
	Transaction []BlockTransaction
}

type BlockTransaction struct {
	CurrHash    string
	CurrHashSrt string
	Output      uint64
	TimeStamp   string
}
