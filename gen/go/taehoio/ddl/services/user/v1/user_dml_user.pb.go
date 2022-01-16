package userddlv1

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/xissy/kubeflake"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	userInsertStmt = `
		INSERT INTO user (
			id, created_at, updated_at, deleted_at, provider, identifier, password_hash, nickname
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?
		)
	`

	userUpdateStmt = `
		UPDATE user SET
			id = ?, created_at = ?, updated_at = ?, deleted_at = ?, provider = ?, identifier = ?, password_hash = ?, nickname = ?
		WHERE id = ?
	`
)

type UserRecorder interface {
	Get(db *sql.DB, id uint64) (*User, error)
	Save(db *sql.DB) error

	FindOneByProviderAndIdentifier(
		db *sql.DB,

		provider interface{},

		identifier interface{},
	) (*User, error)
}

func (m *User) Get(db *sql.DB, id uint64) (*User, error) {
	stmt, err := db.Prepare("SELECT * FROM user WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var mm User

	var createdAt sql.NullTime

	var updatedAt sql.NullTime

	var deletedAt sql.NullTime

	var passwordHash sql.NullString

	if err = stmt.QueryRow(id).Scan(

		&mm.Id,

		&createdAt,

		&updatedAt,

		&deletedAt,

		&mm.Provider,

		&mm.Identifier,

		&passwordHash,

		&mm.Nickname,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if createdAt.Valid {

		mm.CreatedAt = timestamppb.New(createdAt.Time)

	}

	if updatedAt.Valid {

		mm.UpdatedAt = timestamppb.New(updatedAt.Time)

	}

	if deletedAt.Valid {

		mm.DeletedAt = timestamppb.New(deletedAt.Time)

	}

	if passwordHash.Valid {

		mm.PasswordHash = &wrapperspb.StringValue{Value: passwordHash.String}

	}

	return &mm, nil
}

func (m *User) FindOneByProviderAndIdentifier(
	db *sql.DB,

	provider interface{},

	identifier interface{},
) (*User, error) {
	stmt, err := db.Prepare("SELECT * FROM user WHERE provider=? AND identifier=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var mm User

	var createdAt sql.NullTime

	var updatedAt sql.NullTime

	var deletedAt sql.NullTime

	var passwordHash sql.NullString

	if err = stmt.QueryRow(provider, identifier).Scan(

		&mm.Id,

		&createdAt,

		&updatedAt,

		&deletedAt,

		&mm.Provider,

		&mm.Identifier,

		&passwordHash,

		&mm.Nickname,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if createdAt.Valid {

		mm.CreatedAt = timestamppb.New(createdAt.Time)

	}

	if updatedAt.Valid {

		mm.UpdatedAt = timestamppb.New(updatedAt.Time)

	}

	if deletedAt.Valid {

		mm.DeletedAt = timestamppb.New(deletedAt.Time)

	}

	if passwordHash.Valid {

		mm.PasswordHash = &wrapperspb.StringValue{Value: passwordHash.String}

	}

	return &mm, nil
}

func (m *User) Save(db *sql.DB) error {
	if m.Id == 0 {
		m.Id = kubeflake.Must(kubeflake.New())
	}

	shouldInsert := true
	mm, err := m.Get(db, m.Id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if mm != nil {
		shouldInsert = false
	}

	if shouldInsert {
		if err := m.insert(db); err != nil {
			return err
		}
	} else {
		if err := m.update(db); err != nil {
			return err
		}
	}

	mm, err = m.Get(db, m.Id)
	if err != nil {
		return err
	}

	proto.Merge(m, mm)

	return nil
}

func (m *User) insert(db *sql.DB) error {
	currentAt := timestamppb.Now()

	var passwordHash sql.NullString
	if m.PasswordHash != nil {
		passwordHash.Scan(m.PasswordHash.GetValue())
	}

	_, err := db.Exec(
		userInsertStmt,

		m.Id,

		currentAt.AsTime(),

		currentAt.AsTime(),

		nil,

		m.Provider,

		m.Identifier,

		passwordHash,

		m.Nickname,
	)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1062: Duplicate entry") {
			return errors.New("duplicated entry")
		}
		return err
	}

	return nil
}

func (m *User) update(db *sql.DB) error {
	currentAt := timestamppb.Now()

	var passwordHash sql.NullString
	if m.PasswordHash != nil {
		passwordHash.Scan(m.PasswordHash.GetValue())
	}

	_, err := db.Exec(
		userUpdateStmt,

		m.Id,

		currentAt.AsTime(),

		currentAt.AsTime(),

		nil,

		m.Provider,

		m.Identifier,

		passwordHash,

		m.Nickname,

		m.Id,
	)
	if err != nil {
		return err
	}

	return nil
}
