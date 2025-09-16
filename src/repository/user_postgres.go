package repository

import (
	"context"
	"spl-access/ent"
	"spl-access/src/dto"
)

type UserPostgres struct {
	conn *ent.Client
}

func NewUserPostgres(
	conn *ent.Client,
	ctx *context.Context,
) *UserPostgres {
	return &UserPostgres{
		conn: conn,
	}
}

func (u *UserPostgres) CheckUsers(ctx context.Context, users []*dto.UserDto, tx ...any) error {
	var transaction *ent.Tx
	if len(tx) > 0 && tx[0] != nil {
		if entTx, ok := tx[0].(*ent.Tx); ok {
			transaction = entTx
		}
	}

	var usersToBulk = make([]*ent.UserCreate, len(users))
	for i, user := range users {
		usersToBulk[i] = u.conn.User.
			Create().
			SetFullName(user.FullName).
			SetRun(user.Run).
			SetExternalID(user.ExternalId)
	}

	if transaction != nil {
		err := transaction.User.
			CreateBulk(usersToBulk...).
			OnConflictColumns("external_id").
			UpdateNewValues().
			Exec(ctx)
		if err != nil {
			transaction.Rollback()
			return err
		}
	} else {
		err := u.conn.User.
			CreateBulk(usersToBulk...).
			OnConflictColumns("external_id").
			UpdateNewValues().
			Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
