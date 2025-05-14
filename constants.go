// Copyright 2020 by Blank-Xu. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sqladapter

const (
	// defaultTableName  if tableName == "", the Adapter will use this default table name.
	defaultTableName = "casbin_rule"

	// maxParameterCount .
	maxParameterCount = 13

	// defaultPlaceholder .
	defaultPlaceholder = "?"
)

type adapterDriverNameIndex int

const (
	_SQLite adapterDriverNameIndex = iota + 1
)

// general SQL for all supported databases.
const (
	sqlCreateTable = `
CREATE TABLE %[1]s(
    p_type VARCHAR(32),
    v0     VARCHAR(255),
    v1     VARCHAR(255),
    v2     VARCHAR(255),
    v3     VARCHAR(255),
    v4     VARCHAR(255),
    v5     VARCHAR(255)
    v6     VARCHAR(255),
    v7     VARCHAR(255),
    v8     VARCHAR(255),
    v9     VARCHAR(255),
    v10    VARCHAR(255),
    v11    VARCHAR(255),
);
CREATE INDEX idx_%[1]s ON %[1]s (p_type,v0,v1);`
	sqlTableExist   = "SELECT 1 FROM %s WHERE 1=0"
	sqlInsertRow    = "INSERT INTO %s (p_type,v0,v1,v2,v3,v4,v5,v6,v7,v8,v9,v10,v11) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)"
	sqlUpdateRow    = "UPDATE %s SET p_type=?,v0=?,v1=?,v2=?,v3=?,v4=?,v5=?,v6=?,v7=?,v8=?,v9=?,v10=?,v11=? WHERE p_type=? AND v0=? AND v1=? AND v2=? AND v3=? AND v4=? AND v5=? AND v6=? AND v7=? AND v8=? AND v9=? AND v10=? AND v11=?"
	sqlDeleteAll    = "DELETE FROM %s"
	sqlDeleteRow    = "DELETE FROM %s WHERE p_type=? AND v0=? AND v1=? AND v2=? AND v3=? AND v4=? AND v5=? AND v6=? AND v7=? AND v8=? AND v9=? AND v10=? AND v11=?"
	sqlDeleteByArgs = "DELETE FROM %s WHERE p_type=?"
	sqlSelectAll    = "SELECT p_type,v0,v1,v2,v3,v4,v5,v6,v7,v8,v9,v10,v11 FROM %s"
	sqlSelectWhere  = "SELECT p_type,v0,v1,v2,v3,v4,v5,v6,v7,v8,v9,v10,v11 FROM %s WHERE "
)

// for SQLite3.
// for SQLite3.
const (
	sqlCreateTableSQLite3 = `
CREATE TABLE IF NOT EXISTS %[1]s(
    p_type VARCHAR(32)  DEFAULT '' NOT NULL,
    v0     VARCHAR(255) DEFAULT '' NOT NULL,
    v1     VARCHAR(255) DEFAULT '' NOT NULL,
    v2     VARCHAR(255) DEFAULT '' NOT NULL,
    v3     VARCHAR(255) DEFAULT '' NOT NULL,
    v4     VARCHAR(255) DEFAULT '' NOT NULL,
    v5     VARCHAR(255) DEFAULT '' NOT NULL,
    v6     VARCHAR(255) DEFAULT '' NOT NULL,
    v7     VARCHAR(255) DEFAULT '' NOT NULL,
    v8     VARCHAR(255) DEFAULT '' NOT NULL,
    v9     VARCHAR(255) DEFAULT '' NOT NULL,
    v10    VARCHAR(255) DEFAULT '' NOT NULL,
    v11    VARCHAR(255) DEFAULT '' NOT NULL,
    CHECK (TYPEOF("p_type") = "text" AND LENGTH("p_type") <= 32),
    CHECK (TYPEOF("v0") = "text" AND LENGTH("v0") <= 255),
    CHECK (TYPEOF("v1") = "text" AND LENGTH("v1") <= 255),
    CHECK (TYPEOF("v2") = "text" AND LENGTH("v2") <= 255),
    CHECK (TYPEOF("v3") = "text" AND LENGTH("v3") <= 255),
    CHECK (TYPEOF("v4") = "text" AND LENGTH("v4") <= 255),
    CHECK (TYPEOF("v5") = "text" AND LENGTH("v5") <= 255),
    CHECK (TYPEOF("v6") = "text" AND LENGTH("v6") <= 255),
    CHECK (TYPEOF("v7") = "text" AND LENGTH("v7") <= 255),
    CHECK (TYPEOF("v8") = "text" AND LENGTH("v8") <= 255),
    CHECK (TYPEOF("v9") = "text" AND LENGTH("v9") <= 255),
    CHECK (TYPEOF("v10") = "text" AND LENGTH("v10") <= 255),
    CHECK (TYPEOF("v11") = "text" AND LENGTH("v11") <= 255)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_all_%[1]s ON %[1]s (p_type,v0,v1,v2,v3,v4,v5,v6,v7,v8,v9,v10,v11);`
	sqlTruncateTableSQLite3 = "DROP TABLE IF EXISTS %[1]s;" + sqlCreateTableSQLite3
)
