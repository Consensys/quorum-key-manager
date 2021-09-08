package client

import (
	"context"
	"fmt"

	"github.com/consensys/quorum-key-manager/src/infra/postgres"
)

func QuerySearchIDs(ctx context.Context, client postgres.Client, table, idCol, whereCond string, whereArgs []interface{}, isDeleted bool, limit, offset uint64) ([]string, error) {
	var ids = []string{}
	var err error
	var query string

	switch {
	case limit != 0 || offset != 0:
		query = fmt.Sprintf("SELECT (array_agg(%s ORDER BY created_at ASC))[%d:%d] FROM %s WHERE %s", table, idCol, offset+1, offset+limit, whereCond)
	default:
		query = fmt.Sprintf("SELECT array_agg(%s ORDER BY created_at ASC) FROM secrets WHERE %s", idCol, whereCond)
	}

	if isDeleted {
		query = fmt.Sprintf("%s AND deleted_at is NOT NULL", query)
	} else {
		query = fmt.Sprintf("%s AND deleted_at is NULL", query)
	}

	err = client.Query(ctx, &ids, query, whereArgs...)
	if err != nil {
		return nil, err
	}

	return ids, nil
}
