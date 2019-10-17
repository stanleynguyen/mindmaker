package persistence

import "github.com/stanleynguyen/mindmaker/domain"

// Persistence abstraction for databases
type Persistence interface {
	InsertBucket(chatID int64, name string) error
	UpdateDefaultBucket(chatID int64, bucketName string) error
	GetDefaultBucket(chatID int64) (string, error)
	Exists(chatID int64, bucketName string) (bool, error)
	DefaultWasSet(chatID int64) (bool, error)
	DeleteBucket(chatID int64, name string) error
	GetAllBuckets(chatID int64) ([]domain.Bucket, error)
	InsertOption(chatID int64, bucketName string, option domain.Option) error
	ReadAllOptions(chatID int64, bucketName string) ([]domain.Option, error)
	DeleteOption(chatID int64, bucketName string, optionIdx int64) ([]domain.Option, error)
}
