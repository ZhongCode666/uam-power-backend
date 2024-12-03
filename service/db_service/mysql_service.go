package dbservice

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

// MySQLService 结构体，包含一个数据库连接对象
type MySQLService struct {
	db *sql.DB
}

// NewMySQLService 创建一个新的 MySQLService 实例
// dsn 是数据源名称
// 返回一个 MySQLService 实例和一个错误
func NewMySQLService(dsn string) (*MySQLService, error) {
	db, err := sql.Open("mysql", dsn) // 打开数据库连接
	if err != nil {
		return nil, err // 如果打开数据库连接失败，返回错误
	}
	if err := db.Ping(); err != nil {
		return nil, err // 如果数据库连接不可用，返回错误
	}
	db.SetMaxIdleConns(5)                   // 设置最大空闲连接数为 5
	db.SetConnMaxLifetime(30 * time.Second) // 设置连接的最大生命周期为 30 秒
	return &MySQLService{db: db}, nil       // 返回 MySQLService 实例
}

// ExecuteCmd 执行给定的 SQL 命令
// 返回受影响的行数和一个错误
func (s *MySQLService) ExecuteCmd(sql string) (int, error) {
	result, err := s.db.Exec(sql) // 执行 SQL 命令
	if err != nil {
		return 0, err // 如果执行失败，返回错误
	}
	id, _ := result.RowsAffected() // 获取受影响的行数
	return int(id), nil            // 返回受影响的行数
}

// QueryRow 执行给定的查询并返回结果的第一行
// query 是要执行的 SQL 查询
// args 是查询的参数
// 返回一个包含列名和对应值的 map 和一个错误
func (s *MySQLService) QueryRow(query string, args ...interface{}) (map[string]interface{}, error) {
	// 执行查询
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err // 如果查询失败，返回错误
	}

	// 获取列名
	columns, err := rows.Columns()
	if err != nil {
		return nil, err // 如果获取列名失败，返回错误
	}

	// 移动到第一行
	if !rows.Next() {
		return nil, sql.ErrNoRows // 如果没有结果，返回 sql.ErrNoRows
	}

	// 创建一个切片，用于保存列的值
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i] // 将每列的指针存储到 valuePtrs 中
	}

	// 扫描查询结果到 values 中
	if err := rows.Scan(valuePtrs...); err != nil {
		return nil, err // 如果扫描失败，返回错误
	}

	// 将结果转换为 map
	result := make(map[string]interface{})
	for i, colName := range columns {
		val := values[i]

		// 如果值是 []byte 类型，将其转换为 string
		if b, ok := val.([]byte); ok {
			result[colName] = string(b)
		} else {
			result[colName] = val
		}
	}

	return result, nil // 返回结果 map
}

// QueryRows 执行给定的查询并返回所有结果行
// query 是要执行的 SQL 查询
// args 是查询的参数
// 返回一个包含列名和对应值切片的 map 和一个错误
func (s *MySQLService) QueryRows(query string, args ...interface{}) (map[string][]interface{}, error) {
	// 执行查询
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	// 获取列名
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// 初始化结果 map
	result := make(map[string][]interface{})
	for _, colName := range columns {
		result[colName] = []interface{}{}
	}

	// 遍历所有行
	for rows.Next() {
		// 创建用于存储每一列值的切片
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// 扫描每行数据到 values 中
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		// 将每列值追加到相应的列名切片中
		for i, colName := range columns {
			val := values[i]
			// 转换 []byte 为 string
			if b, ok := val.([]byte); ok {
				result[colName] = append(result[colName], string(b))
			} else {
				result[colName] = append(result[colName], val)
			}
		}
	}

	// 检查是否有错误
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
