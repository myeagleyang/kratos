package server

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

// CmdServer the service command.
var CmdServer = &cobra.Command{
	Use:   "server",
	Short: "Generate the proto server implementations",
	Long:  "Generate the proto server implementations. Example: kratos proto server api/xxx.proto --target-dir=internal/service",
	Run:   run,
}
var (
	targetDir string
	mode      string
)

func init() {
	CmdServer.Flags().StringVarP(&targetDir, "target-dir", "t", "internal", "generate target directory")
	CmdServer.Flags().StringVarP(&mode, "mode", "m", "service", "[service | interface | admin]")
}

func run(_ *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Please specify the proto file. Example: kratos proto server api/xxx.proto")
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
	serviceDir := filepath.Join(targetDir, "service")
	bizDir := filepath.Join(targetDir, "biz")
	dataDir := filepath.Join(targetDir, "data")

	if _, err := os.Stat(serviceDir); os.IsNotExist(err) {
		fmt.Printf("service directory: %s does not exsit\n", serviceDir)
		return
	}
	if _, err := os.Stat(bizDir); os.IsNotExist(err) {
		fmt.Printf("biz directory: %s does not exsit\n", bizDir)
		return
	}
	for _, s := range res {
		to1 := filepath.Join(serviceDir, strings.ToLower(s.Service)+".go")
		if _, err := os.Stat(to1); !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "%s already exists: %s\n", s.Service, to1)
			continue
		}
		to2 := filepath.Join(bizDir, strings.ToLower(s.Service)+".go")
		if _, err := os.Stat(to2); !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "%s already exists: %s\n", s.Service, to2)
			continue
		}
		to3 := filepath.Join(dataDir, strings.ToLower(s.Service)+".go")
		if _, err := os.Stat(to3); !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "%s already exists: %s\n", s.Service, to3)
			continue
		}
		b1, b2, b3, err := s.execute()
		if err != nil {
			log.Fatal(err)
		}
		if err := os.WriteFile(to1, b1, 0o644); err != nil {
			log.Fatal(err)
		}
		fmt.Println(to1)

		if err := os.WriteFile(to2, b2, 0o644); err != nil {
			log.Fatal(err)
		}
		fmt.Println(to2)

		if err := os.WriteFile(to3, b3, 0o644); err != nil {
			log.Fatal(err)
		}
		fmt.Println(to3)

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
