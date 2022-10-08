package pageable

import (
	"fmt"
	"github.com/siriusa51/gorm-pageable/dal/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListPageQuery(t *testing.T) {
	items := make([]*types.User, 0)
	size := 20
	for i := 0; i < size; i++ {
		items = append(items, &types.User{
			Name:        fmt.Sprintf("user-%d", i),
			Description: fmt.Sprintf("desc-%d", i),
		})
	}

	p, err := ListPageQuery[*types.User](&ListPageParameter[*types.User]{
		PageNow:    1,
		RawPerPage: 0,
		List:       items,
	})

	assert.NoError(t, err)
	assert.Equal(t, size, len(p.Raws))

	p, err = ListPageQuery[*types.User](&ListPageParameter[*types.User]{
		PageNow:    1,
		RawPerPage: 10,
		List:       items,
	})

	assert.NoError(t, err)
	assert.Equal(t, 10, len(p.Raws))

	p, err = ListPageQuery[*types.User](&ListPageParameter[*types.User]{
		PageNow:    2,
		RawPerPage: 10,
		List:       items,
	})

	assert.NoError(t, err)
	assert.Equal(t, 10, len(p.Raws))

	p, err = ListPageQuery[*types.User](&ListPageParameter[*types.User]{
		PageNow:    3,
		RawPerPage: 10,
		List:       items,
	})

	assert.NoError(t, err)
	assert.Equal(t, 0, len(p.Raws))
}
