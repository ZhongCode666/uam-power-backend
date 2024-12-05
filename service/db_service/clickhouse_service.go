package dbservice

import (
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"strings"
	"sync"
	"time"
	"uam-power-backend/utils"
)

// ClickHouse 封装类，包含缓冲区功能
type ClickHouse struct {
	conn        *sql.DB                    // 数据库连接
	batchSize   int                        // 批量大小
	flushTimer  *time.Ticker               // 刷新定时器
	buffer      map[string][][]interface{} // 表名到数据的缓冲区映射
	mu          sync.Mutex                 // 互斥锁，确保线程安全
	flushPeriod time.Duration              // 刷新周期
	columns     []string                   // 列名
}

// NewClickHouse 创建一个新的 ClickHouse 实例并初始化缓冲区
func NewClickHouse(
	Host string, Port int, Username string, Password string, BatchSize int, FlushPeriod int,
	database string, InitTask bool, Columns []string,
) (*ClickHouse, error) {
	// 创建 ClickHouse 数据库连接
	conn := sql.OpenDB(clickhouse.Connector(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", Host, Port)}, // 设置地址
		Auth: clickhouse.Auth{
			Database: database, // 设置数据库名
			Username: Username, // 设置用户名
			Password: Password, // 设置密码
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 10, // 最大执行时间（秒）
		},
		DialTimeout: 10 * time.Second, // 拨号超时时间
		ReadTimeout: 10 * time.Second, // 读取超时时间
		Debug:       false,            // 是否开启调试模式
	}))

	// 配置连接池参数
	conn.SetConnMaxLifetime(10 * time.Second) // 设置连接的最大生命周期
	conn.SetMaxOpenConns(1000)                // 设置最大打开连接数

	// 测试连接
	if err := conn.Ping(); err != nil {
		utils.MsgError("        [ClickHouse]ping failed: >" + err.Error())
		return nil, fmt.Errorf("failed to ping ClickHouse: %w", err) // 如果连接失败，返回错误
	}

	// 初始化 ClickHouse 实例
	ch := &ClickHouse{
		conn:        conn,                                                     // 数据库连接
		batchSize:   BatchSize,                                                // 批量大小
		buffer:      make(map[string][][]interface{}),                         // 初始化缓冲区
		flushTimer:  time.NewTicker(time.Duration(FlushPeriod) * time.Second), // 刷新定时器
		flushPeriod: time.Duration(FlushPeriod) * time.Second,                 // 刷新周期
		columns:     Columns,                                                  // 列名
	}

	// 启动定时刷新任务
	if InitTask {
		go ch.startFlushWorker() // 启动刷新任务的协程
	}
	return ch, nil // 返回 ClickHouse 实例
}

// Add 向缓冲区添加一条数据
func (ch *ClickHouse) Add(table string, columns []string, row []interface{}) error {
	ch.mu.Lock()         // 加锁，确保线程安全
	defer ch.mu.Unlock() // 在函数结束时解锁

	// 如果表不存在，初始化缓冲区
	if _, exists := ch.buffer[table]; !exists {
		ch.buffer[table] = make([][]interface{}, 0)
	}

	// 添加数据到缓冲区
	ch.buffer[table] = append(ch.buffer[table], row)

	// 如果缓冲区达到批量大小，立即写入
	if len(ch.buffer[table]) >= ch.batchSize {
		if err := ch.flushTable(table, columns); err != nil {
			return fmt.Errorf("failed to flush table %s: %w", table, err) // 返回错误信息
		}
	}

	return nil // 返回 nil 表示成功
}

// ExecuteCmd 执行给定的 SQL 命令
func (ch *ClickHouse) ExecuteCmd(sql string) error {
	// 执行 SQL 命令
	if _, err := ch.conn.Exec(sql); err != nil {
		utils.MsgError("        [ClickHouse]execute failed: >" + err.Error())
		return err // 如果执行失败，返回错误
	}
	return nil // 返回 nil 表示成功
}

// flushTable 将某个表的缓冲区写入数据库
func (ch *ClickHouse) flushTable(table string, columns []string) error {
	// 如果缓冲区为空，直接返回
	if len(ch.buffer[table]) == 0 {
		return nil
	}

	// 获取缓冲区数据
	batch := ch.buffer[table]
	// 清空缓冲区
	ch.buffer[table] = make([][]interface{}, 0)

	// 批量插入数据
	return ch.Insert(table, columns, batch)
}

// Insert 批量插入数据
func (ch *ClickHouse) Insert(table string, columns []string, values [][]interface{}) error {
	// 将列名数组转换为字符串
	columnsStr := "(" + joinStrings(columns, ",") + ")"
	// 构建插入 SQL 语句
	query := fmt.Sprintf("INSERT INTO %s %s VALUES %s", table, columnsStr, valuesToString(values))
	// 执行 SQL 语句
	return ch.ExecuteCmd(query)
}

// InsertOne 向指定表中插入一条数据
func (ch *ClickHouse) InsertOne(table string, columns []string, values []interface{}) error {
	// 将列名数组转换为字符串
	columnsStr := "(" + joinStrings(columns, ",") + ")"
	// 构建插入 SQL 语句
	query := fmt.Sprintf("INSERT INTO %s %s VALUES %s", table, columnsStr, valuesToStringSingle(values))
	// 执行 SQL 语句
	return ch.ExecuteCmd(query)
}

// startFlushWorker 启动定时刷新任务
func (ch *ClickHouse) startFlushWorker() {
	// 遍历定时器的通道
	for range ch.flushTimer.C {
		ch.mu.Lock() // 加锁，确保线程安全
		// 遍历缓冲区中的所有表
		for table := range ch.buffer {
			columns := ch.getColumnsForTable() // 获取表的列名
			// 刷新表的数据
			if err := ch.flushTable(table, columns); err != nil {
				fmt.Printf("Failed to flush table %s: %v\n", table, err) // 打印错误信息
			}
		}
		ch.mu.Unlock() // 解锁
	}
}

// Close 停止定时刷新任务并清空所有缓冲区
func (ch *ClickHouse) Close() error {
	// 停止刷新定时器
	ch.flushTimer.Stop()
	// 加锁，确保线程安全
	ch.mu.Lock()
	// 在函数结束时解锁
	defer ch.mu.Unlock()

	// 写入所有剩余数据
	for table := range ch.buffer {
		// 获取表的列名
		columns := ch.getColumnsForTable()
		// 刷新表的数据
		if err := ch.flushTable(table, columns); err != nil {
			utils.MsgError("        [ClickHouse]flush table failed: >" + err.Error())
			// 打印错误信息
			fmt.Printf("Failed to flush table %s: %v\n", table, err)
		}
	}

	// 关闭数据库连接
	return ch.conn.Close()
}

// getColumnsForTable 获取表对应的列名（可以根据需要动态实现）
func (ch *ClickHouse) getColumnsForTable() []string {
	// 示例：返回固定列名，实际可通过查询数据库元数据动态实现
	return ch.columns
}

// 工具函数：拼接字符串
// joinStrings 将字符串数组 strs 中的元素用分隔符 sep 拼接成一个字符串
func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep // 在每个元素之间添加分隔符
		}
		result += s // 添加字符串元素
	}
	return result // 返回拼接后的字符串
}

// valuesToString 将二维数组 values 转换为 SQL 插入语句中的值字符串
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

// valuesToStringSingle 将一维数组 values 转换为 SQL 插入语句中的值字符串
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
