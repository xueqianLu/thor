package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
)

var (
	restUrl     = flag.String("url", "http://13.228.149.45:10005", "rest url")
	report      = flag.String("report", "/root/node/report.csv", "report file")
	nodeNum     = flag.Int("number", 22, "node number")
	blockHeight = flag.Int("height", 360, "block height")
)

type nodeBenefit = struct {
	nodeName string
	benefit  string
}

func main() {
	flag.Parse()
	height := bestBlock(*restUrl).Number
	if height > uint32(*blockHeight) {
		height = uint32(*blockHeight)
	}
	f, err := os.OpenFile(*report, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return
	}
	defer func() {
		f.Close()
	}()

	write := csv.NewWriter(f)
	nodeBenefitList := make([]nodeBenefit, 0)
	head := []string{"block"}
	for i := 0; i < *nodeNum; i++ {
		nodeBenefit := nodeBenefit{}
		nodeBenefit.nodeName = "node" + strconv.Itoa(i)
		nodeBenefit.benefit = "0x" + fmt.Sprintf("%040d", i+10)
		nodeBenefitList = append(nodeBenefitList, nodeBenefit)
		head = append(head, nodeBenefit.nodeName)
	}

	write.Write(head)
	for i := uint32(0); i < height; i++ {

		record := make([]string, 0)
		record = append(record, strconv.Itoa(int(i)))
		for _, node := range nodeBenefitList {
			acc := accountInfo(*restUrl, node.benefit, strconv.Itoa(int(i)))
			energy := big.Int(acc.Energy)
			//处理金额
			energyFloat, err := strconv.ParseFloat(energy.Text(10), 64)
			if err != nil {
				log.Fatalf("strconv.ParseFloat err: %v", err)
				return
			}
			totalReward := big.NewFloat(energyFloat)
			totalReward.Quo(totalReward, big.NewFloat(1e18))
			totalRewardStr := totalReward.Text('f', 4)

			record = append(record, totalRewardStr)
			log.Printf("%s has benefit %v at block \t%d", node.nodeName, totalRewardStr, i)
		}
		write.Write(record)
	}
	write.Flush()
}
