package label

func splitIdInGroups(id string, maxGroupSize int) []string {
	groups := make([]string, 0, len(id)/maxGroupSize+1)
	for i := 0; i < len(id); {
		endBound := min(i+maxGroupSize, len(id))
		groups = append(groups, id[i:endBound])
		i = endBound
	}
	return groups
}
