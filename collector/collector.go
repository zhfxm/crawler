package collector

type DataCell struct {
	Data map[string]interface{}
}

type Storage interface {
	Save(datas ...*DataCell) error
}
