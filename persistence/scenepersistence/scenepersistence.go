package scenepersistence

import (
	"context"

	"github.com/paulwrubel/photolum/config"
	"github.com/sirupsen/logrus"
)

type Scene struct {
	SceneName  string
	CameraName string
}

var entity = "scene"

func Save(plData *config.PhotolumData, baseLog *logrus.Entry, scene *Scene) error {
	event := "save"
	log := baseLog.WithFields(logrus.Fields{
		"entity": entity,
		"event":  event,
	})
	log.Trace("database event initiated")

	tag, err := plData.DB.Exec(context.Background(), `
		INSERT INTO scenes (
			scene_name,
			camera_name
		) VALUES ($1,$2)`,
		scene.SceneName,
		scene.CameraName,
	)
	if err != nil || tag.RowsAffected() != 1 {
		return err
	}

	log.Trace("database event completed")
	return nil
}

func Get(plData *config.PhotolumData, baseLog *logrus.Entry, sceneName string) (*Scene, error) {
	event := "get"
	log := baseLog.WithFields(logrus.Fields{
		"entity": entity,
		"event":  event,
	})
	log.Trace("database event initiated")

	scene := &Scene{}
	err := plData.DB.QueryRow(context.Background(), `
		SELECT 
			scene_name,
			camera_name
		FROM scenes
		WHERE scene_name = $1`, sceneName).Scan(
		&scene.SceneName,
		&scene.CameraName,
	)
	if err != nil {
		return nil, err
	}

	log.Trace("database event completed")
	return scene, nil
}

func Update(plData *config.PhotolumData, baseLog *logrus.Entry, scene *Scene) error {
	return nil
}

func Delete(plData *config.PhotolumData, baseLog *logrus.Entry, scene *Scene) error {
	return nil
}

func DoesExist(plData *config.PhotolumData, baseLog *logrus.Entry, sceneName string) (bool, error) {
	event := "exist"
	log := baseLog.WithFields(logrus.Fields{
		"entity": entity,
		"event":  event,
	})
	log.Trace("database event initiated")

	var count int
	err := plData.DB.QueryRow(context.Background(), `
		SELECT count(*)
		FROM scenes
		WHERE scene_name = $1`, sceneName).Scan(&count)
	if err != nil {
		return false, err
	}

	log.Trace("database event completed")
	return count == 1, nil
}
