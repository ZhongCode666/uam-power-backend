package dbservice

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type MySQLService struct {
	db *sql.DB
}

func NewMySQLService(dsn string) (*MySQLService, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	// 配置数据库连接池
	// 设置最大打开连接数
	db.SetMaxOpenConns(50) // 根据实际需求调整
	// 设置最大空闲连接数
	db.SetMaxIdleConns(10) // 根据实际需求调整
	// 设置连接的最大生命周期
	db.SetConnMaxLifetime(30 * time.Minute) // 连接最大生命周期为30分钟
	return &MySQLService{db: db}, nil
}

func (s *MySQLService) ExecuteCmd(sql string) (int, error) {
	result, err := s.db.Exec(sql)
	if err != nil {
		return 0, err
	}
	id, _ := result.RowsAffected()
	return int(id), nil
}

func (s *MySQLService) QueryRow(query string, args ...interface{}) (map[string]interface{}, error) {
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

	// 移动到第一行
	if !rows.Next() {
		return nil, sql.ErrNoRows // 如果没有结果，返回 sql.ErrNoRows
	}

	// 创建一个切片，用于保存列的值
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// 扫描查询结果到 values 中
	if err := rows.Scan(valuePtrs...); err != nil {
		return nil, err
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

	return result, nil
}

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
