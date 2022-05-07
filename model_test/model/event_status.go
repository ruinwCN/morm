package model

import (
	"time"
)

type Status struct {
	Id         int64     `morm:"id" json:"id" morm_table:"event_status"`
	Name       string    `morm:"name" json:"name"`
	CreateTime time.Time `morm:"create_time" morm_date:"true" json:"create_time"`
}

func (s *Status) TableName() string {
	return "event_status"
}
