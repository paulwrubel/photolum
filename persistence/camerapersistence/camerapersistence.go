package camerapersistence

import (
	"context"

	"github.com/paulwrubel/photolum/config"
	"github.com/sirupsen/logrus"
)

type Camera struct {
	CameraName     string
	EyeLocation    []float64
	TargetLocation []float64
	UpVector       []float64
	VerticalFOV    float64
	Aperture       float64
	FocusDistance  float64
}

var entity = "camera"

func Save(plData *config.PhotolumData, baseLog *logrus.Entry, camera *Camera) error {
	event := "save"
	log := baseLog.WithFields(logrus.Fields{
		"entity": entity,
		"event":  event,
	})
	log.Trace("database event initiated")

	tag, err := plData.DB.Exec(context.Background(), `
		INSERT INTO cameras (
			camera_name,
			eye_location,
			target_location,
			up_vector,
			vertical_fov,
			aperture,
			focus_distance
		) VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		camera.CameraName,
		camera.EyeLocation,
		camera.TargetLocation,
		camera.UpVector,
		camera.VerticalFOV,
		camera.Aperture,
		camera.FocusDistance,
	)
	if err != nil || tag.RowsAffected() != 1 {
		return err
	}

	log.Trace("database event completed")
	return nil
}

func Get(plData *config.PhotolumData, baseLog *logrus.Entry, cameraName string) (*Camera, error) {
	event := "get"
	log := baseLog.WithFields(logrus.Fields{
		"entity": entity,
		"event":  event,
	})
	log.Trace("database event initiated")

	camera := &Camera{}
	err := plData.DB.QueryRow(context.Background(), `
		SELECT 
			camera_name,
			eye_location,
			target_location,
			up_vector,
			vertical_fov,
			aperture,
			focus_distance
		FROM cameras
		WHERE camera_name = $1`, cameraName).Scan(
		&camera.CameraName,
		&camera.EyeLocation,
		&camera.TargetLocation,
		&camera.UpVector,
		&camera.VerticalFOV,
		&camera.Aperture,
		&camera.FocusDistance,
	)
	if err != nil {
		return nil, err
	}

	log.Trace("database event completed")
	return camera, nil
}

func Update(plData *config.PhotolumData, baseLog *logrus.Entry, camera *Camera) error {
	return nil
}

func Delete(plData *config.PhotolumData, baseLog *logrus.Entry, camera *Camera) error {
	return nil
}

func DoesExist(plData *config.PhotolumData, baseLog *logrus.Entry, cameraName string) (bool, error) {
	event := "exist"
	log := baseLog.WithFields(logrus.Fields{
		"entity": entity,
		"event":  event,
	})
	log.Trace("database event initiated")

	var count int
	err := plData.DB.QueryRow(context.Background(), `
		SELECT count(*)
		FROM cameras
		WHERE camera_name = $1`, cameraName).Scan(&count)
	if err != nil {
		return false, err
	}

	log.Trace("database event completed")
	return count == 1, nil
}
