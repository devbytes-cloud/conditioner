package cmd

import (
	"testing"

	"github.com/devbytes-cloud/conditioner/pkg/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericiooptions"
)

func TestComplete(t *testing.T) {
	streams := genericiooptions.IOStreams{}
	o := NewConditionOptions(streams)

	// Setup Cobra command with flags
	c := &cobra.Command{}
	c.Flags().String("type", "Ready", "")
	c.Flags().String("status", "true", "")
	c.Flags().String("reason", "KubeletReady", "")
	c.Flags().String("message", "kubelet is posting ready status", "")
	c.Flags().Bool("remove", false, "remove the condition") // Make sure to define the 'remove' flag

	err := o.Complete(c, []string{"test-node"}, &config.Config{})
	assert.NoError(t, err)

	// Assert the fields are set correctly
	assert.Equal(t, "Ready", string(o.condition.Type))
	assert.Equal(t, corev1.ConditionTrue, o.condition.Status)
	assert.Equal(t, "KubeletReady", o.condition.Reason)
	assert.Equal(t, "kubelet is posting ready status", o.condition.Message)
}
