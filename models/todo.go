package models

type Todo struct {
	ID          string `gorm:"primaryKey"`
	MeetingID   string
	Title       string
	Description string
	CompletedAt *string // 可为 null，用指针
	CreatedAt   string
	UpdatedAt   string
	List        string
	Assignee    *string // 可为 null，用指针
}
