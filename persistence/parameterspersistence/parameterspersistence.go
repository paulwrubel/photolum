package parameterspersistence

import (
	"context"
	"fmt"

	"github.com/paulwrubel/photolum/config"
	"github.com/paulwrubel/photolum/enumeration/filetype"
	"github.com/paulwrubel/photolum/tracing/shading"
	"github.com/sirupsen/logrus"
)

type Parameters struct {
	ParametersName           string
	ImageWidth               int32
	ImageHeight              int32
	FileType                 filetype.FileType
	GammaCorrection          float64
	TextureGamma             float64
	UseScalingTruncation     bool
	SamplesPerRound          int32
	RoundCount               int32
	TileWidth                int32
	TileHeight               int32
	MaxBounces               int32
	UseBVH                   bool
	BackgroundColorMagnitude float64
	BackgroundColor          *shading.Color
	TMin                     float64
	TMax                     float64
}

var entity = "parameters"

func Save(plData *config.PhotolumData, baseLog *logrus.Entry, parameters *Parameters) error {
	event := "save"
	log := baseLog.WithFields(logrus.Fields{
		"entity": entity,
		"event":  event,
	})
	log.Trace("database event initiated")

	tag, err := plData.DB.Exec(context.Background(), fmt.Sprintf(`
		INSERT INTO parameters (
			parameters_name,
			image_width,
			image_height,
			file_type,
			gamma_correction,
			texture_gamma,
			use_scaling_truncation,
			samples_per_round,
			round_count,
			tile_width,
			tile_height,
			max_bounces,
			use_bvh,
			background_color_magnitude,
			background_color,
			t_min,
			t_max
		) VALUES (%s,%d,%d,%s,%f,%f,%t,%d,%d,%d,%d,%d,%t,%f,{%f,%f,%f},%f,%f)`,
		parameters.ParametersName,
		parameters.ImageWidth,
		parameters.ImageHeight,
		parameters.FileType,
		parameters.GammaCorrection,
		parameters.TextureGamma,
		parameters.UseScalingTruncation,
		parameters.SamplesPerRound,
		parameters.RoundCount,
		parameters.TileWidth,
		parameters.TileHeight,
		parameters.MaxBounces,
		parameters.UseBVH,
		parameters.BackgroundColorMagnitude,
		parameters.BackgroundColor.Red,
		parameters.BackgroundColor.Green,
		parameters.BackgroundColor.Blue,
		parameters.TMin,
		parameters.TMax,
	))
	if err != nil || tag.RowsAffected() != 1 {
		return err
	}

	log.Trace("database event completed")
	return nil
}

func Get(plData *config.PhotolumData, baseLog *logrus.Entry, parametersName string) (*Parameters, error) {
	event := "get"
	log := baseLog.WithFields(logrus.Fields{
		"entity": entity,
		"event":  event,
	})
	log.Trace("database event initiated")

	parameters := &Parameters{
		BackgroundColor: &shading.Color{},
	}
	err := plData.DB.QueryRow(context.Background(), fmt.Sprintf(`
		SELECT 
			parameters_name,
			image_width,
			image_height,
			file_type,
			gamma_correction,
			texture_gamma,
			use_scaling_truncation,
			samples_per_round,
			round_count,
			tile_width,
			tile_height,
			max_bounces,
			use_bvh,
			background_color_magnitude,
			background_color[1],
			background_color[2],
			background_color[3],
			t_min,
			t_max
		FROM parameters
		WHERE parameters_name = %s`, parametersName)).Scan(
		&parameters.ParametersName,
		&parameters.ImageWidth,
		&parameters.ImageHeight,
		&parameters.FileType,
		&parameters.GammaCorrection,
		&parameters.TextureGamma,
		&parameters.UseScalingTruncation,
		&parameters.SamplesPerRound,
		&parameters.RoundCount,
		&parameters.TileWidth,
		&parameters.TileHeight,
		&parameters.MaxBounces,
		&parameters.UseBVH,
		&parameters.BackgroundColorMagnitude,
		&parameters.BackgroundColor.Red,
		&parameters.BackgroundColor.Green,
		&parameters.BackgroundColor.Blue,
		&parameters.TMin,
		&parameters.TMax,
	)
	if err != nil {
		return nil, err
	}

	log.Trace("database event completed")
	return parameters, nil
}

func Update(plData *config.PhotolumData, baseLog *logrus.Entry, parameters *Parameters) error {
	return nil
}

func Delete(plData *config.PhotolumData, baseLog *logrus.Entry, parameters *Parameters) error {
	return nil
}

func DoesExist(plData *config.PhotolumData, baseLog *logrus.Entry, parametersName string) (bool, error) {
	event := "exist"
	log := baseLog.WithFields(logrus.Fields{
		"entity": entity,
		"event":  event,
	})
	log.Trace("database event initiated")

	var count int
	err := plData.DB.QueryRow(context.Background(), fmt.Sprintf(`
		SELECT count(*)
		FROM parameters
		WHERE parameters_name = %s`, parametersName)).Scan(&count)
	if err != nil {
		return false, err
	}

	log.Trace("database event completed")
	return count == 1, nil
}
