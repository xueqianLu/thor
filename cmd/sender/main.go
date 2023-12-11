package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/vechain/thor/api/transactions"
	"github.com/vechain/thor/genesis"
	"github.com/vechain/thor/tx"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	accountFile  = flag.String("account", "/root/account.json", "account file")
	accountIndex = flag.Int("index", 0, "account index")
	url          = flag.String("url", "", "node rpc url")
)

type AccountInfo struct {
	Address string `json:"address"`
	Private string `json:"private"`
}

func main() {

}

func senTx(url string) {
	var blockRef = tx.NewBlockRef(0)
	var chainTag = repo.ChainTag()
	var expiration = uint32(10)
	var gas = uint64(21000)

	tx := new(tx.Builder).
		BlockRef(blockRef).
		ChainTag(chainTag).
		Expiration(expiration).
		Gas(gas).
		Build()
	sig, err := crypto.Sign(tx.SigningHash().Bytes(), genesis.DevAccounts()[0].PrivateKey)
	if err != nil {
		log.Fatalf("crypto.Sign: %v", err)
	}
	tx = tx.WithSignature(sig)
	rlpTx, err := rlp.EncodeToBytes(tx)
	if err != nil {
		log.Fatalf("rlp.EncodeToBytes: %v", err)
	}

	res := httpPost(url+"/transactions", transactions.RawTx{Raw: hexutil.Encode(rlpTx)})
	var txObj map[string]string
	if err = json.Unmarshal(res, &txObj); err != nil {
		log.Fatalf("json.Unmarshal: %v", err)
	}
}

func httpPost(url string, obj interface{}) []byte {
	data, err := json.Marshal(obj)
	if err != nil {
		log.Fatalf("json.Marshal: %v", err)
	}
	res, err := http.Post(url, "application/x-www-form-urlencoded", bytes.NewReader(data))
	if err != nil {
		log.Fatalf("http.Post: %v", err)
	}
	r, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatalf("ioutil.ReadAll: %v", err)
	}
	return r
}
