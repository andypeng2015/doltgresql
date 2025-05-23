// Copyright 2024 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pgcatalog

import (
	"io"

	"github.com/dolthub/go-mysql-server/sql"

	"github.com/dolthub/doltgresql/server/tables"
	pgtypes "github.com/dolthub/doltgresql/server/types"
)

// PgLargeobjectName is a constant to the pg_largeobject name.
const PgLargeobjectName = "pg_largeobject"

// InitPgLargeobject handles registration of the pg_largeobject handler.
func InitPgLargeobject() {
	tables.AddHandler(PgCatalogName, PgLargeobjectName, PgLargeobjectHandler{})
}

// PgLargeobjectHandler is the handler for the pg_largeobject table.
type PgLargeobjectHandler struct{}

var _ tables.Handler = PgLargeobjectHandler{}

// Name implements the interface tables.Handler.
func (p PgLargeobjectHandler) Name() string {
	return PgLargeobjectName
}

// RowIter implements the interface tables.Handler.
func (p PgLargeobjectHandler) RowIter(ctx *sql.Context) (sql.RowIter, error) {
	// TODO: Implement pg_largeobject row iter
	return emptyRowIter()
}

// Schema implements the interface tables.Handler.
func (p PgLargeobjectHandler) Schema() sql.PrimaryKeySchema {
	return sql.PrimaryKeySchema{
		Schema:     pgLargeobjectSchema,
		PkOrdinals: nil,
	}
}

// pgLargeobjectSchema is the schema for pg_largeobject.
var pgLargeobjectSchema = sql.Schema{
	{Name: "loid", Type: pgtypes.Oid, Default: nil, Nullable: false, Source: PgLargeobjectName},
	{Name: "pageno", Type: pgtypes.Int32, Default: nil, Nullable: false, Source: PgLargeobjectName},
	{Name: "data", Type: pgtypes.Bytea, Default: nil, Nullable: false, Source: PgLargeobjectName},
}

// pgLargeobjectRowIter is the sql.RowIter for the pg_largeobject table.
type pgLargeobjectRowIter struct {
}

var _ sql.RowIter = (*pgLargeobjectRowIter)(nil)

// Next implements the interface sql.RowIter.
func (iter *pgLargeobjectRowIter) Next(ctx *sql.Context) (sql.Row, error) {
	return nil, io.EOF
}

// Close implements the interface sql.RowIter.
func (iter *pgLargeobjectRowIter) Close(ctx *sql.Context) error {
	return nil
}
