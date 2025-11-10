package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
)

type GrpcFactory struct {
	once     sync.Once
	err      error
	dns      string
	instance *grpc.ClientConn
}

func NewGrpc(dns string) *GrpcFactory {
	return &GrpcFactory{
		dns: dns,
	}
}

func (f *GrpcFactory) Get() (*grpc.ClientConn, error) {
	f.once.Do(func() {
		f.instance, f.err = grpc.NewClient(f.dns, grpc.WithTransportCredentials(insecure.NewCredentials()))
	})
	return f.instance, f.err
}
