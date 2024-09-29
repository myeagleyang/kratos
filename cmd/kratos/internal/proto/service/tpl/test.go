package tpl

var TestTemplate = `package test

import (
	"gitlab.wwgame.com/chaoshe/blind_box/app/{{ .ServiceLower }}/{{ .Mode }}/internal/conf"
	"gitlab.wwgame.com/wwgame/kratos/v2/config"
	"gitlab.wwgame.com/wwgame/kratos/v2/config/file"
	"os"
	"testing"
)
var testConfData *conf.Bootstrap
var testRegisterConf *conf.Registry

func TestMain(m *testing.M) {
	// 初始化配置
	c := config.New(
		config.WithSource(
			file.NewSource("../configs/"),
		),
	)
	defer c.Close()
	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	testConfData = &bc
	var rc conf.Registry
	if err := c.Scan(&rc); err != nil {
		panic(err)
	}
	testRegisterConf = &rc
	exitCode := m.Run()
	// 退出
	os.Exit(exitCode)
}

`
