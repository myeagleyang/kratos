package server

import (
	"bytes"
	"html/template"
)

//nolint:lll
var serviceTemplate = `
{{- /* delete empty line */ -}}
package service

import (
	{{- if .UseContext }}
	"context"
	{{- end }}
	{{- if .UseIO }}
	"io"
	{{- end }}
	
	pb "{{ .Package }}/api/{{ .ServiceLower }}/{{ .Mode }}/v1"
	"{{ .Package }}/app/{{ .ServiceLower }}/{{ .Mode }}/internal/biz"
	"gitlab.wwgame.com/wwgame/kratos/v2/log"
	{{- if .GoogleEmpty }}
	"google.golang.org/protobuf/types/known/emptypb"
	{{- end }}
)

type {{ .Service }}Service struct {
	pb.Unimplemented{{ .Service }}Server

	usecase *biz.{{ .Service }}UseCase
	log     *log.Helper
}

func New{{ .Service }}Service(uc *biz.{{ .Service }}UseCase, logger log.Logger) *{{ .Service }}Service {
	return &{{ .Service }}Service{usecase: uc, log: log.NewHelper(log.With(logger, "module", "{{ .Mode }}/{{ .ServiceLower }}"))}
}

{{- $s1 := "google.protobuf.Empty" }}
{{ range .Methods }}
{{- if eq .Type 1 }}
func (s *{{ .Service }}Service) {{ .Name }}(ctx context.Context, req {{ if eq .Request $s1 }}*emptypb.Empty{{ else }}*pb.{{ .Request }}{{ end }}) ({{ if eq .Reply $s1 }}*emptypb.Empty{{ else }}*pb.{{ .Reply }}{{ end }}, error) {
	return s.usecase.{{ .Name }}(ctx, req)
}

{{- else if eq .Type 2 }}
func (s *{{ .Service }}Service) {{ .Name }}(conn pb.{{ .Service }}_{{ .Name }}Server) error {
	for {
		req, err := conn.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		
		err = conn.Send(&pb.{{ .Reply }}{})
		if err != nil {
			return err
		}
	}
}

{{- else if eq .Type 3 }}
func (s *{{ .Service }}Service) {{ .Name }}(conn pb.{{ .Service }}_{{ .Name }}Server) error {
	for {
		req, err := conn.Recv()
		if err == io.EOF {
			return conn.SendAndClose(&pb.{{ .Reply }}{})
		}
		if err != nil {
			return err
		}
	}
}

{{- else if eq .Type 4 }}
func (s *{{ .Service }}Service) {{ .Name }}(req {{ if eq .Request $s1 }}*emptypb.Empty
{{ else }}*pb.{{ .Request }}{{ end }}, conn pb.{{ .Service }}_{{ .Name }}Server) error {
	for {
		err := conn.Send(&pb.{{ .Reply }}{})
		if err != nil {
			return err
		}
	}
}

{{- end }}
{{- end }}
`

var bizTemplate = `
{{- /* delete empty line */ -}}
package biz

import (
	{{- if .UseContext }}
	"context"
	{{- end }}
	{{- if .UseIO }}
	"io"
	{{- end }}
	
	pb "{{ .Package }}/api/{{ .ServiceLower }}/{{ .Mode }}/v1"
	"gitlab.wwgame.com/wwgame/kratos/v2/log"
	{{- if .GoogleEmpty }}
	"google.golang.org/protobuf/types/known/emptypb"
	{{- end }}
)

type {{ .Service }}Repo interface {

}


type {{ .Service }}UseCase struct {
	repo {{ .Service }}Repo
	log     *log.Helper
}

func New{{ .Service }}UseCase(repo {{ .Service }}Repo, logger log.Logger) *{{ .Service }}UseCase {
	uc := &{{ .Service }}UseCase{repo: repo, log: log.NewHelper(logger)}

	return uc
}

{{- $s1 := "google.protobuf.Empty" }}
{{ range .Methods }}

func (s *{{ .Service }}UseCase) {{ .Name }}(ctx context.Context, req {{ if eq .Request $s1 }}*emptypb.Empty{{ else }}*pb.{{ .Request }}{{ end }}) ({{ if eq .Reply $s1 }}*emptypb.Empty{{ else }}*pb.{{ .Reply }}{{ end }}, error) {
	return {{ if eq .Reply $s1 }}&emptypb.Empty{}{{ else }}&pb.{{ .Reply }}{}{{ end }}, nil
}

{{- end }}
`

var dataTemplate = `
{{- /* delete empty line */ -}}
package data

import (
	{{- if .UseContext }}
	//"context"
	{{- end }}
	{{- if .UseIO }}
	"io"
	{{- end }}

	//pb "{{ .Package }}/api/{{ .ServiceLower }}/{{ .Mode }}/v1"
	"gitlab.wwgame.com/chaoshe/blind_box/app/{{ .ServiceLower }}/service/internal/biz"
	"gitlab.wwgame.com/wwgame/kratos/v2/log"
	{{- if .GoogleEmpty }}
	"google.golang.org/protobuf/types/known/emptypb"
	{{- end }}
)

type {{ .ServiceLower }}Repo struct {
	data *Data
	log  *log.Helper
}

func New{{ .Service }}Repo(data *Data, logger log.Logger) biz.{{ .Service }}Repo {
	return &{{ .ServiceLower }}Repo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

`

var clientTemplate = `
package client

import (
	"context"
	"gitlab.wwgame.com/chaoshe/blind_box/api/{{ .ServiceLower }}/service/v1"
	"gitlab.wwgame.com/wwgame/kratos/v2/middleware/metadata"
	"gitlab.wwgame.com/wwgame/kratos/v2/middleware/recovery"
	"gitlab.wwgame.com/wwgame/kratos/v2/registry"
	"gitlab.wwgame.com/wwgame/kratos/v2/transport/grpc"
)

func New{{ .Service }}Client(r registry.Discovery) v1.{{ .Service }}Client {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///blind_box.{{ .ServiceLower }}.service"),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			metadata.Client(),
			recovery.Recovery(),
		),
	)
	if err != nil {
		panic(err)
	}
	return v1.New{{ .Service }}Client(conn)
}

`

type MethodType uint8

const (
	unaryType          MethodType = 1
	twoWayStreamsType  MethodType = 2
	requestStreamsType MethodType = 3
	returnsStreamsType MethodType = 4
)

// Service is a proto service.
type Service struct {
	Package      string
	Mode         string
	Service      string
	ServiceLower string
	Methods      []*Method
	GoogleEmpty  bool

	UseIO      bool
	UseContext bool
}

// Method is a proto method.
type Method struct {
	Service string
	Name    string
	Request string
	Reply   string

	// type: unary or stream
	Type MethodType
}

func (s *Service) execute() ([]byte, []byte, []byte, error) {
	const empty = "google.protobuf.Empty"
	svcBuf := new(bytes.Buffer)
	bizBuf := new(bytes.Buffer)
	dataBuf := new(bytes.Buffer)
	for _, method := range s.Methods {
		if (method.Type == unaryType && (method.Request == empty || method.Reply == empty)) ||
			(method.Type == returnsStreamsType && method.Request == empty) {
			s.GoogleEmpty = true
		}
		if method.Type == twoWayStreamsType || method.Type == requestStreamsType {
			s.UseIO = true
		}
		if method.Type == unaryType {
			s.UseContext = true
		}
	}

	tmpl, err := template.New("service").Parse(serviceTemplate)
	if err != nil {
		return nil, nil, nil, err
	}
	if err := tmpl.Execute(svcBuf, s); err != nil {
		return nil, nil, nil, err
	}

	tmpl, err = template.New("biz").Parse(bizTemplate)
	if err != nil {
		return nil, nil, nil, err
	}
	if err := tmpl.Execute(bizBuf, s); err != nil {
		return nil, nil, nil, err
	}

	tmpl, err = template.New("data").Parse(dataTemplate)
	if err != nil {
		return nil, nil, nil, err
	}
	if err := tmpl.Execute(dataBuf, s); err != nil {
		return nil, nil, nil, err
	}

	return svcBuf.Bytes(), bizBuf.Bytes(), dataBuf.Bytes(), nil
}
