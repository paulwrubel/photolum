package renderpersistence

import (
	"context"
	"sync"
	"time"

	"github.com/paulwrubel/photolum/config"
	"github.com/paulwrubel/photolum/enumeration/renderstatus"
	"github.com/sirupsen/logrus"
)

type Render struct {
	RenderName      string
	ParametersName  string
	SceneName       string
	RenderStatus    string
	CompletedRounds uint32
	RoundProgress   float64
	StartTimestamp  time.Time
	EndTimestamp    *time.Time
	ImageData       []byte
}

var entity = "render"

func Save(plData *config.PhotolumData, baseLog *logrus.Entry, render *Render) error {
	event := "save"
	log := baseLog.WithFields(logrus.Fields{
		"entity": entity,
		"event":  event,
	})
	log.Trace("database event initiated")

	tag, err := plData.DB.Exec(context.Background(), `
		INSERT INTO renders (
			render_name,
			parameters_name,
			scene_name,
			render_status,
			completed_rounds,
			round_progress,
			start_timestamp,
			end_timestamp,
			image_data
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		render.RenderName,
		render.ParametersName,
		render.SceneName,
		render.RenderStatus,
		render.CompletedRounds,
		render.RoundProgress,
		render.StartTimestamp,
		render.EndTimestamp,
		render.ImageData,
	)
	if err != nil || tag.RowsAffected() != 1 {
		return err
	}

	log.Trace("database event completed")
	return nil
}

func Get(plData *config.PhotolumData, baseLog *logrus.Entry, renderName string) (*Render, error) {
	event := "get"
	log := baseLog.WithFields(logrus.Fields{
		"entity": entity,
		"event":  event,
	})
	log.Trace("database event initiated")

	render := &Render{}
	err := plData.DB.QueryRow(context.Background(), `
		SELECT 
			render_name,
			parameters_name,
			scene_name,
			render_status,
			completed_rounds,
			round_progress,
			start_timestamp,
			end_timestamp,
			image_data
		FROM renders
		WHERE render_name = $1`, renderName).Scan(
		&render.RenderName,
		&render.ParametersName,
		&render.SceneName,
		&render.RenderStatus,
		&render.CompletedRounds,
		&render.RoundProgress,
		&render.StartTimestamp,
		&render.EndTimestamp,
		&render.ImageData,
	)
	if err != nil {
		return nil, err
	}

	log.Trace("database event completed")
	return render, nil
}

func Update(plData *config.PhotolumData, baseLog *logrus.Entry, render *Render) error {
	event := "update"
	log := baseLog.WithFields(logrus.Fields{
		"entity": entity,
		"event":  event,
	})
	log.Trace("database event initiated")

	tag, err := plData.DB.Exec(context.Background(), `
		UPDATE renders 
		SET 
			render_name = $1,
			parameters_name = $2,
			scene_name = $3,
			render_status = $4,
			completed_rounds = $5,
			round_progress = $6,
			start_timestamp = $7,
			end_timestamps = $8,
			image_data = $9
		WHERE render_name = $1`,
		render.RenderName,
		render.ParametersName,
		render.SceneName,
		render.RenderStatus,
		render.CompletedRounds,
		render.RoundProgress,
		render.StartTimestamp,
		render.EndTimestamp,
		render.ImageData,
	)
	if err != nil || tag.RowsAffected() != 1 {
		return err
	}

	log.Trace("database event completed")
	return nil
}

func Delete(plData *config.PhotolumData, baseLog *logrus.Entry, render *Render) error {
	return nil
}

func DoesExist(plData *config.PhotolumData, baseLog *logrus.Entry, renderName string) (bool, error) {
	event := "exist"
	log := baseLog.WithFields(logrus.Fields{
		"entity": entity,
		"event":  event,
	})
	log.Trace("database event initiated")

	var count int
	err := plData.DB.QueryRow(context.Background(), `
		SELECT count(*)
		FROM renders
		WHERE render_name = $1`, renderName).Scan(&count)
	if err != nil {
		return false, err
	}

	log.Trace("database event completed")
	return count == 1, nil
}

func UpdateRenderStatus(plData *config.PhotolumData, baseLog *logrus.Entry, renderName string, renderStatus renderstatus.RenderStatus) error {
	event := "update render_status"
	log := baseLog.WithFields(logrus.Fields{
		"entity":     entity,
		"event":      event,
		"new_status": string(renderStatus),
	})
	log.Trace("database event initiated")

	tag, err := plData.DB.Exec(context.Background(), `
		UPDATE renders 
		SET render_status = $2
		WHERE render_name = $1`,
		renderName,
		string(renderStatus),
	)
	if err != nil || tag.RowsAffected() != 1 {
		return err
	}

	log.Trace("database event completed")
	return nil
}

func UpdateCompletedRounds(plData *config.PhotolumData, baseLog *logrus.Entry, renderName string, completedRounds uint32) error {
	event := "update completed_rounds"
	log := baseLog.WithFields(logrus.Fields{
		"entity": entity,
		"event":  event,
	})
	log.Trace("database event initiated")

	tag, err := plData.DB.Exec(context.Background(), `
		UPDATE renders 
		SET completed_rounds = $2
		WHERE render_name = $1`,
		renderName,
		completedRounds,
	)
	if err != nil || tag.RowsAffected() != 1 {
		return err
	}

	log.Trace("database event completed")
	return nil
}

func UpdateRoundProgress(plData *config.PhotolumData, baseLog *logrus.Entry, renderName string, roundProgress float64, wg *sync.WaitGroup) error {
	//event := "update render_status"
	// log := baseLog.WithFields(logrus.Fields{
	// 	"entity": entity,
	// 	"event":  event,
	// })
	//log.Trace("database event initiated")

	if wg != nil {
		defer wg.Done()
	}

	tag, err := plData.DB.Exec(context.Background(), `
		UPDATE renders 
		SET round_progress = $2
		WHERE render_name = $1`,
		renderName,
		roundProgress,
	)
	if err != nil || tag.RowsAffected() != 1 {
		return err
	}

	//log.Trace("database event completed")
	return nil
}

func UpdateEndTimestamp(plData *config.PhotolumData, baseLog *logrus.Entry, renderName string, endTime *time.Time) error {
	event := "update end_timestamp"
	log := baseLog.WithFields(logrus.Fields{
		"entity": entity,
		"event":  event,
	})
	log.Trace("database event initiated")

	tag, err := plData.DB.Exec(context.Background(), `
		UPDATE renders 
		SET end_timestamp = $2
		WHERE render_name = $1`,
		renderName,
		endTime,
	)
	if err != nil || tag.RowsAffected() != 1 {
		return err
	}

	log.Trace("database event completed")
	return nil
}

func UpdateImageData(plData *config.PhotolumData, baseLog *logrus.Entry, renderName string, imageData []byte) error {
	event := "update image_data"
	log := baseLog.WithFields(logrus.Fields{
		"entity": entity,
		"event":  event,
	})
	log.Trace("database event initiated")

	tag, err := plData.DB.Exec(context.Background(), `
		UPDATE renders 
		SET image_data = $2
		WHERE render_name = $1`,
		renderName,
		imageData,
	)
	if err != nil || tag.RowsAffected() != 1 {
		return err
	}

	log.Trace("database event completed")
	return nil
}
