package normalise

import (
	"fmt"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario/public"
)

var allNodesString = "all"

func desugarNodeSet(nodeset []string, nodegroups []*public.NodeGroup) ([]string, error) {
	if len(nodeset) == 1 && nodeset[0] == allNodesString {
		result := make([]string, 0, len(nodegroups))

		for _, group := range nodegroups {
			result = append(result, group.Groupname)
		}

		return result, nil
	}

	groupnameset := make(map[string]bool)
	for _, group := range nodegroups {
		groupnameset[group.Groupname] = true
	}

	had := make(map[string]bool)

	for _, name := range nodeset {
		if had[name] {
			return nil, fmt.Errorf("duplicate node group name %s in task", name)
		}

		if !groupnameset[name] {
			if name == allNodesString {
				return nil, fmt.Errorf("can only use '%s' as the only nodegroup", allNodesString)
			}

			return nil, fmt.Errorf("group name %s doesn't exist", name)
		}

		had[name] = true
	}

	return nodeset, nil
}
