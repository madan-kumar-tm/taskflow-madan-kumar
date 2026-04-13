package services

import "taskflow/config"

type StatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type AssigneeCount struct {
	AssigneeID *string `json:"assignee_id"`
	Count      int64   `json:"count"`
}

type ProjectStats struct {
	StatusCounts   []StatusCount   `json:"status_counts"`
	AssigneeCounts []AssigneeCount `json:"assignee_counts"`
}

func GetProjectStats(projectID string) (ProjectStats, error) {
	stats := ProjectStats{}

	statusRows := []StatusCount{}
	if err := config.DB.Raw(`
		SELECT status, COUNT(*) AS count
		FROM tasks
		WHERE project_id = $1
		GROUP BY status
	`, projectID).Scan(&statusRows).Error; err != nil {
		return stats, err
	}

	assigneeRows := []AssigneeCount{}
	if err := config.DB.Raw(`
		SELECT assignee_id, COUNT(*) AS count
		FROM tasks
		WHERE project_id = $1
		GROUP BY assignee_id
	`, projectID).Scan(&assigneeRows).Error; err != nil {
		return stats, err
	}

	stats.StatusCounts = statusRows
	stats.AssigneeCounts = assigneeRows
	return stats, nil
}
