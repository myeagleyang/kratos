module gitlab.wwgame.com/wwgame/kratos/contrib/registry/kubernetes/v2

go 1.16

require (
	gitlab.wwgame.com/wwgame/kratos/v2 v2.6.2
	github.com/json-iterator/go v1.1.12
	k8s.io/api v0.24.3
	k8s.io/apimachinery v0.24.3
	k8s.io/client-go v0.24.3
)

replace gitlab.wwgame.com/wwgame/kratos/v2 => ../../../
