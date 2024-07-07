package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/raeinsoltani/gorello/back/models"
	"github.com/raeinsoltani/gorello/back/repository"
)

type WorkspaceHandler struct {
	WorkspaceRepo         repository.Workspace
	UserWorkspaceRoleRepo repository.UserWorkspaceRole
	UserRepo              repository.User
}

func NewWorkspaceHandler(workspaceRepo repository.Workspace, userWorkspaceRoleRepo repository.UserWorkspaceRole, userRepo repository.User) *WorkspaceHandler {
	return &WorkspaceHandler{
		WorkspaceRepo:         workspaceRepo,
		UserWorkspaceRoleRepo: userWorkspaceRoleRepo,
		UserRepo:              userRepo,
	}
}

type WorkspaceCreateDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *WorkspaceHandler) CreateWorkspace(c echo.Context) error {
	authUsername, ok := c.Get("username").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, "User not authenticated")
	}

	WorkspaceCreateDTO := new(WorkspaceCreateDTO)
	if err := c.Bind(WorkspaceCreateDTO); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(WorkspaceCreateDTO); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	workspace := models.Workspace{
		Name:        WorkspaceCreateDTO.Name,
		Description: WorkspaceCreateDTO.Description,
	}

	err := h.WorkspaceRepo.Create(&workspace)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	user, err := h.UserRepo.FindByUsername(authUsername)
	userWorkspaceRole := models.UserWorkspaceRole{
		User_id:      user.ID,
		Workspace_id: workspace.ID,
		Role:         1,
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	err = h.UserWorkspaceRoleRepo.Create(&userWorkspaceRole)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, userWorkspaceRole)
}

func (h *WorkspaceHandler) GetWorkspaces(c echo.Context) error {
	authUsername, ok := c.Get("username").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, "User not authenticated")
	}

	user, err := h.UserRepo.FindByUsername(authUsername)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userWorkspaceRoles, err := h.UserWorkspaceRoleRepo.FindByUserID(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	workspaces := make([]models.Workspace, 0)
	for _, userWorkspaceRole := range userWorkspaceRoles {
		workspace, err := h.WorkspaceRepo.FindByID(userWorkspaceRole.Workspace_id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		workspaces = append(workspaces, *workspace)
	}

	return c.JSON(http.StatusOK, workspaces)
}

func (h *WorkspaceHandler) GetWorkspaceDescription(c echo.Context) error {
	authUsername, ok := c.Get("username").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, "User not authenticated")
	}

	workspaceId, err := strconv.ParseUint(c.Param("workspaceId"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, err := h.UserRepo.FindByUsername(authUsername)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userWorkspaceRoles, err := h.UserWorkspaceRoleRepo.FindByUserID(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	for _, userWorkspaceRole := range userWorkspaceRoles {
		if userWorkspaceRole.Workspace_id == uint(workspaceId) {
			workspace, err := h.WorkspaceRepo.FindByID(uint(workspaceId))
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
			return c.JSON(http.StatusOK, workspace.Description)
		}
	}

	return c.JSON(http.StatusForbidden, "Access denied to the workspace")
}

func (h *WorkspaceHandler) UpdateWorkspace(c echo.Context) error {
	_, ok := c.Get("username").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, "User not authenticated")
	}

	workspaceId, err := strconv.ParseUint(c.Param("workspaceId"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	workspace, err := h.WorkspaceRepo.FindByID(uint(workspaceId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if workspace == nil {
		return c.JSON(http.StatusNotFound, "Workspace not found")
	}

	workspaceUpdateDTO := new(WorkspaceCreateDTO)
	if err := c.Bind(workspaceUpdateDTO); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(workspaceUpdateDTO); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if workspaceUpdateDTO.Name != "" {
		workspace.Name = workspaceUpdateDTO.Name
	}
	if workspaceUpdateDTO.Description != "" {
		workspace.Description = workspaceUpdateDTO.Description
	}

	err = h.WorkspaceRepo.Update(workspace)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, workspace)
}

func (h *WorkspaceHandler) DeleteWorkspace(c echo.Context) error {
	_, ok := c.Get("username").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, "User not authenticated")
	}

	workspaceId, err := strconv.ParseUint(c.Param("workspaceId"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	workspace, err := h.WorkspaceRepo.FindByID(uint(workspaceId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if workspace == nil {
		return c.JSON(http.StatusNotFound, "Workspace not found")
	}

	err = h.WorkspaceRepo.Delete(workspace.Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusNoContent, nil)
}
