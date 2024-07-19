package sqldb

import (
	"database/sql"
	"errors"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"go.uber.org/zap"
)

type DBer interface {
	Insert(t TableData) error
}

type Sqldb struct {
	db *sql.DB
	options
}

func newSqldb() *Sqldb {
	return &Sqldb{}
}

type TableData struct {
	TableName   string
	ColumnNames []Filed       //标题字段名
	Args        []interface{} //数据
	DataCount   int           //数据量
}

type Filed struct {
	Title string
	Type  string
}

func New(opts ...Option) (*Sqldb, error) {
	options := defaultOption
	for _, opt := range opts {
		opt(&options)
	}
	d := &Sqldb{}
	d.options = options
	if err := d.OpenDB(); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *Sqldb) OpenDB() error {
	db, err := sql.Open("mysql", d.sqlUrl)
	if err != nil {
		return err
	}
	db.SetMaxIdleConns(2048)
	db.SetMaxOpenConns(2048)
	if err = db.Ping(); err != nil {
		return err
	}
	d.db = db
	return nil
}

func (d *Sqldb) Insert(t TableData) error {
	if len(t.ColumnNames) == 0 {
		return errors.New("empty columns")
	}
	sql := `INSERT INTO ` + t.TableName + ` (`

	for _, v := range t.ColumnNames {
		sql += v.Title + ","
	}

	sql = sql[:len(sql) - 1] + `) VALUES `

	blank := ",(" + strings.Repeat(",?", len(t.ColumnNames))[1:] + ")"

	sql += strings.Repeat(blank, t.DataCount)[1:] + `;`

	d.logger.Info("insert table", zap.String("sql", sql))

	_, err := d.db.Exec(sql, t.Args...)
	return err
}