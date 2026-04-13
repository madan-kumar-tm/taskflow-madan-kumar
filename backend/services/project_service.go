package services

import (
	"taskflow/config"
	"taskflow/models"
)

func CreateProject(ownerID, name string, description *string) (models.Project, error) {
	var project models.Project
	query := `
		INSERT INTO projects (name, description, owner_id)
		VALUES ($1, $2, $3)
		RETURNING id, name, description, owner_id, CAST(created_at AS TEXT) AS created_at
	`
	err := config.DB.Raw(query, name, description, ownerID).Scan(&project).Error
	return project, err
}

func ListProjects(ownerID string, limit, offset int) ([]models.Project, error) {
	projects := []models.Project{}
	query := `
		SELECT DISTINCT p.id, p.name, p.description, p.owner_id, CAST(p.created_at AS TEXT) AS created_at
		FROM projects p
		LEFT JOIN tasks t ON t.project_id = p.id
		WHERE p.owner_id = $1 OR t.assignee_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	err := config.DB.Raw(query, ownerID, limit, offset).Scan(&projects).Error
	return projects, err
}

func GetProject(ownerID, projectID string) (models.Project, error) {
	var project models.Project
	query := `
		SELECT id, name, description, owner_id, CAST(created_at AS TEXT) AS created_at
		FROM projects
		WHERE id = $1
	`
	err := config.DB.Raw(query, projectID).Scan(&project).Error
	return project, err
}

func DeleteProject(ownerID, projectID string) error {
	query := `DELETE FROM projects WHERE id = $1 AND owner_id = $2`
	return config.DB.Exec(query, projectID, ownerID).Error
}

func UpdateProject(ownerID, projectID string, name *string, description *string) (models.Project, error) {
	var project models.Project
	query := `
		UPDATE projects
		SET
			name = COALESCE($3, name),
			description = COALESCE($4, description)
		WHERE id = $1 AND owner_id = $2
		RETURNING id, name, description, owner_id, CAST(created_at AS TEXT) AS created_at
	`
	err := config.DB.Raw(query, projectID, ownerID, name, description).Scan(&project).Error
	return project, err
}

// ProjectAccessible returns true if the user owns the project or has a task in it.
func ProjectAccessible(userID, projectID string) (bool, error) {
	var count int64
	check := `
		SELECT COUNT(*) FROM projects p
		LEFT JOIN tasks t ON t.project_id = p.id
		WHERE p.id = $1 AND (p.owner_id = $2 OR t.assignee_id = $2)
	`
	err := config.DB.Raw(check, projectID, userID).Scan(&count).Error
	return count > 0, err
}
