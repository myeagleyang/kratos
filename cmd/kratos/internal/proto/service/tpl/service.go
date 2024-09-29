package tpl

var (
	ServiceTemplate = `package service

import (
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(New{{ .Service }}Service)
`
	ServiceImplTemplate = `package service

import (
	"context"
	
	pb "gitlab.wwgame.com/chaoshe/blind_box/api/{{ .ServiceLower }}/{{ .Mode }}/v1"
	"gitlab.wwgame.com/chaoshe/blind_box/app/{{ .ServiceLower }}/{{ .Mode }}/internal/biz"
	"gitlab.wwgame.com/chaoshe/blind_box/app/{{ .ServiceLower }}/{{ .Mode }}/internal/conf"
	"gitlab.wwgame.com/wwgame/kratos/v2/log"
)

type {{ .Service }}Service struct {
	pb.Unimplemented{{ .Service }}Server
	cfg *conf.Service
	log *log.Helper

	usecase *biz.{{ .Service }}UseCase
}

func New{{ .Service }}Service(uc *biz.{{ .Service }}UseCase, cfg *conf.Service, logger log.Logger) *{{ .Service }}Service {
	return &{{ .Service }}Service{usecase: uc, cfg: cfg, log: log.NewHelper(log.With(logger, "module", "{{ .Mode }}/{{ .ServiceLower }}"))}
}

{{ range .Methods }}

func (s *{{ .Service }}Service) {{ .Name }}(ctx context.Context, req *pb.{{ .Request }}) (*pb.{{ .Reply }}, error) {
	return s.usecase.{{ .Name }}(ctx, req)
}

{{- end }}
`
)
