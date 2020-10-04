package texturepersistence

import (
	"context"

	"github.com/paulwrubel/photolum/config"
	"github.com/sirupsen/logrus"
)

type Texture struct {
	TextureName string
	TextureType string
	Color       []float64
	Gamma       *float64
	Magnitude   *float64
	ImageData   []byte
}

var entity = "texture"

func Save(plData *config.PhotolumData, baseLog *logrus.Entry, texture *Texture) error {
	event := "save"
	log := baseLog.WithFields(logrus.Fields{
		"entity": entity,
		"event":  event,
	})
	log.Trace("database event initiated")

	tag, err := plData.DB.Exec(context.Background(), `
		INSERT INTO textures (
			texture_name,
			texture_type,
			color,
			gamma,
			magnitude,
			image_data
		) VALUES ($1,$2,$3,$4,$5,$6)`,
		texture.TextureName,
		texture.TextureType,
		texture.Color,
		texture.Gamma,
		texture.Magnitude,
		texture.ImageData,
	)
	if err != nil || tag.RowsAffected() != 1 {
		return err
	}

	log.Trace("database event completed")
	return nil
}

func Get(plData *config.PhotolumData, baseLog *logrus.Entry, textureName string) (*Texture, error) {
	event := "get"
	log := baseLog.WithFields(logrus.Fields{
		"entity": entity,
		"event":  event,
	})
	log.Trace("database event initiated")

	texture := &Texture{}
	err := plData.DB.QueryRow(context.Background(), `
		SELECT 
			texture_name,
			texture_type,
			color,
			gamma,
			magnitude,
			image_data
		FROM textures
		WHERE texture_name = $1`, textureName).Scan(
		&texture.TextureName,
		&texture.TextureType,
		&texture.Color,
		&texture.Gamma,
		&texture.Magnitude,
		&texture.ImageData,
	)
	if err != nil {
		return nil, err
	}

	log.Trace("database event completed")
	return texture, nil
}

func Update(plData *config.PhotolumData, baseLog *logrus.Entry, texture *Texture) error {
	return nil
}

func Delete(plData *config.PhotolumData, baseLog *logrus.Entry, texture *Texture) error {
	return nil
}

func DoesExist(plData *config.PhotolumData, baseLog *logrus.Entry, textureName string) (bool, error) {
	event := "exist"
	log := baseLog.WithFields(logrus.Fields{
		"entity": entity,
		"event":  event,
	})
	log.Trace("database event initiated")

	var count int
	err := plData.DB.QueryRow(context.Background(), `
		SELECT count(*)
		FROM textures
		WHERE texture_name = $1`, textureName).Scan(&count)
	if err != nil {
		return false, err
	}

	log.Trace("database event completed")
	return count == 1, nil
}
