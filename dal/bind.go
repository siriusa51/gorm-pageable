package query

import (
	"github.com/siriusa51/gorm-pageable/dal/types"
	"gorm.io/gen"
)

func Bind(generator *gen.Generator) {
	generator.ApplyBasic(types.User{})
}
