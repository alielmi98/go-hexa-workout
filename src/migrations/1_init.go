package migrations

import (
	"log"

	"github.com/alielmi98/go-hexa-workout/constants"
	user_models "github.com/alielmi98/go-hexa-workout/internal/user/core/models"
	workout_models "github.com/alielmi98/go-hexa-workout/internal/workout/core/models"
	"github.com/alielmi98/go-hexa-workout/pkg/db"

	"gorm.io/gorm"
)

func Up_1() {
	database := db.GetDb()

	createTables(database)
}

func createTables(database *gorm.DB) {
	tables := []interface{}{}

	// Account
	tables = addNewTable(database, user_models.User{}, tables)

	// Workout
	tables = addNewTable(database, workout_models.Workout{}, tables)
	tables = addNewTable(database, workout_models.WorkoutExercise{}, tables)
	tables = addNewTable(database, workout_models.ScheduledWorkouts{}, tables)
	tables = addNewTable(database, workout_models.WorkoutReport{}, tables)

	err := database.Migrator().CreateTable(tables...)
	if err != nil {
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Migration, err.Error())
	}
	log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Migration, "tables created")
}

func addNewTable(database *gorm.DB, model interface{}, tables []interface{}) []interface{} {
	if !database.Migrator().HasTable(model) {
		tables = append(tables, model)
	}
	return tables
}

func Down_1() {

}
