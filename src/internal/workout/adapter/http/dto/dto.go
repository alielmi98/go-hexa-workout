package dto

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
