package model

import (
	"fmt"
	"github.com/cocoup/go-smart/core/stores/sqlx"
	"log"
	"testing"
)

func testConn() sqlx.SqlConn {
	conf := sqlx.Config{
		IP:           "127.0.0.1",
		Port:         "3306",
		DB:           "my_db",
		Name:         "root",
		Password:     "root",
		Option:       "charset=utf8mb4&parseTime=True&loc=Local",
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		LogMode:      "info",
	}
	conn, err := sqlx.NewConn(conf)
	if nil != err {
		log.Fatal(err)
	}

	return conn
}

func Test_defaultBookModel_FindOne(t *testing.T) {
	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "FindOne_case1",
			args: args{id: 1},
		},
	}

	sqlConn := testConn()
	bookModel := NewBookModel(sqlConn)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := bookModel.FindOne(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(gotData)
		})
	}
}

func Test_defaultBookModel_FindByFilter(t *testing.T) {
	type args struct {
		filter map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "FindByFilter_case1",
			args: args{filter: map[string]interface{}{
				"name": "语文",
			}},
		},
		{
			name: "FindByFilter_case1",
			args: args{filter: map[string]interface{}{
				"name": "数学",
			}},
		},
	}

	sqlConn := testConn()
	bookModel := NewBookModel(sqlConn)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := bookModel.FindByFilter(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(gotData)
		})
	}
}
