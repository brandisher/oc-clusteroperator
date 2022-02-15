package main

import (
	"github.com/brandisher/oc-clusteroperator/pkg/cmd"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"os"
)

func main() {
	flags := pflag.NewFlagSet("oc-clusteroperator", pflag.ExitOnError)
	pflag.CommandLine = flags

	root := cmd.NewCmdClusterOperator(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}