package utils

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

type Factory struct {
	msSql   MsSql
	schemas []Schema
	dbIndex []DbIndex
}

func NewFactory(msSql MsSql) Factory {
	var indexs []DbIndex
	msSql.Query(QueryIndexs()).Scan(&indexs)
	// fmt.Println(indexs)
	return Factory{msSql, []Schema{}, indexs}
}

func (f Factory) TestToGrom(tabaleName string) {
	text := f.toGorm(tabaleName)
	fmt.Println(text)
}

func (f Factory) GenGormFile() {
	var tables []string
	f.msSql.Query(QueryTables()).Scan(&tables)
	lo.ForEach(tables, func(t string, i int) {
		gormText := f.toGorm(t)
		createFile(fmt.Sprintf("./project/schema/%s.go", strings.ToLower(t)), gormText)
	})
}

func (f Factory) toGorm(tableName string) []string {
	var gormText []string

	f.toSchema(tableName)

	maxLenName := f.maxLenName()
	maxLenDataType := f.maxLenDataType()

	gormText = append(gormText, fmt.Sprintf("package %s\n", packageName))
	_, haveTimeType := lo.Find(f.schemas, func(t Schema) bool {
		return t.DataType == "time.Time"
	})
	if haveTimeType {
		gormText = append(gormText, "import \"time\" \n")
	}

	gormText = append(gormText, fmt.Sprintf("type %s struct {\n", toCamelCase(tableName)))

	colPk := ""
	lo.ForEach(f.schemas, func(s Schema, i int) {
		pk := ternary(s.IsPk, "primaryKey", "")
		colPk = ternary(s.IsPk, s.ColumnName, colPk)
		nameSpace := space(maxLenName - s.LenName)
		typeSpace := space(maxLenDataType - s.LenDataType)
		if s.IsRelation {
			gormText = append(gormText, fmt.Sprintf("	%s%s%s%s `gorm:\"foreignKey:%s;references:%s\"`\n",
				s.Name,
				nameSpace,
				s.DataType,
				typeSpace,
				toCamelCase(colPk),
				toCamelCase(colPk)))
		} else {
			gormText = append(gormText, fmt.Sprintf("	%s%s%s%s `gorm:\"column:%s;type:%s;%s\"`\n",
				s.Name,
				nameSpace,
				s.DataType,
				typeSpace,
				s.ColumnName,
				s.SQLDataType,
				pk))
		}
	})

	gormText = append(gormText, "}\n")

	gormText = append(gormText, fmt.Sprintf("func (%s) TableName() string {\n", toCamelCase(tableName)))
	gormText = append(gormText, fmt.Sprintf("	return \"%s\"\n", tableName))
	gormText = append(gormText, "}\n")

	// fmt.Println(gormText)

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
	var colPkName string
	f.msSql.Query(QueryColums(tableName)).Scan(&columns)
	lo.ForEach(columns, func(t Column, i int) {
		colPkName = ternary(isPk(pks, t.ColumnName), t.ColumnName, colPkName)
		schema := Schema{
			Name:        toCamelCase(t.ColumnName),
			LenName:     len(toCamelCase(t.ColumnName)),
			IsPk:        isPk(pks, t.ColumnName),
			ColumnName:  t.ColumnName,
			DataType:    toGoType(t.DataType),
			SQLDataType: t.DataType,
			LenDataType: len(toGoType(t.DataType)),
			IsRelation:  false,
		}
		f.schemas = append(f.schemas, schema)
	})

	var foreignKeys []ForeignKey
	f.msSql.Query(QueryForeignKey(tableName)).Scan(&foreignKeys)
	lo.ForEach(foreignKeys, func(fk ForeignKey, i int) {
		index, _ := lo.Find(f.dbIndex, func(t DbIndex) bool {
			return t.TableName == fk.TableName && t.ColumnName == colPkName
		})

		fmt.Println(index)

		// check IsUniqueConstraint 0 is false and 1 is true
		dataType := fmt.Sprintf("%v%v", ternary(!index.IsUniqueConstraint, "[]", ""), toCamelCase(fk.TableName))
		// dataType := fmt.Sprintf("%v", toCamelCase(fk.TableName))

		schema := Schema{
			Name:        toCamelCase(fk.TableName),
			LenName:     len(toCamelCase(fk.TableName)),
			IsPk:        false,
			ColumnName:  fk.TableName,
			DataType:    dataType,
			SQLDataType: fk.TableName,
			LenDataType: len(dataType),
			IsRelation:  true,
		}
		f.schemas = append(f.schemas, schema)
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
