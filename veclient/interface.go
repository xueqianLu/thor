package veclient

import (
	"github.com/vechain/thor/block"
	pb "github.com/xueqianLu/vehackcenter/hackcenter"
)

type HackCenterClient interface {
	SubmitBlock(blk *block.Block) (*pb.SubmitBlockResponse, error)
	SubBroadcastTask() error
	SubscribeBlock() error
}

type P2pCenterClient interface {
	//RegisterNode()
	SubscribeBlock() error
	BroadcastBlock(blk *block.Block) error
}
