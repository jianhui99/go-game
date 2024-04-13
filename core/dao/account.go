package dao

import (
	"context"
	"core/models/entity"
	"core/repo"
)

type AccountDao struct {
	repo *repo.Manager
}

func (d *AccountDao) SaveAccount(ctx context.Context, account *entity.Account) error {
	collection := d.repo.Mongo.Db.Collection("account")
	_, err := collection.InsertOne(ctx, account)
	return err
}

func NewAccountDao(m *repo.Manager) *AccountDao {
	return &AccountDao{
		repo: m,
	}
}
