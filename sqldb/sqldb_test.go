package sqldb_test

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zhfxm/simple-crawler/log"
	"github.com/zhfxm/simple-crawler/sqldb"
	"go.uber.org/zap/zapcore"
)

func TestSqlBuild(t *testing.T) {
	plugin := log.NewStdoutPlugin(zapcore.InfoLevel)
	logger := log.NewLogger(plugin)
	logger.Info("logger init end")
	
	dbUrl := "root:root@tcp(localhost:3306)/cloud_crawler?charset=utf8"

	db, _ := sqldb.New(
		sqldb.WidthSqlurl(dbUrl),
		sqldb.WithLogger(logger),
	)

	tb := sqldb.TableData{}
	tb.DataCount = 3
	tb.TableName = "sql_test"
	tb.ColumnNames = []sqldb.Filed{{Title: "name", Type: "VARCHAR(255)"}, {Title: "address", Type: "VARCHAR(255)"}}
	tb.Args = []interface{}{"zh1", "addr1","zh2", "addr2","zh3", "addr3"}
	err := db.Insert(tb)
	fmt.Println(err)
}