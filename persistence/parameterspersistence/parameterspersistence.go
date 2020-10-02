package parameterspersistence

import (
	"github.com/google/uuid"
	"github.com/paulwrubel/photolum/config"
	"github.com/paulwrubel/photolum/enumeration/filetype"
	"github.com/paulwrubel/photolum/tracing/shading"
	"github.com/sirupsen/logrus"
)

type Parameters struct {
	ParametersID             uuid.UUID
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
	BackgroundColor          shading.Color
	TMin                     float64
	TMax                     float64
}

var entity = "parameters"

func Save(plData *config.PhotolumData, baseLog logrus.Entry, parameters *Parameters) error {
	event := "save"
	log := baseLog.WithFields(logrus.Fields{
		"entity": entity,
		"event":  event,
	})
	log.Trace("database event initiated")

	log.Trace("database event completed")
	return nil
}
