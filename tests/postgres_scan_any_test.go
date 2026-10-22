// Copyright 2026 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"xorm.io/xorm/schemas"
)

// TestPostgresFindMapKeepsDriverTypes checks that map rows hold what the
// driver returns: nil for NULL, int64 for int4, float64 for float8 and
// time.Time for timestamptz.
func TestPostgresFindMapKeepsDriverTypes(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	if testEngine.Dialect().URI().DBType != schemas.POSTGRES {
		t.Skip("postgres only")
	}

	const tableName = "find_map_driver_types"
	dropTestTable(tableName)
	defer dropTestTable(tableName)

	_, err := testEngine.Exec(fmt.Sprintf(
		"CREATE TABLE %s (id INT4 NOT NULL PRIMARY KEY, i INT4 NULL, f FLOAT8 NULL, s TEXT NULL, ts TIMESTAMPTZ NULL)", qualifiedTableName(tableName)))
	assert.NoError(t, err)
	_, err = testEngine.Exec(fmt.Sprintf(
		"INSERT INTO %s (id, i, f, s, ts) VALUES (1, 2, 2.5, 'x', '2024-01-02 03:04:05.123+00'), (2, NULL, NULL, NULL, NULL)", qualifiedTableName(tableName)))
	assert.NoError(t, err)

	var rows []map[string]any
	assert.NoError(t, testEngine.Table(tableName).Asc("id").Find(&rows))
	if assert.Len(t, rows, 2) {
		assert.Equal(t, int64(2), rows[0]["i"])
		assert.Equal(t, float64(2.5), rows[0]["f"])
		assert.Equal(t, "x", rows[0]["s"])
		if ts, ok := rows[0]["ts"].(time.Time); assert.True(t, ok, "%T", rows[0]["ts"]) {
			assert.True(t, ts.Equal(time.Date(2024, 1, 2, 3, 4, 5, 123000000, time.UTC)), ts)
		}
		assert.Nil(t, rows[1]["i"])
		assert.Nil(t, rows[1]["f"])
		assert.Nil(t, rows[1]["s"])
		assert.Nil(t, rows[1]["ts"])
	}

	res, err := testEngine.QueryInterface(fmt.Sprintf("SELECT i, f FROM %s WHERE id = 2", qualifiedTableName(tableName)))
	assert.NoError(t, err)
	if assert.Len(t, res, 1) {
		assert.Nil(t, res[0]["i"])
		assert.Nil(t, res[0]["f"])
	}
}
