package main

import (
	"costco/internal/costco"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

// @title Costco API
// @version 1.0
// @description API documentation for Costco service
// @host localhost:8080
// @BasePath /v2
func main() {
	klog.InitFlags(nil)
	defer klog.Flush()

	var costcoCmd = &cobra.Command{
		Use:   "costco",
		Short: "Private container registry for Kubernetes",
		Long:  `Costco is an easily deployable private container registry for Kubernetes.`,
		Args:  cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			r := costco.Init(cmd)
			r.RegisterEndpoints()
			r.Start()
		},
	}

	costcoCmd.Flags().Bool("debug", false, "debug mode for costco")
	costcoCmd.Flags().Bool("verbose", false, "verbose mode for costco")
	costcoCmd.Flags().String("namespace", "costco", "namespace for costco")

	if err := costcoCmd.Execute(); err != nil {
		klog.Fatal(err)
	}
}
