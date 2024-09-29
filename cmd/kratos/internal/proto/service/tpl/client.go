package tpl

var ClientTemplate = `package client

import (
	"github.com/google/wire"
	consulAPI "github.com/hashicorp/consul/api"
	"gitlab.wwgame.com/chaoshe/blind_box/app/{{ .ServiceLower }}/service/internal/conf"
	"gitlab.wwgame.com/wwgame/kratos/contrib/registry/consul/v2"
	"gitlab.wwgame.com/wwgame/kratos/v2/registry"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewDiscovery,)

func NewDiscovery(conf *conf.Registry) registry.Discovery {
	c := consulAPI.DefaultConfig()
	c.Address = conf.Consul.Address
	c.Scheme = conf.Consul.Scheme
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(false))
	return r
}

`

var IDClientTemplate = `package client

import (
	"context"
	"gitlab.wwgame.com/chaoshe/blind_box/api/idgenerator/service/v1"
	"gitlab.wwgame.com/wwgame/kratos/v2/middleware/metadata"
	"gitlab.wwgame.com/wwgame/kratos/v2/middleware/recovery"
	"gitlab.wwgame.com/wwgame/kratos/v2/registry"
	"gitlab.wwgame.com/wwgame/kratos/v2/transport/grpc"
)

func NewIDGeneratorClient(r registry.Discovery) v1.IDGeneratorClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///blind_box.idgenerator.service"),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			metadata.Client(),
			recovery.Recovery(),
		),
	)
	if err != nil {
		panic(err)
	}
	return v1.NewIDGeneratorClient(conn)
}
`
