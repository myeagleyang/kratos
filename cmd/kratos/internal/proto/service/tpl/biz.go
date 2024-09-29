package tpl

var (
	BizTemplate = `package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(New{{ .Service }}UseCase)
`
	BizImplTemplate = `package biz

import (
	"context"
	
	pb "gitlab.wwgame.com/chaoshe/blind_box/api/{{ .ServiceLower }}/{{ .Mode }}/v1"
	"gitlab.wwgame.com/chaoshe/blind_box/app/{{ .ServiceLower }}/{{ .Mode }}/internal/conf"
	"gitlab.wwgame.com/chaoshe/blind_box/pkg/common"
	"gitlab.wwgame.com/wwgame/kratos/v2/log"
)

type {{ .Service }}Repo interface {

}

type {{ .Service }}UseCase struct {
	log *log.Helper
	cfg *conf.Service

	repo {{ .Service }}Repo
}

func New{{ .Service }}UseCase(repo {{ .Service }}Repo, cfg *conf.Service, logger log.Logger) *{{ .Service }}UseCase {
	uc := &{{ .Service }}UseCase{repo: repo, cfg: cfg, log: log.NewHelper(logger)}
	return uc
}

{{ range .Methods }}

func (s *{{ .Service }}UseCase) {{ .Name }}(ctx context.Context, req *pb.{{ .Request }}) (reply *pb.{{ .Reply }}, err error) {
	reply = &pb.{{ .Reply }}{}	
	curUser, err := common.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}
	_ = curUser
	
	return 
}

{{- end }}
`
)
