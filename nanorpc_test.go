package nanorpc

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

const SERVER = "http://[::1]:7076"

var wallet string

func loadConfig() {
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
	loadConfig()
	os.Exit(m.Run())
}

func TestAccounts(t *testing.T) {
	v, err := Accounts(SERVER, wallet)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("accounts:")
	for _, acc := range v {
		t.Log(acc)
	}
}

func TestAccountInfo(t *testing.T) {
	accounts, err := Accounts(SERVER, wallet)
	if err != nil {
		t.Fatal(err)
	}
	a, err := AccountInfo(SERVER, accounts[0])
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("account info:")
	t.Logf("%+v", a)
}

func TestBalance(t *testing.T) {
	accounts, err := Accounts(SERVER, wallet)
	if err != nil {
		t.Fatal(err)
	}
	b, p, err := Balance(SERVER, accounts[0])
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("account balance:")
	t.Logf("%s (%s)", b, p)
}

func getBlockHash(t *testing.T) string {
	accounts, err := Accounts(SERVER, wallet)
	if err != nil {
		t.Fatal(err)
	}
	a, err := AccountInfo(SERVER, accounts[0])
	if err != nil {
		t.Fatal(err)
	}
	return a.Frontier
}

func TestBlockAccount(t *testing.T) {
	acc, err := BlockAccount(SERVER, getBlockHash(t))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("block account:")
	t.Logf(acc)
}

func TestBlockInfo(t *testing.T) {
	b, err := BlockInfo(SERVER, getBlockHash(t))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("block info:")
	t.Logf("%+v", b)
}
