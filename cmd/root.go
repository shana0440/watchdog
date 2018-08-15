package cmd

import (
	"io/ioutil"
	"log"
	"strings"

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
		ignoreGit()
		log.Println("ignores: ", strings.Join(ignores, ", "))
		d := dog.NewDirectory(".", ignores)
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
	rootCmd.Flags().BoolVarP(&silent, "silent", "s", false, "clear screen before command execute, if you have better name to describe this function, please let me know")
}

func ignoreGit() {
	ignores = append(ignores, ".git")
	bytes, err := ioutil.ReadFile(".gitignore")
	if err == nil {
		ignores = append(ignores, strings.Split(strings.Trim(string(bytes), "\n"), "\n")...)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
