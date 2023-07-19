package secret

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	querySecret "gitlab.com/mildd/flow-gateway/internal/app/query/secret"
	"gitlab.com/mildd/flow-gateway/internal/domain/secret"
)

type PGRepository struct {
	db      *sqlx.DB
	builder goqu.DialectWrapper
}

func NewPGRepository(db *sqlx.DB) PGRepository {
	return PGRepository{
		db:      db,
		builder: goqu.Dialect("postgres"),
	}
}

func (r PGRepository) GetSecretQueryModel(_ context.Context, id string) (*querySecret.Secret, error) {
	var qUser querySecret.Secret
	err := r.db.Get(
		&qUser, `SELECT json_agg(ap.codename) 
		FROM user_grouos ug
		JOIN auth_groups_permissions agp ON ug.group_id = agp.group_id
		JOIN auth_permission ap ON agp.permission_id = ap.id
		WHERE ug.user_id = 2
		GROUP BY ug.id`,
		id,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, secret.ErrorNotFound
		}
		return nil, err
	}
	return &qUser, nil
}

func (r PGRepository) GetSecretById(_ context.Context, id string) (*secret.Secret, error) {
	secretObj := secret.NewEmptySecret()
	err := r.db.Get(
		secretObj,
		`SELECT json_agg(ap.codename) 
		FROM user_grouos ug
		JOIN auth_groups_permissions agp ON ug.group_id = agp.group_id
		JOIN auth_permission ap ON agp.permission_id = ap.id
		WHERE ug.user_id = 2
		GROUP BY ug.id`,
		id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, secret.ErrorNotFound
		}
		return nil, err
	}
	return secretObj, nil
}
