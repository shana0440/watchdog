package cmd

import (
	"log"

	"github.com/shana0440/watchdog/dog"
	"github.com/spf13/cobra"
)

var command string
var ignores []string
var silent bool

var rootCmd = &cobra.Command{
	Use:   "watchdog",
	Short: "watchdog can help you run script when file change",
	Long: `A file watcher can execute command when file change.
	built with shana0440. source code at https://github.com/shana0440/watchdog`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		d := dog.NewDirectory(".", ignoreSliceToIgnoreMap(ignores))
		c := dog.NewCommand(silent)
		dog, err := dog.NewDog(d, c)
		if err != nil {
			log.Fatalln("can't create watcher", err)
		}
		defer dog.Close()
		err = dog.Run(command)
		if err != nil {
			log.Fatalln("can't watch files", err)
		}
	},
}

func ignoreSliceToIgnoreMap(arr []string) map[string]struct{} {
	ignores := make(map[string]struct{})
	for _, item := range arr {
		ignores[item] = struct{}{}
	}
	return ignores
}

func init() {
	rootCmd.Flags().StringVarP(&command, "command", "c", "", "the command you want execute when file changed (required)")
	rootCmd.MarkFlagRequired("command")
	rootCmd.Flags().StringArrayVarP(&ignores, "ignore", "i", []string{}, "the file or directory you don't want to trigger command")
	rootCmd.Flags().BoolVarP(&silent, "silent", "s", false, "clear screen before command execute, if you have better name to describe this function, please let me know")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
