package nano_rpc

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

const SERVER = "http://[::1]:7076"

var wallet string

func load_config() {
	b, err := ioutil.ReadFile(os.Getenv("HOME") + "/RaiBlocks/config.json")
	if err != nil {
		panic(err)
	}
	m := make(map[string]interface{})
	if err = json.Unmarshal(b, &m); err != nil {
		panic(err)
	}
	wallet = m["wallet"].(string)
}

func TestMain(m *testing.M) {
	load_config()
	os.Exit(m.Run())
}

func TestAccounts(t *testing.T) {
	t.Run("accounts", func(t *testing.T) {
		v, err := Accounts(SERVER, wallet)
		if err != nil {
			t.Fatal(err)
		}
		for _, acc := range v {
			t.Log(acc)
		}
	})
}

func TestBalance(t *testing.T) {
	t.Run("balance", func(t *testing.T) {
		accounts, err := Accounts(SERVER, wallet)
		b, err := Balance(SERVER, accounts[0])
		if err != nil {
			t.Fatal(err)
		}
		t.Log(b)
	})
}
