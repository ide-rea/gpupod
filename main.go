package main

import (
	"flag"
	"k8s.io/klog"
	"os"
	"path"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gpupod",
	Short: "gpupod is a tool to list/watch NVIDIA GPU pod",
	Run: func(cmd *cobra.Command, args []string) {
		listWatchGpuPod()
	},
}

var createdTime bool
var kubeConfig string
var reason bool
var successPod bool
var watch bool
var timeout int

func main() {

	klogFlagset := flag.NewFlagSet("log", flag.ExitOnError)
	klog.InitFlags(klogFlagset)
	rootCmd.Flags().AddGoFlagSet(klogFlagset)

	rootCmd.Flags().BoolVarP(&createdTime, "createdTime", "t", false, "show pod created time(default without created time)")
	home, _ := os.UserHomeDir()
	rootCmd.Flags().StringVarP(&kubeConfig, "kubeconfig", "k", path.Join(home, "./.kube", "config"), "kubernetes config path")
	rootCmd.Flags().BoolVarP(&reason, "reason", "r", false, "show status reason(default without status reason)")
	rootCmd.Flags().BoolVarP(&successPod, "success-pod", "s", false, "list pod include success pod(default not list success pod)")
	rootCmd.Flags().BoolVarP(&watch, "watch", "w", false, "watch gpu pod(default only list pod)")
	rootCmd.Flags().IntVar(&timeout, "request-timeout", 0, "list watch request timeout seconds(default zero)")
	rootCmd.Execute()
}
