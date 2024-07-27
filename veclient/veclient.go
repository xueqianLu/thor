package veclient

import (
	"context"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/inconshreveable/log15"
	"github.com/vechain/thor/block"
	"github.com/vechain/thor/comm"
	pb "github.com/xueqianLu/vehackcenter/hackcenter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"strconv"
	"time"
)

var log = log15.New("pkg", "veClient")

type VeClient struct {
	proposer string
	index    int
	comu     *comm.Communicator
	conn     pb.CenterServiceClient
	nodes    map[string]string
}

func NewP2PCenterClient(proposer string, comu *comm.Communicator) P2pCenterClient {
	serverUrl := os.Getenv("VE_P2P_SERVER_URL")
	if serverUrl == "" {
		log.Error("VE_P2P_SERVER_URL not set")
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
	client.proposer = proposer
	client.conn = pb.NewCenterServiceClient(conn)
	client.comu = comu
	client.nodes = make(map[string]string)
	return client
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
	client.nodes = make(map[string]string)
	return client
}

func (c *VeClient) BroadcastBlock(blk *block.Block) error {
	pbblk := new(pb.Block)
	var err error
	pbblk.Hash = blk.Header().ID().String()
	pbblk.Height = int64(blk.Header().Number())
	pbblk.Timestamp = int64(blk.Header().Timestamp())
	pbblk.Data, err = rlp.EncodeToBytes(blk)
	if err != nil {
		log.Error("BroadcastBlock encode block failed", "err", err)
		return err
	}

	pbblk.Proposer = new(pb.Proposer)
	pbblk.Proposer.Proposer = c.proposer
	pbblk.Proposer.Index = int32(c.index)
	log.Info("In veclient broadcast block", "number", blk.Header().Number(), "id", blk.Header().ID().String())

	_, err = c.conn.BroadcastBlock(context.TODO(), pbblk)
	return err

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
		log.Info("In veclient broadcast hacked block", "number", block.Header().Number(), "id", block.Header().ID().String())
		c.comu.BroadcastBlock(block)
	}
	return nil
}

func (c *VeClient) RegisterNode() {
	srv := c.comu.P2PServer()
	nodes, err := c.conn.RegisterNode(context.TODO(), &pb.NodeRegisterInfo{
		Node: srv.Self().String(),
	})
	if err != nil {
		log.Error("RegisterNode failed", "err", err)
		return
	}

	addNode := func(node string) error {
		n, err := discover.ParseNode(node)
		if err != nil {
			log.Error("parse node failed", "err", err, "node", node)
			return err
		}
		srv.AddStatic(n)
		return nil
	}
	for _, node := range nodes.Nodes {
		if _, ok := c.nodes[node]; !ok {
			if err := addNode(node); err == nil {
				c.nodes[node] = node
			}
		}
	}
	t := time.NewTicker(time.Second * 2)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			nodes, err := c.conn.FetchNode(context.TODO(), &pb.FetchNodeRequest{
				Self: srv.Self().String(),
			})
			if err != nil {
				log.Error("FetchNode failed", "err", err)
				continue
			}
			for _, node := range nodes.Nodes {
				if _, ok := c.nodes[node]; !ok {
					if err := addNode(node); err == nil {
						c.nodes[node] = node
					}
				}
			}
		}
	}
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
		log.Info("In veclient SubscribeBlock", "block", block.Header().Number(), "id", block.Header().ID().String())
		c.comu.PostNewCenterBlockEvent(block)

	}
	return nil
}

func (c *VeClient) SubscribeHackedBlock() error {
	in := new(pb.SubscribeBlockRequest)
	in.Proposer = c.proposer
	sub, err := c.conn.SubscribeMinedBlock(context.TODO(), in)
	if err != nil {
		log.Error("SubscribeMinedBlock failed", "err", err)
		return err
	}
	for {
		msg, err := sub.Recv()
		if err != nil {
			log.Error("SubscribeMinedBlock Recv failed", "err", err)
			return err
		}

		block := new(block.Block)
		err = rlp.DecodeBytes(msg.Data, block)
		if err != nil {
			log.Error("SubscribeMinedBlock decode block failed", "err", err)
			continue
		}
		log.Info("In veclient SubscribeMinedBlock", "block", block.Header().Number(), "id", block.Header().ID().String())
		c.comu.PostNewHackedBlockEvent(block)
	}
	return nil
}

func (c *VeClient) Vote(height int64, vote bool) bool {
	res, err := c.conn.Vote(context.TODO(), &pb.VoteRequest{
		Block: height,
	})
	if err != nil {
		log.Error("Get vote failed", "err", err)
		return vote
	}
	return res.Vote == 1 && vote
}
