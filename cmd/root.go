package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "watchdog",
	Short: "watchdog can help you run script when file change",
	Long: `A file watcher can execute command when file change.
	built with shana0440. source code at https://github.com/shana0440/watchdog`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		log.Println(cmd)
		log.Println(args)
	},
}

func init() {
	var cmd string
	var ignores []string
	var clear bool
	rootCmd.Flags().StringVarP(&cmd, "command", "c", "", "the command you want execute when file changed (required)")
	rootCmd.MarkFlagRequired("command")
	rootCmd.Flags().StringArrayVarP(&ignores, "ignore", "i", []string{}, "the file or directory you don't want to trigger command")
	rootCmd.Flags().BoolVarP(&clear, "slient", "s", false, "clear screen before command execute, if you have better name to describe this function, please let me know")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
