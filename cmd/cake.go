package cmd

import (
	"github.com/fitzix/sniper-bot/consts"
	"github.com/fitzix/sniper-bot/runner"
	"github.com/spf13/cobra"
)

// cakeCmd represents the unicake command
var cakeCmd = &cobra.Command{
	Use:   "cake",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		runner.NewEthRunner().SniperUniCake(consts.ChainTypeBsc)
	},
}

func init() {
	rootCmd.AddCommand(cakeCmd)
}
