package userv1

import (
	"database/sql"

	"github.com/xissy/kubeflake"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserRecorder interface {
	Get(db *sql.DB, id uint64) (*User, error)
	FindOneByEmail(db *sql.DB, email string) (*User, error)
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

	if err = stmt.QueryRow(id).Scan(
		&uu.Id,
		&createdAt,
		&updatedAt,
		&deletedAt,
		&uu.PasswordHash,
		&uu.FullName,
		&uu.Email,
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

func (u *User) FindOneByEmail(db *sql.DB, email string) (*User, error) {
	stmt, err := db.Prepare("SELECT * FROM user WHERE email=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var uu User

	var createdAt sql.NullTime
	var updatedAt sql.NullTime
	var deletedAt sql.NullTime

	if err = stmt.QueryRow(email).Scan(
		&uu.Id,
		&createdAt,
		&updatedAt,
		&deletedAt,
		&uu.PasswordHash,
		&uu.FullName,
		&uu.Email,
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

	_, err := db.Exec(`INSERT INTO user (
			id, 
			created_at, 
			updated_at, 
			deleted_at, 
			password_hash, 
			full_name, 
			email
		) VALUES (?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE    
			updated_at=VALUES(updated_at),
			password_hash=VALUES(password_hash),
			full_name=VALUES(full_name),
			email=VALUES(email)`,
		id,
		createdAt.AsTime(),
		updatedAt.AsTime(),
		nil,
		u.PasswordHash,
		u.FullName,
		u.Email,
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
