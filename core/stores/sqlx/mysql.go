package sqlx

import (
	"golang.org/x/net/context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type PageInfo struct {
	Page     int `json:"page" form:"page"`         // 页码
	PageSize int `json:"pageSize" form:"pageSize"` // 每页大小
}

type Option func(db *gorm.DB) error

func OrderOption(order string) Option {
	return func(db *gorm.DB) error {
		return db.Order(order).Error
	}
}

func PageOption(info PageInfo) Option {
	return func(db *gorm.DB) error {
		return db.Limit(info.PageSize).Offset((info.Page - 1) * info.PageSize).Error
	}
}

type SqlConn interface {
	DB() *gorm.DB
	Raw(sql string, values ...interface{}) *gorm.DB
	//TODO 携带trace相关逻辑
	RawCtx(ctx context.Context, sql string, values ...interface{}) *gorm.DB
	Insert(val interface{}) error
	InsertCtx(ctx context.Context, val interface{}) error
	FindOne(id int64, out interface{}) error
	FindOneByFilter(filter map[string]interface{}, out interface{}, opts ...Option) error
	FindByFilter(filter map[string]interface{}, out interface{}, opts ...Option) error
	//更新记录，包含零值字段
	Save(val interface{}) *gorm.DB
	//更新记录，非零值字段
	Updates(val interface{}) *gorm.DB
	//根据条件更新记录
	UpdateByFilter(model interface{}, filter map[string]interface{}, upVal map[string]interface{}) *gorm.DB
	//删除主键记录
	Delete(model interface{}, id int64) error
	//根据条件删除记录
	DeleteByFilter(model interface{}, filter map[string]interface{}) error
}

type sqlConn struct {
	db *gorm.DB
}

func NewConn(conf SqlConf) (SqlConn, error) {
	sqlConf := mysql.Config{
		DSN:                       conf.DNS(), // DSN data source name
		DefaultStringSize:         255,        // string 类型字段的默认长度
		SkipInitializeWithVersion: false,      // 根据版本自动配置
	}

	gormDB, err := gorm.Open(mysql.New(sqlConf), gormConfig(conf))
	if err != nil {
		return nil, err
	}

	sqlDB, err := gormDB.DB()
	if nil != err {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
	sqlDB.SetMaxOpenConns(conf.MaxOpenConns)

	return &sqlConn{db: gormDB}, nil
}

func (s *sqlConn) DB() *gorm.DB {
	return s.db
}

func (s *sqlConn) Raw(sql string, values ...interface{}) *gorm.DB {
	return s.RawCtx(context.Background(), sql, values...)
}

func (s *sqlConn) RawCtx(ctx context.Context, sql string, values ...interface{}) *gorm.DB {
	return s.db.Exec(sql, values...)
}

func (s *sqlConn) Insert(val interface{}) error {
	return s.InsertCtx(context.Background(), val)
}

func (s *sqlConn) InsertCtx(ctx context.Context, val interface{}) error {
	return s.db.Create(val).Error
}

func (s *sqlConn) FindOne(id int64, out interface{}) error {
	return s.db.First(out, id).Error
}

func (s *sqlConn) FindOneByFilter(filter map[string]interface{}, out interface{}, opts ...Option) error {
	db := s.db.Where(filter)
	for _, opt := range opts {
		err := opt(db)
		if nil != err {
			return err
		}
	}
	return db.Limit(1).Find(out).Error
}

func (s *sqlConn) FindByFilter(filter map[string]interface{}, out interface{}, opts ...Option) error {
	db := s.db.Where(filter)
	for _, opt := range opts {
		err := opt(db)
		if nil != err {
			return err
		}
	}
	return db.Find(out).Error
}

func (s *sqlConn) Save(val interface{}) *gorm.DB {
	return s.db.Save(val)
}

func (s *sqlConn) Updates(val interface{}) *gorm.DB {
	return s.db.Updates(val)
}

func (s *sqlConn) UpdateByFilter(model interface{}, filter map[string]interface{}, upVal map[string]interface{}) *gorm.DB {
	return s.db.Model(model).Where(filter).Updates(upVal)
}

func (s *sqlConn) Delete(model interface{}, id int64) error {
	return s.db.Delete(model, id).Error
}

func (s *sqlConn) DeleteByFilter(model interface{}, filter map[string]interface{}) error {
	return s.db.Where(filter).Delete(model).Error
}
