package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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
	%[1]s clusteroperator manage [NAME]

	# optionally scale the operator up or down
	%[1]s clusteroperator manage [NAME] --scale
	%[1]s clusteroperator unmanage [NAME] --scale
`
)

// ClusterOperatorOptions wraps the cli options
type ClusterOperatorOptions struct {
	configFlags *genericclioptions.ConfigFlags
	genericclioptions.IOStreams
	clusteroperator string
	state           string
	unmanaged       bool
	scale           bool
}

// NewClusterOperatorOptions provides a new instance of ClusterOperatorOptions with default values
func NewClusterOperatorOptions(streams genericclioptions.IOStreams) *ClusterOperatorOptions {
	return &ClusterOperatorOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
		IOStreams:   streams,
	}
}

// NewCmdClusterOperator provides a cobra command wrapping ClusterOperatorOptions
func NewCmdClusterOperator(streams genericclioptions.IOStreams) *cobra.Command {
	// Create the default option set
	o := NewClusterOperatorOptions(streams)

	// Create the command
	cmd := &cobra.Command{
		Use:          "clusteroperator [manage|unmanage] [NAME]",
		Short:        "Set a cluster operator to unmanaged or unmanaged",
		Example:      fmt.Sprintf(example, "oc"),
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.mergeArgs(c, args); err != nil {
				return err
			}
			if err := o.Execute(); err != nil {
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

// mergeArgs takes the user supplied arguments/flags and merges them into the ClusterOperatorOptions struct.
func (o *ClusterOperatorOptions) mergeArgs(c *cobra.Command, args []string) error {
	argCount := len(args)
	if argCount != 2 {
		return fmt.Errorf("expected exactly two arguments, got %d", argCount)
	}
	o.scale, _ = c.Flags().GetBool("scale")

	state := args[0]
	if err := o.setManagementState(state); err != nil {
		return err
	}

	co := args[1]
	if err := o.setClusterOperatorTarget(co); err != nil {
		return err
	}

	return nil
}

func (o *ClusterOperatorOptions) setManagementState(state string) error {
	if state != MANAGED_STATE && state != UNMANAGED_STATE {
		return fmt.Errorf("wanted %s or %s, got %s", MANAGED_STATE, UNMANAGED_STATE, o.state)
	}
	if state == MANAGED_STATE {
		o.unmanaged = false
	}
	if state == UNMANAGED_STATE {
		o.unmanaged = true
	}
	return nil
}

func (o *ClusterOperatorOptions) setClusterOperatorTarget(co string) error {
	o.clusteroperator = co
	return nil
}

func (o *ClusterOperatorOptions) Execute() error {
	cl := o.newRESTClient()

	body := buildRequestBody(o.clusteroperator, o.unmanaged)
	req := buildRequest(body, cl)

	// Execute the patch
	res := req.Do(context.TODO())

	if err := res.Error(); err != nil {
		return err
	}

	return nil
}

func (o *ClusterOperatorOptions) newRESTClient() rest.Interface {
	// Get the config
	config, _ := o.configFlags.ToRESTConfig()

	// Create a client
	clientSet, _ := kubernetes.NewForConfig(config)
	return clientSet.RESTClient()
}

// buildRequestBody takes the shortname of a cluster operator (e.g. "dns")
// and builds the request body.
func buildRequestBody(clusterOperatorName string, unmanaged bool) []byte {
	con := clusterOperatorName
	u := unmanaged

	reqBody := fmt.Sprintf(`[{"op":"add","path":"/spec/overrides","value":[{"kind":"Deployment","group":"apps/v1","name":"%s-operator","namespace":"openshift-%s-operator","unmanaged":%t}]}]`, con, con, u)

	return []byte(reqBody)
}

func buildRequest(requestBody []byte, restClient rest.Interface) *rest.Request {
	rb := requestBody
	cl := restClient

	return cl.Patch("application/json-patch+json").Body(rb).RequestURI("/apis/config.openshift.io/v1/clusterversions/version")
}
