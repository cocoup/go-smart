package mysql

import (
	"github.com/go-sql-driver/mysql"
	gromSql "gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/cocoup/go-smart/tools/gocli/cmd/model/common"
)

type sqlModel struct {
	db     *gorm.DB
	dbName string
}

func NewSqlModel(dsn string) (common.SchemaModel, error) {
	sqlConf := gromSql.Config{DSN: DSN}
	db, err := gorm.Open(gromSql.New(sqlConf))
	if nil != err {
		return nil, err
	}

	dsnConf, err := mysql.ParseDSN(dsn)
	if nil != err {
		return nil, err
	}

	return &sqlModel{
		db:     db,
		dbName: dsnConf.DBName,
	}, nil
}

func (s *sqlModel) GetAllTables() ([]string, error) {
	query := `select table_name from information_schema.tables where table_schema = ?`

	var tables []string
	err := s.db.Raw(query, s.dbName).Scan(&tables).Error
	if err != nil {
		return nil, err
	}

	return tables, nil
}

func (s *sqlModel) GetColumns(table string) (*common.Table, error) {
	query := `SELECT 
				COLUMN_NAME,
				DATA_TYPE,EXTRA, 
				COLUMN_COMMENT,
				COLUMN_DEFAULT,
				IS_NULLABLE,
				ORDINAL_POSITION 
				FROM INFORMATION_SCHEMA.COLUMNS 
				WHERE TABLE_SCHEMA = ? and TABLE_NAME = ? 
 				ORDER BY ORDINAL_POSITION
				`
	var columns []common.Column
	err := s.db.Raw(query, s.dbName, table).Scan(&columns).Error
	if err != nil {
		return nil, err
	}

	return &common.Table{
		Db:      s.dbName,
		Table:   table,
		Columns: columns,
	}, nil
}

// TODO:: 查询索引，添加缓存相关逻辑
func (s *sqlModel) GetIndex(table, column string) ([]*common.Index, error) {
	//querySql := `SELECT s.INDEX_NAME,s.NON_UNIQUE,s.SEQ_IN_INDEX from  STATISTICS s  WHERE  s.TABLE_SCHEMA = ? and s.TABLE_NAME = ? and s.COLUMN_NAME = ?`
	//var reply []*DbIndex
	//err := s.conn.QueryRowsPartial(&reply, querySql, db, table, column)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return reply, nil

	return nil, nil
}
