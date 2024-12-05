package dbservice

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"sync"
	"time"
	"uam-power-backend/utils"
)

// MySQLWithBufferService 结构体，包含数据库连接、数据缓冲区、互斥锁、列名和间隔时间
type MySQLWithBufferService struct {
	db       *sql.DB                    // 数据库连接对象
	data     map[string][][]interface{} // 数据缓冲区，按表名存储
	mu       sync.Mutex                 // 互斥锁，用于保护数据缓冲区
	columns  []string                   // 列名
	interval time.Duration              // 间隔时间，用于定时刷新数据
}

// NewMySQLWithBufferService 创建一个新的 MySQLWithBufferService 实例
// dsn 是数据源名称
// interval 是定时刷新数据的间隔时间（秒）
// columns 是列名列表
// 返回一个 MySQLWithBufferService 实例和一个错误
func NewMySQLWithBufferService(dsn string, interval int, columns []string) (*MySQLWithBufferService, error) {
	db, err := sql.Open("mysql", dsn) // 打开数据库连接
	if err != nil {
		return nil, err // 如果打开数据库连接失败，返回错误
	}
	if err := db.Ping(); err != nil {
		utils.MsgError("          [MySQLWithBufferService]ping failed: >" + err.Error())
		return nil, err // 如果数据库连接不可用，返回错误
	}
	db.SetConnMaxLifetime(30 * time.Second) // 设置连接的最大生命周期为 30 秒
	db.SetMaxIdleConns(5)                   // 设置最大空闲连接数为 5
	ser := &MySQLWithBufferService{
		db: db, data: make(map[string][][]interface{}),
		interval: time.Duration(interval) * time.Second, columns: columns} // 初始化 MySQLWithBufferService 实例
	ser.Start()                                                              // 启动定时刷新数据的协程
	utils.MsgSuccess("          [MySQLWithBufferService]init successfully!") // 输出初始化成功的消息
	return ser, nil                                                          // 返回 MySQLWithBufferService 实例
}

// Add 向指定表的缓冲区添加一行数据
// table 是表名
// row 是要添加的数据行
func (b *MySQLWithBufferService) Add(table string, row []interface{}) {
	b.mu.Lock()         // 加锁以保护数据缓冲区
	defer b.mu.Unlock() // 在函数结束时解锁

	// 初始化表的缓冲区（如果尚未存在）
	if _, exists := b.data[table]; !exists {
		b.data[table] = [][]interface{}{}
	}

	// 将数据追加到缓冲区
	b.data[table] = append(b.data[table], row)
	utils.MsgSuccess("          [MySQLWithBufferService]Add successfully!") // 输出添加成功的消息
}

// flushTable 批量插入指定表的数据
// table 是表名
func (b *MySQLWithBufferService) flushTable(table string) {
	rows, exists := b.data[table] // 获取表的数据行
	if !exists {
		utils.MsgError("          [MySQLWithBufferService]flushTable not exists " + table) // 如果表不存在，输出错误消息
		return
	}
	err := b.InsertMany(table, rows) // 批量插入数据行
	if err != nil {
		utils.MsgError("          [MySQLWithBufferService]flushTable insert to mysql error " + err.Error()) // 如果插入失败，输出错误消息
		return
	}

	delete(b.data, table)                                                        // 删除已插入的数据
	utils.MsgError("          [MySQLWithBufferService]flushTable successfully!") // 输出成功消息
}

// Start 启动定时刷新数据的协程
func (b *MySQLWithBufferService) Start() {
	go func() {
		for {
			time.Sleep(b.interval) // 等待指定的间隔时间
			b.mu.Lock()            // 加锁以保护数据缓冲区
			for table := range b.data {
				b.flushTable(table) // 刷新每个表的数据
			}
			b.mu.Unlock()                                                              // 解锁
			utils.MsgError("          [MySQLWithBufferService]to mysql successfully!") // 输出成功消息
		}
	}()
}

// InsertMany 批量插入数据到指定表
// table 是表名
// rows 是要插入的数据行
// 返回一个错误，如果插入失败
func (b *MySQLWithBufferService) InsertMany(table string, rows [][]interface{}) error {
	// 构建批量插入的 SQL 语句
	insertSql := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s;",
		table,
		strings.Join(b.columns, ", "),
		valuesToString(rows),
	)

	// 使用 b.db.Exec 执行批量插入
	_, err := b.db.Exec(insertSql)
	if err != nil {
		utils.MsgError("          [MySQLWithBufferService]insert to mysql error " + err.Error()) // 如果插入失败，输出错误消息
		return fmt.Errorf("插入表 %s 失败: %w", table, err)
	}
	return nil
}
