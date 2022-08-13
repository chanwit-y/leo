package utils

func QueryTables() string {
	return `
		SELECT tbl.name AS table_name
		FROM sys.tables tbl
		WHERE tbl.is_ms_shipped = 0
		AND tbl.type = 'U'
		ORDER BY tbl.name;	
	`
}

func QueryColums() string {
	return `
		SELECT c.name                                                   AS column_name,
			CASE typ.is_assembly_type
				WHEN 1 THEN TYPE_NAME(c.user_type_id)
				ELSE TYPE_NAME(c.system_type_id)
			END                                                             AS data_type,
			ISNULL(COLUMNPROPERTY(c.object_id, c.name, 'charmaxlen'), 0)    AS character_maximum_length,
			ISNULL(OBJECT_DEFINITION(c.default_object_id), '')              AS column_default,
			c.is_nullable                                                   AS is_nullable,
			COLUMNPROPERTY(c.object_id, c.name, 'IsIdentity')               AS is_identity,
			OBJECT_NAME(c.object_id)                                        AS table_name,
			ISNULL(OBJECT_NAME(c.default_object_id), '')                    AS constraint_name,
			ISNULL(convert(tinyint, CASE
			WHEN c.system_type_id IN (48, 52, 56, 59, 60, 62, 106, 108, 122, 127) THEN c.precision
			END), 0) AS numeric_precision,
			ISNULL(convert(int, CASE
			WHEN c.system_type_id IN (40, 41, 42, 43, 58, 61) THEN NULL
			ELSE ODBCSCALE(c.system_type_id, c.scale) END), 0) AS numeric_scale
		FROM sys.columns c
			INNER JOIN sys.tables t ON c.object_id = t.object_id
			INNER JOIN sys.types typ ON c.user_type_id = typ.user_type_id
		WHERE OBJECT_NAME(c.object_id) = @P1 
		AND t.is_ms_shipped = 0
		ORDER BY table_name, COLUMNPROPERTY(c.object_id, c.name, 'ordinal');  
	`
}

func QueryForeignKey() string {
	return `
		SELECT OBJECT_NAME(fkc.constraint_object_id) AS constraint_name,
			parent_table.name                       AS table_name,
			referenced_table.name                   AS referenced_table_name,
			SCHEMA_NAME(referenced_table.schema_id) AS referenced_schema_name,
			parent_column.name                      AS column_name,
			referenced_column.name                  AS referenced_column_name,
			fk.delete_referential_action            AS delete_referential_action,
			fk.update_referential_action            AS update_referential_action,
			fkc.constraint_column_id                AS ordinal_position
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
		AND OBJECT_SCHEMA_NAME(fkc.parent_object_id) = @P1
		ORDER BY table_name, constraint_name, ordinal_position
	`
}

func QueryIndexs() string {
	return `
		SELECT DISTINCT
			ind.name AS index_name,
			ind.is_unique AS is_unique,
			ind.is_unique_constraint AS is_unique_constraint,
			ind.is_primary_key AS is_primary_key,
			ind.type_desc as clustering,
			col.name AS column_name,
			ic.key_ordinal AS seq_in_index,
			ic.is_descending_key AS is_descending,
			t.name AS table_name
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
		ORDER BY table_name, index_name, seq_in_index	
	`
}
