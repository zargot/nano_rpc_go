package nano_rpc

import (
	"bytes"
	"encoding/json"
	"github.com/zargot/algo"
	"io/ioutil"
	"net/http"
)

type rpc struct {
	Action string
}

type balance_rpc struct {
	rpc
	Account string
}

const CTYPE = "application/json"

func Balance(url string, acc string) string {
	var err error

	var response *http.Response
	var reqdata []byte
	req := balance_rpc{rpc{"account_balance"}, acc}
	if reqdata, err = json.Marshal(req); err != nil {
		panic(err)
	}
	reqdata = algo.Transform(reqdata, algo.ToLower)
	reqstream := bytes.NewReader(reqdata)
	if response, err = http.Post(url, CTYPE, reqstream); err != nil {
		panic(err)
	}

	var resdata []byte
	res := make(map[string]string)
	if resdata, err = ioutil.ReadAll(response.Body); err != nil {
		panic(err)
	}
	if err = json.Unmarshal(resdata, &res); err != nil {
		panic(err)
	}
	if err, ok := res["error"]; ok {
		panic(err)
	}
	return res["balance"]
}
