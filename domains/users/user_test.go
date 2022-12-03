package users

import (
	"context"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"gopkg.in/stretchr/testify.v1/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	svc UserServicer
)

func TestMain(m *testing.M) {
	code := 1
	defer func() {
		cleanup()
		os.Exit(code)
	}()
	d, err := setupTestDB()
	if err != nil {
		panic(err)
	}
	db = d
	s, err := New(db)
	if err != nil {
		panic(err)
	}
	svc = s
	code = m.Run()
}

func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	user := &User{
		Name:      gofakeit.Name(),
		HighScore: 2000,
	}

	err := svc.Create(ctx, user)
	require.NoError(t, err)
	assert.NotEmpty(t, user.ID)

	err = svc.Delete(ctx, user.ID)
	require.NoError(t, err)
}

func TestFindUsers(t *testing.T) {
	ctx := context.Background()
	ids := []string{}
	defer func() {
		for _, v := range ids {
			require.NoError(t, svc.Delete(ctx, v))
		}
	}()
	users := []*User{
		{
			Name:      gofakeit.Name(),
			HighScore: 2300,
		},
		{
			Name:      gofakeit.Name(),
			HighScore: 1500,
		},
	}
	for _, v := range users {
		require.NoError(t, svc.Create(ctx, v))
		require.NotEmpty(t, v.ID)
		ids = append(ids, v.ID)
	}

	res, err := svc.Find(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, res)
	require.Len(t, res, 2)
}

func TestFindByID(t *testing.T) {
	ctx := context.Background()
	ids := []string{}
	defer func() {
		for _, v := range ids {
			require.NoError(t, svc.Delete(ctx, v))
		}
	}()

	users := []*User{
		{
			Name:      gofakeit.Name(),
			HighScore: 2300,
		},
		{
			Name:      gofakeit.Name(),
			HighScore: 1500,
		},
	}
	for _, v := range users {
		require.NoError(t, svc.Create(ctx, v))
		require.NotEmpty(t, v.ID)
		ids = append(ids, v.ID)
	}

	res, err := svc.FindByID(ctx, ids[1])
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, ids[1], res.ID)
	assert.Equal(t, uint32(1500), res.HighScore)
}

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()
	ids := []string{}
	defer func() {
		for _, v := range ids {
			require.NoError(t, svc.Delete(ctx, v))
		}
	}()

	users := []*User{
		{
			Name:      gofakeit.Name(),
			HighScore: 2300,
		},
		{
			Name:      gofakeit.Name(),
			HighScore: 1500,
		},
	}
	for _, v := range users {
		require.NoError(t, svc.Create(ctx, v))
		require.NotEmpty(t, v.ID)
		ids = append(ids, v.ID)
	}

	newUser := &User{
		Name:      "John",
		HighScore: 2500,
	}

	err := svc.Update(ctx, ids[1], newUser)
	require.NoError(t, err)
	assert.Equal(t, ids[1], newUser.ID)

	res, err := svc.FindByID(ctx, ids[1])
	require.NoError(t, err)
	require.NotEmpty(t, res)

	assert.Equal(t, newUser.Name, res.Name)
	assert.Equal(t, newUser.HighScore, res.HighScore)
	assert.Equal(t, newUser.CreatedAt, res.CreatedAt)
	assert.Equal(t, newUser.UpdatedAt, res.UpdatedAt)
}

func setupTestDB() (*gorm.DB, error) {
	env := os.Getenv("ENVIRONMENT")
	dsn := "root:@tcp(127.0.0.1:3306)/sybo?charset=utf8mb4&parseTime=True&loc=Local"
	if env == "cicd" {
		dsn = "user:password@tcp(127.0.0.1:33306)/scores?charset=utf8mb4&parseTime=True&loc=Local"
	}
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func cleanup() {
	db.Exec("DELETE FROM users")
}
