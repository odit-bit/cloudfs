package repo

import (
	"context"
	"database/sql"

	"github.com/odit-bit/cloudfs/internal/blob"
)

var _ blob.TokenStorer = (*pgShareToken)(nil)

type pgShareToken struct {
	db *sql.DB
}

func (p *pgShareToken) migrate(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS object_tokens (
			token_key VARCHAR (255) PRIMARY KEY,
			bucket VARCHAR (255) NOT NULL,
			filename VARCHAR (255) UNIQUE NOT NULL,
			valid_until timestamp NOT NULL
		)
		;
	`
	_, err := p.db.ExecContext(ctx, query)
	return err
}

func NewPGShareToken(ctx context.Context, db *sql.DB) (*pgShareToken, error) {
	pg := pgShareToken{db: db}
	if err := pg.migrate(ctx); err != nil {
		return nil, err
	}
	return &pg, nil
}

// Delete implements blob.TokenStorer.
func (p *pgShareToken) Delete(ctx context.Context, filename string) error {
	query := `
	DELETE FROM object_tokens
	WHERE filename = $1
	;
`
	_, err := p.db.ExecContext(ctx, query, filename)
	return err
}

// Get implements blob.TokenStorer.
func (p *pgShareToken) Get(ctx context.Context, tokenKey string) (*blob.Token, bool, error) {
	query := `
		SELECT * FROM object_tokens
		WHERE token_key = $1
		;
	`

	tkn := &blob.Token{}
	err := p.db.QueryRowContext(ctx, query, tokenKey).Scan(
		&tkn.Key,
		&tkn.Bucket,
		&tkn.Filename,
		&tkn.Expire,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
	}
	return tkn, true, nil
}

// GetByBucket implements blob.TokenStorer.
func (p *pgShareToken) GetByFilename(ctx context.Context, filename string) (*blob.Token, bool, error) {
	query := `
		SELECT * FROM object_tokens
		WHERE filename = $1
		;
	`
	tkn := &blob.Token{}
	err := p.db.QueryRowContext(ctx, query, filename).Scan(
		&tkn.Key,
		&tkn.Bucket,
		&tkn.Filename,
		&tkn.Expire,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, err
	}
	return tkn, true, nil
}

// Put implements blob.TokenStorer.
func (p *pgShareToken) Put(ctx context.Context, token *blob.Token) blob.OpErr {
	query := `
		INSERT INTO object_tokens
		VALUES ($1, $2, $3, $4)
		;
	`
	_, err := p.db.ExecContext(ctx, query, token.Key, token.Bucket, token.Filename, token.Expire)
	if err != nil {
		// var pgErr *pgconn.PgError
		// if errors.As(err, &pgErr) {
		// 	if pgErr.Code == "23505" {
		// 		return &blob.RecordIsExist{Err: fmt.Errorf("unique constraint violation: %w", err)}
		// 	}
		// }
		return blob.NewException(err)
	}
	return nil
}
