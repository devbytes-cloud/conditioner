package cmd

import (
	"context"
	"fmt"
	"os/user"

	"github.com/devbytes-cloud/conditioner/pkg/config"
	"github.com/devbytes-cloud/conditioner/pkg/jsonpatch"
	"github.com/spf13/cobra"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"k8s.io/client-go/kubernetes"
)

var (
	example = `
# Add a new condition to a node
kubectl conditioner my-node --type Ready --status true --reason KubeletReady --message "kubelet is posting ready status"

# Update an existing condition on a node
kubectl conditioner my-node --type DiskPressure --status false --reason KubeletHasNoDiskPressure --message "kubelet has sufficient disk space available"

# Remove a condition from a node
kubectl conditioner my-node --type NetworkUnavailable --remove
`

	long = `The 'conditioner' command allows you to add, update, or remove status conditions on nodes. 
You need to provide the node name as an argument and use flags to specify the details of the condition. 
The '--type' flag is required and it specifies the type of condition you wish to interact with. 
The '--status' flag sets the status for the specific status condition and it can be 'true', 'false', or left blank for 'unknown'. 
The '--reason' flag sets the reason for the specific status condition. 
The '--message' flag sets the message for the specific status condition. 
If you wish to remove the condition from the node entirely, use the '--remove' flag.`
)

// ConditionOptions is a struct that holds the configuration for the condition command.
type ConditionOptions struct {
	// client is a pointer to a Clientset object that is used to interact with the Kubernetes API.
	client *kubernetes.Clientset

	// configFlags holds the configuration flags for the command.
	configFlags *genericclioptions.ConfigFlags

	// IOStreams provides the standard names for iostreams. This is useful for embedding and for unit testing.
	genericiooptions.IOStreams

	// nodeName is the name of the node that the command is being run against.
	nodeName string

	// remove is a boolean that indicates whether the condition should be removed.
	remove bool

	// condition is a pointer to a NodeCondition object that represents the condition to be added or updated.
	condition *corev1.NodeCondition

	// args is a slice of strings that contains the arguments that were passed to the command.
	args []string
}

// NewConditionOptions is a function that creates a new ConditionOptions.
func NewConditionOptions(streams genericiooptions.IOStreams) *ConditionOptions {
	return &ConditionOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
		IOStreams:   streams,
	}
}

func NewCmdCondition(streams genericiooptions.IOStreams) *cobra.Command {
	o := NewConditionOptions(streams)

	cmd := &cobra.Command{
		Use:          "conditioner [node name] [flags]",
		Short:        "Manipulate status conditions on a specified node.",
		Long:         long,
		Example:      example,
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			o.args = args
			if len(o.args) != 1 {
				return fmt.Errorf("must provide a node to be conditioned")
			}

			o.nodeName = o.args[0]
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			fs := config.FS{}
			conf, err := config.Read(fs)
			if err != nil {
				return err
			}

			if err := o.Complete(c, args, conf); err != nil {
				return err
			}

			if err := o.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringP("status", "", "", "Status for the specific status condition [true, false]")
	cmd.Flags().StringP("reason", "r", "", "Reason for the specific status condition")
	cmd.Flags().StringP("message", "", "", "Message for the specific status condition")
	cmd.Flags().StringP("type", "", "", "(required): type of condition you wish to interact with")
	cmd.Flags().BoolP("remove", "x", false, "If you wish to remove the condition from the node entirely")

	if err := cmd.MarkFlagRequired("type"); err != nil {
		panic(fmt.Sprintf("failed to mark %s flag required: %s", "message", err.Error()))
	}

	o.configFlags.AddFlags(cmd.Flags())

	return cmd
}

// Complete sets all information required for updating the current context
// It retrieves the restConfig from the configFlags and creates a new Kubernetes client.
// It also sets the condition status, reason, message, type, and remove flag from the command flags.
func (o *ConditionOptions) Complete(cmd *cobra.Command, _ []string, config *config.Config) error {
	// Get the restConfig from the configFlags
	restConfig, err := o.configFlags.ToRawKubeConfigLoader().ClientConfig()
	if err != nil {
		return err
	}

	// Create a new Kubernetes client
	o.client, err = kubernetes.NewForConfig(restConfig)
	if err != nil {
		return err
	}

	o.condition = &corev1.NodeCondition{}

	status, err := cmd.Flags().GetString("status")
	if err != nil {
		return err
	}

	if status == "true" {
		o.condition.Status = corev1.ConditionTrue
	} else if status == "false" {
		o.condition.Status = corev1.ConditionFalse
	} else {
		o.condition.Status = corev1.ConditionUnknown
	}

	o.condition.Reason, err = cmd.Flags().GetString("reason")
	if err != nil {
		return err
	}

	o.condition.Message, err = cmd.Flags().GetString("message")
	if err != nil {
		return err
	}

	if config.WhoAmI {
		u, err := user.Current()
		if err != nil {
			return err
		}

		o.condition.Message = fmt.Sprintf("%s: %s", u.Username, o.condition.Message)
	}

	// Get the type from the command flags and set the condition type
	conditionType, err := cmd.Flags().GetString("type")
	if err != nil {
		return err
	}

	if len(config.AllowList) != 0 {
		if ok := allowedType(conditionType, config.AllowList); !ok {
			return fmt.Errorf("condition %s is not in allow-list %v", conditionType, config.AllowList)
		}
	}

	o.condition.Type = corev1.NodeConditionType(conditionType)

	o.remove, err = cmd.Flags().GetBool("remove")
	if err != nil {
		return err
	}

	return nil
}

// Run handles the condition applying or removal on nodes.
func (o *ConditionOptions) Run() error {
	node, err := o.client.CoreV1().Nodes().Get(context.Background(), o.nodeName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	oldConditions, index := findConditionType(node.Status.Conditions, o.condition.Type)

	if index == -1 && o.remove {
		return fmt.Errorf("condition type of %s does not exist", o.condition.Type)
	}

	patch := jsonpatch.GenerateJsonPath(index, o.remove, oldConditions, o.condition)

	jsonPath := []interface{}{patch}
	bytePatch, err := json.Marshal(jsonPath)
	if err != nil {
		return err
	}

	if _, err := o.client.CoreV1().Nodes().Patch(context.Background(), node.Name, types.JSONPatchType, bytePatch, metav1.PatchOptions{}, "status"); err != nil {
		return err
	}

	fmt.Printf("condition status %s has been %sed on node %s\n", o.condition.Type, patch.OP, node.Name)

	return nil
}

// findConditionType is a function that searches for a specific condition type in a slice of NodeCondition objects.
// If a match is found, the function returns a pointer to the matching NodeCondition object and its index in the slice.
// If no match is found, the function returns nil and -1.
func findConditionType(conditions []corev1.NodeCondition, conditionType corev1.NodeConditionType) (*corev1.NodeCondition, int) {
	for k, v := range conditions {
		if v.Type == conditionType {
			return &v, k
		}
	}

	return nil, -1
}

// allowedType checks if a given condition type is in the list of allowed types.
func allowedType(conditionType string, allowedTypes []string) bool {
	for _, v := range allowedTypes {
		if v == conditionType {
			return true
		}
	}
	return false
}
