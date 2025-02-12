package helper

import (
	"rentjoy/internal/dto/homepage"
	"rentjoy/internal/models"
)

func ACTModelToDTO(models []models.ActivityType) []homepage.Activity {
	var dtoActivities []homepage.Activity
	for _, activity := range models {
		dtoActivities = append(dtoActivities, homepage.Activity{
			ID:           activity.ID,
			ActivityName: activity.ActivityName,
		})
	}

	return dtoActivities
}
