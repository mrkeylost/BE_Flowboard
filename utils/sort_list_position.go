package utils

import (
	"github.com/google/uuid"
	"github.com/mrkeylost/Flowboard_Backend/model"
)

func SortListByPos(list []model.List, order []uuid.UUID) []model.List {
	ordered := make([]model.List, 0, len(order))

	listMap := make(map[uuid.UUID]model.List)
	for _, item := range list {
		listMap[item.PublicID] = item
	}

	for _, id := range order {
		if item, ok := listMap[id]; ok {
			ordered = append(ordered, item)
		}
	}

	return ordered
}
