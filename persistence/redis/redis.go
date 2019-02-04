package redis

import (
	"encoding/json"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/stanleynguyen/mindmaker/domain"
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
	return r.Client.Set(name, "[]", 0).Err()
}

// UpdateDefaultBucket change the default bucket for chat to draw from
func (r *Redis) UpdateDefaultBucket(chatID int64, bucketName string) error {
	return r.Client.Set(strconv.Itoa(int(chatID)), bucketName, 0).Err()
}

// GetDefaultBucket get the default bucket of current chat
func (r *Redis) GetDefaultBucket(chatID int64) (string, error) {
	return r.Client.Get(strconv.Itoa(int(chatID))).Result()
}

// DeleteBucket remove a bucket from database
func (r *Redis) DeleteBucket(bucketName string) error {
	return r.Client.Del(bucketName).Err()
}

// InsertOption put an option inside a bucket
func (r *Redis) InsertOption(bucketName string, option domain.Option) error {
	options, err := r.ReadAllOptions(bucketName)
	if err != nil {
		return err
	}

	options = append(options, option)
	strVal, err := json.Marshal(options)
	if err != nil {
		return err
	}

	return r.Client.Set(bucketName, strVal, 0).Err()
}

// ReadAllOptions query all options inside a bucket
func (r *Redis) ReadAllOptions(bucketName string) ([]domain.Option, error) {
	return r.getOptions(bucketName)
}

// DeleteOption remove an option from a bucket
func (r *Redis) DeleteOption(bucketName string, optionIdx int) ([]domain.Option, error) {
	options, err := r.getOptions(bucketName)
	if err != nil {
		return nil, err
	}

	options = append(options[:optionIdx], options[optionIdx+1:]...)
	strVal, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}

	err = r.Client.Set(bucketName, strVal, 0).Err()
	if err != nil {
		return nil, err
	}

	return options, nil
}

func (r *Redis) getOptions(bucketName string) ([]domain.Option, error) {
	rsStr, err := r.Client.Get(bucketName).Result()
	if err != nil {
		return nil, err
	}

	var options []domain.Option
	err = json.Unmarshal([]byte(rsStr), &options)
	if err != nil {
		return nil, err
	}

	return options, nil
}
