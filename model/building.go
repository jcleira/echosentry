package model

import "time"

type Building struct {
	ID        uint      `swaggerignore:"true" xorm:"autoincr pk"`
	Name      string    `xorm:"name"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}
