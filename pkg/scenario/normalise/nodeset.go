package normalise

import (
	"fmt"
	"strings"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario/public"
)

// This string can be used in the configuration to select all nodes there are
const allNodesString = "all"

// desugarNodeSet takes a list of node names and verifies that all these nodes exist
// and that no duplicates exist in the nodeset. (to make it actually a set)
func desugarNodeSet(nodeset []string, nodegroups []*public.NodeGroup) ([]string, error) {
	if len(nodeset) == 1 && nodeset[0] == allNodesString {
		result := make([]string, 0, len(nodegroups))

		for _, group := range nodegroups {
			result = append(result, group.GroupName)
		}

		return result, nil
	}

	groupnameset := make(map[string]bool)
	for _, group := range nodegroups {
		groupnameset[group.GroupName] = true
	}

	had := make(map[string]bool)
	newnodeset := make([]string, 0, len(nodeset))

	for _, name := range nodeset {
		name = strings.TrimSpace(name)

		if had[name] {
			return nil, fmt.Errorf("duplicate node group name %s in task", name)
		}

		if !groupnameset[name] {
			if name == allNodesString {
				return nil, fmt.Errorf("can only use '%s' as the only nodegroup", allNodesString)
			}

			return nil, fmt.Errorf("group name %s doesn't exist", name)
		}

		newnodeset = append(newnodeset, name)
		had[name] = true
	}

	return nodeset, nil
}
