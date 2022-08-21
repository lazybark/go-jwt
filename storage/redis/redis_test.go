package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRestoringIdFromPrevSession(t *testing.T) {
	id := 15
	rdb2, err := NewRedisStorage("localhost:6379", "", false, 5)
	require.NoError(t, err)
	err = rdb2.Flush()
	require.NoError(t, err)

	//There is no record about uid - set to 0 while creating rdb
	n, err := rdb2.GenerateUserId()
	require.NoError(t, err)
	assert.Equal(t, 1, n)

	err = rdb2.db.Set("last_uid", id, 0).Err()
	require.NoError(t, err)

	//New rdb should see last_uid record and use it as own during Init()
	rdb, err := NewRedisStorage("localhost:6379", "", false, 5)
	require.NoError(t, err)

	err = rdb.Init() //Init adds 1 user and makes last_uid+1
	require.NoError(t, err)

	n, err = rdb.GenerateUserId()
	require.NoError(t, err)

	assert.Equal(t, id+2, n)
}
