package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/fitzix/sniper-bot/consts"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "gen default config",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init config file")
		f, err := os.OpenFile("config.yml", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		if _, err := f.WriteString(consts.DefaultConfig); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
