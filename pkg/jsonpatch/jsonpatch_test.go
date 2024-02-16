package jsonpatch

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestOpType(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		index  int
		remove bool
		want   string
	}{
		{-1, false, "add"},
		{0, false, "replace"},
		{1, false, "replace"},
		{0, true, "remove"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			assert.Equal(tt.want, opType(tt.index, tt.remove))
		})
	}
}

func TestPathType(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		index int
		want  string
	}{
		{-1, basePath + "/-"},
		{0, basePath + "/0"},
		{1, basePath + "/1"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			assert.Equal(tt.want, pathType(tt.index))
		})
	}
}

func TestGenerateJsonPath(t *testing.T) {
	assert := assert.New(t)

	// Use a fixed time for testing
	now := metav1.Time{Time: time.Now()}

	tests := []struct {
		name          string
		index         int
		remove        bool
		oldConditions corev1.NodeCondition
		newConditions corev1.NodeCondition
		want          JsonPatch
	}{
		{
			name:          "Add condition",
			index:         -1,
			remove:        false,
			oldConditions: corev1.NodeCondition{},
			newConditions: corev1.NodeCondition{Type: "Ready", Status: "True", Reason: "NodeReady", Message: "Node is ready"},
			want: JsonPatch{
				OP:   "add",
				Path: basePath + "/-",
				Value: &corev1.NodeCondition{
					Type:               "Ready",
					Status:             "True",
					LastHeartbeatTime:  now,
					LastTransitionTime: now,
					Reason:             "NodeReady",
					Message:            "Node is ready",
				},
			},
		},
		{
			name:          "Update condition",
			index:         1,
			remove:        false,
			oldConditions: corev1.NodeCondition{Type: "Ready", Status: "False", Reason: "NodeNotReady", Message: "Node is not ready"},
			newConditions: corev1.NodeCondition{Type: "Ready", Status: "True", Reason: "NodeReady", Message: "Node is now ready"},
			want: JsonPatch{
				OP:   "replace",
				Path: basePath + "/1",
				Value: &corev1.NodeCondition{
					Type:               "Ready",
					Status:             "True",
					LastHeartbeatTime:  now,
					LastTransitionTime: now,
					Reason:             "NodeReady",
					Message:            "Node is now ready",
				},
			},
		},
		{
			name:          "Remove condition",
			index:         0,
			remove:        true,
			oldConditions: corev1.NodeCondition{Type: "DiskPressure", Status: "False", Reason: "DiskOK", Message: "Disk has sufficient space"},
			newConditions: corev1.NodeCondition{},
			want: JsonPatch{
				OP:    "remove",
				Path:  basePath + "/0",
				Value: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateJsonPath(tt.index, tt.remove, &tt.oldConditions, &tt.newConditions)
			assert.Equal(tt.want.OP, got.OP)
			assert.Equal(tt.want.Path, got.Path)
			if tt.want.Value != nil {
				assert.Equal(tt.want.Value.Type, got.Value.Type)
				assert.Equal(tt.want.Value.Status, got.Value.Status)
				assert.Equal(tt.want.Value.Reason, got.Value.Reason)
				assert.Equal(tt.want.Value.Message, got.Value.Message)
			} else {
				assert.Nil(got.Value)
			}
		})
	}
}
