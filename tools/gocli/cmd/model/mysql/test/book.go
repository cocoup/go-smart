package model

import (
	"github.com/cocoup/go-smart/core/stores/sqlx"
	"gorm.io/gorm"
)

type (
	// BookModel defines a model for user
	BookModel interface {
		Insert(data *Book) error
		FindOne(id int64) (*Book, error)
		FindOneByFilter(filter map[string]interface{}, opts ...sqlx.Option) (*Book, error)
		FindByFilter(filter map[string]interface{}, opts ...sqlx.Option) (*[]Book, error)
		Save(*Book) *gorm.DB
		Updates(*Book) *gorm.DB
		UpdateByFilter(filter map[string]interface{}, upVal map[string]interface{}) *gorm.DB
		Delete(id int64) error
		DeleteByFilter(filter map[string]interface{}) error
	}

	defaultBookModel struct {
		conn  sqlx.SqlConn
		table string
	}

	// User defines an data structure for mysql
	Book struct {
		Id   int64
		Name string
	}
)

// NewBookModel creates an instance for UserModel
func NewBookModel(conn sqlx.SqlConn) BookModel {
	return &defaultBookModel{
		conn:  conn,
		table: "`book`",
	}
}

func (d *defaultBookModel) Insert(data *Book) error {
	return d.conn.Insert(data)
}

func (d *defaultBookModel) FindOne(id int64) (data *Book, err error) {
	data = &Book{}
	err = d.conn.FindOne(id, data)
	return
}

func (d *defaultBookModel) FindOneByFilter(filter map[string]interface{}, opts ...sqlx.Option) (data *Book, err error) {
	data = &Book{}
	err = d.conn.FindOneByFilter(filter, data, opts...)
	return
}

func (d *defaultBookModel) FindByFilter(filter map[string]interface{}, opts ...sqlx.Option) (datas *[]Book, err error) {
	datas = &[]Book{}
	err = d.conn.FindByFilter(filter, datas, opts...)
	return
}

func (d *defaultBookModel) Save(data *Book) *gorm.DB {
	return d.conn.Save(data)
}

func (d *defaultBookModel) Updates(data *Book) *gorm.DB {
	return d.conn.Updates(data)
}

func (d *defaultBookModel) UpdateByFilter(filter map[string]interface{}, upVal map[string]interface{}) *gorm.DB {
	return d.conn.UpdateByFilter(Book{}, filter, upVal)
}

func (d *defaultBookModel) Delete(id int64) error {
	return d.conn.Delete(Book{}, id)
}

func (d *defaultBookModel) DeleteByFilter(filter map[string]interface{}) error {
	return d.conn.DeleteByFilter(Book{}, filter)
}
