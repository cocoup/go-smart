package mysql

const (
	category          = "model"
	modelTemplateFile = "model.tpl"
)

//type (
//	// {{.model}}Model defines a model for user
//	{{.model}}}Model interface {
//		Insert(data *{{.model}}) error
//		FindOne(id int64) (*{{model}}, error)
//		FindOneByFilter(filter map[string]interface{}, opts ...Option) (*{{model}}, error)
//		FindByFilter(filter map[string]interface{}, opts ...Option)([]*{{,model}} error)
//		Save(*{{.model}}) *gorm.DB
//		Updates(*{{.model}}) *gorm.DB
//		UpdateByFilter(filter map[string]interface{}, upVal map[string]interface{}) *gorm.DB
//		Delete(id int64) error
//		DeleteByFilter(filter map[string]interface{}) error
//	}
//
//	default{{.model}}Model struct {
//		conn  sqlx.SqlConn
//		table string
//	}
//
//	// User defines an data structure for mysql
//	{{.model}} struct {
//	{{- range .fields }}
//	{{.Name }}  {{.DataType}} {{if ne "" .Comment}} //{{.Comment}}{{end}}
//	{{- end }}
//	}
//)
//
//// New{{.model}}Model creates an instance for UserModel
//func New{{.model}}Model(conn sqlx.SqlConn) {{.model}}Model {
//	return &default{{.model}}Model{
//		conn:  conn,
//		table: "`{{.lowerModel}}`",
//	}
//}

//
//func (d *default{{.model}}Model) Insert(data *{{.model}}) error {
//	return d.conn.Insert(data)
//}
//
//func (d *default{{.model}}Model) FindOne(id int64) (data *{{.model}}, err error) {
//	err = d.conn.FindOne(id, data)
//	return
//}
//
//func (d *default{{.model}}Model) FindOneByFilter(filter map[string]interface{}, opts ...sqlx.Option) (data *{{.model}}, err error) {
//	err = d.conn.FindOneByFilter(filter, data, opts...)
//	return
//}
//
//func (d *default{{.model}}Model) FindByFilter(filter map[string]interface{}, opts ...sqlx.Option) (datas []*{{.model}}, err error) {
//	err = d.conn.FindByFilter(filter, datas, opts...)
//	return
//}
//
//func (d *default{{.model}}Model) Save(data *{{.model}}) *gorm.DB {
//	return d.conn.Save(data)
//}
//
//func (d *default{{.model}}Model) Updates(data *{{.model}}) *gorm.DB {
//	return d.conn.Updates(data)
//}
//
//func (d *default{{.model}}Model) UpdateByFilter(filter map[string]interface{}, upVal map[string]interface{}) *gorm.DB {
//	return d.conn.UpdateByFilter({{.model}}{}, filter, upVal)
//}
//
//func (d *default{{.model}}Model) Delete(id int64) error {
//	return d.conn.Delete({{.model}}{}, id)
//}
//
//func (d *default{{.model}}Model) DeleteByFilter(filter map[string]interface{}) error {
//	return d.conn.DeleteByFilter({{.model}}{}, filter)
//}
