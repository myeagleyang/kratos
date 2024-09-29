module gitlab.wwgame.com/wwgame/kratos/contrib/metrics/datadog/v2

go 1.16

require (
	github.com/DataDog/datadog-go v4.8.3+incompatible
	gitlab.wwgame.com/wwgame/kratos/v2 v2.6.2
)

require github.com/Microsoft/go-winio v0.5.2 // indirect

replace gitlab.wwgame.com/wwgame/kratos/v2 => ../../../
