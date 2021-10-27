package utils

import (
	"fmt"
	"strings"
)

const DEFAULT_PAGE_SIZE uint = 10

type PageQuery struct {
	PageNum  uint   `json:"page"`
	PageSize uint   `json:"size"`
	Sort     string `json:"sort"`

	// 当sort 只有一个属性名字时有效
	Order string `json:"order"` // asc desc
}

func (q *PageQuery) EffecitvePageSize() uint {
	if q.PageSize > 0 {
		return q.PageSize
	}
	return DEFAULT_PAGE_SIZE
}

func (q *PageQuery) EffecitvePageNum() uint {
	if q.PageNum > 0 {
		return q.PageNum
	}
	return 1
}

func (q *PageQuery) Offset() uint {
	return q.EffecitvePageSize() * (q.EffecitvePageNum() - 1)
}

func (q *PageQuery) SortSQL() string {
	if q.Sort == "" {
		return ""
	}
	if !isSingleSortProperty(q.Sort) || !strings.EqualFold("ASC", q.Order) && !strings.EqualFold("DESC", q.Order) {
		return buildOrderSql(q.Sort)
	}
	return fmt.Sprintf("`%s` %s", q.Sort, strings.ToUpper(q.Order))
}

func isSingleSortProperty(sortString string) bool {
	return !strings.Contains(sortString, ",") && !strings.Contains(sortString, ";")
}

func buildOrderSql(sortString string) string {
	sortSQL := strings.Builder{}

	pairs := strings.Split(sortString, ";")
	pairsCount := len(pairs)
	for idx, p := range pairs {
		pair := strings.Split(p, ",")
		if pair[0] == "" {
			continue
		}
		sortSQL.WriteString("`" + pair[0] + "`")
		if len(pair) > 1 && (strings.EqualFold(pair[1], "ASC") || strings.EqualFold(pair[1], "DESC")) {
			sortSQL.WriteString(" " + strings.ToUpper(pair[1]))
		}
		if idx < pairsCount-1 {
			sortSQL.WriteString(",")
		}
	}
	return sortSQL.String()
}
