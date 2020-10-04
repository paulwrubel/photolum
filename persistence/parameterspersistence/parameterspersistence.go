package parameterspersistence

import (
	"context"

	"github.com/paulwrubel/photolum/config"
	"github.com/sirupsen/logrus"
)

type Parameters struct {
	ParametersName           string
	ImageWidth               uint32
	ImageHeight              uint32
	FileType                 string
	GammaCorrection          float64
	UseScalingTruncation     bool
	SamplesPerRound          uint32
	RoundCount               uint32
	TileWidth                uint32
	TileHeight               uint32
	MaxBounces               uint32
	UseBVH                   bool
	BackgroundColorMagnitude float64
	BackgroundColor          [3]float64
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

	tag, err := plData.DB.Exec(context.Background(), `
		INSERT INTO parameters (
			parameters_name,
			image_width,
			image_height,
			file_type,
			gamma_correction,
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
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`,
		parameters.ParametersName,
		parameters.ImageWidth,
		parameters.ImageHeight,
		parameters.FileType,
		parameters.GammaCorrection,
		parameters.UseScalingTruncation,
		parameters.SamplesPerRound,
		parameters.RoundCount,
		parameters.TileWidth,
		parameters.TileHeight,
		parameters.MaxBounces,
		parameters.UseBVH,
		parameters.BackgroundColorMagnitude,
		parameters.BackgroundColor,
		parameters.TMin,
		parameters.TMax,
	)
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

	parameters := &Parameters{}
	err := plData.DB.QueryRow(context.Background(), `
		SELECT 
			parameters_name,
			image_width,
			image_height,
			file_type,
			gamma_correction,
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
		FROM parameters
		WHERE parameters_name = $1`, parametersName).Scan(
		&parameters.ParametersName,
		&parameters.ImageWidth,
		&parameters.ImageHeight,
		&parameters.FileType,
		&parameters.GammaCorrection,
		&parameters.UseScalingTruncation,
		&parameters.SamplesPerRound,
		&parameters.RoundCount,
		&parameters.TileWidth,
		&parameters.TileHeight,
		&parameters.MaxBounces,
		&parameters.UseBVH,
		&parameters.BackgroundColorMagnitude,
		&parameters.BackgroundColor,
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
	err := plData.DB.QueryRow(context.Background(), `
		SELECT count(*)
		FROM parameters
		WHERE parameters_name = $1`, parametersName).Scan(&count)
	if err != nil {
		return false, err
	}

	log.Trace("database event completed")
	return count == 1, nil
}
