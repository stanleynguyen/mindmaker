package redis

import (
	"github.com/go-redis/redis"
)

// Redis abstraction for redis database
type Redis struct {
	Client *redis.Client
}

// NewInstance get new instance of Redis database
func NewInstance(connectionString string) (*Redis, error) {
	opt, err := redis.ParseURL(connectionString)
	if err != nil {
		return nil, err
	}

	return &Redis{redis.NewClient(opt)}, nil
}

// InsertBucket create a new options bucket
func (r *Redis) InsertBucket(name string) error {
	return nil
}

// UpdateDefaultBucket change the default bucket for chat to draw from
func (r *Redis) UpdateDefaultBucket(chatID, bucketName string) error {
	return nil
}

// DeleteBucket remove a bucket from database
func (r *Redis) DeleteBucket(bucketName string) error {
	return nil
}

// InsertOption put an option inside a bucket
func (r *Redis) InsertOption(bucketName, option string) error {
	return nil
}

// ReadAllOptions query all options inside a bucket
func (r *Redis) ReadAllOptions(bucketName string) ([]string, error) {
	return []string{}, nil
}

// DeleteOption remove an option from a bucket
func (r *Redis) DeleteOption(bucketName string, optionIdx int) error {
	return nil
}
