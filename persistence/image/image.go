package image

import (
	"bytes"
	"context"
	"fmt"
	_image "image"
	"image/png"
	"time"

	"github.com/google/uuid"
	"github.com/paulwrubel/photolum/config"
)

type Image struct {
	SceneID           string
	createdTimestamp  time.Time
	modifiedTimestamp time.Time
	accessedTimestamp time.Time
	ImageData         *_image.RGBA64
}

func Create(plData *config.PhotolumData, image *Image) (string, error) {
	fmt.Println("Creating image row in DB...")
	if image.SceneID == "" {
		newSceneID, err := uuid.NewRandom()
		if err != nil {
			return "", err
		}
		image.SceneID = newSceneID.String()
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	tx, err := plData.DB.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	stmt, err := tx.Prepare(`INSERT INTO image (scene_id, image) VALUES (?, ?)`)
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	_, err = stmt.Exec(image.SceneID, image.ImageData)
	if err != nil {
		return "", err
	}
	err = tx.Commit()
	if err != nil {
		return "", err
	}

	fmt.Println("Successfully created image row in DB")
	return image.SceneID, nil
}

func Retrieve(plData *config.PhotolumData, sceneID string) (*Image, error) {
	fmt.Println("Retrieving image row in DB...")
	// first, update accessed timestamp
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	tx, err := plData.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	timestampStmt, err := tx.Prepare(`UPDATE image SET accessed_timestamp = datetime() WHERE scene_id = ?`)
	if err != nil {
		return nil, err
	}
	defer timestampStmt.Close()
	_, err = timestampStmt.Exec(sceneID)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// now retrieve row
	ctx, cancelFunc = context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	tx, err = plData.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	stmt, err := tx.Prepare(`SELECT created_timestamp, modified_timestamp, accessed_timestamp, image_data FROM image WHERE scene_id = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var createdTimestamp time.Time
	var modifiedTimestamp time.Time
	var accessedTimestamp time.Time
	var imgBytes []byte
	err = stmt.QueryRow(sceneID).Scan(&createdTimestamp, &modifiedTimestamp, &accessedTimestamp, &imgBytes)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	imgBytesReader := bytes.NewReader(imgBytes)
	img, err := png.Decode(imgBytesReader)
	if err != nil {
		return nil, err
	}
	rgba64Image := img.(*_image.RGBA64)

	retrievedImageRow := &Image{
		SceneID:           sceneID,
		createdTimestamp:  createdTimestamp,
		modifiedTimestamp: modifiedTimestamp,
		accessedTimestamp: accessedTimestamp,
		ImageData:         rgba64Image,
	}
	fmt.Println("Successfully retrieved image row in DB")
	return retrievedImageRow, nil
}

func Update(plData *config.PhotolumData, image *Image) error {
	fmt.Println("Updating image row in DB...")
	// update row (and timestamp)
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	tx, err := plData.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`UPDATE image SET modified_timestamp = datetime(), image_data = ? WHERE scene_id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	var imgBytes []byte
	imgBytesBuffer := bytes.NewBuffer(imgBytes)
	err = png.Encode(imgBytesBuffer, image.ImageData)
	_, err = stmt.Exec(imgBytesBuffer.Bytes(), image.SceneID)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	fmt.Println("Successfully updated image row in DB")
	return nil
}

func Delete(plData *config.PhotolumData, image *Image) error {
	fmt.Println("Deleting image row in DB...")
	// update row (and timestamp)
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	tx, err := plData.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`DELETE FROM image WHERE scene_id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(image.SceneID)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	fmt.Println("Successfully deleted image row in DB...")
	return nil
}

func DoesExist(plData *config.PhotolumData, sceneID string) (bool, error) {
	fmt.Println("Checking image existance in DB...")
	// first, update accessed timestamp
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	tx, err := plData.DB.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	stmt, err := tx.Prepare(`SELECT count(*) FROM image WHERE scene_id = ?`)
	if err != nil {
		return false, err
	}
	defer stmt.Close()
	var count int
	err = stmt.QueryRow(sceneID).Scan(&count)
	if err != nil {
		return false, err
	}
	err = tx.Commit()
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

func RetrieveAll(plData *config.PhotolumData) ([]*Image, error) {
	fmt.Println("Retrieving all images row in DB...")
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()
	tx, err := plData.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	stmt, err := tx.Prepare(`SELECT scene_id, created_timestamp, modified_timestamp, accessed_timestamp, image_data FROM image`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var sceneID string
	var createdTimestamp time.Time
	var modifiedTimestamp time.Time
	var accessedTimestamp time.Time
	var imgBytes []byte

	totalImageRows := []*Image{}

	rows, err := stmt.Query(sceneID)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&sceneID, &createdTimestamp, &modifiedTimestamp, &accessedTimestamp, &imgBytes)
		if err != nil {
			return nil, err
		}
		imgBytesReader := bytes.NewReader(imgBytes)
		img, err := png.Decode(imgBytesReader)
		if err != nil {
			return nil, err
		}
		rgba64Image := img.(*_image.RGBA64)
		totalImageRows = append(totalImageRows, &Image{
			SceneID:           sceneID,
			createdTimestamp:  createdTimestamp,
			modifiedTimestamp: modifiedTimestamp,
			accessedTimestamp: accessedTimestamp,
			ImageData:         rgba64Image,
		})
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully retrieved all image rows in DB...")
	return totalImageRows, nil
}
