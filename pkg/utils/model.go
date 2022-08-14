package utils

type Schema struct {
	Name        string
	LenName     int
	DataType    string
	IsPk        bool
	IsRelation  bool
	ColumnName  string
	SQLDataType string
	LenDataType int
}

type Column struct {
	ColumnName             string
	DataType               string
	CharacterMaximumLength int
	ColumnDefault          string
	IsNullable             int
	IsIdentity             int
	TableName              string
	ConstraintName         string
	NumericPrecision       int
	NumericScale           int
}

type ForeignKey struct {
	ConstraintName          string
	TableName               string
	ReferencedTableName     string
	ReferencedSchemaName    string
	ColumnName              string
	ReferencedColumnName    string
	DeleteReferentialAction int
	UpdateReferentialAction int
	OrdinalPosition         int
}

type DbIndex struct {
	IndexName          string
	IsUnique           bool
	IsUniqueConstraint bool
	IsPrimaryKey       bool
	Clustering         string
	ColumnName         string
	SeqInIndex         int
	IsDescending       bool
	TableName          string
}
