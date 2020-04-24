package normalization

import (
	"fmt"
	"strings"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"
)

// This string can be used in the configuration to select all nodes there are
const allNodesString = "all"

// desugarNodeGroups takes a list of node names and verifies that all these nodes exist
// and that no duplicates exist in the nodeset. (to make it actually a set)
func desugarNodeGroups(nodeSet []string, nodeGroups []*controlplane.NodeGroup) ([]string, error) {
	onlyHasAllNodesString := len(nodeSet) == 1 && nodeSet[0] == allNodesString
	if onlyHasAllNodesString {
		return desugarAll(nodeGroups), nil
	}

	returnedNodes, err := removeDuplicates(nodeSet, nodeGroups)
	if err != nil {
		return nil, err
	}
	return returnedNodes, nil
}

func removeDuplicates(nodeSet []string, nodeGroups []*controlplane.NodeGroup) ([]string, error) {
	groupNameSet := make(map[string]bool)
	for _, group := range nodeGroups {
		groupNameSet[group.GroupName] = true
	}

	alreadySeen := make(map[string]bool)
	returnedNodes := make([]string, 0, len(nodeSet))

	for _, name := range nodeSet {
		name = strings.TrimSpace(name)

		// If the current groupName is not an existing groupName, error
		if !groupNameSet[name] {
			if name == allNodesString {
				return nil, fmt.Errorf("can only use '%s' as the only nodegroup", allNodesString)
			}

			return nil, fmt.Errorf("group name %s doesn't exist", name)
		}

		if !alreadySeen[name] {
			returnedNodes = append(returnedNodes, name)
			alreadySeen[name] = true
		}
	}
	return returnedNodes, nil
}

func desugarAll(nodeGroups []*controlplane.NodeGroup) []string {
	result := make([]string, 0, len(nodeGroups))

	for _, group := range nodeGroups {
		result = append(result, group.GroupName)
	}

	return result
}
