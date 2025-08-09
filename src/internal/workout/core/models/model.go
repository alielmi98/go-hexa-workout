package models

import (
	"database/sql"
	"time"

	"github.com/alielmi98/go-hexa-workout/constants"
	"gorm.io/gorm"
)

type Workout struct {
	Id          int    `gorm:"primarykey"`
	UserId      int    `gorm:"not null"`
	Name        string `gorm:"type:string;size:100;not null"`
	Description string `gorm:"type:string;size:255;null"`
	Comments    string `gorm:"type:string;size:255;null"`

	CreatedAt  time.Time      `gorm:"type:TIMESTAMP with time zone;not null"`
	ModifiedAt sql.NullTime   `gorm:"type:TIMESTAMP with time zone;null"`
	DeletedAt  sql.NullTime   `gorm:"type:TIMESTAMP with time zone;null"`
	CreatedBy  int            `gorm:"not null"`
	ModifiedBy *sql.NullInt64 `gorm:"null"`
	DeletedBy  *sql.NullInt64 `gorm:"null"`
}

type WorkoutExercise struct {
	Id          int     `gorm:"primarykey"`
	WorkoutId   int     `gorm:"not null"`
	Name        string  `gorm:"type:string;size:100;not null"`
	Description string  `gorm:"type:string;size:255;null"`
	Repetitions int     `gorm:"not null"`
	Sets        int     `gorm:"not null"`
	Weight      float64 `gorm:"not null"`

	CreatedBy  int            `gorm:"not null"`
	ModifiedBy *sql.NullInt64 `gorm:"null"`
	DeletedBy  *sql.NullInt64 `gorm:"null"`
	CreatedAt  time.Time      `gorm:"type:TIMESTAMP with time zone;not null"`
	ModifiedAt sql.NullTime   `gorm:"type:TIMESTAMP with time zone;null"`
	DeletedAt  sql.NullTime   `gorm:"type:TIMESTAMP with time zone;null"`
}

type ScheduledWorkouts struct {
	Id            int       `gorm:"primarykey" `
	WorkoutId     int       `gorm:"not null" `
	ScheduledTime time.Time `gorm:"type:TIMESTAMP with time zone;not null"`
	Status        string    `gorm:"type:string;size:20;not null"`

	CreatedAt  time.Time      `gorm:"type:TIMESTAMP with time zone;not null"`
	ModifiedAt sql.NullTime   `gorm:"type:TIMESTAMP with time zone;null"`
	DeletedAt  sql.NullTime   `gorm:"type:TIMESTAMP with time zone;null"`
	CreatedBy  int            `gorm:"not null" `
	ModifiedBy *sql.NullInt64 `gorm:"null" `
	DeletedBy  *sql.NullInt64 `gorm:"null" `
}

type WorkoutReport struct {
	Id        int    `gorm:"primarykey"`
	WorkoutId int    `gorm:"not null"`
	UserId    int    `gorm:"not null"`
	Details   string `gorm:"type:string;size:255;not null"`

	CreatedAt  time.Time      `gorm:"type:TIMESTAMP with time zone;not null"`
	ModifiedAt sql.NullTime   `gorm:"type:TIMESTAMP with time zone;null"`
	DeletedAt  sql.NullTime   `gorm:"type:TIMESTAMP with time zone;null"`
	CreatedBy  int            `gorm:"not null"`
	ModifiedBy *sql.NullInt64 `gorm:"null"`
	DeletedBy  *sql.NullInt64 `gorm:"null"`
}

func (m *Workout) BeforeCreate(tx *gorm.DB) (err error) {
	value := tx.Statement.Context.Value(constants.UserIdKey)
	var userId = -1
	if value != nil {
		userId = int(value.(float64))
	}
	m.CreatedAt = time.Now().UTC()
	m.CreatedBy = userId
	return
}
func (m *Workout) BeforeUpdate(tx *gorm.DB) (err error) {
	value := tx.Statement.Context.Value(constants.UserIdKey)
	var userId = &sql.NullInt64{Valid: false}
	if value != nil {
		userId = &sql.NullInt64{Valid: true, Int64: int64(value.(float64))}
	}
	m.ModifiedAt = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	m.ModifiedBy = userId
	return
}
func (m *Workout) BeforeDelete(tx *gorm.DB) (err error) {
	value := tx.Statement.Context.Value(constants.UserIdKey)
	var userId = &sql.NullInt64{Valid: false}
	if value != nil {
		userId = &sql.NullInt64{Valid: true, Int64: int64(value.(float64))}
	}
	m.DeletedAt = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	m.DeletedBy = userId
	return
}

// GORM hooks for ScheduledWorkouts
func (m *ScheduledWorkouts) BeforeCreate(tx *gorm.DB) (err error) {
	value := tx.Statement.Context.Value(constants.UserIdKey)
	var userId = -1
	if value != nil {
		userId = int(value.(float64))
	}
	m.CreatedAt = time.Now().UTC()
	m.CreatedBy = userId
	return
}

func (m *ScheduledWorkouts) BeforeUpdate(tx *gorm.DB) (err error) {
	value := tx.Statement.Context.Value(constants.UserIdKey)
	var userId = &sql.NullInt64{Valid: false}
	if value != nil {
		userId = &sql.NullInt64{Valid: true, Int64: int64(value.(float64))}
	}
	m.ModifiedAt = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	m.ModifiedBy = userId
	return
}

func (m *ScheduledWorkouts) BeforeDelete(tx *gorm.DB) (err error) {
	value := tx.Statement.Context.Value(constants.UserIdKey)
	var userId = &sql.NullInt64{Valid: false}
	if value != nil {
		userId = &sql.NullInt64{Valid: true, Int64: int64(value.(float64))}
	}
	m.DeletedAt = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	m.DeletedBy = userId
	return
}

// GORM hooks for WorkoutExercise
func (m *WorkoutExercise) BeforeCreate(tx *gorm.DB) (err error) {
	value := tx.Statement.Context.Value(constants.UserIdKey)
	var userId = -1
	if value != nil {
		userId = int(value.(float64))
	}
	m.CreatedAt = time.Now().UTC()
	m.CreatedBy = userId
	return
}

func (m *WorkoutExercise) BeforeUpdate(tx *gorm.DB) (err error) {
	value := tx.Statement.Context.Value(constants.UserIdKey)
	var userId = &sql.NullInt64{Valid: false}
	if value != nil {
		userId = &sql.NullInt64{Valid: true, Int64: int64(value.(float64))}
	}
	m.ModifiedAt = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	m.ModifiedBy = userId
	return
}

func (m *WorkoutExercise) BeforeDelete(tx *gorm.DB) (err error) {
	value := tx.Statement.Context.Value(constants.UserIdKey)
	var userId = &sql.NullInt64{Valid: false}
	if value != nil {
		userId = &sql.NullInt64{Valid: true, Int64: int64(value.(float64))}
	}
	m.DeletedAt = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	m.DeletedBy = userId
	return
}

// GORM hooks for WorkoutReport
func (m *WorkoutReport) BeforeCreate(tx *gorm.DB) (err error) {
	value := tx.Statement.Context.Value(constants.UserIdKey)
	var userId = -1
	if value != nil {
		userId = int(value.(float64))
	}
	m.CreatedAt = time.Now().UTC()
	m.CreatedBy = userId
	return
}

func (m *WorkoutReport) BeforeUpdate(tx *gorm.DB) (err error) {
	value := tx.Statement.Context.Value(constants.UserIdKey)
	var userId = &sql.NullInt64{Valid: false}
	if value != nil {
		userId = &sql.NullInt64{Valid: true, Int64: int64(value.(float64))}
	}
	m.ModifiedAt = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	m.ModifiedBy = userId
	return
}

func (m *WorkoutReport) BeforeDelete(tx *gorm.DB) (err error) {
	value := tx.Statement.Context.Value(constants.UserIdKey)
	var userId = &sql.NullInt64{Valid: false}
	if value != nil {
		userId = &sql.NullInt64{Valid: true, Int64: int64(value.(float64))}
	}
	m.DeletedAt = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	m.DeletedBy = userId
	return
}
