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
	PrevHash    string
	Height      uint32
	CurrHash    string
	Nonce       uint64
	CurrHashSrt string
	Miner       string
	MinerSrt    string
	TimeStamp   string
	Output      uint64
	Transaction []BlockTransaction
	Difficulty  uint8
	Signature   string
}

type BlockTransaction struct {
	CurrHash    string
	CurrHashSrt string
	Output      uint64
	TimeStamp   string
	RandBytes   string
	PrevBlock   string
	Sender      string
	Receiver    string
	Value       uint64
	ToStorage   uint64
	Signature   string
}
