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

type MySQLWithBufferService struct {
	db       *sql.DB
	data     map[string][][]interface{}
	mu       sync.Mutex
	columns  []string
	interval time.Duration
}

func NewMySQLWithBufferService(dsn string, interval int, columns []string) (*MySQLWithBufferService, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(30 * time.Second)
	db.SetMaxIdleConns(5) // 最大空闲连接数
	ser := &MySQLWithBufferService{
		db: db, data: make(map[string][][]interface{}),
		interval: time.Duration(interval) * time.Second, columns: columns}
	ser.Start()
	utils.MsgSuccess("          [MySQLWithBufferService]init successfully!")
	return ser, nil
}

func (b *MySQLWithBufferService) Add(table string, row []interface{}) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// 初始化表的缓冲区（如果尚未存在）
	if _, exists := b.data[table]; !exists {
		b.data[table] = [][]interface{}{}
	}

	// 将数据追加到缓冲区
	b.data[table] = append(b.data[table], row)
	utils.MsgSuccess("          [MySQLWithBufferService]Add successfully!")
}

// flushTable 批量插入指定表的数据
func (b *MySQLWithBufferService) flushTable(table string) {
	rows, exists := b.data[table]
	if !exists {
		utils.MsgError("          [MySQLWithBufferService]flushTable not exists " + table)
		return
	}
	err := b.InsertMany(table, rows)
	if err != nil {
		utils.MsgError("          [MySQLWithBufferService]flushTable insert to mysql error " + err.Error())
		return
	}

	delete(b.data, table)
	utils.MsgError("          [MySQLWithBufferService]flushTable successfully!")
}

func (b *MySQLWithBufferService) Start() {
	go func() {
		for {
			time.Sleep(b.interval)
			b.mu.Lock()
			for table := range b.data {
				b.flushTable(table)
			}
			b.mu.Unlock()
			utils.MsgError("          [MySQLWithBufferService]to mysql successfully!")
		}
	}()
}

func (b *MySQLWithBufferService) InsertMany(table string, rows [][]interface{}) error {
	insertSql := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s;",
		table,
		strings.Join(b.columns, ", "),
		valuesToString(rows),
	)
	// 收集所有的值

	// 使用 b.db.Exec 执行批量插入
	_, err := b.db.Exec(insertSql)
	if err != nil {
		return fmt.Errorf("failed to insert into table %s: %w", table, err)
	}
	return nil
}
