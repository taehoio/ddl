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

var (
	ErrDuplicatedEntry = errors.New("duplicated entry")
)

const (
	insertStmt = `
		INSERT INTO user (
			id,
			created_at,
			updated_at,
			deleted_at,
			provider,
			identifier,
			password_hash,
			nickname
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	updateStmt = `
		UPDATE user SET
			created_at=?,
			updated_at=?,
			deleted_at=?,
			provider=?,
			identifier=?,
			password_hash=?,
			nickname=?
		WHERE id=?`
)

type UserRecorder interface {
	Get(db *sql.DB, id uint64) (*User, error)
	FindOneByProvideAndIdentifier(db *sql.DB, provider Provider, identifier string) (*User, error)
	Save(db *sql.DB) error
}

func (u *User) Get(db *sql.DB, id uint64) (*User, error) {
	stmt, err := db.Prepare("SELECT * FROM user WHERE id=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var uu User

	var createdAt sql.NullTime
	var updatedAt sql.NullTime
	var deletedAt sql.NullTime
	var passwordHash sql.NullString

	if err = stmt.QueryRow(id).Scan(
		&uu.Id,
		&createdAt,
		&updatedAt,
		&deletedAt,
		&uu.Provider,
		&uu.Identifier,
		&passwordHash,
		&uu.Nickname,
	); err != nil {
		return nil, err
	}

	if createdAt.Valid {
		uu.CreatedAt = timestamppb.New(createdAt.Time)
	}
	if updatedAt.Valid {
		uu.UpdatedAt = timestamppb.New(updatedAt.Time)
	}
	if deletedAt.Valid {
		uu.DeletedAt = timestamppb.New(deletedAt.Time)
	}
	if passwordHash.Valid {
		uu.PasswordHash = &wrapperspb.StringValue{Value: passwordHash.String}
	}

	return &uu, nil
}

func (u *User) FindOneByProvideAndIdentifier(db *sql.DB, provider Provider, identifier string) (*User, error) {
	stmt, err := db.Prepare("SELECT * FROM user WHERE provider=? AND identifier=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var uu User

	var createdAt sql.NullTime
	var updatedAt sql.NullTime
	var deletedAt sql.NullTime
	var passwordHash sql.NullString

	if err = stmt.QueryRow(provider.Number(), identifier).Scan(
		&uu.Id,
		&createdAt,
		&updatedAt,
		&deletedAt,
		&uu.Provider,
		&uu.Identifier,
		&passwordHash,
		&uu.Nickname,
	); err != nil {
		return nil, err
	}

	if createdAt.Valid {
		uu.CreatedAt = timestamppb.New(createdAt.Time)
	}
	if updatedAt.Valid {
		uu.UpdatedAt = timestamppb.New(updatedAt.Time)
	}
	if deletedAt.Valid {
		uu.DeletedAt = timestamppb.New(deletedAt.Time)
	}
	if passwordHash.Valid {
		uu.PasswordHash = &wrapperspb.StringValue{Value: passwordHash.String}
	}

	return &uu, nil
}

func (u *User) Save(db *sql.DB) error {
	if u.Id == 0 {
		u.Id = kubeflake.Must(kubeflake.New())
	}

	shouldInsert := true
	uu, err := u.Get(db, u.Id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if uu != nil {
		shouldInsert = false
	}

	if shouldInsert {
		if err := u.insert(db); err != nil {
			return err
		}
	} else {
		if err := u.update(db); err != nil {
			return err
		}
	}

	uu, err = u.Get(db, u.Id)
	if err != nil {
		return err
	}

	proto.Merge(u, uu)

	return nil
}

func (u *User) insert(db *sql.DB) error {
	currentAt := timestamppb.Now()

	var passwordHash sql.NullString
	if u.PasswordHash != nil {
		passwordHash.Scan(u.PasswordHash.GetValue())
	}

	_, err := db.Exec(
		insertStmt,
		u.Id,
		currentAt.AsTime(),
		currentAt.AsTime(),
		nil,
		u.Provider.Number(),
		u.Identifier,
		passwordHash,
		u.Nickname,
	)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1062: Duplicate entry") {
			return ErrDuplicatedEntry
		}
		return err
	}

	return nil
}

func (u *User) update(db *sql.DB) error {
	currentAt := timestamppb.Now()

	var passwordHash sql.NullString
	if u.PasswordHash != nil {
		passwordHash.Scan(u.PasswordHash.GetValue())
	}

	_, err := db.Exec(
		updateStmt,
		u.CreatedAt.AsTime(),
		currentAt.AsTime(),
		nil,
		u.Provider.Number(),
		u.Identifier,
		passwordHash,
		u.Nickname,
		u.Id,
	)
	if err != nil {
		return err
	}

	return nil
}
