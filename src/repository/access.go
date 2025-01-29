package repository

import (
	"context"
	"os"
	"spl-access/ent"
	entAccess "spl-access/ent/access"
	"spl-access/src/dto"
	"spl-access/src/model"
)

type AccessRepository struct {
	conn *ent.Client
	ctx  *context.Context

	userRepository *UserRepository
}

func NewAccessRepository(
	conn *ent.Client,
	ctx *context.Context,
	userRepository *UserRepository,
) *AccessRepository {
	return &AccessRepository{
		conn:           conn,
		ctx:            ctx,
		userRepository: userRepository,
	}
}

func (a *AccessRepository) UpdateOrCreateAccess(access dto.AccessArrayDto) error {

	// check users (create non existing)
	var userMap = make(map[string]dto.UserDto)
	for _, access := range access.Data {
		if _, exists := userMap[access.Run]; !exists {
			userMap[access.Run] = dto.UserDto{
				Run:        access.Run,
				ExternalId: access.ExternalId,
				FullName:   access.FullName,
			}
		}
	}

	var users []*dto.UserDto
	for _, user := range userMap {
		users = append(users, &user)
	}

	tx, err := a.conn.Tx(*a.ctx)
	if err != nil {
		return err
	}

	err = a.userRepository.CheckUsers(users, tx)
	if err != nil {
		return err
	}

	var accessToBulk = make([]*ent.AccessCreate, len(access.Data))
	for i, access := range access.Data {
		accessCreate := a.conn.Access.
			Create().
			SetRun(access.Run).
			SetEntryAt(access.EntryAt).
			SetLocation(entAccess.Location(access.Location))

		if access.ExitAt.IsZero() {
			accessCreate.SetNillableExitAt(nil)
		} else {
			accessCreate.SetExitAt(access.ExitAt)
		}
		accessToBulk[i] = accessCreate

	}

	err = tx.Access.
		CreateBulk(accessToBulk...).
		OnConflictColumns("run", "location", "entry_at").
		UpdateNewValues().
		Exec(*a.ctx)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (a *AccessRepository) GetAccess() (*[]model.Access, error) {
	sqlFile, err := os.ReadFile("src/repository/sql/access_query.sql")
	if err != nil {
		return nil, err
	}
	query, err := a.conn.QueryContext(*a.ctx, string(sqlFile))
	if err != nil {
		return nil, err
	}
	defer query.Close()

	var accesses []model.Access
	for query.Next() {
		var access model.Access
		if err := query.Scan(
			&access.Run,
			&access.FullName,
			&access.EntryAt,
			&access.Location,
		); err != nil {
			return nil, err
		}
		accesses = append(accesses, access)
	}

	if err := query.Err(); err != nil {
		return nil, err
	}

	return &accesses, nil
}
