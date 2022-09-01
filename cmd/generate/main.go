package main

import (
	"github.com/siriusa51/gorm-pageable/dal"
	"gorm.io/driver/sqlite"
	"gorm.io/gen"
	"gorm.io/gorm"
)

// generate code
func main() {
	// specify the output directory (default: "./query")
	// ### if you want to query without context constrain, set mode gen.WithoutContext ###
	g := gen.NewGenerator(gen.Config{
		OutPath: "./dal/model",
		Mode:    gen.WithDefaultQuery | gen.WithQueryInterface,
		// if you want the nullable field generation property to be pointer type, set FieldNullable true
		FieldNullable: true,
		// if you want to assign field which has default value in `Create` API, set FieldCoverable true, reference: https://gorm.io/docs/create.html#Default-Values
		FieldCoverable: true,
		// if you want generate field with unsigned integer type, set FieldSignable true
		FieldSignable: true,
		// if you want to generate index tags from database, set FieldWithIndexTag true
		FieldWithIndexTag: true,
		// if you want to generate type tags from database, set FieldWithTypeTag true
		FieldWithTypeTag: true,
		// if you need unit tests for query code, set WithUnitTest true
		WithUnitTest: true,
	})

	db, _ := gorm.Open(sqlite.Open(":memory:"))
	g.UseDB(db)

	query.Bind(g)

	g.Execute()
}
