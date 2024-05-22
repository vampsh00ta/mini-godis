package repository

import (
	"sync"
)

type RedisCache struct {
	*sync.Mutex
	storage map[string]int
}

type Repository interface {
}

func New() Repository {
	m := &sync.Mutex{}
	cache := RedisCache{
		m,
		make(map[string]int),
	}
	return &cache
}
