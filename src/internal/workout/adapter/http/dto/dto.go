package dto

import (
	"time"

	"github.com/alielmi98/go-hexa-workout/internal/workout/core/usecase/dto"
)

type CreateWorkoutRequest struct {
	Name        string `json:"name" binding:"required,min=3"`
	Description string `json:"description"`
	Comments    string `json:"comments"`
}

type UpdateWorkoutRequest struct {
	Name        string `json:"name" binding:"required,min=3"`
	Description string `json:"description"`
	Comments    string `json:"comments"`
}
type WorkoutResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Comments    string `json:"comments"`
}

func ToWorkoutResponse(from dto.WorkoutResponse) WorkoutResponse {
	return WorkoutResponse{
		Id:          from.Id,
		Name:        from.Name,
		Description: from.Description,
		Comments:    from.Comments,
	}
}
func ToUpdateWorkoutRequest(from UpdateWorkoutRequest) dto.UpdateWorkoutRequest {
	return dto.UpdateWorkoutRequest{
		Name:        from.Name,
		Description: from.Description,
		Comments:    from.Comments,
	}
}

func ToCreateWorkoutRequest(from CreateWorkoutRequest) dto.CreateWorkoutRequest {
	return dto.CreateWorkoutRequest{
		Name:        from.Name,
		Description: from.Description,
		Comments:    from.Comments,
	}
}

// WorkoutExercise
type CreateWorkoutExerciseRequest struct {
	WorkoutId   int     `json:"workout_id" binding:"required"`
	Name        string  `json:"name" binding:"required,min=3"`
	Description string  `json:"description"`
	Reps        int     `json:"reps" binding:"required"`
	Sets        int     `json:"sets" binding:"required"`
	Weight      float64 `json:"weight" binding:"required"`
}

type UpdateWorkoutExerciseRequest struct {
	WorkoutId   int     `json:"workout_id" binding:"required"`
	Name        string  `json:"name" binding:"required,min=3"`
	Description string  `json:"description"`
	Reps        int     `json:"reps" binding:"required"`
	Sets        int     `json:"sets" binding:"required"`
	Weight      float64 `json:"weight" binding:"required"`
}
type WorkoutExerciseResponse struct {
	Id          int     `json:"id"`
	WorkoutId   int     `json:"workout_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Reps        int     `json:"reps"`
	Sets        int     `json:"sets"`
	Weight      float64 `json:"weight"`
}

func ToWorkoutExerciseResponse(from dto.WorkoutExerciseResponse) WorkoutExerciseResponse {
	return WorkoutExerciseResponse{
		Id:          from.Id,
		WorkoutId:   from.WorkoutId,
		Name:        from.Name,
		Description: from.Description,
		Reps:        from.Repetitions,
		Sets:        from.Sets,
		Weight:      from.Weight,
	}
}
func ToCreateWorkoutExerciseRequest(from CreateWorkoutExerciseRequest) dto.CreateWorkoutExerciseRequest {
	return dto.CreateWorkoutExerciseRequest{
		Name:        from.Name,
		WorkoutId:   from.WorkoutId,
		Description: from.Description,
		Repetitions: from.Reps,
		Sets:        from.Sets,
		Weight:      from.Weight,
	}
}

func ToUpdateWorkoutExerciseRequest(from UpdateWorkoutExerciseRequest) dto.UpdateWorkoutExerciseRequest {
	return dto.UpdateWorkoutExerciseRequest{
		Name:        from.Name,
		WorkoutId:   from.WorkoutId,
		Description: from.Description,
		Repetitions: from.Reps,
		Sets:        from.Sets,
		Weight:      from.Weight,
	}
}

// ScheduledWorkoutss
type ScheduledWorkoutsResponse struct {
	Id            int    `json:"id"`
	WorkoutId     int    `json:"workout_id"`
	ScheduledTime string `json:"scheduled_time"` //ScheduledTime
	Status        string `json:"status"`
}

type CreateScheduledWorkoutsRequest struct {
	WorkoutId     int       `json:"workout_id" binding:"required"`
	ScheduledTime time.Time `json:"scheduled_time" binding:"required"`
	Status        string    `json:"status" binding:"required"`
}

type UpdateScheduledWorkoutsRequest struct {
	ScheduledTime time.Time `json:"scheduled_time" binding:"required"`
	Status        string    `json:"status" binding:"required"`
}

func ToScheduledWorkoutsResponse(from dto.ScheduledWorkoutsResponse) ScheduledWorkoutsResponse {
	return ScheduledWorkoutsResponse{
		Id:            from.Id,
		WorkoutId:     from.WorkoutId,
		Status:        from.Status,
		ScheduledTime: from.ScheduledTime,
	}
}

func ToCreateScheduledWorkoutsRequest(from CreateScheduledWorkoutsRequest) dto.CreateScheduledWorkoutsRequest {
	return dto.CreateScheduledWorkoutsRequest{
		WorkoutId:     from.WorkoutId,
		ScheduledTime: from.ScheduledTime,
		Status:        from.Status,
	}
}

func ToUpdateScheduledWorkoutsRequest(from UpdateScheduledWorkoutsRequest) dto.UpdateScheduledWorkoutsRequest {
	return dto.UpdateScheduledWorkoutsRequest{
		ScheduledTime: from.ScheduledTime,
		Status:        from.Status,
	}
}

// WorkoutReport
type WorkoutReportResponse struct {
	Id        int    `json:"id"`
	WorkoutId int    `json:"workout_id"`
	UserId    int    `json:"user_id"`
	Details   string `json:"details"`
}

type CreateWorkoutReportRequest struct {
	WorkoutId int    `json:"workout_id" binding:"required"`
	Details   string `json:"details" binding:"required"`
}

type UpdateWorkoutReportRequest struct {
	Details   string `json:"details" binding:"required"`
	WorkoutId int    `json:"workout_id" binding:"required"`
}

func ToWorkoutReportResponse(from dto.WorkoutReportResponse) WorkoutReportResponse {
	return WorkoutReportResponse{
		Id:        from.Id,
		WorkoutId: from.WorkoutId,
		UserId:    from.UserId,
		Details:   from.Details,
	}
}

func ToCreateWorkoutReportRequest(from CreateWorkoutReportRequest) dto.CreateWorkoutReportRequest {
	return dto.CreateWorkoutReportRequest{
		WorkoutId: from.WorkoutId,
		Details:   from.Details,
	}
}

func ToUpdateWorkoutReportRequest(from UpdateWorkoutReportRequest) dto.UpdateWorkoutReportRequest {
	return dto.UpdateWorkoutReportRequest{
		Details:   from.Details,
		WorkoutId: from.WorkoutId,
	}
}
