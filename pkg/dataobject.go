package pkg

import (
	"math/big"
)

type Transaction struct {
	Hash  string `json:"hash"`
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
}

type BlockResponse struct {
	Result struct {
		Number       string        `json:"number"`
		Transactions []Transaction `json:"transactions"`
	} `json:"result"`
}

type MaxTransaction struct {
	Address string     `json:"address"`
	Wei     *big.Int   `json:"wei"`
	Eth     *big.Float `json:"eth"`
	Dollars string     `json:"dollars"`
}
