package camerapersistence

import (
	"context"
	"fmt"

	"github.com/paulwrubel/photolum/config"
	"github.com/paulwrubel/photolum/tracing/geometry"
	"github.com/sirupsen/logrus"
)

type Camera struct {
	CameraName     string
	EyeLocation    geometry.Point
	TargetLocation geometry.Point
	UpVector       geometry.Vector
	VerticalFOV    float64
	AspectRatio    float64
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
			aspect_ratio,
			aperture,
			focus_distance
		) VALUES (?,?,?,?,?,?,?,?)`,
		camera.CameraName,
		[]float64{camera.EyeLocation.X, camera.EyeLocation.Y, camera.EyeLocation.Z},
		[]float64{camera.TargetLocation.X, camera.TargetLocation.Y, camera.TargetLocation.Z},
		[]float64{camera.UpVector.X, camera.UpVector.Y, camera.UpVector.Z},
		camera.VerticalFOV,
		camera.AspectRatio,
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

	camera := &Camera{
		EyeLocation:    geometry.Point{},
		TargetLocation: geometry.Point{},
		UpVector:       geometry.Vector{},
	}
	err := plData.DB.QueryRow(context.Background(), `
		SELECT 
			camera_name,
			eye_location[1],
			eye_location[2],
			eye_location[3],
			target_location[1],
			target_location[2],
			target_location[3],
			up_vector[1],
			up_vector[2],
			up_vector[3],
			vertical_fov,
			aspect_ratio,
			aperture,
			focus_distance
		FROM cameras
		WHERE camera_name = ?`, cameraName).Scan(
		&camera.CameraName,
		&camera.EyeLocation.X,
		&camera.EyeLocation.Y,
		&camera.EyeLocation.Z,
		&camera.TargetLocation.X,
		&camera.TargetLocation.Y,
		&camera.TargetLocation.Z,
		&camera.UpVector.X,
		&camera.UpVector.Y,
		&camera.UpVector.Z,
		&camera.VerticalFOV,
		&camera.AspectRatio,
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
	err := plData.DB.QueryRow(context.Background(), fmt.Sprintf(`
		SELECT count(*)
		FROM cameras
		WHERE camera_name = '%s'`, cameraName)).Scan(&count)
	if err != nil {
		return false, err
	}

	log.Trace("database event completed")
	return count == 1, nil
}
