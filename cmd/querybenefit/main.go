package main

import (
	"flag"
	"log"
	"math/big"
	"strconv"
)

var (
	restUrl = flag.String("url", "http://localhost:8669", "rest url")
)

var (
	nodeBenefitList = []struct {
		nodeName string
		benefit  string
	}{
		{"node0", "0x0000000000000000000000000000000000000010"},
		{"node1", "0x0000000000000000000000000000000000000011"},
		{"node2", "0x0000000000000000000000000000000000000012"},
		{"node3", "0x0000000000000000000000000000000000000013"},
		{"node4", "0x0000000000000000000000000000000000000014"},
		{"node5", "0x0000000000000000000000000000000000000015"},
		{"node6", "0x0000000000000000000000000000000000000016"},
	}
)

func main() {
	flag.Parse()
	height := bestBlock(*restUrl).Number
	for i := uint32(0); i < height; i++ {
		for _, node := range nodeBenefitList {
			acc := accountInfo(*restUrl, node.benefit, strconv.Itoa(int(i)))
			{
				energy := big.Int(acc.Energy)
				log.Printf("%s has benefit %v at block \t%d", node.nodeName, energy.Text(10), i)
			}
		}

	}

}
