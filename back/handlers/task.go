package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/raeinsoltani/gorello/back/models"
	"github.com/raeinsoltani/gorello/back/repository"
)

type TaskHandler struct {
	TaskRepo              repository.Task
	UserWorkspaceRoleRepo repository.UserWorkspaceRole
	UserRepo              repository.User
}

func NewTaskHandler(taskRepo repository.Task, userWorkspaceRoleRepo repository.UserWorkspaceRole, userRepo repository.User) *TaskHandler {
	return &TaskHandler{
		TaskRepo:              taskRepo,
		UserWorkspaceRoleRepo: userWorkspaceRoleRepo,
		UserRepo:              userRepo,
	}
}

type TaskCreateDTO struct {
	Title          string `json:"name"`
	Description    string `json:"description"`
	Status         uint   `json:"status"`
	Estimated_time string `json:"estimated_time"`
	Actual_time    string `json:"actual_time"`
	Due_date       string `json:"due_date"`
	Priority       uint   `json:"priority"`
	Assignee_id    uint   `json:"assignee_id"`
	Workspace_id   uint   `json:"workspace_id"`
	Image_url      string `json:"image_url"`
}

func (t *TaskCreateDTO) Validate() error {
	var errFields []string

	if t.Title == "" {
		errFields = append(errFields, "name")
	}

	if t.Workspace_id == 0 {
		errFields = append(errFields, "workspace_id")
	}

	if len(errFields) > 0 {
		return fmt.Errorf("the following fields cannot be empty: %v", strings.Join(errFields, ", "))
	}
	return nil
}

func (h *TaskHandler) CreateTask(c echo.Context) error {
	authUsername, ok := c.Get("username").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, "User not authenticated")
	}

	taskCreateDTO := new(TaskCreateDTO)
	if err := c.Bind(taskCreateDTO); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(taskCreateDTO); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	task := &models.Task{
		Title:          taskCreateDTO.Title,
		Description:    taskCreateDTO.Description,
		Status:         taskCreateDTO.Status,
		Estimated_time: taskCreateDTO.Estimated_time,
		Actual_time:    taskCreateDTO.Actual_time,
		Due_date:       taskCreateDTO.Due_date,
		Priority:       taskCreateDTO.Priority,
		Assignee_id:    taskCreateDTO.Assignee_id,
		Workspace_id:   taskCreateDTO.Workspace_id,
		Image_url:      taskCreateDTO.Image_url,
	}

	workspace_id, err := strconv.ParseUint(c.Param("workspaceId"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	task.Workspace_id = uint(workspace_id)

	user, err := h.UserRepo.FindByUsername(authUsername)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	workspaces, err := h.UserWorkspaceRoleRepo.FindByUserID(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	found := false
	for _, workspace := range workspaces {
		if uint64(workspace.User_id) == uint64(user.ID) {
			found = true
			break
		}
	}

	if !found {
		return c.JSON(http.StatusBadRequest, "Access denied to this workspace")
	}

	err = h.TaskRepo.Create(task)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetTasks(c echo.Context) error {
	workspace_id, err := strconv.ParseUint(c.Param("workspaceId"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	tasks, err := h.TaskRepo.FindByWorkspaceID(uint(workspace_id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) GetTask(c echo.Context) error {
	task_id, err := strconv.ParseUint(c.Param("taskId"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	task, err := h.TaskRepo.FindByID(uint(task_id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(c echo.Context) error {
	taskId, err := strconv.ParseUint(c.Param("taskId"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	taskUpdateDTO := new(TaskCreateDTO)
	if err := c.Bind(taskUpdateDTO); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(taskUpdateDTO); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	task, err := h.TaskRepo.FindByID(uint(taskId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	task.Title = taskUpdateDTO.Title
	task.Description = taskUpdateDTO.Description
	task.Status = taskUpdateDTO.Status
	task.Estimated_time = taskUpdateDTO.Estimated_time
	task.Actual_time = taskUpdateDTO.Actual_time
	task.Priority = taskUpdateDTO.Priority
	task.Image_url = taskUpdateDTO.Image_url

	err = h.TaskRepo.Update(task)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(c echo.Context) error {
	taskId, err := strconv.ParseUint(c.Param("taskId"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = h.TaskRepo.Delete(uint(taskId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusNoContent, fmt.Sprintf("Task with id %d deleted", taskId))
}
