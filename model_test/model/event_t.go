package model

import "time"

type EventT struct {
	Id         int64     `morm:"id" json:"id" morm_table:"event_t"`
	Ndecimal   float32   `morm:"ndecimal" json:"ndecimal"`
	Nbigint    int64     `morm:"nbigint" json:"nbigint"`
	Nchar      byte      `morm:"nchar" json:"nchar"`
	Ntext      string    `morm:"ntext" json:"ntext"`
	Nfloat     float32   `morm:"nfloat" json:"nfloat"`
	CreateTime time.Time `morm:"create_time" morm_date:"true" json:"create_time"`
}

func (s *EventT) TableName() string {
	return "event_t"
}
