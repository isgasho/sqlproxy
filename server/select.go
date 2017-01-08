package server

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	mysql "github.com/siddontang/go-mysql/mysql"
	"github.com/xwb1989/sqlparser"
)

func (h MysqlHandler) handleSelect(selectStatement *sqlparser.Select) (*mysql.Result, error) {

	sqlparser.Rewrite(selectStatement.From, func(origin []byte) []byte {
		s := string(origin)
		if s == "users" {
			s = "customers"
		}
		return []byte(s)
	})
	newSelect := sqlparser.String(selectStatement)
	log.Println(selectStatement, "->", newSelect)

	result, err := h.SelectDB(newSelect)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil

}

func (h MysqlHandler) SelectDB(selectStatement string) (*mysql.Result, error) {
	// 1. Exec mysql query
	rows, err := h.db.Query(selectStatement)
	if err != nil {
		return nil, err
	}

	// 2. Process result
	//字典类型, 构造scanArgs、values两个数组，scanArgs的每个值指向values相应值的地址
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))

	valueList := [][]interface{}{}

	for rows.Next() {

		values := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		// parse records
		err = rows.Scan(scanArgs...)
		valueList = append(valueList, values)

	}

	// 处理结果
	result, err := mysql.BuildSimpleResultset(
		columns,
		valueList,
		false,
	)
	// log.Println([]interface{}{values})
	if err != nil {
		return nil, err
	}

	return &mysql.Result{0, 0, 0, result}, nil
}
