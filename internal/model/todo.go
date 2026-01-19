package model

import "time"

type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
	StatusClosed     Status = "closed"
)

func (s Status) Icon() string {
	switch s {
	case StatusPending:
		return "○"
	case StatusInProgress:
		return "◐"
	case StatusDone:
		return "●"
	case StatusClosed:
		return "✕"
	default:
		return "○"
	}
}

func (s Status) Next() Status {
	switch s {
	case StatusPending:
		return StatusInProgress
	case StatusInProgress:
		return StatusDone
	case StatusDone:
		return StatusPending
	case StatusClosed:
		return StatusPending
	default:
		return StatusPending
	}
}

func (s Status) ToggleClosed() Status {
	if s == StatusClosed {
		return StatusPending
	}
	return StatusClosed
}

type TimeTag string

const (
	TimeTagNone  TimeTag = ""
	TimeTagToday TimeTag = "today"
	TimeTagWeek  TimeTag = "week"
	TimeTagMonth TimeTag = "month"
)

func (t TimeTag) Label() string {
	switch t {
	case TimeTagToday:
		return "today"
	case TimeTagWeek:
		return "this week"
	case TimeTagMonth:
		return "this month"
	default:
		return ""
	}
}

type Todo struct {
	ID        int64
	Title     string
	Category  string
	Detail    string
	Status    Status
	TimeTag   TimeTag
	DueDate   *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
