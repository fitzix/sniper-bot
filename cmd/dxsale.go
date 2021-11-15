package cmd

import (
	"github.com/fitzix/sniper-bot/consts"
	"github.com/fitzix/sniper-bot/runner"
	"github.com/spf13/cobra"
)

// dxsaleCmd represents the dxsale command
var dxsaleCmd = &cobra.Command{
	Use:   "dxsale",
	Short: "sniper dxsale",
	Run: func(cmd *cobra.Command, args []string) {
		runner.NewEthRunner().SniperDxsale(consts.ChainTypeBsc)
	},
}

func init() {
	rootCmd.AddCommand(dxsaleCmd)
}
