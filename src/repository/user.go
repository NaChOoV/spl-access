package repository

import (
	"context"
	"spl-access/ent"
	"spl-access/src/dto"
)

type UserRepository struct {
	conn *ent.Client
	ctx  *context.Context
}

func NewUserRepository(
	conn *ent.Client,
	ctx *context.Context,
) *UserRepository {
	return &UserRepository{
		conn: conn,
		ctx:  ctx,
	}
}

func (u *UserRepository) CheckUsers(users []*dto.UserDto, tx *ent.Tx) error {
	var usersToBulk = make([]*ent.UserCreate, len(users))
	for i, user := range users {
		usersToBulk[i] = u.conn.User.
			Create().
			SetFullName(user.FullName).
			SetRun(user.Run).
			SetExternalID(user.ExternalId)
	}

	if tx != nil {
		err := tx.User.
			CreateBulk(usersToBulk...).
			OnConflictColumns("external_id", "run").
			UpdateNewValues().
			Exec(*u.ctx)
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		err := u.conn.User.
			CreateBulk(usersToBulk...).
			OnConflict().
			UpdateNewValues().
			Exec(*u.ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
