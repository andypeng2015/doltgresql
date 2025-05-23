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

package binary

import (
	"strconv"

	"github.com/cockroachdb/errors"
	"github.com/dolthub/go-mysql-server/sql"

	"github.com/dolthub/doltgresql/server/functions/framework"
	pgtypes "github.com/dolthub/doltgresql/server/types"
)

// These functions can be gathered using the following query from a Postgres 15 instance:
// SELECT * FROM pg_operator o WHERE o.oprname = <OPERATOR> ORDER BY o.oprcode::varchar;
// Replace <OPERATOR> with the desired JSON operator

// initJSON registers the functions to the catalog.
func initJSON() {
	framework.RegisterBinaryFunction(framework.Operator_BinaryJSONExtractJson, json_array_element)
	framework.RegisterBinaryFunction(framework.Operator_BinaryJSONExtractJson, jsonb_array_element)
	framework.RegisterBinaryFunction(framework.Operator_BinaryJSONExtractJson, json_object_field)
	framework.RegisterBinaryFunction(framework.Operator_BinaryJSONExtractJson, jsonb_object_field)
	framework.RegisterBinaryFunction(framework.Operator_BinaryJSONExtractText, json_array_element_text)
	framework.RegisterBinaryFunction(framework.Operator_BinaryJSONExtractText, jsonb_array_element_text)
	framework.RegisterBinaryFunction(framework.Operator_BinaryJSONExtractText, json_object_field_text)
	framework.RegisterBinaryFunction(framework.Operator_BinaryJSONExtractText, jsonb_object_field_text)
	framework.RegisterBinaryFunction(framework.Operator_BinaryJSONExtractPathJson, json_extract_path)
	framework.RegisterBinaryFunction(framework.Operator_BinaryJSONExtractPathJson, jsonb_extract_path)
	framework.RegisterBinaryFunction(framework.Operator_BinaryJSONExtractPathText, json_extract_path_text)
	framework.RegisterBinaryFunction(framework.Operator_BinaryJSONExtractPathText, jsonb_extract_path_text)
	framework.RegisterBinaryFunction(framework.Operator_BinaryJSONContainsRight, jsonb_contains)
	framework.RegisterBinaryFunction(framework.Operator_BinaryJSONContainsLeft, jsonb_contained)
	framework.RegisterBinaryFunction(framework.Operator_BinaryJSONTopLevel, jsonb_exists)
	framework.RegisterBinaryFunction(framework.Operator_BinaryJSONTopLevelAny, jsonb_exists_any)
	framework.RegisterBinaryFunction(framework.Operator_BinaryJSONTopLevelAll, jsonb_exists_all)
	framework.RegisterBinaryFunction(framework.Operator_BinaryMinus, jsonb_delete_text)
	framework.RegisterBinaryFunction(framework.Operator_BinaryMinus, jsonb_delete_text_array)
	framework.RegisterBinaryFunction(framework.Operator_BinaryMinus, jsonb_delete_int32)
}

// json_array_element represents the PostgreSQL function of the same name, taking the same parameters.
var json_array_element = framework.Function2{
	Name:       "json_array_element",
	Return:     pgtypes.Json,
	Parameters: [2]*pgtypes.DoltgresType{pgtypes.Json, pgtypes.Int32},
	Strict:     true,
	Callable: func(ctx *sql.Context, _ [3]*pgtypes.DoltgresType, val1 any, val2 any) (any, error) {
		// TODO: make a bespoke implementation that preserves whitespace
		newVal, err := pgtypes.JsonB.IoInput(ctx, val1.(string))
		if err != nil {
			return nil, err
		}
		var unusedTypes [3]*pgtypes.DoltgresType
		retVal, err := jsonb_array_element.Callable(ctx, unusedTypes, newVal, val2)
		if err != nil {
			return nil, err
		}
		if retVal == nil {
			return "", nil
		}
		return pgtypes.JsonB.IoOutput(ctx, retVal)
	},
}

// jsonb_array_element represents the PostgreSQL function of the same name, taking the same parameters.
var jsonb_array_element = framework.Function2{
	Name:       "jsonb_array_element",
	Return:     pgtypes.JsonB,
	Parameters: [2]*pgtypes.DoltgresType{pgtypes.JsonB, pgtypes.Int32},
	Strict:     true,
	Callable: func(ctx *sql.Context, _ [3]*pgtypes.DoltgresType, val1 any, val2 any) (any, error) {
		array, ok := val1.(pgtypes.JsonDocument).Value.(pgtypes.JsonValueArray)
		if !ok {
			return nil, nil
		}
		idx := val2.(int32)
		if idx < 0 {
			idx += int32(len(array))
		}
		if int(idx) >= len(array) {
			return nil, nil
		}
		return pgtypes.JsonDocument{Value: array[idx]}, nil
	},
}

// json_object_field represents the PostgreSQL function of the same name, taking the same parameters.
var json_object_field = framework.Function2{
	Name:       "json_object_field",
	Return:     pgtypes.Json,
	Parameters: [2]*pgtypes.DoltgresType{pgtypes.Json, pgtypes.Text},
	Strict:     true,
	Callable: func(ctx *sql.Context, _ [3]*pgtypes.DoltgresType, val1 any, val2 any) (any, error) {
		// TODO: make a bespoke implementation that preserves whitespace
		newVal, err := pgtypes.JsonB.IoInput(ctx, val1.(string))
		if err != nil {
			return nil, err
		}
		var unusedTypes [3]*pgtypes.DoltgresType
		retVal, err := jsonb_object_field.Callable(ctx, unusedTypes, newVal, val2)
		if err != nil {
			return nil, err
		}
		if retVal == nil {
			return "", nil
		}
		return pgtypes.JsonB.IoOutput(ctx, retVal)
	},
}

// jsonb_object_field represents the PostgreSQL function of the same name, taking the same parameters.
var jsonb_object_field = framework.Function2{
	Name:       "jsonb_object_field",
	Return:     pgtypes.JsonB,
	Parameters: [2]*pgtypes.DoltgresType{pgtypes.JsonB, pgtypes.Text},
	Strict:     true,
	Callable: func(ctx *sql.Context, _ [3]*pgtypes.DoltgresType, val1 any, val2 any) (any, error) {
		object, ok := val1.(pgtypes.JsonDocument).Value.(pgtypes.JsonValueObject)
		if !ok {
			return nil, nil
		}
		idx, ok := object.Index[val2.(string)]
		if !ok {
			return nil, nil
		}
		return pgtypes.JsonDocument{Value: object.Items[idx].Value}, nil
	},
}

// json_array_element_text represents the PostgreSQL function of the same name, taking the same parameters.
var json_array_element_text = framework.Function2{
	Name:       "json_array_element_text",
	Return:     pgtypes.Text,
	Parameters: [2]*pgtypes.DoltgresType{pgtypes.Json, pgtypes.Int32},
	Strict:     true,
	Callable: func(ctx *sql.Context, _ [3]*pgtypes.DoltgresType, val1 any, val2 any) (any, error) {
		// TODO: make a bespoke implementation that preserves whitespace
		newVal, err := pgtypes.JsonB.IoInput(ctx, val1.(string))
		if err != nil {
			return nil, err
		}
		var unusedTypes [3]*pgtypes.DoltgresType
		return jsonb_array_element_text.Callable(ctx, unusedTypes, newVal, val2)
	},
}

// jsonb_array_element_text represents the PostgreSQL function of the same name, taking the same parameters.
var jsonb_array_element_text = framework.Function2{
	Name:       "jsonb_array_element_text",
	Return:     pgtypes.Text,
	Parameters: [2]*pgtypes.DoltgresType{pgtypes.JsonB, pgtypes.Int32},
	Strict:     true,
	Callable: func(ctx *sql.Context, dt [3]*pgtypes.DoltgresType, val1 any, val2 any) (any, error) {
		doc, err := jsonb_array_element.Callable(ctx, dt, val1, val2)
		if err != nil || doc == nil {
			return nil, err
		}
		switch value := doc.(pgtypes.JsonDocument).Value.(type) {
		case pgtypes.JsonValueString:
			return string(value), nil
		default:
			return pgtypes.JsonB.IoOutput(ctx, pgtypes.JsonDocument{Value: value})
		}
	},
}

// json_object_field_text represents the PostgreSQL function of the same name, taking the same parameters.
var json_object_field_text = framework.Function2{
	Name:       "json_object_field_text",
	Return:     pgtypes.Text,
	Parameters: [2]*pgtypes.DoltgresType{pgtypes.Json, pgtypes.Text},
	Strict:     true,
	Callable: func(ctx *sql.Context, _ [3]*pgtypes.DoltgresType, val1 any, val2 any) (any, error) {
		// TODO: make a bespoke implementation that preserves whitespace
		newVal, err := pgtypes.JsonB.IoInput(ctx, val1.(string))
		if err != nil {
			return nil, err
		}
		var unusedTypes [3]*pgtypes.DoltgresType
		return jsonb_object_field_text.Callable(ctx, unusedTypes, newVal, val2)
	},
}

// jsonb_object_field_tex_jsonb represents the PostgreSQL function of the same name, taking the same parameters.
var jsonb_object_field_text = framework.Function2{
	Name:       "jsonb_object_field_text",
	Return:     pgtypes.Text,
	Parameters: [2]*pgtypes.DoltgresType{pgtypes.JsonB, pgtypes.Text},
	Strict:     true,
	Callable: func(ctx *sql.Context, dt [3]*pgtypes.DoltgresType, val1 any, val2 any) (any, error) {
		doc, err := jsonb_object_field.Callable(ctx, dt, val1, val2)
		if err != nil || doc == nil {
			return nil, err
		}
		switch value := doc.(pgtypes.JsonDocument).Value.(type) {
		case pgtypes.JsonValueString:
			return string(value), nil
		default:
			return pgtypes.JsonB.IoOutput(ctx, pgtypes.JsonDocument{Value: value})
		}
	},
}

// json_extract_path represents the PostgreSQL function of the same name, taking the same parameters.
var json_extract_path = framework.Function2{
	Name:       "json_extract_path",
	Return:     pgtypes.Json,
	Parameters: [2]*pgtypes.DoltgresType{pgtypes.Json, pgtypes.TextArray},
	Strict:     true,
	Callable: func(ctx *sql.Context, _ [3]*pgtypes.DoltgresType, val1 any, val2 any) (any, error) {
		// TODO: make a bespoke implementation that preserves whitespace
		newVal, err := pgtypes.JsonB.IoInput(ctx, val1.(string))
		if err != nil {
			return nil, err
		}
		var unusedTypes [3]*pgtypes.DoltgresType
		retVal, err := jsonb_extract_path.Callable(ctx, unusedTypes, newVal, val2)
		if err != nil {
			return nil, err
		}
		if retVal == nil {
			return "", nil
		}
		return pgtypes.JsonB.IoOutput(ctx, retVal)
	},
}

// jsonb_extract_path represents the PostgreSQL function of the same name, taking the same parameters.
var jsonb_extract_path = framework.Function2{
	Name:       "jsonb_extract_path",
	Return:     pgtypes.JsonB,
	Parameters: [2]*pgtypes.DoltgresType{pgtypes.JsonB, pgtypes.TextArray},
	Strict:     true,
	Callable: func(ctx *sql.Context, _ [3]*pgtypes.DoltgresType, val1 any, val2 any) (any, error) {
		value := val1.(pgtypes.JsonDocument).Value
		paths := val2.([]interface{})
		for _, path := range paths {
			textPath, ok := path.(string)
			if !ok {
				return nil, nil
			}
			switch currentValue := value.(type) {
			case pgtypes.JsonValueObject:
				idx, ok := currentValue.Index[textPath]
				if !ok {
					return nil, nil
				}
				value = currentValue.Items[idx].Value
			case pgtypes.JsonValueArray:
				idx, err := strconv.Atoi(textPath)
				if err != nil {
					// We don't return the error here, a bad parse is treated as an object key which isn't valid
					return nil, nil
				}
				value = currentValue[idx]
			default:
				return nil, nil
			}
		}
		return pgtypes.JsonDocument{Value: value}, nil
	},
}

// json_extract_path_text represents the PostgreSQL function of the same name, taking the same parameters.
var json_extract_path_text = framework.Function2{
	Name:       "json_extract_path_text",
	Return:     pgtypes.Text,
	Parameters: [2]*pgtypes.DoltgresType{pgtypes.Json, pgtypes.TextArray},
	Strict:     true,
	Callable: func(ctx *sql.Context, _ [3]*pgtypes.DoltgresType, val1 any, val2 any) (any, error) {
		// TODO: make a bespoke implementation that preserves whitespace
		newVal, err := pgtypes.JsonB.IoInput(ctx, val1.(string))
		if err != nil {
			return nil, err
		}
		var unusedTypes [3]*pgtypes.DoltgresType
		return jsonb_extract_path_text.Callable(ctx, unusedTypes, newVal, val2)
	},
}

// jsonb_extract_path_text represents the PostgreSQL function of the same name, taking the same parameters.
var jsonb_extract_path_text = framework.Function2{
	Name:       "jsonb_extract_path_text",
	Return:     pgtypes.Text,
	Parameters: [2]*pgtypes.DoltgresType{pgtypes.JsonB, pgtypes.TextArray},
	Strict:     true,
	Callable: func(ctx *sql.Context, dt [3]*pgtypes.DoltgresType, val1 any, val2 any) (any, error) {
		doc, err := jsonb_extract_path.Callable(ctx, dt, val1, val2)
		if err != nil || doc == nil {
			return nil, err
		}
		switch value := doc.(pgtypes.JsonDocument).Value.(type) {
		case pgtypes.JsonValueString:
			return string(value), nil
		default:
			return pgtypes.JsonB.IoOutput(ctx, pgtypes.JsonDocument{Value: value})
		}
	},
}

// jsonb_contains represents the PostgreSQL function of the same name, taking the same parameters.
var jsonb_contains = framework.Function2{
	Name:       "jsonb_contains",
	Return:     pgtypes.Bool,
	Parameters: [2]*pgtypes.DoltgresType{pgtypes.JsonB, pgtypes.JsonB},
	Strict:     true,
	Callable: func(ctx *sql.Context, _ [3]*pgtypes.DoltgresType, val1 any, val2 any) (any, error) {
		return nil, errors.Errorf("JSON contains is not yet supported")
	},
}

// jsonb_contained represents the PostgreSQL function of the same name, taking the same parameters.
var jsonb_contained = framework.Function2{
	Name:       "jsonb_contained",
	Return:     pgtypes.Bool,
	Parameters: [2]*pgtypes.DoltgresType{pgtypes.JsonB, pgtypes.JsonB},
	Strict:     true,
	Callable: func(ctx *sql.Context, dt [3]*pgtypes.DoltgresType, val1 any, val2 any) (any, error) {
		return jsonb_contains.Callable(ctx, dt, val2, val1)
	},
}

// jsonb_exists represents the PostgreSQL function of the same name, taking the same parameters.
var jsonb_exists = framework.Function2{
	Name:       "jsonb_exists",
	Return:     pgtypes.Bool,
	Parameters: [2]*pgtypes.DoltgresType{pgtypes.JsonB, pgtypes.Text},
	Strict:     true,
	Callable: func(ctx *sql.Context, _ [3]*pgtypes.DoltgresType, val1 any, val2 any) (any, error) {
		switch value := val1.(pgtypes.JsonDocument).Value.(type) {
		case pgtypes.JsonValueObject:
			_, ok := value.Index[val2.(string)]
			return ok, nil
		case pgtypes.JsonValueArray:
			str := val2.(string)
			for _, arrayItem := range value {
				itemStr, ok := arrayItem.(pgtypes.JsonValueString)
				if ok && str == string(itemStr) {
					return true, nil
				}
			}
			return false, nil
		case pgtypes.JsonValueString:
			return string(value) == val2.(string), nil
		default:
			return false, nil
		}
	},
}

// jsonb_exists_any represents the PostgreSQL function of the same name, taking the same parameters.
var jsonb_exists_any = framework.Function2{
	Name:       "jsonb_exists_any",
	Return:     pgtypes.Bool,
	Parameters: [2]*pgtypes.DoltgresType{pgtypes.JsonB, pgtypes.TextArray},
	Strict:     true,
	Callable: func(ctx *sql.Context, _ [3]*pgtypes.DoltgresType, val1 any, val2 any) (any, error) {
		keys := val2.([]interface{})
		switch value := val1.(pgtypes.JsonDocument).Value.(type) {
		case pgtypes.JsonValueObject:
			for _, key := range keys {
				if _, ok := value.Index[key.(string)]; ok {
					return true, nil
				}
			}
			return false, nil
		case pgtypes.JsonValueArray:
			// Inefficient but good enough for now
			for _, key := range keys {
				for _, arrayItem := range value {
					itemStr, ok := arrayItem.(pgtypes.JsonValueString)
					if ok && string(itemStr) == key.(string) {
						return true, nil
					}
				}
			}
			return false, nil
		case pgtypes.JsonValueString:
			for _, key := range keys {
				if string(value) == key.(string) {
					return true, nil
				}
			}
			return false, nil
		default:
			return false, nil
		}
	},
}

// jsonb_exists_all represents the PostgreSQL function of the same name, taking the same parameters.
var jsonb_exists_all = framework.Function2{
	Name:       "jsonb_exists_all",
	Return:     pgtypes.Bool,
	Parameters: [2]*pgtypes.DoltgresType{pgtypes.JsonB, pgtypes.TextArray},
	Strict:     true,
	Callable: func(ctx *sql.Context, _ [3]*pgtypes.DoltgresType, val1 any, val2 any) (any, error) {
		keys := val2.([]interface{})
		switch value := val1.(pgtypes.JsonDocument).Value.(type) {
		case pgtypes.JsonValueObject:
			for _, key := range keys {
				if _, ok := value.Index[key.(string)]; !ok {
					return false, nil
				}
			}
			return true, nil
		case pgtypes.JsonValueArray:
			// Inefficient but good enough for now
			for _, key := range keys {
				found := false
				for _, arrayItem := range value {
					itemStr, ok := arrayItem.(pgtypes.JsonValueString)
					if ok && string(itemStr) == key.(string) {
						found = true
						break
					}
				}
				if !found {
					return false, nil
				}
			}
			return true, nil
		case pgtypes.JsonValueString:
			for _, key := range keys {
				if string(value) != key.(string) {
					return false, nil
				}
			}
			return true, nil
		default:
			return false, nil
		}
	},
}

// jsonb_delete_text represents the PostgreSQL function of the same name, taking the same parameters.
var jsonb_delete_text = framework.Function2{
	Name:       "jsonb_delete",
	Return:     pgtypes.JsonB,
	Parameters: [2]*pgtypes.DoltgresType{pgtypes.JsonB, pgtypes.Text},
	Strict:     true,
	Callable: func(ctx *sql.Context, _ [3]*pgtypes.DoltgresType, val1 any, val2 any) (any, error) {
		return nil, errors.Errorf("JSON deletions are not yet supported")
	},
}

// jsonb_delete_text_array represents the PostgreSQL function of the same name, taking the same parameters.
var jsonb_delete_text_array = framework.Function2{
	Name:       "jsonb_delete",
	Return:     pgtypes.JsonB,
	Parameters: [2]*pgtypes.DoltgresType{pgtypes.JsonB, pgtypes.TextArray},
	Strict:     true,
	Callable: func(ctx *sql.Context, _ [3]*pgtypes.DoltgresType, val1 any, val2 any) (any, error) {
		return nil, errors.Errorf("JSON deletions are not yet supported")
	},
}

// jsonb_delete_int32 represents the PostgreSQL function of the same name, taking the same parameters.
var jsonb_delete_int32 = framework.Function2{
	Name:       "jsonb_delete",
	Return:     pgtypes.JsonB,
	Parameters: [2]*pgtypes.DoltgresType{pgtypes.JsonB, pgtypes.Int32},
	Strict:     true,
	Callable: func(ctx *sql.Context, _ [3]*pgtypes.DoltgresType, val1 any, val2 any) (any, error) {
		return nil, errors.Errorf("JSON deletions are not yet supported")
	},
}
