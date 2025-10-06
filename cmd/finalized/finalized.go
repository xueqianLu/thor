package main

import (
	"flag"
	"fmt"
	"github.com/vechain/thor/cmd/utils"
	"log"
	"math/big"
)

var (
	restUrl        = flag.String("url", "http://127.0.0.1:10005", "rest url")
	report         = flag.String("report", "/root/node/collect.csv", "report file")
	validatorCount = flag.Int("validators", 101, "validator count")
	blockHeight    = flag.Int("height", 360, "block height")
)

type BlockInfo struct {
	Number    int
	Signer    string
	SignerIdx int
	Timestamp int64
}

func getHackSignerIdx(validatorCount int) (int, int) {
	offset := 10
	start := validatorCount / 3
	end := start + validatorCount/3 - 1
	return start + offset, end + offset
}

func calcAuraFinalizedTime(epoch int, validatorCount int, blocks map[int]BlockInfo) (int64, bool) {
	epochStart := epoch * 180
	epochEnd := (epoch + 1) * 180
	needValidators := validatorCount / 2
	gotted := make(map[string]bool)
	gotCount := 0
	var epochStartTime int64 = 0
	var finalizedTime int64 = 0
	for i := epochStart; i < epochEnd; i++ {
		blk, ok := blocks[i]
		if !ok {
			continue
		}
		if i == epochStart {
			epochStartTime = blk.Timestamp
		}
		if _, ok := gotted[blk.Signer]; !ok {
			gotted[blk.Signer] = true
			gotCount += 1
		}
		if gotCount >= needValidators {
			finalizedTime = blk.Timestamp
			break
		}
	}
	return finalizedTime - epochStartTime, finalizedTime != 0
}

func calcFobFinalizedTime(epoch int, validatorCount int, blocks map[int]BlockInfo) (int64, bool) {
	epochStart := epoch * 180
	epochEnd := (epoch + 1) * 180
	needValidators := validatorCount/3*2 + 1
	gotted := make(map[string]bool)
	gotCount := 0
	var epochStartTime int64 = 0
	var finalizedTime int64 = 0
	for i := epochStart; i < epochEnd; i++ {
		blk, ok := blocks[i]
		if !ok {
			continue
		}
		if i == epochStart {
			epochStartTime = blk.Timestamp
		}
		if _, ok := gotted[blk.Signer]; !ok {
			gotted[blk.Signer] = true
			gotCount += 1
		}

		if gotCount >= needValidators {
			finalizedTime = blk.Timestamp
			break
		}
	}
	return finalizedTime - epochStartTime, finalizedTime != 0
}

func main() {
	flag.Parse()
	var blocks = make(map[int]BlockInfo) // blocknumber -> BlockInfo

	for i := 0; i < *blockHeight; i++ {
		blk := utils.BlockByNumber(*restUrl, int64(i))
		if blk == nil {
			log.Printf("block %d not found\n", i)
			continue
		}
		number := int(blk.Number)
		signerIdx, _ := new(big.Int).SetString(blk.Beneficiary.String()[2:], 16)
		blocks[number] = BlockInfo{
			Number:    number,
			Signer:    blk.Beneficiary.String(),
			Timestamp: int64(blk.Timestamp),
			SignerIdx: int(signerIdx.Int64()),
		}
	}

	endpeoch := *blockHeight / 180
	for epoch := 1; epoch < endpeoch; epoch++ {
		auraTime, auraOk := calcAuraFinalizedTime(epoch, *validatorCount, blocks)
		auraTimeStr := "none"
		if auraOk && auraTime > 0 {
			auraTimeStr = fmt.Sprintf("%d", auraTime)
		}
		fobTime, fobOk := calcFobFinalizedTime(epoch, *validatorCount, blocks)
		fobTimeStr := "none"
		if fobOk && fobTime > 0 {
			fobTimeStr = fmt.Sprintf("%d", fobTime)
		}
		log.Printf("epoch %d, aura finalized time %s, fob finalized time %s\n", epoch, auraTimeStr, fobTimeStr)
	}

	log.Printf("collect finished")
}
