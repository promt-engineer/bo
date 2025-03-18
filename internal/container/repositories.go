package container

import (
	"backoffice/internal/config"
	"backoffice/internal/constants"
	"backoffice/internal/entities"
	"backoffice/internal/repositories/pgsql"
	"backoffice/internal/repositories/redis"
	r "backoffice/pkg/redis"

	"github.com/sarulabs/di"
	"gorm.io/gorm"
)

func BuildRepositories() []di.Def {
	return []di.Def{
		{
			Name: constants.AccountRepositoryName,
			Build: func(ctn di.Container) (interface{}, error) {
				conn := ctn.Get(constants.PgSQLConnectionName).(*gorm.DB)

				return pgsql.NewAccountRepository(conn), nil
			},
		},
		{
			Name: constants.RefreshTokenRepositoryName,
			Build: func(ctn di.Container) (interface{}, error) {
				conn := ctn.Get(constants.PgSQLConnectionName).(*gorm.DB)

				return pgsql.NewTokenRepository(conn), nil
			},
		},
		{
			Name: constants.RoleRepositoryName,
			Build: func(ctn di.Container) (interface{}, error) {
				conn := ctn.Get(constants.PgSQLConnectionName).(*gorm.DB)

				return pgsql.NewRoleRepository(conn), nil
			},
		},
		{
			Name: constants.PermissionRepositoryName,
			Build: func(ctn di.Container) (interface{}, error) {
				conn := ctn.Get(constants.PgSQLConnectionName).(*gorm.DB)

				return pgsql.NewPermissionRepository(conn), nil
			},
		},
		{
			Name: constants.GameRepositoryName,
			Build: func(ctn di.Container) (interface{}, error) {
				conn := ctn.Get(constants.PgSQLConnectionName).(*gorm.DB)

				return pgsql.NewGameRepository(conn), nil
			},
		},
		{
			Name: constants.OrganizationRepositoryName,
			Build: func(ctn di.Container) (interface{}, error) {
				conn := ctn.Get(constants.PgSQLConnectionName).(*gorm.DB)

				return pgsql.NewOrganizationRepository(conn), nil
			},
		},
		{
			Name: constants.SessionRepositoryName,
			Build: func(ctn di.Container) (interface{}, error) {
				conn := ctn.Get(constants.RedisName).(*r.Client)

				return redis.NewSessionRepository(conn), nil
			},
		},
		{
			Name: constants.FileRepositoryName,
			Build: func(ctn di.Container) (interface{}, error) {
				conn := ctn.Get(constants.RedisName).(*r.Client)

				return redis.NewFileRepository(conn), nil
			},
		},
		{
			Name: constants.CurrencyRepositoryName,
			Build: func(ctn di.Container) (interface{}, error) {
				conn := ctn.Get(constants.PgSQLConnectionName).(*gorm.DB)

				return pgsql.NewCurrencyRepository(conn), nil
			},
		},
		{
			Name: constants.WagerSetRepositoryName,
			Build: func(ctn di.Container) (interface{}, error) {
				conn := ctn.Get(constants.PgSQLConnectionName).(*gorm.DB)

				return pgsql.NewBaseRepository[entities.WagerSet](conn), nil
			},
		},
		{
			Name: constants.CurrencySetRepositoryName,
			Build: func(ctn di.Container) (interface{}, error) {
				conn := ctn.Get(constants.PgSQLConnectionName).(*gorm.DB)

				return pgsql.NewBaseRepository[entities.CurrencySet](conn), nil
			},
		},
		{
			Name: constants.DebugRepositoryName,
			Build: func(ctn di.Container) (interface{}, error) {
				conn := ctn.Get(constants.PgSQLConnectionName).(*gorm.DB)
				cfg := ctn.Get(constants.ConfigName).(*config.Config)

				return pgsql.NewDebugRepository(conn, cfg.PgSQLConfig.Name), nil
			},
		},
	}
}
