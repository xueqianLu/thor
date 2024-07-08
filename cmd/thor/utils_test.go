package main

import (
	"encoding/json"
	"fmt"
	"github.com/vechain/thor/cmd/thor/node"
	"os"
	"testing"
	"time"
)

func TestLoadOrGenerateKey(t *testing.T) {
	validatorCount := 101
	type AuthInfo struct {
		MasterAddress   string `json:"masterAddress"`
		EndorsorAddress string `json:"endorsorAddress"`
		Identity        string `json:"identity"`
	}
	type AccountInfo struct {
		Name    string `json:"name"`
		Private string `json:"private"`
		Address string `json:"address"`
	}

	// generate 101 account.
	accounts := make([]AccountInfo, 0)
	for i := 0; i < validatorCount; i++ {
		path := fmt.Sprintf("keys/account.key.%d", i)
		key, err := loadOrGeneratePrivateKey(path)
		if err != nil {
			panic(err)
		}
		acckey := &node.Master{PrivateKey: key}
		accounts = append(accounts, AccountInfo{
			Name:    fmt.Sprintf("account-%d", i),
			Private: key.D.String(),
			Address: acckey.Address().String(),
		})
	}
	d, _ := json.MarshalIndent(accounts, "", "  ")
	os.WriteFile("accounts.json", d, 0644)

	authority := make([]AuthInfo, 0)
	masters := ""
	for i := 0; i < validatorCount; i++ {
		path := fmt.Sprintf("keys/master.key.%d", i)
		key, err := loadOrGeneratePrivateKey(path)
		if err != nil {
			panic(err)
		}
		master := &node.Master{PrivateKey: key}
		authority = append(authority, AuthInfo{
			MasterAddress:   master.Address().String(),
			EndorsorAddress: "0x782c4C7E8bA047edfbfA0F2815D4D035467C6aFD",
			Identity:        "0x000000000000000068747470733a2f2f656e762e7665636861696e2e6f72672f",
		})
		masters += fmt.Sprintf("master-%d: %s\n", i, master.Address().String())
	}
	os.WriteFile("master.keys", []byte(masters), 0644)

	// generate genesis.json
	gs := ""
	gs += "{\n"
	gs += fmt.Sprintf("  \"launchTime\": %d,\n", time.Now().Unix())
	gs += "  \"gasLimit\": 1000000000000,\n"
	gs += "  \"extraData\": \"Test Chain\",\n"
	gs += "  \"accounts\": [\n"
	for _, acc := range accounts {
		gs += fmt.Sprintf("    {\n")
		gs += fmt.Sprintf("      \"address\": \"%s\",\n", acc.Address)
		gs += fmt.Sprintf("      \"balance\": 20000000000000000000000000000,\n")
		gs += fmt.Sprintf("      \"energy\": 500000000000000000000000000\n")
		gs += "    },\n"
	}
	gs += fmt.Sprintf("    {\n")
	gs += fmt.Sprintf("      \"address\": \"%s\",\n", "0x782c4C7E8bA047edfbfA0F2815D4D035467C6aFD")
	gs += fmt.Sprintf("      \"balance\": 25000000000000000000000000,\n")
	gs += fmt.Sprintf("      \"energy\": 0,\n")
	gs += fmt.Sprintf("      \"code\": \"%s\",\n", "0x6060604052600256")
	gs += fmt.Sprintf("		\"storage\": {\n\"0x0000000000000000000000000000000000000000000000000000000000000001\": \"0x0000000000000000000000000000000000000000000000000000000000000002\"}")
	gs += fmt.Sprintf("    }\n")
	gs += "  ],\n"

	gs += "  \"authority\": [\n"
	for _, auth := range authority {
		gs += fmt.Sprintf("    {\n")
		gs += fmt.Sprintf("      \"masterAddress\": \"%s\",\n", auth.MasterAddress)
		gs += fmt.Sprintf("      \"endorsorAddress\": \"%s\",\n", auth.EndorsorAddress)
		gs += fmt.Sprintf("      \"identity\": \"%s\"\n", auth.Identity)
		if auth != authority[len(authority)-1] {
			gs += "    },\n"
		} else {
			gs += "    }\n"
		}
	}
	gs += "  ],\n"
	gs += "  \"params\": {\n"
	gs += "    \"rewardRatio\": 300000000000000000,\n"
	gs += "    \"baseGasPrice\": 100000000000000,\n"
	gs += "    \"proposerEndorsement\": 25000000000000000000000000,\n"
	gs += "    \"executorAddress\": \"0x0000000000000000000000004578656375746f72\",\n"
	gs += fmt.Sprintf("    \"maxBlockProposers\": %d\n", validatorCount)
	gs += "  },\n"
	gs += "  \"executor\": {\n"
	gs += "    \"approvers\": [\n"
	gs += "      {\n"
	gs += "		\"address\": \"0xE0785611500B582cCE651c17477a408Dd0057D30\",\n"
	gs += "		\"identity\": \"0x00000000000067656e6572616c20707572706f736520626c6f636b636861696e\"\n"
	gs += "		 }\n"
	gs += "    ]\n"
	gs += "  }\n"
	gs += "}\n"

	os.WriteFile("genesis.json", []byte(gs), 0644)
}
