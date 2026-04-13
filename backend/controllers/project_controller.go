package controllers

import (
	"log/slog"
	"net/http"
	"taskflow/services"
	"taskflow/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createProjectRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type updateProjectRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

func CreateProject(c *gin.Context) {
	userID := c.GetString("user_id")
	var req createProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Name == "" {
		fields := map[string]string{}
		if req.Name == "" {
			fields["name"] = "is required"
		}
		if len(fields) == 0 {
			fields["body"] = "Project name is required"
		}
		utils.ValidationFailed(c, fields)
		return
	}

	project, err := services.CreateProject(userID, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create project"})
		return
	}

	c.JSON(http.StatusCreated, project)
}

func ListProjects(c *gin.Context) {
	userID := c.GetString("user_id")
	pg, err := utils.ParsePagination(c)
	if err != nil {
		utils.ValidationFailed(c, map[string]string{"page/limit": err.Error()})
		return
	}

	projects, err := services.ListProjects(userID, pg.Limit, pg.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch projects"})
		return
	}
	c.JSON(http.StatusOK, projects)
}

func GetProject(c *gin.Context) {
	userID := c.GetString("user_id")
	projectID := c.Param("id")

	allowed, err := services.ProjectAccessible(userID, projectID)
	if err != nil {
		slog.Error("project access check failed", "project_id", projectID, "user_id", userID, "err", err)
		utils.NotFound(c)
		return
	}
	if !allowed {
		utils.NotFound(c)
		return
	}

	project, err := services.GetProject(userID, projectID)
	if err != nil || project.ID == "" {
		utils.NotFound(c)
		return
	}

	status := c.Query("status")
	assignee := c.Query("assignee")
	if status != "" && !isValidProjectTaskStatus(status) {
		utils.ValidationFailed(c, map[string]string{"status": "must be one of todo, in_progress, done"})
		return
	}
	if assignee != "" {
		if _, err := uuid.Parse(assignee); err != nil {
			utils.ValidationFailed(c, map[string]string{"assignee": "must be a valid uuid"})
			return
		}
	}

	pg, err := utils.ParsePagination(c)
	if err != nil {
		utils.ValidationFailed(c, map[string]string{"page/limit": err.Error()})
		return
	}

	tasks, err := services.ListTasks(projectID, nullableString(status), nullableString(assignee), pg.Limit, pg.Offset)
	if err != nil {
		slog.Error("list tasks for project failed", "project_id", projectID, "status", status, "assignee", assignee, "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"project": project,
		"tasks":   tasks,
	})
}

func DeleteProject(c *gin.Context) {
	userID := c.GetString("user_id")
	projectID := c.Param("id")
	err := services.DeleteProject(userID, projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete project"})
		return
	}
	c.Status(http.StatusNoContent)
}

func UpdateProject(c *gin.Context) {
	userID := c.GetString("user_id")
	projectID := c.Param("id")

	var req updateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	project, err := services.UpdateProject(userID, projectID, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update project"})
		return
	}
	if project.ID == "" {
		utils.Forbidden(c)
		return
	}
	c.JSON(http.StatusOK, project)
}

func nullableString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func ProjectStats(c *gin.Context) {
	userID := c.GetString("user_id")
	projectID := c.Param("id")

	allowed, err := services.ProjectAccessible(userID, projectID)
	if err != nil || !allowed {
		utils.NotFound(c)
		return
	}

	stats, err := services.GetProjectStats(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func isValidProjectTaskStatus(status string) bool {
	return status == "todo" || status == "in_progress" || status == "done"
}
