package main

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"io/ioutil"
	"log"
	"net/http"
)

//	{
//	 "number": 0,
//	 "id": "0x0000000025b7e3372117c18fcee8ddee8f08ace384639d728c9656e2419cf0e5",
//	 "size": 172,
//	 "parentID": "0xffffffff5465737420436861696e000000000000000000000000000000000000",
//	 "timestamp": 1722054674,
//	 "gasLimit": 1000000000000,
//	 "beneficiary": "0x0000000000000000000000000000000000000000",
//	 "gasUsed": 0,
//	 "totalScore": 0,
//	 "txsRoot": "0x45b0cfc220ceec5b7c1c62c4d4193d38e4eba48e8815729ce75f9c0ab0e4c1c0",
//	 "txsFeatures": 0,
//	 "stateRoot": "0xdc5a406e794a6d3894af3d6b2493933d62442c313464e199c0d13e864c556e44",
//	 "receiptsRoot": "0x45b0cfc220ceec5b7c1c62c4d4193d38e4eba48e8815729ce75f9c0ab0e4c1c0",
//	 "com": false,
//	 "signer": "0x0000000000000000000000000000000000000000",
//	 "isTrunk": true,
//	 "isFinalized": true,
//	 "transactions": []
//	}
type BlockInfo struct {
	Number int64  `json:"number"`
	ID     string `json:"id"`
}

func getChainTag(url string) (byte, error) {
	api := fmt.Sprintf("%s/blocks/0", url)
	res, err := httpGet(api)
	if err != nil {
		return 0x00, err
	}
	var genInfo BlockInfo
	if err = json.Unmarshal(res, &genInfo); err != nil {
		log.Printf("json.Unmarshal genesis info: %v", err)
		return 0x00, err
	}
	id := common.FromHex(genInfo.ID)
	return id[31], nil
}

func httpGet(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Printf("http.Get: %v", err)
		return nil, err
	}
	r, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Printf("ioutil.ReadAll: %v", err)
		return nil, err
	}
	return r, nil
}
