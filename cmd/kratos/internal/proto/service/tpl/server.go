package tpl

var (
	GrpcTemplate = `package server

import (
	pb "gitlab.wwgame.com/chaoshe/blind_box/api/{{ .ServiceLower }}/{{ .Mode }}/v1"
	"gitlab.wwgame.com/chaoshe/blind_box/app/{{ .ServiceLower }}/{{ .Mode }}/internal/conf"
	"gitlab.wwgame.com/chaoshe/blind_box/app/{{ .ServiceLower }}/{{ .Mode }}/internal/service"
	"gitlab.wwgame.com/wwgame/kratos/v2/log"
	"gitlab.wwgame.com/wwgame/kratos/v2/middleware/logging"
	"gitlab.wwgame.com/wwgame/kratos/v2/middleware/metadata"
	"gitlab.wwgame.com/wwgame/kratos/v2/middleware/recovery"
	"gitlab.wwgame.com/wwgame/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, svc *service.{{ .Service }}Service, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
			metadata.Server(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	pb.Register{{ .Service }}Server(srv, svc)
	return srv
}
`
	HttpTemplate = `package server

import (
	"gitlab.wwgame.com/chaoshe/blind_box/app/{{ .ServiceLower }}/{{ .Mode }}/internal/conf"
	"gitlab.wwgame.com/chaoshe/blind_box/app/{{ .ServiceLower }}/{{ .Mode }}/internal/middleware/validate"
	"gitlab.wwgame.com/wwgame/kratos/v2/log"
	"gitlab.wwgame.com/wwgame/kratos/v2/middleware/logging"
	"gitlab.wwgame.com/wwgame/kratos/v2/middleware/ratelimit"
	"gitlab.wwgame.com/wwgame/kratos/v2/middleware/recovery"
	"gitlab.wwgame.com/wwgame/kratos/v2/transport/http"

	"github.com/gorilla/handlers"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
			ratelimit.Server(),
			validate.Validator(),
		),
		http.Filter(
			handlers.CORS(
				handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
				handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"}),
				handlers.AllowedOrigins([]string{"*"}),
			),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	pb.Register{{ .Service }}HTTPServer(srv, svc)
	return srv
}
`
	ServerTemplate = `package server

import (
	"github.com/google/wire"
	consulAPI "github.com/hashicorp/consul/api"
	"gitlab.wwgame.com/chaoshe/blind_box/app/{{ .ServiceLower }}/{{ .Mode }}/internal/conf"
	"gitlab.wwgame.com/wwgame/kratos/contrib/registry/consul/v2"
	"gitlab.wwgame.com/wwgame/kratos/v2/registry"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGRPCServer, NewRegistrar)

func NewRegistrar(conf *conf.Registry) registry.Registrar {
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
)
