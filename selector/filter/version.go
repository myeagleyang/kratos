package filter

import (
	"context"

	"gitlab.wwgame.com/wwgame/kratos/v2/selector"
)

// Version is version filter.
func Version(version string) selector.NodeFilter {
	return func(_ context.Context, nodes []selector.Node) []selector.Node {
		newNodes := make([]selector.Node, 0)
		for _, n := range nodes {
			if n.Version() == version {
				newNodes = append(newNodes, n)
			}
		}
		return newNodes
	}
}
