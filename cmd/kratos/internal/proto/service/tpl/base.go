package tpl

var (
	GenerateTemplate = `package generate

//go:generate kratos proto client .
`
	MakefileTemplate = `include ../../../app_makefile
`
	DockerfileTemplate = `FROM 192.168.8.90:5000/ubuntu

# 配置时区
ENV TZ=Asia/Shanghai

RUN mkdir -p /opt/service
WORKDIR /opt/service

COPY cmd/cmd /opt/service/

CMD ["/opt/service/cmd", "--conf", "./configs"]

`
)
