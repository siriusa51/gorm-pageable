package pageable

type ListPageParameter[T any] struct {
	PageNow    int
	RawPerPage int
	List       []T
}

type ListPageResponse[T any] struct {
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

// ListPageQuery
// main handler of list query
func ListPageQuery[T any](param *ListPageParameter[T]) (*ListPageResponse[T], error) {
	// recovery
	if recovery != nil {
		defer recovery()
	}
	if param.PageNow < 0 {
		param.PageNow = 0
	}

	if param.RawPerPage <= 0 {
		param.RawPerPage = defaultRpp
	}

	count := len(param.List)

	pageCount := count / param.RawPerPage
	if count%param.RawPerPage != 0 {
		pageCount++
	}

	var startRow int
	var endRow int
	if use0Page {
		startRow = max(0, param.RawPerPage*param.PageNow)
		endRow = min(count, param.RawPerPage*(param.PageNow+1))
	} else {
		param.PageNow = max(1, param.PageNow)
		startRow = max(0, param.RawPerPage*(param.PageNow-1))
		endRow = min(count, param.RawPerPage*(param.PageNow))
	}
	empty := count == 0
	lastPage := param.PageNow == pageCount
	result := param.List[startRow:endRow]

	// prepare base response
	return &ListPageResponse[T]{
		PageNow:    param.PageNow,
		PageCount:  pageCount,
		RawPerPage: param.RawPerPage,
		RawCount:   count,
		Raws:       result,
		FirstPage:  param.PageNow == 1,
		LastPage:   lastPage,
		Empty:      empty,
		StartRow:   startRow,
		EndRow:     endRow,
	}, nil
}
