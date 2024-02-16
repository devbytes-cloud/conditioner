package jsonpatch

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// basePath is a constant string that represents the base path in the JSON document where the operations are performed.
// It is used in the formation of the path for the JSON Patch operation.
const (
	basePath string = "/status/conditions"
)

// JsonPatch represents a JSON Patch operation.
// JSON Patch is a format (identified by the media type "application/json-patch+json")
// for expressing a sequence of operations to apply to a JavaScript Object Notation (JSON) document.
// It is suitable for use with HTTP PATCH requests.
type JsonPatch struct {
	// OP is the operation to be performed. It's a string and can be one of "add", "remove", or "replace".
	OP string `json:"op"`
	// Path is the string that contains the location in the JSON document where the operation is performed.
	Path string `json:"path"`
	// Value is the actual value that is used by the operation.
	// It's a pointer to a NodeCondition object from the "k8s.io/api/core/v1" package.
	Value *corev1.NodeCondition `json:"value"`
}

// GenerateJsonPath is a function that generates a JSON Patch operation based on the provided parameters.
// It takes four parameters:
// - index: an integer that represents the index of the condition in the conditions array.
// - remove: a boolean that indicates whether the operation is a remove operation.
// - oldConditions: a pointer to a NodeCondition object that represents the old conditions.
// - newConditions: a pointer to a NodeCondition object that represents the new conditions.
func GenerateJsonPath(index int, remove bool, oldConditions, newConditions *corev1.NodeCondition) JsonPatch {
	jsonPatch := JsonPatch{
		OP:   opType(index, remove),
		Path: pathType(index),
	}

	// If the operation is a remove operation, return the JsonPatch object.
	if remove {
		return jsonPatch
	}

	jsonPatch.Value = &corev1.NodeCondition{
		Type:               newConditions.Type,
		Status:             newConditions.Status,
		LastHeartbeatTime:  metav1.Time{Time: time.Now()},
		LastTransitionTime: metav1.Time{Time: time.Now()},
		Reason:             newConditions.Reason,
		Message:            newConditions.Message,
	}

	// If the index is not -1, update the value of the JsonPatch object based on the oldConditions parameter.
	if index != -1 {
		if newConditions.Message == "" {
			jsonPatch.Value.Message = oldConditions.Reason
		}

		if newConditions.Reason == "" {
			jsonPatch.Value.Reason = oldConditions.Reason
		}

		if newConditions.Status == "" {
			jsonPatch.Value.Status = oldConditions.Status
		}

	}

	// Return the JsonPatch object.
	return jsonPatch
}

// opType is a function that determines the operation type for a JSON Patch operation.
// The function returns a string that represents the operation type.
// If the remove parameter is true, the function returns "remove".
// If the index parameter is -1, the function returns "add".
func opType(index int, remove bool) string {
	if remove {
		return "remove"
	}
	if index == -1 {
		return "add"
	}

	return "replace"
}

// pathType is a function that generates the path for a JSON Patch operation.
// If the index parameter is -1, the function returns a string that ends with "/-".
// This is due to adding elements to the end of the array is denoted with "/-"
// Otherwise, the function returns a string that ends with the index.
func pathType(index int) string {
	if index == -1 {
		return fmt.Sprintf("%s/-", basePath)
	}

	return fmt.Sprintf("%s/%d", basePath, index)
}
