/*
	The package support the slave mode and sends statistic to master node.
*/
package slavenode

import (
	"context"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/iostrovok/conveyor/protobuf/go/nodes"
)

// Max Size of message for Lambda client
const MaxMsgSize = 10 * 1024 * 1024 // max message size 10 MB

func dialOption() []grpc.DialOption {

	kp := keepalive.ClientParameters{
		/*
			After a duration of this time if the client doesn't see any activity it
			pings the server to see if the transport is still alive.
			If set below 10s, a minimum value of 10s will be used instead.
		*/
		Time: 10 * time.Second,

		/*
			After having pinged for keepalive check, the client waits for a duration
			of Timeout and if no activity is seen even after that the connection is
			closed.
		*/
		Timeout: 30 * time.Second,

		/*
			If true, client sends keepalive pings even with no active RPCs. If false,
			when there are no active RPCs, Time and Timeout will be ignored
			and no keepalive pings will be sent.
		*/
		PermitWithoutStream: true,
	}

	options := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(kp),
		grpc.WithDefaultCallOptions(
			grpc.WaitForReady(false),
			grpc.MaxCallRecvMsgSize(MaxMsgSize),
			grpc.MaxCallSendMsgSize(MaxMsgSize),
		),
	}

	return options
}

type SlaveNode struct {
	sync.RWMutex

	conn                    *grpc.ClientConn
	client                  nodes.MasterNodeClient
	lastConnectError        error
	host, clusterID, nodeID string

	data    []string
	isError bool
}

func New(host string) (*SlaveNode, error) {

	s := &SlaveNode{
		host: host,
	}

	err := s.connection()

	return s, err
}

func (s *SlaveNode) connection() error {
	conn, err := grpc.Dial(s.host, dialOption()...)
	if err == nil {
		s.conn = conn
		s.client = nodes.NewMasterNodeClient(s.conn)
	}

	return err
}

func (s *SlaveNode) Send(ctx context.Context, request *nodes.SlaveNodeInfoRequest) (*nodes.SimpleResult, error) {
	return s.client.UpdateNodeInfo(ctx, request)
}
