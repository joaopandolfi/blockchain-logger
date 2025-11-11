package dao

import (
	"logger/remotes/postgres"

	"github.com/joaopandolfi/blackwhale/models/dao"
)

func new() dao.SQLDAO {
	return dao.Sql(postgres.Driver())
}
