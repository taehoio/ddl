// Code generated by protoc-gen-go-ddl. DO NOT EDIT.
// versions:
//  protoc-gen-go-ddl v0.0.1-alpha
//  protoc            (unknown)
// source: taehoio/ddl/services/oneonone/v1/oneonone.proto

package oneononeddlv1

import (
	"database/sql"
	"strings"

	"github.com/xissy/kubeflake"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

//go:generate mockgen -package oneononeddlv1 -self_package "github.com/taehoio/ddl/gen/go/taehoio/ddl/services/oneonone/v1" -source ./oneonone_dml_question.pb.go -destination ./oneonone_dml_question_mock.pb.go -mock_names QuestionRecorder=MockQuestionRecorder "github.com/taehoio/ddl/gen/go/taehoio/ddl/services/oneonone/v1" QuestionRecorder

const (
	questionInsertStmt = `
		INSERT INTO question (
			id, created_at, updated_at, deleted_at, question
		) VALUES (
			?, ?, ?, ?, ?
		)
	`

	questionUpdateStmt = `
		UPDATE question SET
			id = ?, created_at = ?, updated_at = ?, deleted_at = ?, question = ?
		WHERE
			id = ?
	`
)

var (
	_ = timestamppb.Timestamp{}
	_ = wrapperspb.Int32Value{}
)

type QuestionRecorder interface {
	Get(db *sql.DB, id uint64) (*Question, error)
	List(db *sql.DB, lastID *wrapperspb.UInt64Value, asc bool, limit int64) ([]*Question, error)
	Save(db *sql.DB) error
}

func (m *Question) Get(db *sql.DB, id uint64) (*Question, error) {
	stmt, err := db.Prepare("SELECT * FROM question WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var mm Question

	var createdAt sql.NullTime
	var updatedAt sql.NullTime
	var deletedAt sql.NullTime

	if err = stmt.QueryRow(id).Scan(
		&mm.Id,
		&createdAt,
		&updatedAt,
		&deletedAt,
		&mm.Question,
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

func (m *Question) List(db *sql.DB, lastID *wrapperspb.UInt64Value, asc bool, limit int64) ([]*Question, error) {
	q := "SELECT * FROM question"
	if lastID != nil {
		if asc {
			q += " WHERE id > ?"
		} else {
			q += " WHERE id < ?"
		}
	}
	if asc {
		q += " ORDER BY id ASC"
	} else {
		q += " ORDER BY id DESC"
	}
	q += " LIMIT ?"

	stmt, err := db.Prepare(q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var args []interface{}
	if lastID != nil {
		args = append(args, lastID.Value)
	}
	args = append(args, limit)

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var arr []*Question

	for rows.Next() {
		var mm Question

		var createdAt sql.NullTime
		var updatedAt sql.NullTime
		var deletedAt sql.NullTime

		if err = rows.Scan(
			&mm.Id,
			&createdAt,
			&updatedAt,
			&deletedAt,
			&mm.Question,
		); err != nil {
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

		arr = append(arr, &mm)
	}

	return arr, nil
}

func (m *Question) FindByIDs(db *sql.DB, ids []uint64) ([]*Question, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	q := "SELECT * FROM question WHERE id IN ("
	for i := range ids {
		if i > 0 {
			q += ", "
		}
		q += "?"
	}
	q += ")"

	stmt, err := db.Prepare(q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var args []interface{}
	for _, id := range ids {
		args = append(args, id)
	}

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var arr []*Question

	for rows.Next() {
		var mm Question

		var createdAt sql.NullTime
		var updatedAt sql.NullTime
		var deletedAt sql.NullTime

		if err = rows.Scan(
			&mm.Id,
			&createdAt,
			&updatedAt,
			&deletedAt,
			&mm.Question,
		); err != nil {
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

		arr = append(arr, &mm)
	}

	return arr, nil
}

func (m *Question) Save(db *sql.DB) error {
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

func (m *Question) insert(db *sql.DB) error {
	currentAt := timestamppb.Now()

	_, err := db.Exec(
		questionInsertStmt,
		m.Id,
		currentAt.AsTime(),
		currentAt.AsTime(),
		nil,
		m.Question,
	)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1062: Duplicate entry") {
			return ErrDuplicateEntry
		}
		return err
	}

	return nil
}

func (m *Question) update(db *sql.DB) error {
	currentAt := timestamppb.Now()

	_, err := db.Exec(
		questionUpdateStmt,
		m.Id,
		currentAt.AsTime(),
		currentAt.AsTime(),
		nil,
		m.Question,
		m.Id,
	)
	if err != nil {
		return err
	}

	return nil
}
