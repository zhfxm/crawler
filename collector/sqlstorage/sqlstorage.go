package sqlstorage

import (
	"encoding/json"

	"github.com/zhfxm/simple-crawler/collector"
	"github.com/zhfxm/simple-crawler/sqldb"
)

type SqlStore struct {
	dataDocker []*collector.DataCell //分批输出结果缓存
	db         sqldb.DBer
	options
}

func NewStorage(opts ...Option) (*SqlStore, error) {
	options := defaultOption
	for _, opt := range opts {
		opt(&options)
	}
	s := &SqlStore{}
	s.options = options
	db, err := sqldb.New(
		sqldb.WidthSqlurl(s.sqlurl),
		sqldb.WithLogger(s.logger),
	)
	if err != nil {
		return nil, err
	}
	s.db = db
	return s, nil
}

var fields []string = []string{
	"hsNow",
	"hsNowUnit",
	"Hs24",
	"Hs24Unit",
}

func (s *SqlStore) Save(dataCells ...*collector.DataCell) error {
	for _, dataCell := range dataCells {
		s.dataDocker = append(s.dataDocker, dataCell)
		if len(s.dataDocker) >= s.BatchCount {
			return s.Flush()
		}
	}
	return nil
}

func (s *SqlStore) Flush() error {
	if len(s.dataDocker) == 0 {
		return nil
	}
	args := make([]interface{}, 0)
	values := []string{}
	for _, dataCell := range s.dataDocker {
		d := dataCell.Data
		for _, field := range fields {
			v := d[field]
			switch v := v.(type) {
			case nil:
				values = append(values, "")
			case string:
				values = append(values, v)
			default:
				b, err := json.Marshal(v)
				if err != nil {
					values = append(values, "")
				} else {
					values = append(values, string(b))
				}
			}
		}
		for _, v := range values {
			args = append(args, v)
		}
	}
	return s.db.Insert(sqldb.TableData{
		TableName:   "t_test",
		ColumnNames: getFields(),
		Args:        args,
		DataCount:   len(s.dataDocker),
	})
}

func getFields() []sqldb.Filed {
	return []sqldb.Filed{
		{Title: "hsNow", Type: "VARCHAR(255)"},
		{Title: "hsNowUnit", Type: "VARCHAR(255)"},
		{Title: "Hs24", Type: "VARCHAR(255)"},
		{Title: "Hs24Unit", Type: "VARCHAR(255)"},
	}
}