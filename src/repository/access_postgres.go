package repository

import (
	"context"
	"os"
	"spl-access/ent"
	entAccess "spl-access/ent/access"
	"spl-access/src/dto"
	"spl-access/src/model"
)

type PostgresAccess struct {
	conn *ent.Client

	userRepository UserRepository
}

func NewPostgresAccess(
	conn *ent.Client,
	ctx *context.Context,
	userRepository UserRepository,
) *PostgresAccess {
	return &PostgresAccess{
		conn:           conn,
		userRepository: userRepository,
	}
}

func (a *PostgresAccess) UpdateOrCreateAccess(ctx context.Context, access []*dto.AccessDto) error {

	// check users (create non existing)
	var userMap = make(map[string]dto.UserDto)
	for _, accessItem := range access {
		if _, exists := userMap[accessItem.ExternalId]; !exists {
			userMap[accessItem.ExternalId] = dto.UserDto{
				Run:        accessItem.Run,
				ExternalId: accessItem.ExternalId,
				FullName:   accessItem.FullName,
			}
		}
	}

	var users []*dto.UserDto
	for _, user := range userMap {
		users = append(users, &user)
	}

	tx, err := a.conn.Tx(ctx)
	if err != nil {
		return err
	}

	err = a.userRepository.CheckUsers(ctx, users, tx)
	if err != nil {
		return err
	}

	var accessToBulk = make([]*ent.AccessCreate, len(access))
	for i, access := range access {
		accessCreate := a.conn.Access.
			Create().
			SetRun(access.Run).
			SetEntryAt(access.EntryAt).
			SetNillableExitAt(access.ExitAt).
			SetLocation(entAccess.Location(access.Location))

		accessToBulk[i] = accessCreate

	}

	err = tx.Access.
		CreateBulk(accessToBulk...).
		OnConflictColumns("run", "location", "entry_at").
		UpdateNewValues().
		Exec(ctx)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (a *PostgresAccess) GetAccess(ctx context.Context, complete bool) ([]*model.Access, error) {
	filePath := "src/repository/sql/get_access_query.sql"
	if complete {
		filePath = "src/repository/sql/get_access_query_complete.sql"
	}

	sqlFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	query, err := a.conn.QueryContext(ctx, string(sqlFile))
	if err != nil {
		return nil, err
	}
	defer query.Close()

	var accesses []*model.Access
	for query.Next() {
		var access model.Access
		if err := query.Scan(
			&access.ExternalId,
			&access.Run,
			&access.FullName,
			&access.EntryAt,
			&access.ExitAt,
			&access.Location,
		); err != nil {
			return nil, err
		}
		accesses = append(accesses, &access)
	}

	if err := query.Err(); err != nil {
		return nil, err
	}

	return accesses, nil
}
