package utils

import "fmt"

func QueryTables() string {
	return `
		SELECT tbl.name AS TableName
		FROM sys.tables tbl
		WHERE tbl.is_ms_shipped = 0
		AND tbl.type = 'U'
		ORDER BY tbl.name;	
	`
}

func QueryColums(tableName string) string {
	return fmt.Sprintf(`
		SELECT c.name                                                   AS ColumnName,
			CASE typ.is_assembly_type
				WHEN 1 THEN TYPE_NAME(c.user_type_id)
				ELSE TYPE_NAME(c.system_type_id)
			END                                                             AS DataType,
			ISNULL(COLUMNPROPERTY(c.object_id, c.name, 'charmaxlen'), 0)    AS CharacterMaximumLength,
			ISNULL(OBJECT_DEFINITION(c.default_object_id), '')              AS ColumnDefault,
			c.is_nullable                                                   AS IsNullable,
			COLUMNPROPERTY(c.object_id, c.name, 'IsIdentity')               AS IsIdentity,
			OBJECT_NAME(c.object_id)                                        AS TableName,
			ISNULL(OBJECT_NAME(c.default_object_id), '')                    AS ConstraintName,
			ISNULL(convert(tinyint, CASE
			WHEN c.system_type_id IN (48, 52, 56, 59, 60, 62, 106, 108, 122, 127) THEN c.precision
			END), 0) AS NumericPrecision,
			ISNULL(convert(int, CASE
			WHEN c.system_type_id IN (40, 41, 42, 43, 58, 61) THEN NULL
			ELSE ODBCSCALE(c.system_type_id, c.scale) END), 0) AS NumericScale
		FROM sys.columns c
			INNER JOIN sys.tables t ON c.object_id = t.object_id
			INNER JOIN sys.types typ ON c.user_type_id = typ.user_type_id
		WHERE OBJECT_NAME(c.object_id) = '%v' 
		AND t.is_ms_shipped = 0
		ORDER BY TableName, COLUMNPROPERTY(c.object_id, c.name, 'ordinal');  
	`, tableName)
}

func QueryForeignKey(name string) string {
	return fmt.Sprintf(`
		SELECT OBJECT_NAME(fkc.constraint_object_id) AS ConstraintName,
			parent_table.name                       AS TableName,
			referenced_table.name                   AS ReferencedTableName,
			SCHEMA_NAME(referenced_table.schema_id) AS ReferencedSchemaName,
			parent_column.name                      AS ColumnName,
			referenced_column.name                  AS ReferencedColumnName,
			fk.delete_referential_action            AS DeleteReferentialAction,
			fk.update_referential_action            AS UpdateReferentialAction,
			fkc.constraint_column_id                AS OrdinalPosition
		FROM sys.foreign_key_columns AS fkc
			INNER JOIN sys.tables AS parent_table
					ON fkc.parent_object_id = parent_table.object_id
			INNER JOIN sys.tables AS referenced_table
					ON fkc.referenced_object_id = referenced_table.object_id
			INNER JOIN sys.columns AS parent_column
					ON fkc.parent_object_id = parent_column.object_id
					AND fkc.parent_column_id = parent_column.column_id
			INNER JOIN sys.columns AS referenced_column
					ON fkc.referenced_object_id = referenced_column.object_id
					AND fkc.referenced_column_id = referenced_column.column_id
			INNER JOIN sys.foreign_keys AS fk
					ON fkc.constraint_object_id = fk.object_id
					AND fkc.parent_object_id = fk.parent_object_id
		WHERE parent_table.is_ms_shipped = 0
		AND referenced_table.is_ms_shipped = 0
		AND referenced_table.name = '%v'
		--AND OBJECT_SCHEMA_NAME(fkc.parent_object_id) = @P1
		ORDER BY TableName, ConstraintName, OrdinalPosition
	`, name)
}

func QueryIndexs() string {
	return `
		SELECT DISTINCT
			ind.name AS IndexName,
			ind.is_unique AS IsUnique,
			ind.is_unique_constraint AS IsUniqueConstraint,
			ind.is_primary_key AS IsPrimaryKey,
			ind.type_desc as Clustering,
			col.name AS ColumnName,
			ic.key_ordinal AS SeqInIndex,
			ic.is_descending_key AS IsDescending,
			t.name AS TableName
		FROM
			sys.indexes ind
		INNER JOIN sys.index_columns ic
			ON ind.object_id = ic.object_id AND ind.index_id = ic.index_id
		INNER JOIN sys.columns col
			ON ic.object_id = col.object_id AND ic.column_id = col.column_id
		INNER JOIN
			sys.tables t ON ind.object_id = t.object_id
		WHERE t.is_ms_shipped = 0
			AND t.name = @P1 
			AND col.name = @P2 
			AND ind.filter_definition IS NULL
			AND ind.name IS NOT NULL
			AND ind.type_desc IN (
			'CLUSTERED',
			'NONCLUSTERED',
			'CLUSTERED COLUMNSTORE',
			'NONCLUSTERED COLUMNSTORE'
			)
		ORDER BY TableName, IndexName, SeqInIndex	
	`
}

func QueryPrimarykey(name string) string {
	return fmt.Sprintf(`
		SELECT column_name as PRIMARYKEYCOLUMN
		FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS AS TC 
		INNER JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE AS KU
		ON TC.CONSTRAINT_TYPE = 'PRIMARY KEY' 
		AND TC.CONSTRAINT_NAME = KU.CONSTRAINT_NAME 
		AND KU.table_name='%s'
		ORDER BY KU.TABLE_NAME, KU.ORDINAL_POSITION`, name)
}
