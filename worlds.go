package main

import "context"

type WorldDescription struct {
	Id          int64  `sql:"id"`
	Status      int8   `sql:"status"`
	Name        string `sql:"name"`
	Description string `sql:"description"`
}

func GetWorldDescriptions(ctx context.Context, connections *transport) ([]WorldDescription, error) {
	var worldDescriptions = make([]WorldDescription, 0)
	rows, err := connections.ydbClient.Query().QueryResultSet(ctx, "SELECT * FROM world_descriptions")
	if err != nil {
		return worldDescriptions, err
	}
	defer rows.Close(ctx)
	for row, err := range rows.Rows(ctx) {
		if err != nil {
			return worldDescriptions, err
		}
		var wd WorldDescription
		err := row.ScanStruct(&wd)
		if err != nil {
			return worldDescriptions, err
		}
		worldDescriptions = append(worldDescriptions, wd)
	}

	return worldDescriptions, nil
}
