package pageable

import (
	"context"
	"fmt"
	"github.com/siriusa51/gorm-pageable/dal/model"
	"github.com/siriusa51/gorm-pageable/dal/types"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"testing"
)

func TestPageQuery(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db = db.Debug()

	q := model.Use(db)
	_ = db.AutoMigrate(types.User{})

	udo := q.User.WithContext(context.TODO())
	items := make([]*types.User, 0)
	size := 20
	for i := 0; i < size; i++ {
		items = append(items, &types.User{
			Name:        fmt.Sprintf("user-%d", i),
			Description: fmt.Sprintf("desc-%d", i),
		})
	}
	_ = udo.Create(items...)

	p, err := PageQuery[*types.User](&PageParameter{
		PageNow:    1,
		RawPerPage: 0,
		Dao:        udo.As(udo.TableName()),
	})

	assert.NoError(t, err)
	assert.Equal(t, size, len(p.Raws))

	p, err = PageQuery[*types.User](&PageParameter{
		PageNow:    1,
		RawPerPage: 10,
		Dao:        udo.As(udo.TableName()),
	})

	assert.NoError(t, err)
	assert.Equal(t, 10, len(p.Raws))

	p, err = PageQuery[*types.User](&PageParameter{
		PageNow:    1,
		RawPerPage: 5,
		Dao:        udo.As(udo.TableName()),
		ConditionFunc: func(query gen.Dao) gen.Dao {
			return query.Where(q.User.Name.Value("user-1")).Or(q.User.Name.Value("user-2"))
		},
		OrderBy: []field.Expr{q.User.Name.Desc()},
	})

	assert.NoError(t, err)
	assert.Equal(t, 2, len(p.Raws))
}
