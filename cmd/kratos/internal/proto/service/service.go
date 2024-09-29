package service

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/emicklei/proto"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// CmdService the service command.
var CmdService = &cobra.Command{
	Use:   "service",
	Short: "Generate the service file  implementations",
	Long:  "Generate the service file implementations. Example: kratos proto service api/xxx.proto --target-dir=internal",
	Run:   run,
}
var (
	targetDir string
	mode      string
)

func init() {
	CmdService.Flags().StringVarP(&targetDir, "target-dir", "t", "internal", "generate target directory")
	CmdService.Flags().StringVarP(&mode, "mode", "m", "service", "[service | interface | admin]")
}

func run(_ *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Please specify the proto file. Example: kratos proto server api/xxx.proto")
		return
	}
	if mode != "service" {
		fmt.Fprintln(os.Stderr, "generate file only support service mode")
		return
	}
	reader, err := os.Open(args[0])
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		log.Fatal(err)
	}

	var (
		pkg string

		res []*Service
	)
	proto.Walk(definition,
		proto.WithOption(func(o *proto.Option) {
			if o.Name == "go_package" {
				split := strings.Split(o.Constant.Source, ";")
				pkg = split[0]
			}
			// 用 objc_class_prefix 指定项目目录,拼import
			if o.Name == "objc_class_prefix" {
				pkg = o.Constant.Source
			}

		}),
		proto.WithService(func(s *proto.Service) {
			name := serviceName(s.Name)
			low := strings.ToLower(name)
			cs := &Service{
				Package:      pkg,
				Mode:         mode,
				Service:      name,
				ServiceLower: low,
				TargetDir:    targetDir,
			}
			for _, e := range s.Elements {
				r, ok := e.(*proto.RPC)
				if !ok {
					continue
				}
				cs.Methods = append(cs.Methods, &Method{
					Service: serviceName(s.Name), Name: serviceName(r.Name), Request: parametersName(r.RequestType),
					Reply: parametersName(r.ReturnsType), Type: getMethodType(r.StreamsRequest, r.StreamsReturns),
				})
			}
			res = append(res, cs)
		}),
	)
	serviceDir := filepath.Join(targetDir, mode, "internal/service")
	bizDir := filepath.Join(targetDir, mode, "internal/biz")
	dataDir := filepath.Join(targetDir, mode, "internal/data")
	clientDir := filepath.Join(targetDir, mode, "internal/client")
	confDir := filepath.Join(targetDir, mode, "internal/conf")
	serverDir := filepath.Join(targetDir, mode, "internal/server")
	testDir := filepath.Join(targetDir, mode, "test")
	cmdDir := filepath.Join(targetDir, mode, "cmd")
	configsDir := filepath.Join(targetDir, mode, "configs")

	if _, err := os.Stat(cmdDir); os.IsNotExist(err) {
		os.MkdirAll(cmdDir, 0666)
	}
	if _, err := os.Stat(configsDir); os.IsNotExist(err) {
		os.MkdirAll(configsDir, 0666)
	}
	if _, err := os.Stat(bizDir); os.IsNotExist(err) {
		os.MkdirAll(bizDir, 0666)
	}
	if _, err := os.Stat(clientDir); os.IsNotExist(err) {
		os.MkdirAll(clientDir, 0666)
	}
	if _, err := os.Stat(confDir); os.IsNotExist(err) {
		os.MkdirAll(confDir, 0666)
	}
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		os.MkdirAll(dataDir, 0666)
	}
	if _, err := os.Stat(serverDir); os.IsNotExist(err) {
		os.MkdirAll(serverDir, 0666)
	}
	if _, err := os.Stat(serviceDir); os.IsNotExist(err) {
		os.MkdirAll(serviceDir, 0666)
	}
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		os.MkdirAll(testDir, 0666)
	}

	for _, s := range res {
		err = s.execute()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getMethodType(streamsRequest, streamsReturns bool) MethodType {
	if !streamsRequest && !streamsReturns {
		return unaryType
	} else if streamsRequest && streamsReturns {
		return twoWayStreamsType
	} else if streamsRequest {
		return requestStreamsType
	} else if streamsReturns {
		return returnsStreamsType
	}
	return unaryType
}

func parametersName(name string) string {
	return strings.ReplaceAll(name, ".", "_")
}

func serviceName(name string) string {
	return toUpperCamelCase(strings.Split(name, ".")[0])
}

func toUpperCamelCase(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	s = cases.Title(language.Und, cases.NoLower).String(s)
	return strings.ReplaceAll(s, " ", "")
}
