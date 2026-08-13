package app

import (
	"github.com/mythenox/binge-tracker/internal/database"
)

type State struct {
	DB  *database.Queries
	Cfg *Config
}
