package encoding

import (
	"bytes"
	"image/jpeg"
	"image/png"

	"github.com/paulwrubel/photolum/config"
	"github.com/paulwrubel/photolum/enumeration/filetype"
	"github.com/paulwrubel/photolum/enumeration/renderstatus"
	"github.com/paulwrubel/photolum/persistence/renderpersistence.go"
	"github.com/sirupsen/logrus"
)

func RunWorker(plData *config.PhotolumData, log *logrus.Entry, renderName string, encodingChan <-chan *config.TracingPayload) {
	log.Debug("running encoding worker")

	var err error
	for {
		tracingPayload, active := <-encodingChan
		if active {
			log.Debug("encoding new image to render")
			buffer := new(bytes.Buffer)
			switch tracingPayload.FileType {
			case filetype.PNG:
				err = png.Encode(buffer, tracingPayload.Image)
			case filetype.JPEG:
				err = jpeg.Encode(buffer, tracingPayload.Image, nil)
			}
			if err != nil {
				log.WithError(err).Error("error encoding image")
				renderpersistence.UpdateRenderStatus(plData, log, renderName, renderstatus.Error)
			}
			err = renderpersistence.UpdateImageData(plData, log, renderName, buffer.Bytes())
			if err != nil {
				log.WithError(err).Error("error updating render")
				renderpersistence.UpdateRenderStatus(plData, log, renderName, renderstatus.Error)
			}
			log.Debug("image encoding finished")
		} else {
			log.Debug("encoder signalled to exit")
			break
		}
	}
	log.Debug("closing encoding worker")
}
