package repository

import (
	"github.com/hnnsly/library-console/internal/repository/postgres"
	"github.com/hnnsly/library-console/internal/repository/redis"
)

// LibraryRepository фасадный репозиторий для работы с БД и кешем
type LibraryRepository struct {
	postgres.Queries
	redis.Redis
}

// New создает новый LibraryRepository
func New(pg *postgres.Queries, rd *redis.Redis) *LibraryRepository {
	return &LibraryRepository{
		Queries: *pg,
		Redis:   *rd,
	}
}
