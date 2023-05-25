//go:build ignore
// +build ignore

// Copyright (C) 2019-2021 Zilliz. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License
// is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
// or implied. See the License for the specific language governing permissions and limitations under the License.

// This program generates entity/columns_{{FieldType}}.go. Invoked by go generate
package main

import (
	"fmt"
	"os"
	"text/template"
	"time"

	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

var scalarColumnTemplate = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT
// This file is generated by go generate

package entity 

import (
	"fmt"

	"github.com/cockroachdb/errors"
	schema "github.com/milvus-io/milvus-proto/go-api/schemapb"
)
{{ range .Types }}{{with .}}
// Column{{.TypeName}} generated columns type for {{.TypeName}}
type Column{{.TypeName}} struct {
	ColumnBase
	name   string
	values []{{.TypeDef}}
}

// Name returns column name
func (c *Column{{.TypeName}}) Name() string {
	return c.name
}

// Type returns column FieldType
func (c *Column{{.TypeName}}) Type() FieldType {
	return FieldType{{.TypeName}}
}

// Len returns column values length
func (c *Column{{.TypeName}}) Len() int {
	return len(c.values)
}

// Get returns value at index as interface{}.
func (c *Column{{.TypeName}}) Get(idx int) (interface{}, error) {
	var r {{.TypeDef}} // use default value
	if idx < 0 || idx >= c.Len() {
		return r, errors.New("index out of range")
	}
	return c.values[idx], nil
}

// FieldData return column data mapped to schema.FieldData
func (c *Column{{.TypeName}}) FieldData() *schema.FieldData {
	fd := &schema.FieldData{
		Type: schema.DataType_{{.TypeName}},
		FieldName: c.name,
	}
	data := make([]{{.PbType}}, 0, c.Len())
	for i := 0 ;i < c.Len(); i++ {
		data = append(data, {{.PbType}}(c.values[i]))
	}
	fd.Field = &schema.FieldData_Scalars{
		Scalars: &schema.ScalarField{
			Data: &schema.ScalarField_{{.PbName}}Data{
				{{.PbName}}Data: &schema.{{.PbName}}Array{
					Data: data,
				},
			},
		},
	}
	return fd
}

// ValueByIdx returns value of the provided index
// error occurs when index out of range
func (c *Column{{.TypeName}}) ValueByIdx(idx int) ({{.TypeDef}}, error) {
	var r {{.TypeDef}} // use default value
	if idx < 0 || idx >= c.Len() {
		return r, errors.New("index out of range")
	}
	return c.values[idx], nil
}

// AppendValue append value into column
func(c *Column{{.TypeName}}) AppendValue(i interface{}) error {
	v, ok := i.({{.TypeDef}})
	if !ok {
		return fmt.Errorf("invalid type, expected {{.TypeDef}}, got %T", i)
	}
	c.values = append(c.values, v)

	return nil
}

// Data returns column data
func (c *Column{{.TypeName}}) Data() []{{.TypeDef}} {
	return c.values
}

// NewColumn{{.TypeName}} auto generated constructor
func NewColumn{{.TypeName}}(name string, values []{{.TypeDef}}) *Column{{.TypeName}} {
	return &Column{{.TypeName}} {
		name: name,
		values: values,
	}
}
{{end}}{{end}}
`))
var vectorColumnTemplate = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT
// This file is generated by go generate 

package entity 

import (
	"fmt"

	"github.com/cockroachdb/errors"
	schema "github.com/milvus-io/milvus-proto/go-api/schemapb"
)

{{ range .Types }}{{with.}}
// Column{{.TypeName}} generated columns type for {{.TypeName}}
type Column{{.TypeName}} struct {
	ColumnBase
	name   string
	dim    int
	values []{{.TypeDef}}
}

// Name returns column name
func (c *Column{{.TypeName}}) Name() string {
	return c.name
}

// Type returns column FieldType
func (c *Column{{.TypeName}}) Type() FieldType {
	return FieldType{{.TypeName}}
}

// Len returns column data length
func (c * Column{{.TypeName}}) Len() int {
	return len(c.values)
}

// Dim returns vector dimension
func (c *Column{{.TypeName}}) Dim() int {
	return c.dim
}

// Get returns values at index as interface{}.
func (c *Column{{.TypeName}}) Get(idx int) (interface{}, error) {
	if idx < 0 || idx >= c.Len() {
		return nil, errors.New("index out of range")
	}
	return c.values[idx], nil
}

// AppendValue append value into column
func(c *Column{{.TypeName}}) AppendValue(i interface{}) error {
	v, ok := i.({{.TypeDef}})
	if !ok {
		return fmt.Errorf("invalid type, expected {{.TypeDef}}, got %T", i)
	}
	c.values = append(c.values, v)

	return nil
}

// Data returns column data
func (c *Column{{.TypeName}}) Data() []{{.TypeDef}} {
	return c.values
}

// FieldData return column data mapped to schema.FieldData
func (c *Column{{.TypeName}}) FieldData() *schema.FieldData {
	fd := &schema.FieldData{
		Type: schema.DataType_{{.TypeName}},
		FieldName: c.name,
	}

	data := make({{.TypeDef}}, 0, len(c.values)* c.dim)

	for _, vector := range c.values {
		data = append(data, vector...)
	}

	fd.Field = &schema.FieldData_Vectors{
		Vectors: &schema.VectorField{
			Dim: int64(c.dim),
			{{if eq .TypeName "BinaryVector" }}
			Data: &schema.VectorField_BinaryVector{
				BinaryVector: data,
			},
			{{else}}
			Data: &schema.VectorField_FloatVector{
				FloatVector: &schema.FloatArray{
					Data: data,
				},
			},
			{{end}}
		},
	}
	return fd
}

// NewColumn{{.TypeName}} auto generated constructor
func NewColumn{{.TypeName}}(name string, dim int, values []{{.TypeDef}}) *Column{{.TypeName}} {
	return &Column{{.TypeName}} {
		name:   name,
		dim:    dim,
		values: values,
	}
}
{{end}}{{end}}
`))

var scalarColumnTestTemplate = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT
// This file is generated by go generate 

package entity 

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	schema "github.com/milvus-io/milvus-proto/go-api/schemapb"
)
{{ range .Types }}{{with.}}
func TestColumn{{.TypeName}}(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	columnName := fmt.Sprintf("column_{{.TypeName}}_%d", rand.Int())
	columnLen := 8 + rand.Intn(10)

	v := make([]{{.TypeDef}}, columnLen)
	column := NewColumn{{.TypeName}}(columnName, v)

	t.Run("test meta", func(t *testing.T) {
		ft := FieldType{{.TypeName}}
		assert.Equal(t, "{{.TypeName}}", ft.Name())
		assert.Equal(t, "{{.TypeDef}}", ft.String())
		pbName, pbType := ft.PbFieldType()
		assert.Equal(t, "{{.PbName}}", pbName)
		assert.Equal(t, "{{.PbType}}", pbType)
	})

	t.Run("test column attribute", func(t *testing.T) {
		assert.Equal(t, columnName, column.Name())
		assert.Equal(t, FieldType{{.TypeName}}, column.Type())
		assert.Equal(t, columnLen, column.Len())
		assert.EqualValues(t, v, column.Data())
	})

	t.Run("test column field data", func(t *testing.T) {
		fd := column.FieldData()
		assert.NotNil(t, fd)
		assert.Equal(t, fd.GetFieldName(), columnName)
	})

	t.Run("test column value by idx", func(t *testing.T) {
		_, err := column.ValueByIdx(-1)
		assert.NotNil(t, err)
		_, err = column.ValueByIdx(columnLen)
		assert.NotNil(t, err)
		for i := 0; i < columnLen; i++ {
			v, err := column.ValueByIdx(i)
			assert.Nil(t, err)
			assert.Equal(t, column.values[i], v)
		}
	})
}

func TestFieldData{{.TypeName}}Column(t *testing.T) {
	len := rand.Intn(10) + 8
	name := fmt.Sprintf("fd_{{.TypeName}}_%d", rand.Int())
	fd := &schema.FieldData{
		Type: schema.DataType_{{.TypeName}},
		FieldName: name,
	}

	t.Run("normal usage", func(t *testing.T) {
		fd.Field = &schema.FieldData_Scalars{
			Scalars: &schema.ScalarField{
				Data: &schema.ScalarField_{{.PbName}}Data{
					{{.PbName}}Data: &schema.{{.PbName}}Array{
						Data: make([]{{.PbType}}, len),
					},
				},
			},
		}
		column, err := FieldDataColumn(fd, 0, len)
		assert.Nil(t, err)
		assert.NotNil(t, column)
 
		assert.Equal(t, name, column.Name())
		assert.Equal(t, len, column.Len())
		assert.Equal(t, FieldType{{.TypeName}}, column.Type())

		var ev {{.TypeDef}}
		err = column.AppendValue(ev)
		assert.Equal(t, len+1, column.Len())
		assert.Nil(t, err)

		err = column.AppendValue(struct{}{})
		assert.Equal(t, len+1, column.Len())
		assert.NotNil(t, err)
	})

	
	t.Run("nil data", func(t *testing.T) {
		fd.Field = nil
		_, err := FieldDataColumn(fd, 0, len)
		assert.NotNil(t, err)
	})
	
	t.Run("get all data", func(t *testing.T) {
		fd.Field = &schema.FieldData_Scalars{
			Scalars: &schema.ScalarField{
				Data: &schema.ScalarField_{{.PbName}}Data{
					{{.PbName}}Data: &schema.{{.PbName}}Array{
						Data: make([]{{.PbType}}, len),
					},
				},
			},
		}
		column, err := FieldDataColumn(fd, 0, -1)
		assert.Nil(t, err)
		assert.NotNil(t, column)
		
		assert.Equal(t, name, column.Name())
		assert.Equal(t, len, column.Len())
		assert.Equal(t, FieldType{{.TypeName}}, column.Type())
	})
}
{{end}}{{end}}
`))
var vectorColumnTestTemplate = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT
// This file is generated by go generate 

package entity 

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	schema "github.com/milvus-io/milvus-proto/go-api/schemapb"
	"github.com/stretchr/testify/assert"
)
{{range .Types}}{{with .}}
func TestColumn{{.TypeName}}(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	columnName := fmt.Sprintf("column_{{.TypeName}}_%d", rand.Int())
	columnLen := 12 + rand.Intn(10)
	dim := ([]int{64, 128, 256, 512})[rand.Intn(4)]

	v := make([]{{.TypeDef}},0, columnLen)
	dlen := dim
	{{if eq .TypeName "BinaryVector" }}dlen /= 8{{end}}
	
	for i := 0; i < columnLen; i++ {
		entry := make({{.TypeDef}}, dlen)
		v = append(v, entry)
	}
	column := NewColumn{{.TypeName}}(columnName, dim, v)
	
	t.Run("test meta", func(t *testing.T) {
		ft := FieldType{{.TypeName}}
		assert.Equal(t, "{{.TypeName}}", ft.Name())
		assert.Equal(t, "{{.TypeDef}}", ft.String())
		pbName, pbType := ft.PbFieldType()
		assert.Equal(t, "{{.PbName}}", pbName)
		assert.Equal(t, "{{.PbType}}", pbType)
	})

	t.Run("test column attribute", func(t *testing.T) {
		assert.Equal(t, columnName, column.Name())
		assert.Equal(t, FieldType{{.TypeName}}, column.Type())
		assert.Equal(t, columnLen, column.Len())
		assert.Equal(t, dim, column.Dim())
		assert.Equal(t ,v, column.Data())
		
		var ev {{.TypeDef}}
		err := column.AppendValue(ev)
		assert.Equal(t, columnLen+1, column.Len())
		assert.Nil(t, err)
		
		err = column.AppendValue(struct{}{})
		assert.Equal(t, columnLen+1, column.Len())
		assert.NotNil(t, err)
	})

	t.Run("test column field data", func(t *testing.T) {
		fd := column.FieldData()
		assert.NotNil(t, fd)
		assert.Equal(t, fd.GetFieldName(), columnName)

		c, err := FieldDataVector(fd)
		assert.NotNil(t, c)
		assert.NoError(t, err)
	})

	t.Run("test column field data error", func(t *testing.T) {
		fd := &schema.FieldData{
			Type:      schema.DataType_{{.TypeName}},
			FieldName: columnName,
		}
		_, err := FieldDataVector(fd) 
		assert.Error(t, err)
	})

}
{{end}}{{end}}
`))

func main() {
	scalarFieldTypes := []entity.FieldType{
		entity.FieldTypeBool,
		entity.FieldTypeInt8,
		entity.FieldTypeInt16,
		entity.FieldTypeInt32,
		entity.FieldTypeInt64,
		entity.FieldTypeFloat,
		entity.FieldTypeDouble,
		entity.FieldTypeString,
	}
	vectorFieldTypes := []entity.FieldType{
		entity.FieldTypeBinaryVector,
		entity.FieldTypeFloatVector,
	}
	now := time.Now()

	pf := func(ft entity.FieldType) interface{} {
		pbName, pbType := ft.PbFieldType()
		return struct {
			Timestamp time.Time
			TypeName  string
			TypeDef   string
			PbName    string
			PbType    string
		}{
			Timestamp: now,
			TypeName:  ft.Name(),
			TypeDef:   ft.String(),
			PbName:    pbName,
			PbType:    pbType,
		}
	}
	fn := func(fn string, types []entity.FieldType, tmpl *template.Template, pf func(entity.FieldType) interface{}) {
		params := struct {
			Timestamp time.Time
			Types     []interface{}
		}{
			Timestamp: now,
			Types:     make([]interface{}, 0, len(types)),
		}
		for _, ft := range types {
			params.Types = append(params.Types, pf(ft))
		}
		f, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer f.Close()

		tmpl.Execute(f, params)
	}
	fnTest := func(fn string, types []entity.FieldType, tmpl *template.Template, pf func(entity.FieldType) interface{}) {
		params := struct {
			Timestamp time.Time
			Types     []interface{}
		}{
			Timestamp: now,
			Types:     make([]interface{}, 0, len(types)),
		}
		for _, ft := range types {
			params.Types = append(params.Types, pf(ft))
		}
		f, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer f.Close()
		tmpl.Execute(f, params)
	}
	fn("columns_scalar_gen.go", scalarFieldTypes, scalarColumnTemplate, pf)
	fn("columns_vector_gen.go", vectorFieldTypes, vectorColumnTemplate, pf)
	fnTest("columns_scalar_gen_test.go", scalarFieldTypes, scalarColumnTestTemplate, pf)
	fnTest("columns_vector_gen_test.go", vectorFieldTypes, vectorColumnTestTemplate, pf)
}
