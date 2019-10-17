package redis

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
	"github.com/stanleynguyen/mindmaker/domain"
)

// BucketNameSeparator string used for concatting chatID with user defined name
// to form bucket names
const BucketNameSeparator = " - "

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
func (r *Redis) InsertBucket(chatID int64, name string) error {
	nameWChatID := getBucketNameFromChatID(chatID, name)
	return r.Client.Set(nameWChatID, "[]", 0).Err()
}

// UpdateDefaultBucket change the default bucket for chat to draw from
func (r *Redis) UpdateDefaultBucket(chatID int64, bucketName string) error {
	return r.Client.Set(strconv.Itoa(int(chatID)), bucketName, 0).Err()
}

// GetDefaultBucket get the default bucket of current chat
func (r *Redis) GetDefaultBucket(chatID int64) (string, error) {
	return r.Client.Get(strconv.Itoa(int(chatID))).Result()
}

// Exists check if a bucket exists in database
func (r *Redis) Exists(chatID int64, bucketName string) (bool, error) {
	name := getBucketNameFromChatID(chatID, bucketName)
	return r.rowExists(name)
}

// DefaultWasSet check if a default bucket set for a chat
func (r *Redis) DefaultWasSet(chatID int64) (bool, error) {
	return r.rowExists(strconv.Itoa(int(chatID)))
}

func (r *Redis) rowExists(key string) (bool, error) {
	rs, err := r.Client.Exists(key).Result()
	if err != nil {
		return false, err
	}

	if rs > 0 {
		return true, nil
	}

	return false, nil
}

// DeleteBucket remove a bucket from database
func (r *Redis) DeleteBucket(chatID int64, bucketName string) error {
	name := getBucketNameFromChatID(chatID, bucketName)
	return r.Client.Del(name).Err()
}

// GetAllBuckets get all buckets in a chat
func (r *Redis) GetAllBuckets(chatID int64) ([]domain.Bucket, error) {
	bucketNames, err := r.Client.Keys(fmt.Sprintf("%v%v*", chatID, BucketNameSeparator)).Result()
	if err != nil {
		return nil, err
	}

	buckets := []domain.Bucket{}
	for _, name := range bucketNames {
		userReadableName := strings.SplitN(name, BucketNameSeparator, 2)[1]
		b := domain.Bucket{
			ChatID:  chatID,
			Name:    userReadableName,
			Options: nil,
		}
		buckets = append(buckets, b)
	}

	return buckets, nil
}

// InsertOption put an option inside a bucket
func (r *Redis) InsertOption(chatID int64, bucketName string, option domain.Option) error {
	options, err := r.ReadAllOptions(chatID, bucketName)
	if err != nil {
		return err
	}

	options = append(options, option)
	strVal, err := json.Marshal(options)
	if err != nil {
		return err
	}

	name := getBucketNameFromChatID(chatID, bucketName)
	return r.Client.Set(name, strVal, 0).Err()
}

// ReadAllOptions query all options inside a bucket
func (r *Redis) ReadAllOptions(chatID int64, bucketName string) ([]domain.Option, error) {
	return r.getOptions(chatID, bucketName)
}

// DeleteOption remove an option from a bucket
func (r *Redis) DeleteOption(chatID int64, bucketName string, optionIdx int64) ([]domain.Option, error) {
	options, err := r.getOptions(chatID, bucketName)
	if err != nil {
		return nil, err
	}

	options = append(options[:optionIdx], options[optionIdx+1:]...)
	strVal, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}

	name := getBucketNameFromChatID(chatID, bucketName)
	err = r.Client.Set(name, strVal, 0).Err()
	if err != nil {
		return nil, err
	}

	return options, nil
}

func (r *Redis) getOptions(chatID int64, bucketName string) ([]domain.Option, error) {
	name := getBucketNameFromChatID(chatID, bucketName)
	rsStr, err := r.Client.Get(name).Result()
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

func getBucketNameFromChatID(chatID int64, userGivenName string) string {
	return strconv.Itoa(int(chatID)) + BucketNameSeparator + userGivenName
}
