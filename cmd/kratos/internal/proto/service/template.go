package service

import (
	"bytes"
	"fmt"
	"gitlab.wwgame.com/wwgame/kratos/cmd/kratos/v2/internal/proto/service/tpl"
	"html/template"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type MethodType uint8

const (
	unaryType          MethodType = 1
	twoWayStreamsType  MethodType = 2
	requestStreamsType MethodType = 3
	returnsStreamsType MethodType = 4
)

var (
	ModeService   = "service"
	ModeInterface = "interface"
	ModeAdmin     = "admin"
	ModeJob       = "job"
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

	TargetDir string
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

func (s *Service) execute() error {
	const empty = "google.protobuf.Empty"
	bf := new(bytes.Buffer)

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

	// main.go
	tmpl, err := template.New("cmd").Parse(tpl.MainTemplate)
	if err != nil {
		return err
	}
	if err = tmpl.Execute(bf, s); err != nil {
		return err
	}
	filename := filepath.Join(s.TargetDir, s.Mode, "cmd", "main.go")
	fmt.Println(filename)
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		log.Println(err.Error())
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Write(bf.Bytes())
		file.Close()
	} else if err != nil {
		return err
	}

	// wire.go
	bf = new(bytes.Buffer)
	tmpl, err = template.New("cmd").Parse(tpl.WireTemplate)
	if err != nil {
		return err
	}
	if err = tmpl.Execute(bf, s); err != nil {
		return err
	}
	filename = filepath.Join(s.TargetDir, s.Mode, "cmd", "wire.go")
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Write(bf.Bytes())
		file.Close()
	} else if err != nil {
		return err
	}
	fmt.Println(filename)
	// README.md
	filename = filepath.Join(s.TargetDir, s.Mode, "README.md")
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Write([]byte(fmt.Sprintf("# %s Service", s.Service)))
		file.Close()
	} else if err != nil {
		return err
	}
	fmt.Println(filename)
	// generate.go
	filename = filepath.Join(s.TargetDir, s.Mode, "generate.go")
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Write([]byte(tpl.GenerateTemplate))
		file.Close()
	} else if err != nil {
		return err
	}

	// Dockerfile
	filename = filepath.Join(s.TargetDir, s.Mode, "Dockerfile")
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Write([]byte(tpl.DockerfileTemplate))
		file.Close()
	} else if err != nil {
		return err
	}
	// Makefile
	filename = filepath.Join(s.TargetDir, s.Mode, "Makefile")
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Write([]byte(tpl.MakefileTemplate))
		file.Close()
	} else if err != nil {
		return err
	}
	// version
	filename = filepath.Join(s.TargetDir, s.Mode, "version")
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Write([]byte("v0.0.0"))
		file.Close()
	} else if err != nil {
		return err
	}

	// configs
	// config.yaml
	filename = filepath.Join(s.TargetDir, s.Mode, "configs", "config.yaml")
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Write([]byte(tpl.ConfigYaml))
		file.Close()
	} else if err != nil {
		return err
	}
	// registry.yaml
	filename = filepath.Join(s.TargetDir, s.Mode, "configs", "registry.yaml")
	fmt.Println(filename)
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Write([]byte(tpl.RegistryYaml))
		file.Close()
	} else if err != nil {
		return err
	}

	// biz.go
	bf = new(bytes.Buffer)
	tmpl, err = template.New("biz").Parse(tpl.BizTemplate)
	if err != nil {
		return err
	}
	if err = tmpl.Execute(bf, s); err != nil {
		return err
	}
	filename = filepath.Join(s.TargetDir, s.Mode, "internal", "biz", "biz.go")
	fmt.Println(filename)
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Write(bf.Bytes())
		file.Close()
	} else if err != nil {
		return err
	}

	// biz_impl.go
	bf = new(bytes.Buffer)
	tmpl, err = template.New("biz_impl").Parse(tpl.BizImplTemplate)
	if err != nil {
		return err
	}
	if err = tmpl.Execute(bf, s); err != nil {
		return err
	}
	filename = filepath.Join(s.TargetDir, s.Mode, "internal", "biz", s.ServiceLower+".go")
	fmt.Println(filename)
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Write(bf.Bytes())
		file.Close()
	} else if err != nil {
		return err
	}

	// README.md
	filename = filepath.Join(s.TargetDir, s.Mode, "internal", "biz", "README.md")
	fmt.Println(filename)
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Write([]byte("# Biz"))
		file.Close()
	} else if err != nil {
		return err
	}

	// client.go
	bf = new(bytes.Buffer)
	tmpl, err = template.New("client").Parse(tpl.ClientTemplate)
	if err != nil {
		return err
	}
	if err = tmpl.Execute(bf, s); err != nil {
		return err
	}
	filename = filepath.Join(s.TargetDir, s.Mode, "internal", "client", "client.go")
	fmt.Println(filename)
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Write(bf.Bytes())
		file.Close()
	} else if err != nil {
		return err
	}

	// conf.proto
	bf = new(bytes.Buffer)
	tmpl, err = template.New("conf").Parse(tpl.ConfTemplate)
	if err != nil {
		return err
	}
	if err = tmpl.Execute(bf, s); err != nil {
		return err
	}
	filename = filepath.Join(s.TargetDir, s.Mode, "internal", "conf", "conf.proto")
	fmt.Println(filename)
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Write(bf.Bytes())
		file.Close()
	} else if err != nil {
		return err
	}

	// data.go
	bf = new(bytes.Buffer)
	tmpl, err = template.New("data").Parse(tpl.DataTemplate)
	if err != nil {
		return err
	}
	if err = tmpl.Execute(bf, s); err != nil {
		return err
	}
	filename = filepath.Join(s.TargetDir, s.Mode, "internal", "data", "data.go")
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Write(bf.Bytes())
		file.Close()
	} else if err != nil {
		return err
	}

	// data_impl.go
	bf = new(bytes.Buffer)
	tmpl, err = template.New("data_impl").Parse(tpl.DataImplTemplate)
	if err != nil {
		return err
	}
	if err = tmpl.Execute(bf, s); err != nil {
		return err
	}
	filename = filepath.Join(s.TargetDir, s.Mode, "internal", "data", s.ServiceLower+".go")
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Write(bf.Bytes())
		file.Close()
	} else if err != nil {
		return err
	}

	// README.md
	filename = filepath.Join(s.TargetDir, s.Mode, "internal", "data", "README.md")
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Write([]byte("# %s Data"))
		file.Close()
	} else if err != nil {
		return err
	}

	// grpc.go
	bf = new(bytes.Buffer)
	tmpl, err = template.New("grpc").Parse(tpl.GrpcTemplate)
	if err != nil {
		return err
	}
	if err = tmpl.Execute(bf, s); err != nil {
		return err
	}
	filename = filepath.Join(s.TargetDir, s.Mode, "internal", "server", "grpc.go")
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Write(bf.Bytes())
		file.Close()
	} else if err != nil {
		return err
	}

	switch s.Mode {
	case ModeInterface, ModeAdmin:
		// grpc.go
		bf = new(bytes.Buffer)
		tmpl, err = template.New("http").Parse(tpl.HttpTemplate)
		if err != nil {
			return err
		}
		if err = tmpl.Execute(bf, s); err != nil {
			return err
		}
		filename = filepath.Join(s.TargetDir, s.Mode, "internal", "server", "http.go")
		if _, err = os.Stat(filename); os.IsNotExist(err) {
			file, err := os.Create(filename)
			if err != nil {
				return err
			}
			file.Write(bf.Bytes())
			file.Close()
		} else {
			return err
		}
	}

	// server.go
	bf = new(bytes.Buffer)
	tmpl, err = template.New("grpc").Parse(tpl.ServerTemplate)
	if err != nil {
		return err
	}
	if err = tmpl.Execute(bf, s); err != nil {
		return err
	}
	filename = filepath.Join(s.TargetDir, s.Mode, "internal", "server", "server.go")
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Write(bf.Bytes())
		file.Close()
	} else if err != nil {
		return err
	}

	// service.go
	bf = new(bytes.Buffer)
	tmpl, err = template.New("service").Parse(tpl.ServiceTemplate)
	if err != nil {
		return err
	}
	if err = tmpl.Execute(bf, s); err != nil {
		return err
	}
	filename = filepath.Join(s.TargetDir, s.Mode, "internal", "service", "service.go")
	fmt.Println(filename)
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Write(bf.Bytes())
		file.Close()
	} else if err != nil {
		return err
	}
	// service_impl.go
	bf = new(bytes.Buffer)
	tmpl, err = template.New("service_impl").Parse(tpl.ServiceImplTemplate)
	if err != nil {
		return err
	}
	if err = tmpl.Execute(bf, s); err != nil {
		return err
	}
	filename = filepath.Join(s.TargetDir, s.Mode, "internal", "service", s.ServiceLower+".go")
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Write(bf.Bytes())
		file.Close()
	} else if err != nil {
		return err
	}

	// test
	bf = new(bytes.Buffer)
	tmpl, err = template.New("test").Parse(tpl.TestTemplate)
	if err != nil {
		return err
	}
	if err = tmpl.Execute(bf, s); err != nil {
		return err
	}
	filename = filepath.Join(s.TargetDir, s.Mode, "test", "main_test.go")
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Write(bf.Bytes())
		file.Close()
	} else if err != nil {
		return err
	}

	// README.md
	filename = filepath.Join(s.TargetDir, s.Mode, "internal", "service", "README.md")
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Write([]byte("# %s Service"))
		file.Close()
	} else if err != nil {
		return err
	}

	cmd := exec.Command("make", "api", fmt.Sprintf("app=%s", s.ServiceLower))
	log.Println(cmd.String())
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(string(output), err.Error())
	}
	cmd = exec.Command("make", "config", fmt.Sprintf("app=%s", s.ServiceLower))
	log.Println(cmd.String())
	output, err = cmd.Output()
	if err != nil {
		log.Fatal(string(output), err.Error())
	}
	cmd = exec.Command("make", "wire", fmt.Sprintf("app=%s", s.ServiceLower), fmt.Sprintf("mode=%s", s.Mode))
	log.Println(cmd.String())
	output, err = cmd.Output()
	if err != nil {
		log.Fatal(string(output), err.Error())
	}
	return nil
}
