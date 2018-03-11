package nano_rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
)

type Account struct {
	Frontier            string `json:"frontier"`
	OpenBlock           string `json:"open_block"`
	RepresentativeBlock string `json:"representative_block"`
	Balance             string `json:"balance"`
	ModifiedTimestamp   string `json:"modified_timestamp"`
	BlockCount          string `json:"block_count"`
}

type Block struct {
	Type        string
	Previous    string
	Destination string
	Balance     string
	Work        string
	Signature   string
}

type actionRPC struct {
	Action string `json:"action"`
}

type accountRPC struct {
	actionRPC
	Account string `json:"account"`
}

type blockRPC struct {
	actionRPC
	Hash string `json:"hash"`
}

type walletRPC struct {
	actionRPC
	Wallet string `json:"wallet"`
}

func marshal(x interface{}) []byte {
	buf, err := json.Marshal(x)
	if err != nil {
		panic(err)
	}
	return buf
}

func unmarshal(buf []byte) map[string]interface{} {
	res := make(map[string]interface{})
	if err := json.Unmarshal(buf, &res); err != nil {
		panic(err)
	}
	return res
}

func request(url string, req interface{}) (resdata []byte, res map[string]interface{}, err error) {
	reqdata := marshal(req)

	response, err := http.Post(url, "application/json", bytes.NewReader(reqdata))
	if err != nil {
		return
	}
	resdata, err = ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	res = unmarshal(resdata)
	if val, ok := res["error"]; ok {
		panic(val.(string) + ": " + string(reqdata))
	}
	return
}

func rawtonano(rawstr string) uint64 {
	var raw, rawtonano, nano = new(big.Int), new(big.Int), new(big.Int)
	raw, ok := new(big.Int).SetString(rawstr, 10)
	if !ok {
		panic("invalid balance")
	}
	rawtonano.Exp(big.NewInt(10), big.NewInt(24), nil)
	nano.Div(raw, rawtonano)
	if nano.BitLen() > 64 {
		return 9999999999999999999
	}
	return nano.Uint64()
}

func Accounts(url string, wallet string) (accounts []string, err error) {
	req := walletRPC{actionRPC{"account_list"}, wallet}
	_, res, err := request(url, req)
	if err != nil {
		return
	}
	v := res["accounts"].([]interface{})
	accounts = make([]string, len(v))
	for i, x := range v {
		accounts[i] = x.(string)
	}
	return
}

func AccountInfo(url string, acc string) (info Account, err error) {
	req := accountRPC{actionRPC{"account_info"}, acc}
	resdata, _, err := request(url, req)
	if err != nil {
		return
	}
	err = json.Unmarshal(resdata, &info)
	return
}

func Balance(url string, acc string) (balance uint64, pending uint64, err error) {
	req := accountRPC{actionRPC{"account_balance"}, acc}
	_, res, err := request(url, req)
	if err != nil {
		return
	}
	balance = rawtonano(res["balance"].(string))
	pending = rawtonano(res["pending"].(string))
	return
}

func BlockInfo(url string, hash string) (block Block, err error) {
	req := blockRPC{actionRPC{"block"}, hash}
	_, res, err := request(url, req)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(res["contents"].(string)), &block)
	if err != nil {
		return
	}

	i, ok := new(big.Int).SetString(block.Balance, 16)
	if !ok {
		err = fmt.Errorf("internal error")
		return
	}
	block.Balance = i.Text(10)
	return
}
