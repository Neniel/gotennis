package redis

import (
	"context"

	"github.com/Neniel/gotennis/lib/entity"
	"github.com/Neniel/gotennis/lib/util"

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
	mapp, err := rdbr.redisClient.HGetAll(util.CollNameCategories).Result()
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
	err := rdbr.redisClient.HGet(util.CollNameCategories, id).Scan(&category)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (rdbr *RedisReader) GetPlayers(ctx context.Context) ([]entity.Player, error) {
	output := make([]entity.Player, 0)
	mapp, err := rdbr.redisClient.HGetAll(util.CollNamePlayers).Result()
	if err != nil {
		return nil, err
	}
	for _, v := range mapp {
		var p entity.Player
		p.UnmarshalBinary([]byte(v))
		output = append(output, p)
	}
	return output, nil
}

func (rdbr *RedisReader) GetPlayer(ctx context.Context, id string) (*entity.Player, error) {
	var player entity.Player
	err := rdbr.redisClient.HGet(util.CollNamePlayers, id).Scan(&player)
	if err != nil {
		return nil, err
	}
	return &player, nil
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
	cmdResult := rdbw.redisClient.HSet(util.CollNameCategories, category.ID.Hex(), category)
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

func (rdbw *RedisWriter) AddPlayer(ctx context.Context, player *entity.Player) (*entity.Player, error) {
	cmdResult := rdbw.redisClient.HSet(util.CollNameCategories, player.ID.Hex(), player)
	if cmdResult.Err() != nil {
		return nil, cmdResult.Err()
	}
	return player, nil
}

func (rdbw *RedisWriter) UpdatePlayer(ctx context.Context, category *entity.Player) (*entity.Player, error) {
	return nil, nil
}

func (rdbw *RedisWriter) DeletePlayer(ctx context.Context, id string) error {
	return nil
}
