package proto

import (
	"github.com/spf13/cobra"
	"gitlab.wwgame.com/wwgame/kratos/cmd/kratos/v2/internal/proto/add"
	"gitlab.wwgame.com/wwgame/kratos/cmd/kratos/v2/internal/proto/client"
	"gitlab.wwgame.com/wwgame/kratos/cmd/kratos/v2/internal/proto/router"
	"gitlab.wwgame.com/wwgame/kratos/cmd/kratos/v2/internal/proto/server"
	"gitlab.wwgame.com/wwgame/kratos/cmd/kratos/v2/internal/proto/service"
)

// CmdProto represents the proto command.
var CmdProto = &cobra.Command{
	Use:   "proto",
	Short: "Generate the proto files",
	Long:  "Generate the proto files.",
}

func init() {
	CmdProto.AddCommand(add.CmdAdd)
	CmdProto.AddCommand(client.CmdClient)
	CmdProto.AddCommand(server.CmdServer)
	CmdProto.AddCommand(router.CmdServer)
	CmdProto.AddCommand(service.CmdService)
}
