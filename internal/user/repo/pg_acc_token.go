package repo

import (
	"context"
	"database/sql"

	"github.com/odit-bit/cloudfs/internal/user"
)

var _ user.TokenStorer = (*userTokenPG)(nil)

type userTokenPG struct {
	db *sql.DB
}

func NewUserTokenPG(ctx context.Context, db *sql.DB) (*userTokenPG, error) {
	pg := &userTokenPG{db: db}
	if err := pg.migrate(ctx); err != nil {
		return nil, err
	}
	return pg, nil
}

func (pg *userTokenPG) migrate(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS user_tokens (
			token_key VARCHAR (255) UNIQUE PRIMARY KEY,
			user_id VARCHAR (255) UNIQUE NOT NULL,
			valid_until TIMESTAMP NOT NULL
		)
	;
	`
	_, err := pg.db.ExecContext(ctx, query)
	return err
}

func (pg *userTokenPG) Delete(ctx context.Context, token string) error {
	query := `
		DELETE FROM user_tokens
		WHERE token_key = $1
		;
	`
	_, err := pg.db.ExecContext(ctx, query, token)
	return err
}

func (pg *userTokenPG) GetToken(ctx context.Context, token string) (*user.Token, error) {
	query := `
	SELECT * FROM user_tokens
	WHERE token_key = $1
	;
`
	tkn := &user.Token{}
	err := pg.db.QueryRowContext(ctx, query, token).Scan(
		&tkn.Key,
		&tkn.UserID,
		&tkn.Expire,
	)
	if err != nil {
		return nil, err
	}
	return tkn, err
}

func (pg *userTokenPG) GetTokenUserID(ctx context.Context, userID string) (*user.Token, bool) {
	query := `
		SELECT * FROM user_tokens
		WHERE user_id = $1
		;
	`
	tkn := user.Token{}
	err := pg.db.QueryRowContext(ctx, query, userID).Scan(
		&tkn.Key,
		&tkn.UserID,
		&tkn.Expire,
	)
	if err != nil {
		return nil, false
	}
	return &tkn, true
}

func (pg *userTokenPG) PutToken(ctx context.Context, token *user.Token) error {
	// put should idempoten
	query := `
		INSERT INTO user_tokens
		VALUES ($1, $2, $3)
		ON CONFLICT (
				token_key
			)
		DO UPDATE SET 
				valid_until = EXCLUDED.valid_until
		;
	`
	res, err := pg.db.ExecContext(ctx, query, token.Key, token.UserID, token.Expire)
	if err != nil {
		return err
	}
	_ = res
	return nil
}
