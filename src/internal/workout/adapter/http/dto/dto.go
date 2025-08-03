package dto

import "github.com/alielmi98/go-hexa-workout/internal/workout/core/usecase/dto"

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
