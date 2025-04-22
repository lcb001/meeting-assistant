package database

import (
	"meetingagent/models"
)

func GetTodosByMeetingID(meetingID string) ([]models.Todo, error) {
	var todos []models.Todo
	err := DB.Where("meetingID = ?", meetingID).Find(&todos).Error
	return todos, err
}
