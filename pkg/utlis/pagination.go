package utlis

import (
	"fmt"
	"math"
	"strconv"
)

const (
	defaultSize = 10
)

type Pagination struct {
	Size    int    `json:"size,omitempty"`
	Page    int    `json:"page,omitempty"`
	OrderBy string `json:"orderBy,omitempty"`
}

func NewPaginationQuery(size int, page int) *Pagination {
	return &Pagination{
		Size: size,
		Page: page,
	}
}

func (q *Pagination) SetPage(pageQuery string) error {
	if pageQuery == "" {
		q.Size = 0
		return nil
	}
	n, err := strconv.Atoi(pageQuery)
	if err != nil {
		return err
	}
	q.Page = n
	
	return nil
}


func (q *Pagination) SetOrderBy(orderByQuery string) {
	q.OrderBy = orderByQuery
}

func (q *Pagination) GetOffSet() int {
	if q.Page == 0{
		return 0
	}
	return (q.Page - 1) * q.Size
}

func (q *Pagination) GetLimit() int {
	return q.Size
}

func (q *Pagination) GetOrderBy() string {
	return q.OrderBy
}

func (q *Pagination) GetPage() int {
	return q.Page
}

func (q *Pagination) GetSize() int {
	return q.Size
}

func (q *Pagination) GetQueryString() string {
	return fmt.Sprintf("page=%v&size=%v&orderBy=%s", q.GetPage(),q.GetSize(),q.GetOrderBy())
}

func (q *Pagination) GetTotalPages(totalCount int) int {
	d := float64(totalCount)/float64(q.GetSize())
	return int(math.Ceil(d))
}

func (q *Pagination) GetHasMore(totalCount int) bool {
	return q.GetPage() < totalCount/q.GetSize()
}