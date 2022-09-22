package sqlx

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

type Account struct {
	UID      string
	UNO      string
	NickName string
	Sex      int
}

type ExchangeCode struct {
	ID    int64
	Code  string
	Logic string
	Param int
	Use   int
}

type TableTest struct {
	Id        int64
	A         string //测试a
	B         string //bbbbb
	C         string //cc
	D         int64  //ddd
	E         int64  //eee
	CreatedAt time.Time
	UpdatedAt time.Time
}

func testConn() SqlConn {
	conf := Config{
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
	conn, err := NewConn(conf)
	if nil != err {
		log.Fatal(err)
	}

	return conn
}

func Test_sqlConn_Insert(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Insert_case1",
			args: args{val: &TableTest{
				A: "aaa",
				B: "bbb",
				C: "ccc",
				D: 100,
				E: 10,
			}},
		},
		{
			name: "Insert_case2",
			args: args{val: &TableTest{
				A:         "xxx",
				B:         "yyy",
				CreatedAt: time.Unix(time.Now().Unix()+10*24*3600, 0),
			}},
		},
	}
	conn := testConn()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := conn.Insert(tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_sqlConn_FindOne(t *testing.T) {
	type args struct {
		id int64
		//out *ExchangeCode
		out *TableTest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		//{
		//	name: "FindOne-case1",
		//	args: args{
		//		id:  5005,
		//		out: &ExchangeCode{},
		//	},
		//	wantErr: true,
		//},
		//{
		//	name: "FindOne-case2",
		//	args: args{
		//		id:  5505,
		//		out: &ExchangeCode{},
		//	},
		//	wantErr: false,
		//},
		{
			name: "FindOne-case3",
			args: args{
				id:  2,
				out: &TableTest{},
			},
			wantErr: false,
		},
	}
	conn := testConn()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := conn.FindOne(tt.args.id, tt.args.out); (err != nil) != tt.wantErr {
				t.Errorf("FindOne() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_sqlConn_FindOneByFilter(t *testing.T) {
	type args struct {
		filter map[string]interface{}
		out    interface{}
		opts   []Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "FindOne-case1",
			args: args{
				filter: map[string]interface{}{
					"logic": "GAME",
				},
				out: &ExchangeCode{},
			},
			wantErr: false,
		},
		{
			name: "FindOne-case2",
			args: args{
				filter: map[string]interface{}{
					"logic": "GAME",
				},
				out: &ExchangeCode{},
				opts: []Option{
					OrderOption("id desc"),
				},
			},
			wantErr: false,
		},
	}

	conn := testConn()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := conn.FindOneByFilter(tt.args.filter, tt.args.out, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("FindOneByFilter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_sqlConn_FindByFilter(t *testing.T) {
	type args struct {
		filter map[string]interface{}
		out    interface{}
		opts   []Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "FindByFilter_case1",
			args: args{
				filter: map[string]interface{}{},
				out:    &[]ExchangeCode{},
				opts: []Option{
					PageOption(1, 2),
					OrderOption("id asc"),
				},
			},
		},
		{
			name: "FindByFilter_case1",
			args: args{
				filter: map[string]interface{}{},
				out:    &[]ExchangeCode{},
				opts: []Option{
					PageOption(1, 2),
					OrderOption("id desc"),
				},
			},
		},
		{
			name: "FindByFilter_case1",
			args: args{
				filter: map[string]interface{}{
					"logic": "ROOT",
				},
				out: &[]ExchangeCode{},
				opts: []Option{
					PageOption(1, 2),
					OrderOption("id desc"),
				},
			},
		},
	}
	conn := testConn()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := conn.FindByFilter(tt.args.filter, tt.args.out, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("FindByFilter() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				fmt.Println(tt.args.out)
			}
		})
	}
}

func Test_sqlConn_Save(t *testing.T) {
	type args struct {
		val *ExchangeCode
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Save_case1", //ID未赋值更新(会新插入一行)
			args: args{val: &ExchangeCode{Code: "8888", Logic: "xxxx"}},
		},
		{
			name: "Save_case2", //未赋值字段零值更新
			args: args{val: &ExchangeCode{ID: 5505, Code: "6666", Logic: "yyyy"}},
		},
		{
			name: "Save_case3", //不存在记录更新
			args: args{val: &ExchangeCode{ID: 1, Code: "9999", Logic: "zzzz"}},
		},
		{
			name: "Save_case4", //所有字段零值更新，ID除外
			args: args{val: &ExchangeCode{ID: 5503}},
		},
	}
	conn := testConn()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := conn.Save(tt.args.val); !assert.Equal(t, got.Error, nil) {
				t.Errorf("Save() = %v", got.Error)
			}
		})
	}
}

func Test_sqlConn_Updates(t *testing.T) {
	type args struct {
		val *ExchangeCode
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Save_case1", //ID未赋值会报错
			args:    args{val: &ExchangeCode{Code: "8888", Logic: "xxxx"}},
			wantErr: true,
		},
		{
			name:    "Save_case2", //未赋值字段零值不更新
			args:    args{val: &ExchangeCode{ID: 5505, Code: "6666"}},
			wantErr: false,
		},
		{
			name:    "Save_case3", //不存在记录更新
			args:    args{val: &ExchangeCode{ID: 1, Code: "9999", Logic: "zzzz"}},
			wantErr: false,
		},
		{
			name:    "Save_case4", //所有字段零值不更新
			args:    args{val: &ExchangeCode{ID: 5522}},
			wantErr: false,
		},
	}
	conn := testConn()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := conn.Updates(tt.args.val); !assert.Equal(t, got.Error != nil, tt.wantErr) {
				t.Errorf("Save() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func Test_sqlConn_UpdateByFilter(t *testing.T) {
	type args struct {
		model  interface{}
		filter map[string]interface{}
		upVal  map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "UpdateByFilter_case1",
			args: args{
				model: ExchangeCode{},
				filter: map[string]interface{}{
					"code": "9999",
				},
				upVal: map[string]interface{}{
					"code":  "0000",
					"logic": "",
				},
			},
		},
		{
			name: "UpdateByFilter_case2",
			args: args{
				model: ExchangeCode{},
				filter: map[string]interface{}{
					"code": "6666",
				},
				upVal: map[string]interface{}{
					"logic": "8888",
				},
			},
		},
		{
			name: "UpdateByFilter_case1",
			args: args{
				model:  ExchangeCode{},
				filter: map[string]interface{}{},
				upVal: map[string]interface{}{
					"code":  "0000",
					"logic": "",
				},
			},
			wantErr: true,
		},
	}
	conn := testConn()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := conn.UpdateByFilter(tt.args.model, tt.args.filter, tt.args.upVal); !assert.Equal(t, got.Error != nil, tt.wantErr) {
				t.Errorf("UpdateByFilter() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func Test_sqlConn_Delete(t *testing.T) {
	type args struct {
		model interface{}
		id    int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Delete_case1", //ID未赋值，默认为0
			args: args{
				model: ExchangeCode{},
			},
			wantErr: false,
		},
		{
			name: "Delete_case1", //正常删除
			args: args{
				model: ExchangeCode{},
				id:    1,
			},
			wantErr: false,
		},
	}
	conn := testConn()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := conn.Delete(tt.args.model, tt.args.id); !assert.Equal(t, err != nil, tt.wantErr) {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_sqlConn_DeleteByFilter(t *testing.T) {
	type args struct {
		model  interface{}
		filter map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "DeleteByFilter_case1",
			args: args{
				model:  ExchangeCode{},
				filter: map[string]interface{}{},
			},
			wantErr: true,
		},
		{
			name: "DeleteByFilter_case2",
			args: args{
				model: ExchangeCode{},
				filter: map[string]interface{}{
					"code": "8888",
				},
			},
		},
	}
	conn := testConn()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := conn.DeleteByFilter(tt.args.model, tt.args.filter); !assert.Equal(t, err != nil, tt.wantErr) {
				t.Errorf("DeleteByFilter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
