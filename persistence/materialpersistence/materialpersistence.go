package materialpersistence

import (
	"context"

	"github.com/paulwrubel/photolum/config"
	"github.com/sirupsen/logrus"
)

type Material struct {
	MaterialName           string
	MaterialType           string
	ReflectanceTextureName *string
	EmittanceTextureName   *string
	Fuzziness              *float64
	RefractiveIndex        *float64
}

var entity = "material"

func Save(plData *config.PhotolumData, baseLog *logrus.Entry, material *Material) error {
	event := "save"
	log := baseLog.WithFields(logrus.Fields{
		"entity": entity,
		"event":  event,
	})
	log.Trace("database event initiated")

	tag, err := plData.DB.Exec(context.Background(), `
		INSERT INTO materials (
			material_name,
			material_type,
			reflectance_texture_name,
			emittance_texture_name,
			fuzziness,
			refractive_index
		) VALUES ($1,$2,$3,$4,$5,$6)`,
		material.MaterialName,
		material.MaterialType,
		material.ReflectanceTextureName,
		material.EmittanceTextureName,
		material.Fuzziness,
		material.RefractiveIndex,
	)
	if err != nil || tag.RowsAffected() != 1 {
		return err
	}

	log.Trace("database event completed")
	return nil
}

func Get(plData *config.PhotolumData, baseLog *logrus.Entry, materialName string) (*Material, error) {
	event := "get"
	log := baseLog.WithFields(logrus.Fields{
		"entity": entity,
		"event":  event,
	})
	log.Trace("database event initiated")

	material := &Material{}
	err := plData.DB.QueryRow(context.Background(), `
		SELECT 
			material_name,
			material_type,
			reflectance_texture_name,
			emittance_texture_name,
			fuzziness,
			refractive_index
		FROM materials
		WHERE material_name = $1`, materialName).Scan(
		&material.MaterialName,
		&material.MaterialType,
		&material.ReflectanceTextureName,
		&material.EmittanceTextureName,
		&material.Fuzziness,
		&material.RefractiveIndex,
	)
	if err != nil {
		return nil, err
	}

	log.Trace("database event completed")
	return material, nil
}

func Update(plData *config.PhotolumData, baseLog *logrus.Entry, material *Material) error {
	return nil
}

func Delete(plData *config.PhotolumData, baseLog *logrus.Entry, material *Material) error {
	return nil
}

func DoesExist(plData *config.PhotolumData, baseLog *logrus.Entry, materialName string) (bool, error) {
	event := "exist"
	log := baseLog.WithFields(logrus.Fields{
		"entity": entity,
		"event":  event,
	})
	log.Trace("database event initiated")

	var count int
	err := plData.DB.QueryRow(context.Background(), `
		SELECT count(*)
		FROM materials
		WHERE material_name = $1`, materialName).Scan(&count)
	if err != nil {
		return false, err
	}

	log.Trace("database event completed")
	return count == 1, nil
}
