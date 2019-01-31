package persistence

// Persistence abstraction for databases
type Persistence interface {
	InsertBucket(name string) error
	UpdateDefaultBucket(chatID, bucketName string) error
	DeleteBucket(name string) error
	InsertOption(bucketName, option string) error
	ReadAllOptions(bucketName string) ([]string, error)
	DeleteOption(bucketName string, optionIdx int) error
}
