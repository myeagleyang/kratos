module gitlab.wwgame.com/wwgame/kratos/contrib/config/kubernetes/v2

go 1.16

require (
	gitlab.wwgame.com/wwgame/kratos/v2 v2.6.2
	k8s.io/api v0.26.3
	k8s.io/apimachinery v0.26.3
	k8s.io/client-go v0.26.3
)

replace gitlab.wwgame.com/wwgame/kratos/v2 => ../../../
