package handlers

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/raeinsoltani/gorello/back/models"
	"github.com/raeinsoltani/gorello/back/repository"
	"github.com/raeinsoltani/gorello/back/utils"
)

type UserHandler struct {
	UserRepo repository.User
}

func NewUserHandler(userRepo repository.User) *UserHandler {
	return &UserHandler{UserRepo: userRepo}
}

type UserResponseDTO struct {
	ID        uint   `json:"id"`
	Firstname string `json:"firstname"`
	Username  string `json:"username"`
}

type UserRegisterDTO struct {
	Username string `json:"username" validate:"required,alphanum,min=3,max=25"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserUpdateDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *UserHandler) Register(c echo.Context) error {
	userRegisterDTO := new(UserRegisterDTO)
	if err := c.Bind(userRegisterDTO); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(userRegisterDTO); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user := models.User{
		Username: userRegisterDTO.Username,
		Email:    userRegisterDTO.Email,
		Password: utils.HashPassword(userRegisterDTO.Password),
	}

	existingUser, err := h.UserRepo.FindByUsername(user.Username)
	if err != nil {
		log.Printf("error finding user: %s", err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}
	if existingUser != nil {
		return c.JSON(http.StatusBadRequest, "Username already in use")
	}

	if err := h.UserRepo.Create(&user); err != nil {
		log.Printf("error creating user: %s", err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) Login(c echo.Context) error {
	userLoginDTO := new(UserLoginDTO)
	if err := c.Bind(userLoginDTO); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, err := h.UserRepo.FindByUsername(userLoginDTO.Username)
	if err != nil {
		log.Printf("error finding user: %s", err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	if user == nil || !utils.CheckPasswordHash(userLoginDTO.Password, user.Password) {
		return c.JSON(http.StatusUnauthorized, "Invalid username or password")
	}

	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		log.Printf("failed to generate token: %s", err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

func (h *UserHandler) GetUser(c echo.Context) error {
	username := c.Param("username")
	authUsername := c.Get("username")
	log.Printf("Auth Username: %s", authUsername)
	if authUsername != username {
		return c.JSON(http.StatusForbidden, "Access denied")
	}

	user, err := h.UserRepo.FindByUsername(username)
	if err != nil {
		log.Printf("error finding user: %s", err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}
	if user == nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	username := c.Param("username")
	authUsername := c.Get("username")
	if authUsername != username {
		return c.JSON(http.StatusForbidden, "Access denied")
	}

	if err := h.UserRepo.Delete(username); err != nil {
		log.Printf("error deleting user: %s", err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *UserHandler) SearchUsers(c echo.Context) error {
	keyword := c.QueryParam("keyword")
	if keyword == "" {
		return c.JSON(http.StatusBadRequest, "Keyword is required")
	}

	users, err := h.UserRepo.FindByKeyWord(keyword)
	if err != nil {
		log.Printf("error searching for users: %s", err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	username := c.Param("username")
	authUsername := c.Get("username")
	if authUsername != username {
		return c.JSON(http.StatusForbidden, "Access denied")
	}

	userUpdateDTO := new(UserUpdateDTO)
	if err := c.Bind(userUpdateDTO); err != nil {
		return c.JSON(http.StatusBadRequest, "Error binding user data")
	}

	user, err := h.UserRepo.FindByUsername(username)
	if err != nil {
		log.Printf("error fetching user: %s", err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}
	if user == nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	if userUpdateDTO.Password != "" {
		user.Password = utils.HashPassword(userUpdateDTO.Password)
	}
	user.Email = userUpdateDTO.Email

	if err := h.UserRepo.Update(user); err != nil {
		log.Printf("error updating user: %s", err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	users, err := h.UserRepo.FindAll()
	if err != nil {
		log.Printf("error fetching users: %s", err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	userList := make([]map[string]string, 0)

	for _, user := range users {
		userData := map[string]string{
			"username": user.Username,
			"email":    user.Email,
		}
		userList = append(userList, userData)
	}

	return c.JSON(http.StatusOK, userList)
}
