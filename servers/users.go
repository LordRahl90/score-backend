package servers

import (
	"errors"
	"net/http"

	"sybo/domains/users"
	"sybo/requests"
	"sybo/responses"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (s *Server) userRoutes() {
	s.Router.GET("users", s.findAllUsers)
	user := s.Router.Group("user")
	{
		user.POST("", s.createUser)
		user.GET(":id", s.findUser)
		user.PUT(":id", s.updateUser)
		user.DELETE(":id", s.deleteUser)
	}
}

func (s *Server) findAllUsers(ctx *gin.Context) {
	res, err := s.userService.Find(ctx)
	if err != nil {
		serverError(ctx, http.StatusInternalServerError, err)
		return
	}
	var response []responses.User
	for _, v := range res {
		response = append(response, responses.User{
			ID:        v.ID,
			Name:      v.Name,
			HighScore: v.HighScore,
		})
	}

	success(ctx, map[string]interface{}{"users": response})
}

func (s *Server) findUser(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := s.userService.FindByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			serverError(ctx, http.StatusNotFound, err)
			return
		}
		serverError(ctx, http.StatusInternalServerError, err)
		return
	}
	response := responses.User{
		ID:        res.ID,
		Name:      res.Name,
		HighScore: res.HighScore,
	}

	success(ctx, response)
}

func (s *Server) createUser(ctx *gin.Context) {
	var req requests.User
	if err := ctx.ShouldBindJSON(&req); err != nil {
		serverError(ctx, http.StatusBadRequest, err)
		return
	}

	ent := &users.User{
		Name:      req.Name,
		HighScore: req.HighScore,
	}
	err := s.userService.Create(ctx.Request.Context(), ent)
	if err != nil {
		serverError(ctx, http.StatusInternalServerError, err)
		return
	}

	res := responses.User{
		ID:        ent.ID,
		Name:      ent.Name,
		HighScore: ent.HighScore,
	}
	created(ctx, res)
}

func (s *Server) updateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	var req requests.User
	if err := ctx.ShouldBindJSON(&req); err != nil {
		serverError(ctx, http.StatusBadRequest, err)
		return
	}
	ent := users.User{
		Name:      req.Name,
		HighScore: req.HighScore,
	}
	err := s.userService.Update(ctx, id, &ent)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			serverError(ctx, http.StatusNotFound, err)
			return
		}
		serverError(ctx, http.StatusInternalServerError, err)
		return
	}
	res := responses.User{
		ID:        ent.ID,
		Name:      ent.Name,
		HighScore: ent.HighScore,
	}
	success(ctx, res)
}

func (s *Server) deleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	err := s.userService.Delete(ctx.Request.Context(), id)
	if err != nil {
		serverError(ctx, http.StatusInternalServerError, err)
		return
	}
	success(ctx, "")
}
