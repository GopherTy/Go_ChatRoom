package cmd

import (
	"chatroom/cmd/daemon"
	"chatroom/configure"
	"chatroom/logger"
	"chatroom/utils"
	"log"

	"github.com/spf13/cobra"
)

func init() {
	var filename string
	basePaht := utils.BasePath()
	cmd := &cobra.Command{
		Use:   "daemon",
		Short: "run as daemon",
		Run: func(cmd *cobra.Command, args []string) {
			// load configure
			cnf := configure.Single()
			e := cnf.Load(filename)
			if e != nil {
				log.Fatalln(e)
			}
			e = cnf.Format(basePaht)
			if e != nil {
				log.Fatalln(e)
			}

			// init logger
			e = logger.Init(basePaht, &cnf.Logger)
			if e != nil {
				log.Fatalln(e)
			}

			// run
			daemon.Run()
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&filename, "config",
		"c",
		utils.Abs(basePaht, "chatroom.jsonnet"),
		"configure file",
	)
	rootCmd.AddCommand(cmd)
}
