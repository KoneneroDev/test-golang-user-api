package postgres

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
	"time"
)

func newTestStorage(t *testing.T) (*Storage, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	cleanup := func() {
		db.Close()
	}
	return &Storage{db: db}, mock, cleanup
}

func TestStorageCreateUser(t *testing.T) {
	t.Run("success save user to db", func(t *testing.T) {
		//given
		storage, mock, cleanup := newTestStorage(t)
		defer cleanup()

		user := &UserDto{
			ID:        uuid.New(),
			Firstname: "Ivan",
			Lastname:  "Ivanov",
			Email:     "ivan@gmail.com",
			Age:       30,
			Created:   time.Now(),
		}

		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO users (id, firstname, lastname, email, age, created)`)).
			WithArgs(user.ID, user.Firstname, user.Lastname, user.Email, user.Age, user.Created).
			WillReturnResult(sqlmock.NewResult(1, 1))
		//when
		err := storage.CreateUser(user)
		//then
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("returns error when insert fails", func(t *testing.T) {
		//given
		storage, mock, cleanup := newTestStorage(t)
		defer cleanup()

		user := &UserDto{
			ID:        uuid.New(),
			Firstname: "Ivan",
			Lastname:  "Ivanov",
			Email:     "ivan@gmail.com",
			Age:       30,
			Created:   time.Now(),
		}

		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO users (id, firstname, lastname, email, age, created)`)).
			WithArgs(user.ID, user.Firstname, user.Lastname, user.Email, user.Age, user.Created).
			WillReturnError(fmt.Errorf("insert error"))

		//when
		err := storage.CreateUser(user)

		//then
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to insert user")
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestStorageGetUser(t *testing.T) {
	t.Run("success save user to db", func(t *testing.T) {
		//given
		storage, mock, cleanup := newTestStorage(t)
		defer cleanup()

		id := uuid.New()
		created := time.Now()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, firstname, lastname, email, age, created FROM users WHERE id = $1`)).
			WithArgs(id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "firstname", "lastname", "email", "age", "created"}).
				AddRow(id, "Ivan", "Ivanov", "ivan@gmail.com", 30, created))

		//when
		user, err := storage.GetUser(id)
		//then
		require.NoError(t, err)
		assert.Equal(t, "Ivan", user.Firstname)
		assert.Equal(t, "Ivanov", user.Lastname)
		assert.Equal(t, "ivan@gmail.com", user.Email)
		assert.Equal(t, 30, user.Age)
	})

	t.Run("returns error when user not found", func(t *testing.T) {
		// given
		storage, mock, cleanup := newTestStorage(t)
		defer cleanup()

		id := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, firstname, lastname, email, age, created FROM users WHERE id = $1`)).
			WithArgs(id).
			WillReturnError(sql.ErrNoRows)

		// when
		user, err := storage.GetUser(id)
		// then
		require.Nil(t, user)
		require.Error(t, err)
		require.Contains(t, err.Error(), "user not found")
	})

}

func TestStorageEditUser(t *testing.T) {
	t.Run("success edit user to db", func(t *testing.T) {
		//given
		storage, mock, cleanup := newTestStorage(t)
		defer cleanup()

		user := &UserDto{
			ID:        uuid.New(),
			Firstname: "Ivan",
			Lastname:  "Ivanov",
			Email:     "ivan@gmail.com",
			Age:       30,
		}

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET firstname = $1, lastname = $2, email = $3, age = $4 WHERE id = $5`)).
			WithArgs(user.Firstname, user.Lastname, user.Email, user.Age, user.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		//when
		updated, err := storage.EditUser(user)
		//then
		require.NoError(t, err)
		assert.Equal(t, user, updated)
	})

	t.Run("returns error when no rows updated", func(t *testing.T) {
		// given
		storage, mock, cleanup := newTestStorage(t)
		defer cleanup()

		user := &UserDto{
			ID:        uuid.New(),
			Firstname: "Ivan",
			Lastname:  "Ivanov",
			Email:     "ivan@gmail.com",
			Age:       30,
		}

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET firstname = $1, lastname = $2, email = $3, age = $4 WHERE id = $5`)).
			WithArgs(user.Firstname, user.Lastname, user.Email, user.Age, user.ID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		// when
		updated, err := storage.EditUser(user)
		// then
		require.Nil(t, updated)
		require.Error(t, err)
		require.Contains(t, err.Error(), "user not found")
	})
}

func TestStorageDeleteUser(t *testing.T) {
	t.Run("success delete user to db", func(t *testing.T) {
		//given
		storage, mock, cleanup := newTestStorage(t)
		defer cleanup()

		id := uuid.New()

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM users WHERE id = $1`)).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		//when
		err := storage.DeleteUser(id)
		//then
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("returns error when no rows deleted", func(t *testing.T) {
		// given
		storage, mock, cleanup := newTestStorage(t)
		defer cleanup()

		id := uuid.New()

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM users WHERE id = $1`)).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 0))

		// when
		err := storage.DeleteUser(id)
		// then
		require.Error(t, err)
		require.Contains(t, err.Error(), "user not found")
	})
}
