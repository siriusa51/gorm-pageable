package pageable

import (
	"fmt"
	"runtime/debug"
	"time"
)

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

type number interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

func max[T number](a T, b T) T {
	if a > b {
		return a
	}
	return b
}

func min[T number](a T, b T) T {
	if a > b {
		return b
	}
	return a
}
