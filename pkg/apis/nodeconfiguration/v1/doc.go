//go:generate sh -c "cd ../../../../ && make crd_gen"

// Package v1 is the v1 version of the API.
// +groupName=apate.opendc.org
// +k8s:deepcopy-gen=package,register
// +k8s:openapi-gen=true
package v1
