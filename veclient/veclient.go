package veclient

import (
	"context"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/inconshreveable/log15"
	"github.com/vechain/thor/block"
	pb "github.com/xueqianLu/vehackcenter/hackcenter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

var log = log15.New("pkg", "veClient")

type VeClient struct {
	conn pb.CenterServiceClient
}

func getVeClient() *VeClient {
	serverUrl := os.Getenv("VE_HACK_SERVER_URL")
	if serverUrl == "" {
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
	client.conn = pb.NewCenterServiceClient(conn)
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
	signer, _ := blk.Header().Signer()
	pbblk.Proposer.Proposer = signer.String()
	pbblk.Proposer.Index = 0 // todo: change to real index

	return c.conn.SubmitBlock(context.TODO(), pbblk)
}

func (c *VeClient) SubBroadcastTask() error {
	proposer := ""
	sub, err := c.conn.SubBroadcastTask(context.TODO(), &pb.SubBroadcastTaskRequest{
		Proposer: proposer,
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
	}
	return nil
}

func (c *VeClient) SubscribeBlock() error {
	sub, err := c.conn.SubscribeBlock(nil, nil)
	if err != nil {
		log.Error("SubscribeBlock failed", "err", err)
		return err
	}
	for {
		block, err := sub.Recv()
		if err != nil {
			log.Error("SubscribeBlock Recv failed", "err", err)
			return err
		}
		log.Info("SubscribeBlock Recv", "block", block)

	}
	return nil
}
