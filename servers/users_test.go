package servers

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"sybo/domains/users"
	"sybo/requests"
	"sybo/responses"

	"github.com/brianvoe/gofakeit/v6"
	"gopkg.in/stretchr/testify.v1/assert"
	"gopkg.in/stretchr/testify.v1/require"
	"gorm.io/gorm"
)

func TestCreate(t *testing.T) {
	req := requests.User{
		Name:      "John",
		HighScore: 1337,
	}

	b, err := json.Marshal(req)
	require.NoError(t, err)
	require.NotNil(t, b)

	res := handleRequest(t, http.MethodPost, "/user", b)
	require.Equal(t, http.StatusCreated, res.Code)

	var response responses.User
	err = json.Unmarshal(res.Body.Bytes(), &response)
	require.NoError(t, err)
	require.NotEmpty(t, response)

	assert.Equal(t, req.Name, response.Name)
	assert.Equal(t, req.HighScore, response.HighScore)

	require.NoError(t, server.userService.Delete(context.Background(), response.ID))
}

func TestCreateWithBadJson(t *testing.T) {
	b := []byte(`
	""{"name":"John","highscore":1337}
	`)
	res := handleRequest(t, http.MethodPost, "/user", b)
	require.Equal(t, http.StatusBadRequest, res.Code)
}

func TestCreateWithNoHighScoreField(t *testing.T) {
	b := []byte(`
		{"name":"John","higsscore":1337}
	`)
	res := handleRequest(t, http.MethodPost, "/user", b)
	require.Equal(t, http.StatusBadRequest, res.Code)
}

func TestCreateWithNoNameField(t *testing.T) {
	b := []byte(`
		{"names":"John","highscore":1337}
	`)
	res := handleRequest(t, http.MethodPost, "/user", b)
	require.Equal(t, http.StatusBadRequest, res.Code)
}

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()
	ent := &users.User{
		Name:      "John",
		HighScore: 3500,
	}
	err := server.userService.Create(ctx, ent)
	require.NoError(t, err)

	req := requests.User{
		Name:      ent.Name,
		HighScore: 2345,
	}
	b, err := json.Marshal(req)
	require.NoError(t, err)
	require.NotNil(t, b)

	res := handleRequest(t, http.MethodPut, "/user/"+ent.ID, b)
	require.Equal(t, http.StatusOK, res.Code)

	val, err := server.userService.FindByID(ctx, ent.ID)
	require.NoError(t, err)
	require.NotEmpty(t, val)

	assert.Equal(t, req.Name, val.Name)
	assert.Equal(t, req.HighScore, val.HighScore)

	// attempt to update with invalid ID
	res = handleRequest(t, http.MethodPut, "/user/"+gofakeit.UUID(), b)
	require.Equal(t, http.StatusNotFound, res.Code)

	require.NoError(t, server.userService.Delete(context.Background(), ent.ID))
}

func TestUpdateWithBadJson(t *testing.T) {
	b := []byte(`
	""{"name":"John","highscore":1337}
	`)
	res := handleRequest(t, http.MethodPut, "/user/"+gofakeit.UUID(), b)
	require.Equal(t, http.StatusBadRequest, res.Code)
}

func TestFindNonExistingUser(t *testing.T) {
	res := handleRequest(t, http.MethodGet, "/user/"+gofakeit.UUID(), nil)
	require.Equal(t, http.StatusNotFound, res.Code)
}

func TestFind(t *testing.T) {
	ctx := context.Background()
	ids := []string{}
	defer func() {
		for _, v := range ids {
			require.NoError(t, server.userService.Delete(ctx, v))
		}
	}()
	users := []*users.User{
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
		require.NoError(t, server.userService.Create(ctx, v))
		require.NotEmpty(t, v.ID)
		ids = append(ids, v.ID)
	}

	res := handleRequest(t, http.MethodGet, "/user/"+ids[0], nil)
	require.Equal(t, http.StatusOK, res.Code)
	var response responses.User
	err := json.Unmarshal(res.Body.Bytes(), &response)
	require.NoError(t, err)
	require.NotEmpty(t, response)

	assert.Equal(t, users[0].Name, response.Name)
	assert.Equal(t, users[0].HighScore, response.HighScore)
}

func TestFindEmptyAllUsers(t *testing.T) {
	res := handleRequest(t, http.MethodGet, "/users", nil)
	require.Equal(t, http.StatusOK, res.Code)

	var response map[string]interface{}

	err := json.Unmarshal(res.Body.Bytes(), &response)
	require.NoError(t, err)
	require.NotEmpty(t, response)

	users := response["users"]
	assert.Nil(t, users)
}

func TestFindAllUsers(t *testing.T) {
	ctx := context.Background()
	ids := []string{}
	defer func() {
		for _, v := range ids {
			require.NoError(t, server.userService.Delete(ctx, v))
		}
	}()
	users := []*users.User{
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
		require.NoError(t, server.userService.Create(ctx, v))
		require.NotEmpty(t, v.ID)
		ids = append(ids, v.ID)
	}

	res := handleRequest(t, http.MethodGet, "/users", nil)
	require.Equal(t, http.StatusOK, res.Code)

	var response map[string]interface{}

	err := json.Unmarshal(res.Body.Bytes(), &response)
	require.NoError(t, err)
	require.NotEmpty(t, response)

	resUsers := response["users"]
	assert.NotNil(t, resUsers)
	assert.Len(t, resUsers, 2)
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	ids := []string{}
	defer func() {
		for _, v := range ids {
			require.NoError(t, server.userService.Delete(ctx, v))
		}
	}()
	users := []*users.User{
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
		require.NoError(t, server.userService.Create(ctx, v))
		require.NotEmpty(t, v.ID)
		ids = append(ids, v.ID)
	}

	res := handleRequest(t, http.MethodDelete, "/user/"+ids[0], nil)
	require.Equal(t, http.StatusOK, res.Code)

	val, err := server.userService.FindByID(ctx, ids[0])
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	require.Empty(t, val)
}
