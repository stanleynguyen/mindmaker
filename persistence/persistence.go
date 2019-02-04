package persistence

import "github.com/stanleynguyen/mindmaker/domain"

// Persistence abstraction for databases
type Persistence interface {
	InsertBucket(name string) error
	UpdateDefaultBucket(chatID int64, bucketName string) error
	GetDefaultBucket(chatID int64) (string, error)
	Exists(bucketName string) (bool, error)
	DefaultWasSet(chatID int64) (bool, error)
	DeleteBucket(name string) error
	GetAllBuckets(chatID int64) ([]domain.Bucket, error)
	InsertOption(bucketName string, option domain.Option) error
	ReadAllOptions(bucketName string) ([]domain.Option, error)
	DeleteOption(bucketName string, optionIdx int) ([]domain.Option, error)
}
