/*
Copyright © 2024 kube-devops@cisco.com
Copyright apply to this source code.
Check LICENSE for detail.

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Create a container with namespace and cgroups limit",
	Long: `Create a container with namespace and cgroups, for example:

	mydocker run -ti [image] [command]
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")
	},
}

type runCmdOptions {
	enableTTY	bool
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	runCmd.Flags().BoolVarP()
}
