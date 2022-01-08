package userddlv1

import (
	"database/sql"

	"github.com/xissy/kubeflake"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
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

	if err = stmt.QueryRow(provider.Number(), identifier).Scan(
		&uu.Id,
		&createdAt,
		&updatedAt,
		&deletedAt,
		&uu.Provider,
		&uu.Identifier,
		&uu.PasswordHash,
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

	return &uu, nil
}

func (u *User) Save(db *sql.DB) error {
	id := u.Id
	if id == 0 {
		id = kubeflake.Must(kubeflake.New())
	}

	currentAt := timestamppb.Now()

	createdAt := u.CreatedAt
	if createdAt == nil {
		createdAt = currentAt
	}

	updatedAt := currentAt

	var passwordHash sql.NullString
	if u.PasswordHash != nil {
		passwordHash.Scan(u.PasswordHash.GetValue())
	}

	_, err := db.Exec(`INSERT INTO user (
			id, 
			created_at, 
			updated_at, 
			deleted_at, 
			provider,
			identifier,
			password_hash, 
			nickname
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE    
			updated_at=VALUES(updated_at),
			provider=VALUES(provider),
			identifier=VALUES(identifier),
			password_hash=VALUES(password_hash),
			nickname=VALUES(nickname)`,
		id,
		createdAt.AsTime(),
		updatedAt.AsTime(),
		nil,
		u.Provider.Number(),
		u.Identifier,
		passwordHash,
		u.Nickname,
	)
	if err != nil {
		return err
	}

	uu, err := u.Get(db, id)
	if err != nil {
		return err
	}

	proto.Merge(u, uu)

	return nil
}
