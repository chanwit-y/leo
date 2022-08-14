package utils

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

type Factory struct {
	msSql   MsSql
	schemas []Schema
}

func NewFactory(msSql MsSql) Factory {
	return Factory{msSql, []Schema{}}
}

func (f Factory) CreateGorm(tableName string) []string {
	var gormText []string

	gormText = append(gormText, fmt.Sprintf("package %s\n", packageName))
	gormText = append(gormText, fmt.Sprintf("type %s struct {\n", toCamelCase(tableName)))

	f.toSchema(tableName)

	maxLenName := f.maxLenName()
	maxLenDataType := f.maxLenDataType()

	for _, s := range f.schemas {
		pk := ternary(s.IsPk, "primaryKey", "")
		nameSpace := space(maxLenName - s.LenName)
		typeSpace := space(maxLenDataType - s.LenDataType)
		gormText = append(gormText, fmt.Sprintf("	%s%s%s%s `gorm:\"column:%s;type:%s;%s\"`\n",
			s.Name,
			nameSpace,
			s.DataType,
			typeSpace,
			s.ColumnName,
			s.SQLDataType,
			pk))
	}
	fmt.Println(gormText)

	return gormText
}

func (f Factory) maxLenName() int {
	return lo.MaxBy(f.schemas, func(s1, s2 Schema) bool {
		return s1.LenName > s2.LenName
	}).LenName
}

func (f Factory) maxLenDataType() int {
	return lo.MaxBy(f.schemas, func(s1, s2 Schema) bool {
		return s1.LenDataType > s2.LenDataType
	}).LenDataType
}

func (f *Factory) toSchema(tableName string) {
	var pks []string
	f.msSql.Query(QueryPrimarykey(tableName)).Scan(&pks)

	var columns []Column
	f.msSql.Query(QueryColums(tableName)).Scan(&columns)
	lo.ForEach(columns, func(t Column, i int) {
		schema := Schema{
			Name:        toCamelCase(t.ColumnName),
			LenName:     len(toCamelCase(t.ColumnName)),
			IsPk:        isPk(pks, t.ColumnName),
			ColumnName:  t.ColumnName,
			DataType:    toGoType(t.DataType),
			SQLDataType: t.DataType,
			LenDataType: len(toGoType(t.DataType)),
		}
		f.schemas = append(f.schemas, schema)
		// schemas = append(schemas, fmt.Sprintf("	%s %s `gorm:\"column:%s;type:%s;%s\"`\n", columName, dataType, t.ColumnName, t.DataType, pk))
	})
}

// func maxlen(perv, next int) int {
// 	if perv < next {
// 		return next
// 	}
// 	return perv
// }

func isPk(pks []string, colName string) bool {
	_, f := lo.Find(pks, func(t string) bool {
		return t == colName
	})

	return f
}

func toGoType(dataType string) string {

	switch strings.ToLower(dataType) {
	case "nvarchar":
		return "string"
	case "bigint":
		return "int64"
	case "int":
		return "int64"
	case "bit":
		return "bool"
	case "datetime":
		return "time.Time"
	case "time":
		return "time.Time"
	case "decimal":
		return "float64"
	case "varchar":
		return "string"
	default:
		return ""
	}
}
