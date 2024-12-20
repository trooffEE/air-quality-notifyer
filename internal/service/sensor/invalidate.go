package sensor

import (
	"air-quality-notifyer/internal/db/models"
	"database/sql"
	"errors"
	"fmt"
)

func (s *Service) invalidateSensor(incomingSensor AirqualitySensorScriptScrapped, districtId int64) {
	savedSensor, err := s.repo.FindSensorByApiId(incomingSensor.Id)
	if errors.Is(err, sql.ErrNoRows) {
		dbModel := models.AirqualitySensor{
			ApiId:      incomingSensor.Id,
			DistrictId: districtId,
			Address:    incomingSensor.Address,
			Lat:        incomingSensor.Lat,
			Lon:        incomingSensor.Lon,
		}
		err = s.repo.SaveSensor(dbModel)
		if err != nil {
			fmt.Printf("Failed to save air quality sensor: %v\n", err)
		}
	}

	if err != nil {
		fmt.Println(err)
	}

	//TODO
	//ids, err := s.repo.GetAllApiIds()
	//if err != nil {
	//	fmt.Println(err)
	//}
}
