package common

var p2m = map[string]string{
	"int8":        "bigint",
	"numeric":     "bigint",
	"float8":      "double",
	"float4":      "float",
	"int2":        "smallint",
	"int4":        "integer",
	"timestamptz": "timestamp",
}

type (
	// Model gets table information from information_schema
	SchemaModel interface {
		GetAllTables() ([]string, error)
		GetColumns(table string) (*Table, error)
		GetIndex(table, column string) ([]*Index, error)
	}

	Column struct {
		Name       string `json:"name" gorm:"column:COLUMN_NAME"`
		DataType   string `json:"dataType" gorm:"column:DATA_TYPE"`
		Extra      string `gorm:"column:EXTRA"`
		Comment    string `json:"comment" gorm:"column:COLUMN_COMMENT"`
		IsNullAble string `gorm:"column:IS_NULLABLE"`
	}

	// DbIndex defines index of columns in information_schema.statistic
	Index struct {
		IndexName  string `db:"INDEX_NAME"`
		NonUnique  int    `db:"NON_UNIQUE"`
		SeqInIndex int    `db:"SEQ_IN_INDEX"`
	}

	// Table describes mysql table which contains database name, table name, columns, keys
	Table struct {
		Db      string
		Table   string
		Columns []Column
		// Primary key not included
		UniqueIndex map[string][]*Column
		PrimaryKey  *Column
		NormalIndex map[string][]*Column
	}

	// 结构体字段描述
	Field struct {
		Name     string
		DataType string
		Comment  string
	}
)
