package service

import (
	gmModel "dnf/biz/gm/model"
	logModel "dnf/biz/log/model"
	"dnf/mods/game_db"
	"errors"
	"github.com/localhostjason/webserver/server/util/uv"
)

func GetReChargeLogList(q *RechargeQ, pi *uv.PagingIn, order *uv.Order) ([]logModel.RechargeLog, *uv.PagingOut, error) {
	dbx := game_db.DBPools.Get(gmModel.WebServer)
	if dbx == nil {
		return nil, nil, errors.New("webserver database not connected")
	}

	tx := q.FilterQuery(dbx)
	var lst = make([]logModel.RechargeLog, 0)

	po, err := uv.PagingFind(tx, &lst, pi, order)
	return lst, po, err
}

func GetOperateLogList(q *OperateQ, pi *uv.PagingIn, order *uv.Order) ([]logModel.OperateLog, *uv.PagingOut, error) {
	dbx := game_db.DBPools.Get(gmModel.WebServer)
	if dbx == nil {
		return nil, nil, errors.New("webserver database not connected")
	}

	tx := q.FilterQuery(dbx)
	var lst = make([]logModel.OperateLog, 0)

	po, err := uv.PagingFind(tx, &lst, pi, order)
	return lst, po, err
}
