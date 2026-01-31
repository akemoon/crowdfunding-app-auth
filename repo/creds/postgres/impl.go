package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	"github.com/akemoon/crowdfunding-app-auth/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	constraintCredsUniqueEmail = "credentials_email_unique"
)

type CredsRepo struct {
	db *sql.DB
}

func NewCredsRepo(db *sql.DB) *CredsRepo {
	return &CredsRepo{
		db: db,
	}
}

//go:embed sql/create_creds.sql
var createCredsSQL string

func (r *CredsRepo) CreateCreds(ctx context.Context, c domain.Creds) (uuid.UUID, error) {

	var userID uuid.UUID

	err := r.db.QueryRowContext(ctx, createCredsSQL,
		c.Email,
		c.PasswordHash,
	).Scan(
		&userID,
	)
	if err != nil {
		pgErr := asPostgresError(err)
		if pgErr != nil {
			mappedErr := mapPostgresError(pgErr)
			return uuid.Nil, fmt.Errorf("%w: %s", mappedErr, pgErr.Detail)
		}
		return uuid.Nil, fmt.Errorf("%w: %s", domain.ErrInternal, err)
	}

	return userID, nil
}

//go:embed sql/delete_creds_by_user_id.sql
var deleteCredsByUserIDSQL string

func (r *CredsRepo) DeleteCredsByUserID(ctx context.Context, userID uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, deleteCredsByUserIDSQL, userID)
	if err != nil {
		return fmt.Errorf("%w: %s", domain.ErrInternal, err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %s", domain.ErrInternal, err)
	}
	if n == 0 {
		return domain.ErrCredsNotFound
	}

	return nil
}

//go:embed sql/get_creds_by_email.sql
var getCredsByEmailSQL string

func (r *CredsRepo) GetCredsByEmail(ctx context.Context, email string) (domain.Creds, error) {
	var c domain.Creds

	err := r.db.QueryRowContext(ctx, getCredsByEmailSQL, email).Scan(
		&c.UserID,
		&c.Email,
		&c.PasswordHash,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Creds{}, domain.ErrCredsNotFound
		}
		return domain.Creds{}, fmt.Errorf("%w: %s", domain.ErrInternal, err)
	}

	return c, nil
}

func asPostgresError(err error) *pgconn.PgError {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr
	}
	return nil
}

func mapPostgresError(err *pgconn.PgError) error {
	if err.Code == "23505" {
		switch err.ConstraintName {
		case constraintCredsUniqueEmail:
			return domain.ErrEmailExists
		default:
			return domain.ErrUnknownConflict
		}
	}

	return domain.ErrInternal
}
