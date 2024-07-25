package veclient

import (
	"context"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/inconshreveable/log15"
	"github.com/vechain/thor/block"
	"github.com/vechain/thor/comm"
	pb "github.com/xueqianLu/vehackcenter/hackcenter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"strconv"
)

var log = log15.New("pkg", "veClient")

type VeClient struct {
	proposer string
	index    int
	comu     *comm.Communicator
	conn     pb.CenterServiceClient
}

func NewClient(proposer string, comu *comm.Communicator) *VeClient {
	serverUrl := os.Getenv("VE_HACK_SERVER_URL")
	hackIndex := os.Getenv("VE_HACK_CLIENT_INDEX")
	if serverUrl == "" || hackIndex == "" {
		log.Error("VE_HACK_SERVER_URL or VE_HACK_CLIENT_INDEX not set")
		return nil
	}
	client := new(VeClient)
	conn, err := grpc.Dial(serverUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(1024*1024*1024),
			grpc.MaxCallSendMsgSize(1024*1024*1024)),
	)
	if err != nil {
		log.Error("veClient connect failed", "err", err)
	}
	client.index, _ = strconv.Atoi(hackIndex)
	client.proposer = proposer
	client.conn = pb.NewCenterServiceClient(conn)
	client.comu = comu
	return client
}

func (c *VeClient) SubmitBlock(blk *block.Block) (*pb.SubmitBlockResponse, error) {
	pbblk := new(pb.Block)
	var err error
	pbblk.Hash = blk.Header().ID().String()
	pbblk.Height = int64(blk.Header().Number())
	pbblk.Timestamp = int64(blk.Header().Timestamp())
	pbblk.Data, err = rlp.EncodeToBytes(blk)
	if err != nil {
		log.Error("SubmitBlock encode block failed", "err", err)
		return nil, err
	}

	pbblk.Proposer = new(pb.Proposer)
	pbblk.Proposer.Proposer = c.proposer
	pbblk.Proposer.Index = int32(c.index)
	log.Info("In veclient SubmitBlock", "number", blk.Header().Number())

	return c.conn.SubmitBlock(context.TODO(), pbblk)
}

func (c *VeClient) SubBroadcastTask() error {
	sub, err := c.conn.SubBroadcastTask(context.TODO(), &pb.SubBroadcastTaskRequest{
		Proposer: c.proposer,
	})
	if err != nil {
		log.Error("SubBroadcastTask failed", "err", err)
		return err
	}
	for {
		task, err := sub.Recv()
		if err != nil {
			log.Error("SubBroadcastTask Recv failed", "err", err)
			return err
		}
		log.Info("SubBroadcastTask Recv", "task", task)
		block := new(block.Block)
		err = rlp.DecodeBytes(task.Data, block)
		if err != nil {
			log.Error("SubBroadcastTask decode block failed", "err", err)
			continue
		}
		log.Info("In veclient broadcast block", "number", block.Header().Number())
		c.comu.BroadcastBlock(block)
	}
	return nil
}

func (c *VeClient) SubscribeBlock() error {
	in := new(pb.SubscribeBlockRequest)
	in.Proposer = c.proposer
	sub, err := c.conn.SubscribeBlock(context.TODO(), in)
	if err != nil {
		log.Error("SubscribeBlock failed", "err", err)
		return err
	}
	for {
		msg, err := sub.Recv()
		if err != nil {
			log.Error("SubscribeBlock Recv failed", "err", err)
			return err
		}

		block := new(block.Block)
		err = rlp.DecodeBytes(msg.Data, block)
		if err != nil {
			log.Error("SubscribeBlock decode block failed", "err", err)
			continue
		}
		log.Info("In veclient Recv new block", "block", block.Header().ID())
		c.comu.PostNewCenterBlockEvent(block)

	}
	return nil
}
