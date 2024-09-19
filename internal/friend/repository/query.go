package repository

import (
	"fmt"
	"go-service/internal/friend/domain"
)

func BuildQuery(filter domain.FriendFilter, buildParam func(int) string) (qr string, params []interface{}) {
	selectClause := fmt.Sprintf("select * from ((select * from users a left join relations b on a.id = b.user_id1 where b.user_id2 = %s) union (select * from users a join relations b on a.id = b.user_id2 where b.user_id1 = %s))", buildParam(1), buildParam(1))

	params = append(params, *filter.UserId)
	whereClause := "where"
	orderByClause := ""
	limitClause := ""
	if filter.Q != nil {
		params = append(params, *filter.Q)
		whereClause = fmt.Sprintf("%s b.name like CONCAT('%%',%s::text,'%%')", whereClause, buildParam(len(params)))
	}

	if len(filter.Sorts) > 0 {
		for _, v := range filter.Sorts {
			if len(v) > 0 {
				sortType := "ASC"
				if v[0] == '-' {
					sortType = "DESC"
				}
				orderByClause = fmt.Sprintf(" %s %s %s,", orderByClause, v, sortType)
			}

		}

		orderByClause = "order by" + orderByClause[:len(orderByClause)-2]
	}

	if filter.Page != nil && filter.Limit != nil {
		offset := *filter.Page * *filter.Limit
		limitClause = fmt.Sprintf("LIMIT %v OFFSET %v ", filter.Limit, offset)
	}
	return fmt.Sprintf("%s %s %s %s", selectClause, whereClause, orderByClause, limitClause), params

}
