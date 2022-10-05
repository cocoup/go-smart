package redis

import (
	"golang.org/x/net/context"
	"reflect"
	"testing"
	"time"
)

func testConn() RedisConn {
	conf := Config{
		Addrs:    []string{":26380", "26381"},
		Password: "123456",
		DB:       0,
		Master:   "mymaster",
	}
	return NewConn(conf, Sentinel())
}

func Test_redisConn_Set(t *testing.T) {
	type args struct {
		key  string
		val  interface{}
		args []int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{key: "key1", val: 10, args: []int{100}},
		},
	}
	r := testConn()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := r.Set(tt.args.key, tt.args.val, tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_redisConn_Eval(t *testing.T) {
	type args struct {
		script string
		keys   []string
		args   []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantVal interface{}
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				script: "return {KEYS[1],KEYS[2],ARGV[1],ARGV[2]}",
				keys:   []string{"keykey1", "keykey2"},
				args:   []interface{}{"小明", 20},
			},
			wantVal: []interface{}{"keykey1", "keykey2", "小明", "20"},
		},
	}

	r := testConn()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := r.Eval(tt.args.script, tt.args.keys, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("Eval() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func Test_redisConn_Ping(t *testing.T) {
	tests := []struct {
		name    string
		wantVal string
		wantErr bool
	}{
		{
			name:    "test1",
			wantVal: "PONG",
		},
	}
	r := testConn()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := r.Ping()
			if (err != nil) != tt.wantErr {
				t.Errorf("Ping() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVal != tt.wantVal {
				t.Errorf("Ping() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func Test_redisConn_Pipelined(t *testing.T) {
	type args struct {
		fn func(Pipeliner) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{fn: func(pipe Pipeliner) error {
				pipe.Set(context.Background(), "test", "xxx", 10*time.Second)
				pipe.Exec(context.Background())
				return nil
			}},
		},
	}
	r := testConn()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := r.Pipelined(tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("Pipelined() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_redisConn_TxPipelined(t *testing.T) {
	type args struct {
		fn func(Pipeliner) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{fn: func(pipe Pipeliner) error {
				pipe.Set(context.Background(), "test", "xxx", 10*time.Second)
				return nil
			}},
		},
	}
	r := testConn()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := r.TxPipelined(tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("Pipelined() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
