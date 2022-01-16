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
	userRoleInsertStmt = `
		INSERT INTO user_role (
			id, created_at, updated_at, deleted_at, user_id, role
		) VALUES (
			?, ?, ?, ?, ?, ?
		)
	`

	userRoleUpdateStmt = `
		UPDATE user_role SET
			id = ?, created_at = ?, updated_at = ?, deleted_at = ?, user_id = ?, role = ?
		WHERE
			id = ?
	`
)

var (
	_ = timestamppb.Timestamp{}
	_ = wrapperspb.Int32Value{}
)

type UserRoleRecorder interface {
	Get(db *sql.DB, id uint64) (*UserRole, error)
	Save(db *sql.DB) error
	FindOneByUserId(db *sql.DB, user_id interface{}) (*UserRole, error)
}

func (m *UserRole) Get(db *sql.DB, id uint64) (*UserRole, error) {
	stmt, err := db.Prepare("SELECT * FROM user_role WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var mm UserRole

	var createdAt sql.NullTime
	var updatedAt sql.NullTime
	var deletedAt sql.NullTime

	if err = stmt.QueryRow(id).Scan(
		&mm.Id,
		&createdAt,
		&updatedAt,
		&deletedAt,
		&mm.UserId,
		&mm.Role,
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

	return &mm, nil
}

func (m *UserRole) FindOneByUserId(db *sql.DB, user_id interface{}) (*UserRole, error) {
	stmt, err := db.Prepare("SELECT * FROM user_role WHERE user_id=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var mm UserRole

	var createdAt sql.NullTime
	var updatedAt sql.NullTime
	var deletedAt sql.NullTime

	if err = stmt.QueryRow(user_id).Scan(
		&mm.Id,
		&createdAt,
		&updatedAt,
		&deletedAt,
		&mm.UserId,
		&mm.Role,
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

	return &mm, nil
}

func (m *UserRole) Save(db *sql.DB) error {
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

func (m *UserRole) insert(db *sql.DB) error {
	currentAt := timestamppb.Now()

	_, err := db.Exec(
		userRoleInsertStmt,
		m.Id,
		currentAt.AsTime(),
		currentAt.AsTime(),
		nil,
		m.UserId,
		m.Role,
	)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1062: Duplicate entry") {
			return errors.New("duplicated entry")
		}
		return err
	}

	return nil
}

func (m *UserRole) update(db *sql.DB) error {
	currentAt := timestamppb.Now()

	_, err := db.Exec(
		userRoleUpdateStmt,
		m.Id,
		currentAt.AsTime(),
		currentAt.AsTime(),
		nil,
		m.UserId,
		m.Role,
		m.Id,
	)
	if err != nil {
		return err
	}

	return nil
}