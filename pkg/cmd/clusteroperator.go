package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	MANAGED_STATE   = "manage"
	UNMANAGED_STATE = "unmanage"
)

var (
	example = `
	# set a cluster operator to unmanaged
	%[1]s clusteroperator unmanage [NAME]

	# set a cluster operator to unmanaged
	%[1]s clusteroperator unmanage [NAME]

	# optionally scale the operator up or down
	%[1]s clusteroperator manage [NAME] --scale
	%[1]s clusteroperator unmanage [NAME] --scale
`
)

// ClusterOperatorOptions wraps the cli options
type ClusterOperatorOptions struct {
	configFlags *genericclioptions.ConfigFlags
	genericclioptions.IOStreams
	state string
	scale bool
}

// NewClusterOperatorOptions provides a new instance of ClusterOperatorOptions with default values
func NewClusterOperatorOptions(streams genericclioptions.IOStreams) *ClusterOperatorOptions {
	return &ClusterOperatorOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
		IOStreams:   streams,
		state:       "manage",
		scale:       false,
	}
}

// NewCmdClusterOperator provides a cobra command wrapping ClusterOperatorOptions
func NewCmdClusterOperator(streams genericclioptions.IOStreams) *cobra.Command {
	// Create the default option set
	o := NewClusterOperatorOptions(streams)

	// Create the command
	cmd := &cobra.Command{
		Use:     "clusteroperator [manage|unmanage] [NAME]",
		Short:   "Set a cluster operator to managed or unmanaged",
		Example: fmt.Sprintf(example, "oc"),
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Validate(); err != nil {
				return err
			}
			return nil
		},
	}

	// Add the scale flag
	cmd.Flags().BoolVar(&o.scale, "scale", o.scale, "if true, scale the operator deployment appropriately")
	o.configFlags.AddFlags(cmd.Flags())

	return cmd
}

func (o *ClusterOperatorOptions) Validate() error {
	if o.state != MANAGED_STATE && o.state != UNMANAGED_STATE {
		return fmt.Errorf("wanted %s or %s, got %s", MANAGED_STATE, UNMANAGED_STATE, o.state)
	}
	return nil
}
