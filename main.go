package main

import (
	"os"
	"path"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gpupod",
	Short: "gpupod is a tool to list/watch pod with nvidia gpu resources.",
	Long: "gpupod get cluster gpu usage info, which pod is occupy gpu resource\n" +
		"and the pod info include namespace, image etc.",
	Run: func(cmd *cobra.Command, args []string) {
		listWatchGpuPod()
	},
}

var createdTime bool
var kubeConfig string
var reason bool
var successPod bool
var watch bool

func main() {
	home, _ := os.UserHomeDir()
	rootCmd.Flags().BoolVarP(&createdTime, "createdTime", "t", false, "with pod created time")
	rootCmd.Flags().StringVarP(&kubeConfig, "kubeconfig", "k", path.Join(home, "./.kube", "config"), "kubernetes config path")
	rootCmd.Flags().BoolVarP(&reason, "reason", "r", false, "with pod created time")
	rootCmd.Flags().BoolVarP(&successPod, "success-pod", "s", false, "with success pod")
	rootCmd.Flags().BoolVarP(&watch, "watch", "w", false, "watch gpu pod ")
	rootCmd.Execute()
}
