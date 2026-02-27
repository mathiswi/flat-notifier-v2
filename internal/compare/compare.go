package compare

import "github.com/mathiswitte/flat-notifier-go/internal/model"

// FindNewFlats returns the flats from found that are not in existingIDs.
func FindNewFlats(found []model.FlatInfo, existingIDs []string) []model.FlatInfo {
	existing := make(map[string]struct{}, len(existingIDs))
	for _, id := range existingIDs {
		existing[id] = struct{}{}
	}

	var newFlats []model.FlatInfo
	for _, flat := range found {
		if _, ok := existing[flat.FlatID]; !ok {
			newFlats = append(newFlats, flat)
		}
	}
	return newFlats
}
