package nano_rpc

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type actionRPC struct {
	Action string `json:"action"`
}

type accountRPC struct {
	actionRPC
	Account string `json:"account"`
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

func request(url string, req interface{}) (res map[string]interface{}, err error) {
	reqdata := marshal(req)

	response, err := http.Post(url, "application/json", bytes.NewReader(reqdata))
	if err != nil {
		return
	}
	resdata, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	res = unmarshal(resdata)
	if val, ok := res["error"]; ok {
		panic(val.(string) + ": " + string(reqdata))
	}
	return
}

func Accounts(url string, wallet string) (accounts []string, err error) {
	req := walletRPC{actionRPC{"account_list"}, wallet}
	res, err := request(url, req)
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

func Balance(url string, acc string) (balance string, err error) {
	req := accountRPC{actionRPC{"account_balance"}, acc}
	res, err := request(url, req)
	if err != nil {
		return
	}
	balance = res["balance"].(string)
	return
}
