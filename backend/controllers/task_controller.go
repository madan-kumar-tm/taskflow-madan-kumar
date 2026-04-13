package controllers

import (
	"log/slog"
	"net/http"
	"taskflow/services"
	"taskflow/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createTaskRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
	Priority    *string `json:"priority"`
	AssigneeID  *string `json:"assignee_id"`
	DueDate     *string `json:"due_date"`
}

type updateTaskRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
	Priority    *string `json:"priority"`
	AssigneeID  *string `json:"assignee_id"`
	DueDate     *string `json:"due_date"`
}

func CreateTask(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		utils.ValidationFailed(c, map[string]string{"project_id": "is required"})
		return
	}
	creatorID := c.GetString("user_id")
	var req createTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Title == "" {
		fields := map[string]string{}
		if req.Title == "" {
			fields["title"] = "is required"
		}
		if len(fields) == 0 {
			fields["body"] = "Task title is required"
		}
		utils.ValidationFailed(c, fields)
		return
	}

	priority := "medium"
	if req.Priority != nil && *req.Priority != "" {
		priority = *req.Priority
	}
	if !isValidPriority(priority) {
		utils.ValidationFailed(c, map[string]string{"priority": "must be one of low, medium, high"})
		return
	}

	status := "todo"
	if req.Status != nil && *req.Status != "" {
		status = normalizeTaskStatus(*req.Status)
	}
	if !isValidStatus(status) {
		utils.ValidationFailed(c, map[string]string{"status": "must be one of todo, in_progress, done"})
		return
	}
	if req.AssigneeID != nil && *req.AssigneeID != "" && !isValidUUID(*req.AssigneeID) {
		utils.ValidationFailed(c, map[string]string{"assignee_id": "must be a valid uuid"})
		return
	}

	task, err := services.CreateTask(projectID, creatorID, req.Title, req.Description, status, priority, req.AssigneeID, req.DueDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create task"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func ListTasks(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		utils.ValidationFailed(c, map[string]string{"project_id": "is required"})
		return
	}

	userID := c.GetString("user_id") // 👈 logged-in user

	status := c.Query("status")
	if status != "" {
		status = normalizeTaskStatus(status)
		if !isValidStatus(status) {
			utils.ValidationFailed(c, map[string]string{"status": "must be one of todo, in_progress, done"})
			return
		}
	}

	pg, err := utils.ParsePagination(c)
	if err != nil {
		utils.ValidationFailed(c, map[string]string{"page/limit": err.Error()})
		return
	}

	//ONLY fetch tasks assigned to logged-in user
	tasks, err := services.ListTasks(projectID, nullable(status), &userID, pg.Limit, pg.Offset)
	if err != nil {
		slog.Error("list tasks failed", "project_id", projectID, "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func UpdateTask(c *gin.Context) {
	taskID := c.Param("id")
	var req updateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationFailed(c, map[string]string{"body": "invalid json"})
		return
	}

	status := req.Status
	if status != nil {
		normalized := normalizeTaskStatus(*status)
		status = &normalized
		if !isValidStatus(*status) {
			utils.ValidationFailed(c, map[string]string{"status": "must be one of todo, in_progress, done"})
			return
		}
	}
	if req.Priority != nil && *req.Priority != "" && !isValidPriority(*req.Priority) {
		utils.ValidationFailed(c, map[string]string{"priority": "must be one of low, medium, high"})
		return
	}
	if req.AssigneeID != nil && *req.AssigneeID != "" && !isValidUUID(*req.AssigneeID) {
		utils.ValidationFailed(c, map[string]string{"assignee_id": "must be a valid uuid"})
		return
	}

	task, err := services.UpdateTask(taskID, status, req.AssigneeID, req.Title, req.Description, req.Priority, req.DueDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update task"})
		return
	}
	if task.ID == "" {
		utils.NotFound(c)
		return
	}
	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	userID := c.GetString("user_id")
	taskID := c.Param("id")

	ok, err := services.CanDeleteTask(userID, taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not verify permissions"})
		return
	}
	if !ok {
		utils.Forbidden(c)
		return
	}

	if err := services.DeleteTask(taskID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete task"})
		return
	}
	c.Status(http.StatusNoContent)
}

func nullable(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func normalizeTaskStatus(status string) string {
	switch status {
	case "in progress":
		return "in_progress"
	case "completed":
		return "done"
	default:
		return status
	}
}

func isValidStatus(status string) bool {
	return status == "todo" || status == "in_progress" || status == "done"
}

func isValidPriority(priority string) bool {
	return priority == "low" || priority == "medium" || priority == "high"
}

func isValidUUID(value string) bool {
	_, err := uuid.Parse(value)
	return err == nil
}
