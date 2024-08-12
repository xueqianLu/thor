package main

import (
	"flag"
	"github.com/vechain/thor/cmd/utils"
	"log"
)

var (
	restUrl     = flag.String("url", "http://13.228.149.45:10005", "rest url")
	report      = flag.String("report", "/root/node/collect.csv", "report file")
	blockHeight = flag.Int("height", 360, "block height")
)

func main() {
	flag.Parse()
	var history = make(map[int]map[string]int)
	list := make([]int, 0)

	for i := 0; i < *blockHeight; i++ {
		blk := utils.BlockByNumber(*restUrl, int64(i))
		if blk == nil {
			continue
		}
		epoch := int(blk.Number) / 180
		signer := blk.Beneficiary.String()
		if _, ok := history[epoch]; !ok {
			history[epoch] = make(map[string]int)
			history[epoch][signer] = 1
			list = append(list, epoch)
		} else {
			history[epoch][signer]++
		}
	}
	log.Printf("collect finished")
	for _, epoch := range list {
		for signer, count := range history[epoch] {
			log.Printf("epoch %d, signer %s block %d\n", epoch, signer, count)
		}
	}
	return
}
