package redis

import (
	"context"

	"github.com/Neniel/gotennis/entity"

	"github.com/go-redis/redis"
)

type RedisReader struct {
	redisClient *redis.Client
}

func NewRedisReader(client *redis.Client) *RedisReader {
	return &RedisReader{
		redisClient: client,
	}
}

func (rdbr *RedisReader) GetCategories(ctx context.Context) ([]entity.Category, error) {
	output := make([]entity.Category, 0)
	mapp, err := rdbr.redisClient.HGetAll("categories").Result()
	if err != nil {
		return nil, err
	}
	for _, v := range mapp {
		var c entity.Category
		c.UnmarshalBinary([]byte(v))
		output = append(output, c)
	}
	return output, nil
}

func (rdbr *RedisReader) GetCategory(ctx context.Context, id string) (*entity.Category, error) {
	var category entity.Category
	err := rdbr.redisClient.HGet("categories", id).Scan(&category)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

type RedisWriter struct {
	redisClient *redis.Client
}

func NewRedisWriter(client *redis.Client) *RedisWriter {
	return &RedisWriter{
		redisClient: client,
	}
}

func (rdbw *RedisWriter) AddCategory(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	cmdResult := rdbw.redisClient.HSet("categories", category.ID.Hex(), category)
	if cmdResult.Err() != nil {
		return nil, cmdResult.Err()
	}
	return category, nil
}

func (rdbw *RedisWriter) UpdateCategory(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	return nil, nil
}

func (rdbw *RedisWriter) DeleteCategory(ctx context.Context, id string) error {
	return nil
}
