package context

import (
	"github.com/forbole/juno/v4/modules/registrar"
	"github.com/forbole/juno/v4/node"
	"google.golang.org/grpc"
)

type Context struct {
	registrar.Context
	Node           node.Node
	GRPCConnection *grpc.ClientConn
}

func NewContext(registrarContext registrar.Context, node node.Node, grpcConnection *grpc.ClientConn) Context {
	return Context{
		Context:        registrarContext,
		Node:           node,
		GRPCConnection: grpcConnection,
	}
}
