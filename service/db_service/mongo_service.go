package dbservice

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// MongoDBClient 是MongoDB操作类
type MongoDBClient struct {
	client *mongo.Client
	db     *mongo.Database
}

// NewMongoDBClient 创建MongoDBClient实例
func NewMongoDBClient(URI string, DB string) (*MongoDBClient, error) {
	// 设置连接超时
	clientOptions := options.Client().ApplyURI(URI).SetMaxPoolSize(10).SetMinPoolSize(5).SetMaxConnIdleTime(10 * time.Minute)

	// 使用 mongo.Connect() 来创建新的 MongoDB 客户端
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	// 尝试连接
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	// 获取指定数据库
	db := client.Database(DB)
	return &MongoDBClient{client: client, db: db}, nil
}

// Close 关闭MongoDB连接
func (mongoDb *MongoDBClient) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return mongoDb.client.Disconnect(ctx)
}

// InsertOne 插入一条数据
func (mongoDb *MongoDBClient) InsertOne(collection string, document interface{}) (*mongo.InsertOneResult, error) {
	coll := mongoDb.db.Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := coll.InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// FindOne 查找一条数据
func (mongoDb *MongoDBClient) FindOne(collection string, filter interface{}) (bson.M, error) {
	coll := mongoDb.db.Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var re bson.M
	result := coll.FindOne(ctx, filter)
	err := result.Decode(&re)
	if err != nil {
		return nil, err
	}
	if result.Err() != nil {
		return nil, result.Err()
	}
	return re, nil
}

// UpdateOne 更新一条数据
func (mongoDb *MongoDBClient) UpdateOne(collection string, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	coll := mongoDb.db.Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteOne 删除一条数据
func (mongoDb *MongoDBClient) DeleteOne(collection string, filter interface{}) (*mongo.DeleteResult, error) {
	coll := mongoDb.db.Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// FindAll 查找所有数据
func (mongoDb *MongoDBClient) FindAll(collection string, filter interface{}) ([]bson.M, error) {
	coll := mongoDb.db.Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// CreateCollection 创建新集合
func (mongoDb *MongoDBClient) CreateCollection(collectionName string) error {
	// 检查集合是否已存在
	collections, err := mongoDb.db.ListCollectionNames(context.Background(), bson.M{})
	if err != nil {
		return err
	}

	for _, col := range collections {
		if col == collectionName {
			return fmt.Errorf("collection %s already exists", collectionName)
		}
	}

	// 创建新集合
	err = mongoDb.db.CreateCollection(context.Background(), collectionName)
	if err != nil {
		return err
	}

	//fmt.Println("Collection created:", collectionName)
	return nil
}

func (mongoDb *MongoDBClient) DropCollection(collectionName string) error {
	// 获取集合对象
	collection := mongoDb.db.Collection(collectionName)

	// 删除集合
	err := collection.Drop(context.Background())
	if err != nil {
		return fmt.Errorf("failed to drop collection: %v", err)
	}

	//fmt.Println("Collection dropped:", collectionName)
	return nil
}
