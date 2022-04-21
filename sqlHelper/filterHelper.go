package sqlHelper

import (
	"github.com/gin-gonic/gin"
	"github.com/pro-assistance/pro-assister/sqlHelper/filter"
	"github.com/pro-assistance/pro-assister/sqlHelper/paginator"
	"github.com/pro-assistance/pro-assister/sqlHelper/sorter"
)

type QueryFilter struct {
	Col       string
	Value     string
	Filter    *filter.Filter
	Sorter    *sorter.Sorter
	Paginator *paginator.Paginator
}

func (i *SQLHelper) CreateQueryFilter(c *gin.Context) (*QueryFilter, error) {
	col := c.Query("col")
	value := c.Query("value")
	filterItem, err := filter.NewFilter(c)
	if err != nil {
		return nil, err
	}
	sorterItem, err := sorter.NewSorter(c)
	if err != nil {
		return nil, err
	}
	paginatorItem, err := paginator.NewPaginator(c)
	if err != nil {
		return nil, err
	}
	return &QueryFilter{Col: col, Value: value, Filter: filterItem, Sorter: sorterItem, Paginator: paginatorItem}, nil
}