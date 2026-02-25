package test

import (
	"bluebell/dao"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newCommunityDaoWithMock(t *testing.T) (*dao.CommunityDao, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	cleanup := func() {
		_ = db.Close()
	}
	return dao.NewCommunityDao(sqlxDB), mock, cleanup
}

func TestNewCommunityDao(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	c := dao.NewCommunityDao(sqlxDB)
	assert.NotNil(t, c)
}

func TestGetCommunityDetail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c, mock, cleanup := newCommunityDaoWithMock(t)
		defer cleanup()

		now := time.Date(2026, 2, 25, 0, 0, 0, 0, time.UTC)
		rows := sqlmock.NewRows([]string{"community_id", "community_name", "introduction", "create_time"}).
			AddRow(int64(1), "go", "go intro", now)
		mock.ExpectQuery(regexp.QuoteMeta("select community_id, community_name, introduction, create_time from community where community_id = ? ")).
			WithArgs(int64(1)).
			WillReturnRows(rows)

		detail, err := c.GetCommunityDetail(1)
		require.NoError(t, err)
		require.NotNil(t, detail)
		assert.Equal(t, int64(1), detail.ID)
		assert.Equal(t, "go", detail.Name)
		assert.Equal(t, "go intro", detail.Introduction)
		assert.Equal(t, now, detail.CreateTime)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not_found", func(t *testing.T) {
		c, mock, cleanup := newCommunityDaoWithMock(t)
		defer cleanup()

		mock.ExpectQuery(regexp.QuoteMeta("select community_id, community_name, introduction, create_time from community where community_id = ? ")).
			WithArgs(int64(100)).
			WillReturnError(sql.ErrNoRows)

		detail, err := c.GetCommunityDetail(100)
		assert.Nil(t, detail)
		assert.ErrorIs(t, err, dao.ErrorCommunityNotExist)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query_error", func(t *testing.T) {
		c, mock, cleanup := newCommunityDaoWithMock(t)
		defer cleanup()

		expectedErr := errors.New("db timeout")
		mock.ExpectQuery(regexp.QuoteMeta("select community_id, community_name, introduction, create_time from community where community_id = ? ")).
			WithArgs(int64(101)).
			WillReturnError(expectedErr)

		detail, err := c.GetCommunityDetail(101)
		assert.Nil(t, detail)
		assert.ErrorIs(t, err, expectedErr)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetCommunityList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c, mock, cleanup := newCommunityDaoWithMock(t)
		defer cleanup()

		rows := sqlmock.NewRows([]string{"community_id", "community_name"}).
			AddRow(int64(1), "go").
			AddRow(int64(2), "java")
		mock.ExpectQuery(regexp.QuoteMeta("select community_id, community_name from community")).
			WillReturnRows(rows)

		list, err := c.GetCommunityList()
		require.NoError(t, err)
		require.Len(t, list, 2)
		assert.Equal(t, int64(1), list[0].ID)
		assert.Equal(t, "go", list[0].Name)
		assert.Equal(t, int64(2), list[1].ID)
		assert.Equal(t, "java", list[1].Name)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not_found", func(t *testing.T) {
		c, mock, cleanup := newCommunityDaoWithMock(t)
		defer cleanup()

		mock.ExpectQuery(regexp.QuoteMeta("select community_id, community_name from community")).
			WillReturnError(sql.ErrNoRows)

		list, err := c.GetCommunityList()
		assert.NoError(t, err)
		assert.Nil(t, list)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query_error", func(t *testing.T) {
		c, mock, cleanup := newCommunityDaoWithMock(t)
		defer cleanup()

		expectedErr := errors.New("db timeout")
		mock.ExpectQuery(regexp.QuoteMeta("select community_id, community_name from community")).
			WillReturnError(expectedErr)

		list, err := c.GetCommunityList()
		assert.Nil(t, list)
		assert.ErrorIs(t, err, expectedErr)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
