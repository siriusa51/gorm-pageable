package pageable

import (
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type PageParameter struct {
	PageNow    int
	RawPerPage int
	// for conditional queries
	Condition [][]gen.Condition
	// for sorting results
	OrderBy []field.Expr
	Dao     gen.Dao
}

type PageResponse[T any] struct {
	dao gen.Dao
	// current page of query
	PageNow int
	// total page of the query
	PageCount int
	// total raw of query
	RawCount int
	// if the result is empty
	Empty bool
	// rpp
	RawPerPage int
	// query result list
	Raws []T
	// if the result is the first page
	FirstPage bool
	// if the result is the last page
	LastPage bool
	// The number of first record the resultSet
	StartRow int
	// The number of last record the resultSet
	EndRow int
}

// PageQuery
// main handler of query
func PageQuery[T any](param *PageParameter) (*PageResponse[T], error) {
	// recovery
	if recovery != nil {
		defer recovery()
	}
	// get limit and offSet
	var limit, offset int
	if !use0Page {
		limit, offset = getLimitOffset(param.PageNow-1, param.RawPerPage)
		if param.PageNow < 1 {
			param.PageNow = 1
		}
	} else {
		limit, offset = getLimitOffset(param.PageNow, param.RawPerPage)
	}

	query := param.Dao
	if param.Condition != nil && len(param.Condition) > 0 {
		query = query.Where(param.Condition[0]...)
		for _, cond := range param.Condition[1:] {
			query = query.Or(cond...)
		}
	}

	// get total count of the table
	count64, err := query.Count()
	if err != nil {
		return nil, err
	}

	if param.OrderBy != nil && len(param.OrderBy) > 0 {
		query = query.Order(param.OrderBy...)
	}

	count := int(count64)

	// get page count
	if param.RawPerPage <= 0 {
		param.RawPerPage = count
		limit = count
	}

	query = query.Limit(limit)

	result, err := query.Offset(offset).Find()
	if err != nil {
		return nil, err
	}

	pageCount := count / param.RawPerPage
	if count%param.RawPerPage != 0 {
		pageCount++
	}

	startRow, endRow, empty, lastPage := 0, 0, (param.PageNow > pageCount) || count == 0, param.PageNow == pageCount
	if !empty {
		startRow = param.PageNow * param.RawPerPage
		if !lastPage {
			endRow = (param.PageNow+1)*param.RawPerPage - 1
		} else {
			endRow = count
		}
	}

	// prepare base response
	return &PageResponse[T]{
		dao:        param.Dao,
		PageNow:    param.PageNow,
		PageCount:  pageCount,
		RawPerPage: param.RawPerPage,
		RawCount:   count,
		Raws:       result.([]T),
		FirstPage:  param.PageNow == 1,
		LastPage:   lastPage,
		Empty:      empty,
		StartRow:   startRow,
		EndRow:     endRow,
	}, nil
}

// GetNextPage
// return next page`s Response
// 	func getResultSet (page int,rowsPerPage int)(*pageable.Response,error){
// 	//your empty result set
// 		resultSet := make([]*user,0,30)
// 		//prepare a handler to query
// 		handler := DB.
// 			Module(&user{}).
// 			Where(&user{Active:true})
// 		//use PageQuery to get data (this page)
// 		resp,err := pageable.PageQuery(page,rowsPerPage,handler,&resultSet)
// 		// handle error
// 		f err != nil {
// 			panic(err)
// 		}
//		//get next page
//		resp,err := resp.GetNextPage()	//Response of next page
// 	}
func (r *PageResponse[T]) GetNextPage() (*PageResponse[T], error) {
	return PageQuery[T](&PageParameter{PageNow: r.PageNow + 1, RawPerPage: r.RawPerPage, Dao: r.dao})
}

// GetLastPage
// return last page`s Response
func (r *PageResponse[T]) GetLastPage() (*PageResponse[T], error) {
	return PageQuery[T](&PageParameter{PageNow: r.PageNow - 1, RawPerPage: r.RawPerPage, Dao: r.dao})
}

// GetEndPage
// return end page`s Response
func (r *PageResponse[T]) GetEndPage() (*PageResponse[T], error) {
	return PageQuery[T](&PageParameter{PageNow: r.PageCount, RawPerPage: r.RawPerPage, Dao: r.dao})
}

// GetFirstPage
// return first page`s Response
func (r *PageResponse[T]) GetFirstPage() (*PageResponse[T], error) {
	p := 1
	if use0Page {
		p = 0
	}

	return PageQuery[T](&PageParameter{PageNow: p, RawPerPage: r.RawPerPage, Dao: r.dao})
}
