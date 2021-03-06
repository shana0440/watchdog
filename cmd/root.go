package cmd

import (
	"log"
	"strings"

	"github.com/shana0440/watchdog/dog"
	"github.com/shana0440/watchdog/helper"
	"github.com/spf13/cobra"
)

var command string
var ignores []string
var matchs []string
var silent bool

var rootCmd = &cobra.Command{
	Use:   "watchdog",
	Short: "watchdog can help you run script when file change",
	Long: `A file watcher can execute command when file change.
	built with shana0440. source code at https://github.com/shana0440/watchdog`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ignores = append(ignores, helper.IgnoreGit()...)
		log.Println("ignores: ", strings.Join(ignores, ", "))
		d := dog.NewDirectory(".", ignores, matchs)
		c := dog.NewCommand(silent)
		w := dog.NewWatch()
		dog := dog.NewDog(d, c, w)
		defer dog.Close()
		err := dog.Run(command)
		if err != nil {
			log.Fatalln("can't watch files", err)
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&command, "command", "c", "", "the command you want execute when file changed (required)")
	rootCmd.MarkFlagRequired("command")
	rootCmd.Flags().StringArrayVarP(&ignores, "ignore", "i", []string{}, "the file or directory you don't want to trigger command")
	rootCmd.Flags().StringArrayVarP(&matchs, "match", "m", []string{}, "the glob pattern you want to trigger command, watchdog will watch whole directory if not specify")
	rootCmd.Flags().BoolVarP(&silent, "silent", "s", false, "clear screen before command execute, if you have better name to describe this function, please let me know")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
