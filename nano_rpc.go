package nano_rpc

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type action_rpc struct {
	Action string
}

type account_rpc struct {
	action_rpc
	Account string
}

type wallet_rpc struct {
	action_rpc
	Wallet string
}

func marshal(x interface{}) string {
	buf, err := json.Marshal(x)
	if err != nil {
		panic(err)
	}
	return strings.ToLower(string(buf))
}

func unmarshal(buf []byte) map[string]interface{} {
	res := make(map[string]interface{})
	if err := json.Unmarshal(buf, &res); err != nil {
		panic(err)
	}
	return res
}

func request(url string, req interface{}) (res map[string]interface{}, err error) {
	reqstr := marshal(req)

	buf := bytes.NewBufferString(reqstr)
	response, err := http.Post(url, "application/json", buf)
	if err != nil {
		return
	}
	resdata, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	res = unmarshal(resdata)
	if val, ok := res["error"]; ok {
		panic(val)
	}
	return
}

func Accounts(url string, wallet string) (accounts []string, err error) {
	req := wallet_rpc{action_rpc{"account_list"}, wallet}
	res, err := request(url, req)
	v := res["accounts"].([]interface{})
	accounts = make([]string, len(v))
	for i, x := range v {
		accounts[i] = x.(string)
	}
	return
}

func Balance(url string, acc string) (balance string, err error) {
	req := account_rpc{action_rpc{"account_balance"}, acc}
	res, err := request(url, req)
	if err == nil {
		balance = res["balance"].(string)
	}
	return
}
