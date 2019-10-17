package postgres

import (
	"database/sql"
	"errors"
	"log"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	pg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // to allow migrations from files
	_ "github.com/lib/pq"                                // init for postgresql db
	"github.com/stanleynguyen/mindmaker/domain"
)

// Postgres abstraction for Postgresql database
type Postgres struct {
	Conn *sql.DB
}

// NewInstance create new connection to Postgres DB
// this initialization also attempts to migrate the DB to latest version
func NewInstance(connectionString string) (*Postgres, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	driver, err := pg.WithInstance(db, &pg.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance("file://./migrations", "postgres", driver)
	if err != nil {
		return nil, err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	return &Postgres{db}, nil
}

// InsertBucket create a new options bucket
func (p *Postgres) InsertBucket(chatID int64, name string) error {
	strChatID := strconv.FormatInt(chatID, 10)
	chatExists, err := p.chatExists(strChatID)
	if err != nil {
		logDBErr(err)
		return err
	}
	if !chatExists {
		err := p.insertChat(strChatID)
		if err != nil {
			logDBErr(err)
			return err
		}
	}
	query := `INSERT INTO buckets(id, chat_id) VALUES ($1, $2) RETURNING id`
	_, err = p.Conn.Exec(query, name, strChatID)
	return err
}

// UpdateDefaultBucket change the default bucket for chat to draw from
func (p *Postgres) UpdateDefaultBucket(chatID int64, bucketName string) error {
	query := `UPDATE chats SET default_bucket = $1 WHERE id = $2`
	stmt, err := p.Conn.Prepare(query)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(bucketName, strconv.FormatInt(chatID, 10))
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

// GetDefaultBucket get the default bucket of current chat
func (p *Postgres) GetDefaultBucket(chatID int64) (string, error) {
	query := `SELECT default_bucket FROM chats WHERE id = $1`
	var bucketName string
	err := p.Conn.QueryRow(query, strconv.FormatInt(chatID, 10)).Scan(&bucketName)
	if err != nil {
		return "", err
	} else if bucketName == "" {
		return "", errors.New("No default bucket found")
	}

	return bucketName, nil
}

// Exists check if a bucket exists in database
func (p *Postgres) Exists(chatID int64, bucketName string) (bool, error) {
	query := `SELECT * FROM buckets WHERE chat_id = $1 AND id = $2`
	rows, err := p.Conn.Query(query, strconv.FormatInt(chatID, 10), bucketName)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

// DefaultWasSet check if a default bucket set for a chat
func (p *Postgres) DefaultWasSet(chatID int64) (bool, error) {
	query := `SELECT default_bucket FROM chats WHERE id = $1`
	var bucketName string
	err := p.Conn.QueryRow(query, strconv.FormatInt(chatID, 10)).Scan(&bucketName)
	if err != nil {
		return false, err
	}

	return bucketName != "", nil
}

// DeleteBucket remove a bucket from database
func (p *Postgres) DeleteBucket(chatID int64, bucketName string) error {
	query := `DELETE FROM buckets WHERE chat_id = $1 AND id = $2`
	_, err := p.Conn.Exec(query, strconv.FormatInt(chatID, 10), bucketName)
	return err
}

// GetAllBuckets get all buckets in a chat
func (p *Postgres) GetAllBuckets(chatID int64) ([]domain.Bucket, error) {
	query := `SELECT id FROM buckets WHERE chat_id = $1`
	rows, err := p.Conn.Query(query, strconv.FormatInt(chatID, 10))
	if err != nil {
		return nil, err
	}

	results := make([]domain.Bucket, 0)
	for rows.Next() {
		b := domain.Bucket{}
		err = rows.Scan(&b.Name)
		if err != nil {
			return nil, err
		}
		results = append(results, b)
	}

	return results, nil
}

// InsertOption put an option inside a bucket
func (p *Postgres) InsertOption(chatID int64, bucketName string, option domain.Option) error {
	query := `INSERT INTO options (chat_id, bucket_id, content) VALUES ($1, $2, $3)`
	_, err := p.Conn.Exec(query, strconv.FormatInt(chatID, 10), bucketName, option)
	return err
}

// ReadAllOptions query all options inside a bucket
func (p *Postgres) ReadAllOptions(chatID int64, bucketName string) ([]domain.Option, error) {
	query := `SELECT content FROM options WHERE chat_id = $1 AND bucket_id = $2`
	rows, err := p.Conn.Query(query, strconv.FormatInt(chatID, 10), bucketName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]domain.Option, 0)
	for rows.Next() {
		o := domain.Option("")
		err = rows.Scan(&o)
		if err != nil {
			return nil, err
		}
		results = append(results, o)
	}

	return results, nil
}

// DeleteOption remove an option from a bucket
func (p *Postgres) DeleteOption(chatID int64, bucketName string, optionIdx int64) ([]domain.Option, error) {
	query := `SELECT id, content FROM options WHERE chat_id = $1 AND bucket_id = $2`
	rows, err := p.Conn.Query(query, strconv.FormatInt(chatID, 10), bucketName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	IDs := make([]int64, 0)
	options := make([]domain.Option, 0)
	for rows.Next() {
		idx := int64(0)
		o := domain.Option("")
		err = rows.Scan(&idx, &o)
		if err != nil {
			return nil, err
		}
		IDs = append(IDs, idx)
		options = append(options, o)
	}

	optionID := IDs[optionIdx]
	query = `DELETE FROM options WHERE id = $1`
	_, err = p.Conn.Exec(query, optionID)
	if err != nil {
		return nil, err
	}

	return append(options[:optionIdx], options[optionIdx+1:]...), nil
}

func (p *Postgres) insertChat(id string) error {
	query := `INSERT INTO chats(id) VALUES ($1)`
	_, err := p.Conn.Exec(query, id)
	return err
}

func (p *Postgres) chatExists(id string) (bool, error) {
	rows, err := p.Conn.Query(`SELECT * FROM chats WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		return true, nil
	}

	return false, nil
}

func logDBErr(err error) {
	log.Printf("DB error: %s", err.Error())
}
