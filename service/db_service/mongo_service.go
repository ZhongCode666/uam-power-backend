package dbservice

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// MongoDBClient 是 MongoDB 操作类
type MongoDBClient struct {
	// client 是 MongoDB 客户端
	client *mongo.Client
	// db 是 MongoDB 数据库
	db *mongo.Database
}

// NewMongoDBClient 创建 MongoDBClient 实例
// URI 是 MongoDB 的连接 URI
// DB 是要连接的数据库名称
func NewMongoDBClient(URI string, DB string) (*MongoDBClient, error) {
	// 设置连接选项，包括最大连接池大小、最小连接池大小和最大空闲时间
	clientOptions := options.Client().ApplyURI(URI).SetMaxPoolSize(10).SetMinPoolSize(5).SetMaxConnIdleTime(10 * time.Minute)

	// 使用 mongo.Connect() 来创建新的 MongoDB 客户端
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	// 尝试连接到 MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	// 获取指定的数据库
	db := client.Database(DB)
	return &MongoDBClient{client: client, db: db}, nil
}

// Close 关闭 MongoDB 连接
func (mongoDb *MongoDBClient) Close() error {
	// 设置上下文和超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // 函数结束时取消上下文

	// 断开 MongoDB 客户端连接
	return mongoDb.client.Disconnect(ctx)
}

// InsertOne 插入一条数据
// collection 是集合名称
// document 是要插入的文档
// 返回插入结果和可能的错误
func (mongoDb *MongoDBClient) InsertOne(collection string, document interface{}) (*mongo.InsertOneResult, error) {
	coll := mongoDb.db.Collection(collection)                               // 获取集合对象
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 设置上下文和超时时间
	defer cancel()                                                          // 函数结束时取消上下文

	result, err := coll.InsertOne(ctx, document) // 插入文档
	if err != nil {
		return nil, err // 如果插入失败，返回错误
	}
	return result, nil // 返回插入结果
}

// FindOne 查找一条数据
// collection 是集合名称
// filter 是查询过滤器
// 返回查找到的文档和可能的错误
func (mongoDb *MongoDBClient) FindOne(collection string, filter interface{}) (bson.M, error) {
	coll := mongoDb.db.Collection(collection)                               // 获取集合对象
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 设置上下文和超时时间
	defer cancel()                                                          // 函数结束时取消上下文
	var re bson.M
	result := coll.FindOne(ctx, filter) // 查找一条数据
	err := result.Decode(&re)           // 解码查询结果
	if err != nil {
		return nil, err // 如果解码失败，返回错误
	}
	if result.Err() != nil {
		return nil, result.Err() // 如果查询结果有错误，返回错误
	}
	return re, nil // 返回查询结果
}

// FindOneWithDropRow 查找一条数据并排除指定字段
// collection 是集合名称
// filter 是查询过滤器
// dropRows 是要排除的字段
// 返回查找到的文档和可能的错误
func (mongoDb *MongoDBClient) FindOneWithDropRow(collection string, filter interface{}, dropRows interface{}) (bson.M, error) {
	coll := mongoDb.db.Collection(collection)                               // 获取集合对象
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 设置上下文和超时时间
	defer cancel()                                                          // 函数结束时取消上下文
	var re bson.M
	result := coll.FindOne(ctx, filter, options.FindOne().SetProjection(dropRows)) // 查找一条数据并排除指定字段
	err := result.Decode(&re)                                                      // 解码查询结果
	if err != nil {
		return nil, err // 如果解码失败，返回错误
	}
	if result.Err() != nil {
		return nil, result.Err() // 如果查询结果有错误，返回错误
	}
	return re, nil // 返回查询结果
}

// UpdateOne 更新一条数据
// collection 是集合名称
// filter 是查询过滤器
// update 是更新操作
// 返回更新结果和可能的错误
func (mongoDb *MongoDBClient) UpdateOne(collection string, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	coll := mongoDb.db.Collection(collection)                               // 获取集合对象
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 设置上下文和超时时间
	defer cancel()                                                          // 函数结束时取消上下文

	result, err := coll.UpdateOne(ctx, filter, update) // 更新一条数据
	if err != nil {
		return nil, err // 如果更新失败，返回错误
	}
	return result, nil // 返回更新结果
}

// Update 更新多条数据
// collection 是集合名称
// filter 是查询过滤器
// update 是更新操作
// 返回更新结果和可能的错误
func (mongoDb *MongoDBClient) Update(collection string, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	coll := mongoDb.db.Collection(collection)                               // 获取集合对象
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 设置上下文和超时时间
	defer cancel()                                                          // 函数结束时取消上下文

	result, err := coll.UpdateMany(ctx, filter, update) // 更新多条数据
	if err != nil {
		return nil, err // 如果更新失败，返回错误
	}
	return result, nil // 返回更新结果
}

// DeleteOne 删除一条数据
// collection 是集合名称
// filter 是查询过滤器
// 返回删除结果和可能的错误
func (mongoDb *MongoDBClient) DeleteOne(collection string, filter interface{}) (*mongo.DeleteResult, error) {
	coll := mongoDb.db.Collection(collection)                               // 获取集合对象
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 设置上下文和超时时间
	defer cancel()                                                          // 函数结束时取消上下文

	result, err := coll.DeleteOne(ctx, filter) // 删除一条数据
	if err != nil {
		return nil, err // 如果删除失败，返回错误
	}
	return result, nil // 返回删除结果
}

// Delete 删除多条数据
// collection 是集合名称
// filter 是查询过滤器
// 返回删除结果和可能的错误
func (mongoDb *MongoDBClient) Delete(collection string, filter interface{}) (*mongo.DeleteResult, error) {
	coll := mongoDb.db.Collection(collection)                               // 获取集合对象
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 设置上下文和超时时间
	defer cancel()                                                          // 函数结束时取消上下文

	result, err := coll.DeleteMany(ctx, filter) // 删除多条数据
	if err != nil {
		return nil, err // 如果删除失败，返回错误
	}
	return result, nil // 返回删除结果
}

// FindAll 查找所有数据
// collection 是集合名称
// filter 是查询过滤器
// 返回查找到的文档切片和可能的错误
func (mongoDb *MongoDBClient) FindAll(collection string, filter interface{}) ([]bson.M, error) {
	coll := mongoDb.db.Collection(collection)                               // 获取集合对象
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 设置上下文和超时时间
	defer cancel()                                                          // 函数结束时取消上下文

	cursor, err := coll.Find(ctx, filter) // 查找所有符合过滤器的数据
	if err != nil {
		return nil, err // 如果查找失败，返回错误
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx) // 关闭游标
		if err != nil {

		}
	}(cursor, ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err // 如果解码失败，返回错误
	}
	return results, nil // 返回查找到的文档切片
}

// FindAllWithDrops 查找所有数据并排除指定字段
// collection 是集合名称
// filter 是查询过滤器
// drops 是要排除的字段
// 返回查找到的文档切片和可能的错误
func (mongoDb *MongoDBClient) FindAllWithDrops(collection string, filter interface{}, drops interface{}) ([]bson.M, error) {
	coll := mongoDb.db.Collection(collection)                               // 获取集合对象
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 设置上下文和超时时间
	defer cancel()                                                          // 函数结束时取消上下文

	cursor, err := coll.Find(ctx, filter, options.Find().SetProjection(drops)) // 查找所有符合过滤器的数据并排除指定字段
	if err != nil {
		return nil, err // 如果查找失败，返回错误
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx) // 关闭游标
		if err != nil {

		}
	}(cursor, ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err // 如果解码失败，返回错误
	}
	return results, nil // 返回查找到的文档切片
}

// CreateCollection 创建新集合
// collectionName 是要创建的集合名称
// 返回可能的错误
func (mongoDb *MongoDBClient) CreateCollection(collectionName string) error {
	// 检查集合是否已存在
	collections, err := mongoDb.db.ListCollectionNames(context.Background(), bson.M{})
	if err != nil {
		return err // 如果获取集合名称列表失败，返回错误
	}

	// 遍历集合名称列表，检查集合是否已存在
	for _, col := range collections {
		if col == collectionName {
			return fmt.Errorf("collection %s already exists", collectionName) // 如果集合已存在，返回错误
		}
	}

	// 创建新集合
	err = mongoDb.db.CreateCollection(context.Background(), collectionName)
	if err != nil {
		return err // 如果创建集合失败，返回错误
	}

	//fmt.Println("Collection created:", collectionName)
	return nil // 成功创建集合，返回 nil
}

// DropCollection 删除集合
// collectionName 是要删除的集合名称
// 返回可能的错误
func (mongoDb *MongoDBClient) DropCollection(collectionName string) error {
	// 获取集合对象
	collection := mongoDb.db.Collection(collectionName)

	// 删除集合
	err := collection.Drop(context.Background())
	if err != nil {
		return fmt.Errorf("删除集合失败: %v", err)
	}

	//fmt.Println("Collection dropped:", collectionName)
	return nil
}
