package dto

// Workout
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

// WorkoutExercise
type WorkoutExerciseResponse struct {
	Id          int
	WorkoutId   int
	Name        string
	Description string
	Repetitions int
	Sets        int
	Weight      float64
}
type CreateWorkoutExerciseRequest struct {
	WorkoutId   int
	Name        string
	Description string
	Repetitions int
	Sets        int
	Weight      float64
}

type UpdateWorkoutExerciseRequest struct {
	Name        string
	WorkoutId   int
	Description string
	Repetitions int
	Sets        int
	Weight      float64
}
