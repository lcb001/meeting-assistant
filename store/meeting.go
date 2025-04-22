package store

import "meetingagent/models"

var Meetings []models.Meeting

type Summary struct {
	Content string
	Todos   []string
}

var Summaries map[string]Summary
