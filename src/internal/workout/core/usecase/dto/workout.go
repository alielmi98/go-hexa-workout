package dto

type CreateWorkoutRequest struct {
	Name        string
	Description string
	Comments    string
	UserId      int
}

type UpdateWorkoutRequest struct {
	Name        string
	Description string
	Comments    string
}

type WorkoutResponse struct {
	Id          int
	UserId      int
	Name        string
	Description string
	Comments    string
}
