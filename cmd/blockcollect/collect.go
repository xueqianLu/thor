package main

import (
	"flag"
	"github.com/vechain/thor/cmd/utils"
	"log"
	"strconv"
)

var (
	restUrl     = flag.String("url", "http://13.228.149.45:10005", "rest url")
	report      = flag.String("report", "/root/node/collect.csv", "report file")
	blockHeight = flag.Int("height", 360, "block height")
)

func main() {
	flag.Parse()
	var history = make(map[int]map[string]int)
	var honest = make(map[int]int)
	var hacker = make(map[int]int)
	list := make([]int, 0)

	for i := 0; i < *blockHeight; i++ {
		blk := utils.BlockByNumber(*restUrl, int64(i))
		if blk == nil {
			log.Printf("block %d not found\n", i)
			continue
		}
		epoch := int(blk.Number) / 180
		signer := blk.Beneficiary.String()
		// signer latest 2 bytes.

		if _, ok := history[epoch]; !ok {
			history[epoch] = make(map[string]int)
			history[epoch][signer] = 1
			list = append(list, epoch)
		} else {
			history[epoch][signer]++
		}
		singerIdx, _ := strconv.Atoi(signer[len(signer)-2:])
		singerIdx -= 10
		log.Printf("epoch %d, signer %s, block %d\n", epoch, signer, history[epoch][signer])
		if singerIdx >= 7 && singerIdx <= 15 {
			hacker[epoch]++
		} else {
			honest[epoch]++
		}

	}
	log.Printf("collect finished")
	//for _, epoch := range list {
	//	for signer, count := range history[epoch] {
	//		log.Printf("epoch %d, signer %s block %d\n", epoch, signer, count)
	//	}
	//}
	for epoch, count := range honest {
		log.Printf("epoch %d honest block %d\n", epoch, count)
	}
	for epoch, count := range hacker {
		log.Printf("epoch %d hacker block %d\n", epoch, count)
	}
	return
}
