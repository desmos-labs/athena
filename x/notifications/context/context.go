package context

import (
	"github.com/forbole/juno/v5/modules/registrar"
	"github.com/forbole/juno/v5/node"
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
