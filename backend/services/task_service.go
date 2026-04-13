package services

import (
	"taskflow/config"
	"taskflow/models"
)

func CreateTask(projectID, creatorID, title string, description *string, status string, priority string, assigneeID *string, dueDate *string) (models.Task, error) {
	var task models.Task
	query := `
		INSERT INTO tasks (title, description, status, priority, project_id, assignee_id, due_date, creator_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, title, description, status, priority, project_id, assignee_id, creator_id,
		          CAST(due_date AS TEXT) AS due_date, CAST(created_at AS TEXT) AS created_at, CAST(updated_at AS TEXT) AS updated_at
	`
	err := config.DB.Raw(query, title, description, status, priority, projectID, assigneeID, dueDate, creatorID).Scan(&task).Error
	return task, err
}

func ListTasks(projectID string, status *string, assigneeID *string, limit, offset int) ([]models.Task, error) {
	tasks := []models.Task{}
	q := config.DB.Table("tasks").
		Select("id, title, description, status, priority, project_id, assignee_id, creator_id, CAST(due_date AS TEXT) AS due_date, CAST(created_at AS TEXT) AS created_at, CAST(updated_at AS TEXT) AS updated_at").
		Where("project_id = ?", projectID)

	if status != nil && *status != "" {
		q = q.Where("status = ?", *status)
	}
	if assigneeID != nil && *assigneeID != "" {
		q = q.Where("assignee_id = ?", *assigneeID)
	}

	err := q.Order("created_at DESC").Limit(limit).Offset(offset).Scan(&tasks).Error
	return tasks, err
}

func UpdateTask(taskID string, status *string, assigneeID *string, title *string, description *string, priority *string, dueDate *string) (models.Task, error) {
	var task models.Task
	query := `
		UPDATE tasks
		SET
			status = COALESCE($2, status),
			assignee_id = COALESCE($3, assignee_id),
			title = COALESCE($4, title),
			description = COALESCE($5, description),
			priority = COALESCE($6, priority),
			due_date = COALESCE($7, due_date),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, title, description, status, priority, project_id, assignee_id, creator_id,
		          CAST(due_date AS TEXT) AS due_date, CAST(created_at AS TEXT) AS created_at, CAST(updated_at AS TEXT) AS updated_at
	`
	err := config.DB.Raw(query, taskID, status, assigneeID, title, description, priority, dueDate).Scan(&task).Error
	return task, err
}

func DeleteTask(taskID string) error {
	return config.DB.Exec(`DELETE FROM tasks WHERE id = $1`, taskID).Error
}

// CanDeleteTask returns true if user is project owner or task creator.
func CanDeleteTask(userID, taskID string) (bool, error) {
	var count int64
	query := `
		SELECT COUNT(*)
		FROM tasks t
		JOIN projects p ON p.id = t.project_id
		WHERE t.id = $1 AND (p.owner_id = $2 OR t.creator_id = $2)
	`
	err := config.DB.Raw(query, taskID, userID).Scan(&count).Error
	return count > 0, err
}
