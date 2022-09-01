package pageable

import (
	"fmt"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"runtime/debug"
	"time"
)

type PageParameter struct {
	PageNow    int
	RawPerPage int
	// for conditional queries
	Condition []gen.Condition
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

// getLimitOffset
// (private) get LIMIT and OFFSET keyword in SQL
func getLimitOffset(page, rpp int) (limit, offset int) {
	if page < 0 {
		page = 0
	}
	if rpp < 1 {
		rpp = defaultRpp
	}
	return rpp, page * rpp
}

// SetDefaultRPP
// Set default rpp
func SetDefaultRPP(rpp int) error {
	if rpp < 1 {
		return fmt.Errorf("invalid input rpp")
	}
	defaultRpp = rpp
	return nil
}

// recoveryHandler
// default type of recovery handler
type recoveryHandler func()

// recovery
// handler of panic
var recovery recoveryHandler

var defaultRpp int

var use0Page bool

// SetRecovery
// Set custom recovery
// Here are some sample of the custom recovery
// 	package main
// 	import (
// 		"fmt"
// 		pageable "github.com/BillSJC/gorm-pageable"
// 	)
//
// 	//your recovery
// 	func myRecovery(){
// 		if err := recover ; err != nil {
// 			fmt.Println("something happened")
// 			fmt.Println(err)
// 			//then you can do some logs...
// 		}
// 	}
//
// 	func init(){
// 		//setup your recovery
// 		pageable.SetRecovery(myRecovery)
// 	}
func SetRecovery(handler func()) {
	recovery = handler
}

// defaultRecovery
// print base recover info
func defaultRecovery() {
	if err := recover(); err != nil {
		// print panic info
		fmt.Printf("Panic recovered: %s \n\n Time: %s \n\n Stack Trace: \n\n",
			fmt.Sprint(err),
			time.Now().Format("2006-01-02 15:04:05"),
		)
		// stack
		debug.PrintStack()
	}
}

// init
// use default recovery
func init() {
	// use default rpp
	_ = SetDefaultRPP(25)
	// use 1 as default page
	use0Page = false
}

// Use0AsFirstPage
// the default first page is 1. However,if u want to use 0 as the first page, just follow this step:
// pageable.Use0AsFirstPage()
func Use0AsFirstPage() {
	use0Page = true
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
		query = query.Where(param.Condition...)
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
