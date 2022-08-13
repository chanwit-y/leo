package utils

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

type Index struct {
	IndexName          string
	IsUnique           int
	IsUniqueConstraint int
	IsPrimaryKey       int
	Clustering         string
	ColumnName         string
	SeqInIndex         int
	IsDescending       int
	TableName          string
}
