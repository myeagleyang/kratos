package main

import (
	"log"

	"github.com/spf13/cobra"

	"gitlab.wwgame.com/wwgame/kratos/cmd/kratos/v2/internal/change"
	"gitlab.wwgame.com/wwgame/kratos/cmd/kratos/v2/internal/project"
	"gitlab.wwgame.com/wwgame/kratos/cmd/kratos/v2/internal/proto"
	"gitlab.wwgame.com/wwgame/kratos/cmd/kratos/v2/internal/run"
	"gitlab.wwgame.com/wwgame/kratos/cmd/kratos/v2/internal/upgrade"
)

var rootCmd = &cobra.Command{
	Use:     "kratos",
	Short:   "Kratos: An elegant toolkit for Go microservices.",
	Long:    `Kratos: An elegant toolkit for Go microservices.`,
	Version: release,
}

func init() {
	rootCmd.AddCommand(project.CmdNew)
	rootCmd.AddCommand(proto.CmdProto)
	rootCmd.AddCommand(upgrade.CmdUpgrade)
	rootCmd.AddCommand(change.CmdChange)
	rootCmd.AddCommand(run.CmdRun)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
