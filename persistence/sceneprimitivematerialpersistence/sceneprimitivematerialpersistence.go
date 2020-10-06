package sceneprimitivematerialpersistence

import (
	"context"

	"github.com/paulwrubel/photolum/config"
	"github.com/sirupsen/logrus"
)

type ScenePrimitiveMaterial struct {
	SceneName     string
	PrimitiveName string
	MaterialName  string
}

var entity = "sceneprimitivematerial"

func Save(plData *config.PhotolumData, baseLog *logrus.Entry, scenePrimitiveMaterial *ScenePrimitiveMaterial) error {
	event := "save"
	log := baseLog.WithFields(logrus.Fields{
		"entity": entity,
		"event":  event,
	})
	log.Trace("database event initiated")

	tag, err := plData.DB.Exec(context.Background(), `
		INSERT INTO scene_primitive_materials (
			scene_name,
			primitive_name,
			material_name
		) VALUES ($1,$2,$3)`,
		scenePrimitiveMaterial.SceneName,
		scenePrimitiveMaterial.PrimitiveName,
		scenePrimitiveMaterial.MaterialName,
	)
	if err != nil || tag.RowsAffected() != 1 {
		return err
	}

	log.Trace("database event completed")
	return nil
}

func Get(plData *config.PhotolumData, baseLog *logrus.Entry, sceneName, primitiveName, materialName string) (*ScenePrimitiveMaterial, error) {
	event := "get"
	log := baseLog.WithFields(logrus.Fields{
		"entity": entity,
		"event":  event,
	})
	log.Trace("database event initiated")

	scenePrimitiveMaterial := &ScenePrimitiveMaterial{}
	err := plData.DB.QueryRow(context.Background(), `
		SELECT 
			scene_name,
			primitive_name,
			material_name
		FROM scene_primitive_materials
		WHERE 
			scene_name = $1 AND
			primitive_name = $2 AND
			material_name = $3`, sceneName, primitiveName, materialName).Scan(
		&scenePrimitiveMaterial.SceneName,
		&scenePrimitiveMaterial.PrimitiveName,
		&scenePrimitiveMaterial.MaterialName,
	)
	if err != nil {
		return nil, err
	}

	log.Trace("database event completed")
	return scenePrimitiveMaterial, nil
}

func GetAllInScene(plData *config.PhotolumData, baseLog *logrus.Entry, sceneName string) ([]*ScenePrimitiveMaterial, error) {
	event := "get"
	log := baseLog.WithFields(logrus.Fields{
		"entity": entity,
		"event":  event,
	})
	log.Trace("database event initiated")

	scenePrimitiveMaterials := []*ScenePrimitiveMaterial{}
	rows, err := plData.DB.Query(context.Background(), `
		SELECT 
			scene_name,
			primitive_name,
			material_name
		FROM scene_primitive_materials
		WHERE 
			scene_name = $1`, sceneName)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		spm := &ScenePrimitiveMaterial{}
		err := rows.Scan(
			&spm.SceneName,
			&spm.PrimitiveName,
			&spm.MaterialName,
		)
		if err != nil {
			return nil, err
		}
		scenePrimitiveMaterials = append(scenePrimitiveMaterials, spm)
	}

	log.Trace("database event completed")
	return scenePrimitiveMaterials, nil
}

func Update(plData *config.PhotolumData, baseLog *logrus.Entry, scenePrimitiveMaterial *ScenePrimitiveMaterial) error {
	return nil
}

func Delete(plData *config.PhotolumData, baseLog *logrus.Entry, scenePrimitiveMaterial *ScenePrimitiveMaterial) error {
	return nil
}

func DoesExist(plData *config.PhotolumData, baseLog *logrus.Entry, sceneName, primitiveName, materialName string) (bool, error) {
	event := "exist"
	log := baseLog.WithFields(logrus.Fields{
		"entity": entity,
		"event":  event,
	})
	log.Trace("database event initiated")

	var count int
	err := plData.DB.QueryRow(context.Background(), `
		SELECT count(*)
		FROM scene_primitive_materials
		WHERE 
			scene_name = $1 AND
			primitive_name = $2 AND
			material_name = $3`, sceneName, primitiveName, materialName).Scan(&count)
	if err != nil {
		return false, err
	}

	log.Trace("database event completed")
	return count == 1, nil
}
