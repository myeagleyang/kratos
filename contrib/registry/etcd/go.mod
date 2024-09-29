module gitlab.wwgame.com/wwgame/kratos/contrib/registry/etcd/v2

go 1.16

require (
	gitlab.wwgame.com/wwgame/kratos/v2 v2.6.2
	go.etcd.io/etcd/client/v3 v3.5.8
	google.golang.org/grpc v1.50.1
)

replace gitlab.wwgame.com/wwgame/kratos/v2 => ../../../
