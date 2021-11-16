package cmd

import (
	"github.com/fitzix/sniper-bot/consts"
	"github.com/fitzix/sniper-bot/runner"
	"github.com/spf13/cobra"
)

// cakeCmd represents the unicake command
var uniCmd = &cobra.Command{
	Use:   "uni",
	Short: "sniper on uniswap v2",
	Run: func(cmd *cobra.Command, args []string) {
		runner.NewEthRunner().SniperUniCake(consts.ChainTypeEth)
	},
}

func init() {
	rootCmd.AddCommand(uniCmd)
}
