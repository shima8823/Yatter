package dao

import (
	"strings"
)

func buildQuery(DBName, idColumnName string, since_id, max_id, limit *uint64) (string, []interface{}) {
	queryParts := []string{"select * from " + DBName}
	var args []interface{}

	if since_id != nil || max_id != nil {
		conditions := []string{}

		if since_id != nil {
			conditions = append(conditions, idColumnName+" >= ?")
			args = append(args, *since_id)
		}

		if max_id != nil {
			conditions = append(conditions, idColumnName+" <= ?")
			args = append(args, *max_id)
		}

		queryParts = append(queryParts, "where "+strings.Join(conditions, " AND "))
	}

	queryParts = append(queryParts, "order by create_at desc")

	if limit != nil {
		queryParts = append(queryParts, "limit ?")
		args = append(args, *limit)
	}

	query := strings.Join(queryParts, " ")
	return query, args
}
