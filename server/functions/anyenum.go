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

package functions

import (
	"github.com/dolthub/go-mysql-server/sql"

	"github.com/dolthub/doltgresql/server/functions/framework"
	pgtypes "github.com/dolthub/doltgresql/server/types"
)

// initAnyEnum registers the functions to the catalog.
func initAnyEnum() {
	framework.RegisterFunction(anyenum_in)
	framework.RegisterFunction(anyenum_out)
}

// anyenum_in represents the PostgreSQL function of anyenum type IO input.
var anyenum_in = framework.Function1{
	Name:       "anyenum_in",
	Return:     pgtypes.AnyEnum,
	Parameters: [1]*pgtypes.DoltgresType{pgtypes.Cstring},
	Strict:     true,
	Callable: func(ctx *sql.Context, _ [2]*pgtypes.DoltgresType, val any) (any, error) {
		// TODO
		return val, nil
	},
}

// anyenum_out represents the PostgreSQL function of anyenum type IO output.
var anyenum_out = framework.Function1{
	Name:       "anyenum_out",
	Return:     pgtypes.Cstring,
	Parameters: [1]*pgtypes.DoltgresType{pgtypes.AnyEnum},
	Strict:     true,
	Callable: func(ctx *sql.Context, _ [2]*pgtypes.DoltgresType, val any) (any, error) {
		// TODO
		return val, nil
	},
}
