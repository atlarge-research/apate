package normalise

import (
	"fmt"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario/public"
)

func desugarNodeSet(nodeset []string, nodegroups []*public.NodeGroup) ([]string, error) {
	if len(nodeset) == 1 && nodeset[0] == "all" {
		result := make([]string, 0, len(nodegroups))

		for _, group := range nodegroups {
			result = append(result, group.Groupname)
		}

		return result, nil
	}

	had := make(map[string]bool)

	for _, name := range nodeset {
		if had[name] {
			return nodeset, fmt.Errorf("duplicate node group name %s in task", name)
		}

		had[name] = true
	}

	return nodeset, nil
}
