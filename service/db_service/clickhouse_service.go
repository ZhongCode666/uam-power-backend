package dbservice

import (
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"strings"
	"sync"
	"time"
)

// ClickHouse 封装类，包含缓冲区功能
type ClickHouse struct {
	conn        *sql.DB
	batchSize   int
	flushTimer  *time.Ticker
	buffer      map[string][][]interface{} // 表名到数据的缓冲区映射
	mu          sync.Mutex                 // 互斥锁，确保线程安全
	flushPeriod time.Duration              // 刷新周期
	columns     []string
}

// NewClickHouse 创建一个新的 ClickHouse 实例并初始化缓冲区
func NewClickHouse(
	Host string, Port int, Username string, Password string, BatchSize int, FlushPeriod int,
	database string, InitTask bool, Columns []string,
) (*ClickHouse, error) {
	//host, port, username, password, database string, batchSize int, flushPeriod time.Duration
	conn := sql.OpenDB(clickhouse.Connector(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", Host, Port)},
		Auth: clickhouse.Auth{
			Database: database,
			Username: Username,
			Password: Password,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 30, // 最大执行时间（秒）
		},
		DialTimeout: 10 * time.Second,
		ReadTimeout: 10 * time.Second,
		Debug:       false,
	}))

	// 配置连接池参数
	conn.SetConnMaxLifetime(30 * time.Second)
	conn.SetMaxOpenConns(1000)

	// 测试连接
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping ClickHouse: %w", err)
	}

	// 初始化 ClickHouse 实例
	ch := &ClickHouse{
		conn:        conn,
		batchSize:   BatchSize,
		buffer:      make(map[string][][]interface{}),
		flushTimer:  time.NewTicker(time.Duration(FlushPeriod) * time.Second),
		flushPeriod: time.Duration(FlushPeriod) * time.Second,
		columns:     Columns,
	}

	// 启动定时刷新任务
	if InitTask {
		go ch.startFlushWorker()
	}
	return ch, nil
}

// Add 向缓冲区添加一条数据
func (ch *ClickHouse) Add(table string, columns []string, row []interface{}) error {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	// 如果表不存在，初始化缓冲区
	if _, exists := ch.buffer[table]; !exists {
		ch.buffer[table] = make([][]interface{}, 0)
	}

	// 添加数据到缓冲区
	ch.buffer[table] = append(ch.buffer[table], row)

	// 如果缓冲区达到批量大小，立即写入
	if len(ch.buffer[table]) >= ch.batchSize {
		if err := ch.flushTable(table, columns); err != nil {
			return fmt.Errorf("failed to flush table %s: %w", table, err)
		}
	}

	return nil
}

func (ch *ClickHouse) ExecuteCmd(sql string) error {
	if _, err := ch.conn.Exec(sql); err != nil {
		return err
	}
	return nil
}

// flushTable 将某个表的缓冲区写入数据库
func (ch *ClickHouse) flushTable(table string, columns []string) error {
	if len(ch.buffer[table]) == 0 {
		return nil
	}

	// 获取缓冲区数据
	batch := ch.buffer[table]
	ch.buffer[table] = make([][]interface{}, 0) // 清空缓冲区

	// 批量插入数据
	return ch.Insert(table, columns, batch)
}

// Insert 批量插入数据
func (ch *ClickHouse) Insert(table string, columns []string, values [][]interface{}) error {
	columnsStr := "(" + joinStrings(columns, ",") + ")"
	query := fmt.Sprintf("INSERT INTO %s %s VALUES %s", table, columnsStr, valuesToString(values))
	return ch.ExecuteCmd(query)
}

func (ch *ClickHouse) InsertOne(table string, columns []string, values []interface{}) error {
	columnsStr := "(" + joinStrings(columns, ",") + ")"
	query := fmt.Sprintf("INSERT INTO %s %s VALUES %s", table, columnsStr, valuesToStringSingle(values))
	return ch.ExecuteCmd(query)
}

// StartFlushWorker 启动定时刷新任务
func (ch *ClickHouse) startFlushWorker() {
	for range ch.flushTimer.C {
		ch.mu.Lock()
		for table := range ch.buffer {
			columns := ch.getColumnsForTable()
			if err := ch.flushTable(table, columns); err != nil {
				fmt.Printf("Failed to flush table %s: %v\n", table, err)
			}
		}
		ch.mu.Unlock()
	}
}

// Close 停止定时刷新任务并清空所有缓冲区
func (ch *ClickHouse) Close() error {
	ch.flushTimer.Stop()
	ch.mu.Lock()
	defer ch.mu.Unlock()

	// 写入所有剩余数据
	for table := range ch.buffer {
		columns := ch.getColumnsForTable()
		if err := ch.flushTable(table, columns); err != nil {
			fmt.Printf("Failed to flush table %s: %v\n", table, err)
		}
	}

	return ch.conn.Close()
}

// getColumnsForTable 获取表对应的列名（可以根据需要动态实现）
func (ch *ClickHouse) getColumnsForTable() []string {
	// 示例：返回固定列名，实际可通过查询数据库元数据动态实现
	return ch.columns
}

// 工具函数：拼接字符串
func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}

func valuesToString(values [][]interface{}) string {
	var rows []string

	for _, row := range values {
		var formattedRow []string
		for _, value := range row {
			switch v := value.(type) {
			case string:
				// 字符串需要加单引号，并转义
				formattedRow = append(formattedRow, fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''")))
			case nil:
				// 空值处理
				formattedRow = append(formattedRow, "NULL")
			default:
				// 其他类型直接转换为字符串
				formattedRow = append(formattedRow, fmt.Sprintf("%v", v))
			}
		}
		// 拼接单行数据
		rows = append(rows, fmt.Sprintf("(%s)", strings.Join(formattedRow, ", ")))
	}

	// 拼接所有行
	return strings.Join(rows, ",\n")
}

func valuesToStringSingle(values []interface{}) string {
	var formattedRow []string
	for _, value := range values {
		switch v := value.(type) {
		case string:
			// 字符串需要加单引号，并转义
			formattedRow = append(formattedRow, fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''")))
		case nil:
			// 空值处理
			formattedRow = append(formattedRow, "NULL")
		default:
			// 其他类型直接转换为字符串
			formattedRow = append(formattedRow, fmt.Sprintf("%v", v))
		}
	}
	return fmt.Sprintf("(%s)", strings.Join(formattedRow, ", "))
}
