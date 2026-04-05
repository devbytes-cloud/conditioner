package cmd

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/devbytes-cloud/conditioner/pkg/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

func TestSetNodeNames(t *testing.T) {
	o := NewConditionOptions(genericiooptions.IOStreams{})

	err := o.setNodeNames([]string{"node/worker-01", "worker-02", "nodes/worker-03"})
	require.NoError(t, err)
	assert.Equal(t, []string{"worker-01", "worker-02", "worker-03"}, o.nodeNames)
}

func TestSetNodeNamesRequiresAtLeastOneNode(t *testing.T) {
	o := NewConditionOptions(genericiooptions.IOStreams{})

	err := o.setNodeNames(nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must provide at least one node")
}

func TestSetNodeNamesRejectsEmptyNodeName(t *testing.T) {
	o := NewConditionOptions(genericiooptions.IOStreams{})

	err := o.setNodeNames([]string{"node/"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "node name cannot be empty")
}

func TestReadStdinNames_NonTTY_ReturnsLines(t *testing.T) {
	streams, in, _, _ := genericiooptions.NewTestIOStreams()
	fmt.Fprintln(in, "worker-01")
	fmt.Fprintln(in, "worker-02")

	o := NewConditionOptions(streams)
	names, err := o.readStdinNames()
	require.NoError(t, err)
	assert.Equal(t, []string{"worker-01", "worker-02"}, names)
}

func TestReadStdinNames_SkipsBlankLines(t *testing.T) {
	streams, in, _, _ := genericiooptions.NewTestIOStreams()
	fmt.Fprintln(in, "worker-01")
	fmt.Fprintln(in, "")
	fmt.Fprintln(in, "worker-02")

	o := NewConditionOptions(streams)
	names, err := o.readStdinNames()
	require.NoError(t, err)
	assert.Equal(t, []string{"worker-01", "worker-02"}, names)
}

func TestReadStdinNames_EmptyInput_ErrorsFromSetNodeNames(t *testing.T) {
	streams, _, _, _ := genericiooptions.NewTestIOStreams()

	o := NewConditionOptions(streams)
	names, err := o.readStdinNames()
	require.NoError(t, err)
	assert.Empty(t, names)

	err = o.setNodeNames(names)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must provide at least one node")
}

func TestReadStdinNames_PrefixNormalization(t *testing.T) {
	streams, in, _, _ := genericiooptions.NewTestIOStreams()
	fmt.Fprintln(in, "node/worker-01")
	fmt.Fprintln(in, "nodes/worker-02")

	o := NewConditionOptions(streams)
	names, err := o.readStdinNames()
	require.NoError(t, err)
	// readStdinNames returns raw names; normalization happens inside setNodeNames
	assert.Equal(t, []string{"node/worker-01", "nodes/worker-02"}, names)

	err = o.setNodeNames(names)
	require.NoError(t, err)
	assert.Equal(t, []string{"worker-01", "worker-02"}, o.nodeNames)
}

func TestReadStdinNames_MergeOrder(t *testing.T) {
	streams, in, _, _ := genericiooptions.NewTestIOStreams()
	fmt.Fprintln(in, "stdin-node")

	o := NewConditionOptions(streams)
	stdinNames, err := o.readStdinNames()
	require.NoError(t, err)

	merged := append(stdinNames, "positional-node")
	err = o.setNodeNames(merged)
	require.NoError(t, err)
	assert.Equal(t, []string{"stdin-node", "positional-node"}, o.nodeNames)
}

type errReader struct{}

func (e errReader) Read(_ []byte) (int, error) {
	return 0, fmt.Errorf("read error")
}

func TestReadStdinNames_ScannerError(t *testing.T) {
	streams := genericiooptions.IOStreams{
		In:     io.MultiReader(strings.NewReader("worker-01\n"), errReader{}),
		Out:    io.Discard,
		ErrOut: io.Discard,
	}

	o := NewConditionOptions(streams)
	_, err := o.readStdinNames()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "reading stdin:")
}
