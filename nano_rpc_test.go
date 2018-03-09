package nano_rpc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

const SERVER = "http://[::1]:7076"

func load_config(t *testing.T) (wallet, acc string) {
	var err error
	var b []byte
	if b, err = ioutil.ReadFile(os.Getenv("HOME") + "/RaiBlocks/config.json"); err != nil {
		t.Fatal(err)
	}
	m := make(map[string]interface{})
	if err = json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if str, ok := m["wallet"].(string); ok {
		wallet = str
	}
	if str, ok := m["account"].(string); ok {
		acc = str
	}
	return
}

func TestBalance(t *testing.T) {
	_, acc := load_config(t)
	t.Run("balance", func(t *testing.T) {
		b := Balance(SERVER, acc)
		fmt.Println(b)
	})
}
