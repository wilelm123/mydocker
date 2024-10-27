/*
Copyright © 2024 kube-devops@cisco.com
Copyright apply to this source code.
Check LICENSE for detail.
*/
package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wilelm123/mydocker/pkg/cgroups/subsystems"
	"github.com/wilelm123/mydocker/pkg/command"
	"strings"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Create a container with namespace and cgroups limit",
	Long: `Create a container with namespace and cgroups, for example:

	mydocker run -ti [image] [command]
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("missing container command")
		}
		imageName := args[0]
		cmdArr := args[1:]

		log.Infof(imageName)
		log.Infof(strings.Join(cmdArr, ","))

		if runOptions.enableTTY && runOptions.detach {
			return fmt.Errorf("ti and d parameters should not be specified at the same time")
		}

		resConf := &subsystems.ResourceConfig{
			MemoryLimit: runOptions.memoryLimit,
			CpuSet:      runOptions.cpuLimit,
			CpuShare:    runOptions.cpuShare,
		}
		log.Infof("create TTy %v", runOptions.enableTTY)

		command.Run(runOptions.enableTTY, cmdArr, resConf, runOptions.name, runOptions.volume, imageName, runOptions.environ, runOptions.network, runOptions.portMap)
		return nil
	},
}

type runCmdOptions struct {
	enableTTY   bool
	detach      bool
	memoryLimit string
	cpuShare    string
	cpuLimit    string
	name        string
	volume      string
	environ     []string
	network     string
	portMap     []string
}

var runOptions runCmdOptions

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	runCmd.Flags().BoolVarP(&runOptions.enableTTY, "tty", "t", false, "Enable tty")
	runCmd.Flags().BoolVarP(&runOptions.detach, "detach", "d", false, "detach container")
	runCmd.Flags().StringVarP(&runOptions.memoryLimit, "mem", "m", "", "memory limit")
	runCmd.Flags().StringVarP(&runOptions.cpuShare, "cpushare", "c", "", "cpushare limit")
	runCmd.Flags().StringVarP(&runOptions.cpuLimit, "cpuset", "", "", "cpuset limit")
	runCmd.Flags().StringVarP(&runOptions.volume, "volume", "v", "", "volume")
	runCmd.Flags().StringSliceVarP(&runOptions.environ, "env", "e", []string{}, "set environment")
	runCmd.Flags().StringVarP(&runOptions.network, "net", "n", "", "container network")
	runCmd.Flags().StringSliceVarP(&runOptions.portMap, "port", "p", []string{}, "port mapping")
}
