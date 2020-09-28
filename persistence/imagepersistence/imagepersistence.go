package imagepersistence

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"time"

	"github.com/google/uuid"
	"github.com/paulwrubel/photolum/config"
)

type Image struct {
	SceneID           string
	FileType          string
	createdTimestamp  time.Time
	modifiedTimestamp time.Time
	accessedTimestamp time.Time
	ImageData         image.Image
}

func Save(plData *config.PhotolumData, img *Image) (string, error) {
	fmt.Println("Creating image row in DB...")
	if img.SceneID == "" {
		newSceneID, err := uuid.NewRandom()
		if err != nil {
			return "", err
		}
		img.SceneID = newSceneID.String()
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*60)
	defer cancelFunc()
	tx, err := plData.DB.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	stmt, err := tx.Prepare(`INSERT INTO image (scene_id, file_type, image_data) VALUES (?, ?, ?)`)
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	var imgBytes []byte
	imgBytesBuffer := bytes.NewBuffer(imgBytes)
	switch img.FileType {
	case "png":
		err = png.Encode(imgBytesBuffer, img.ImageData)
	case "jpeg":
		err = jpeg.Encode(imgBytesBuffer, img.ImageData, nil)
	default:
		err = fmt.Errorf("invalid image file to save in db: %s", img.FileType)
	}
	if err != nil {
		return "", err
	}

	_, err = stmt.Exec(img.SceneID, img.FileType, imgBytesBuffer.Bytes())
	if err != nil {
		return "", err
	}
	err = tx.Commit()
	if err != nil {
		return "", err
	}

	fmt.Println("Successfully created image row in DB")
	return img.SceneID, nil
}

func Get(plData *config.PhotolumData, sceneID string, fileType string) (*Image, error) {
	fmt.Println("Retrieving image row in DB...")
	// first, update accessed timestamp
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	tx, err := plData.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	timestampStmt, err := tx.Prepare(`UPDATE image SET accessed_timestamp = datetime() WHERE scene_id = ? AND file_type = ?`)
	if err != nil {
		return nil, err
	}
	defer timestampStmt.Close()
	_, err = timestampStmt.Exec(sceneID, fileType)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// now retrieve row
	ctx, cancelFunc = context.WithTimeout(context.Background(), time.Second*60)
	defer cancelFunc()
	tx, err = plData.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	stmt, err := tx.Prepare(`SELECT created_timestamp, modified_timestamp, accessed_timestamp, image_data FROM image WHERE scene_id = ? AND file_type = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var createdTimestamp time.Time
	var modifiedTimestamp time.Time
	var accessedTimestamp time.Time
	var imgBytes []byte
	err = stmt.QueryRow(sceneID, fileType).Scan(&createdTimestamp, &modifiedTimestamp, &accessedTimestamp, &imgBytes)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	imgBytesReader := bytes.NewReader(imgBytes)
	var img image.Image
	switch fileType {
	case "png":
		img, err = png.Decode(imgBytesReader)
	case "jpeg":
		img, err = jpeg.Decode(imgBytesReader)
	default:
		img = nil
		err = fmt.Errorf("invalid image file type found in db: %s", fileType)
	}
	if err != nil {
		return nil, err
	}

	retrievedImageRow := &Image{
		SceneID:           sceneID,
		FileType:          fileType,
		createdTimestamp:  createdTimestamp,
		modifiedTimestamp: modifiedTimestamp,
		accessedTimestamp: accessedTimestamp,
		ImageData:         img,
	}
	fmt.Println("Successfully retrieved image row in DB")
	return retrievedImageRow, nil
}

func GetEncoded(plData *config.PhotolumData, sceneID string, fileType string) ([]byte, error) {
	fmt.Println("Retrieving encoded image in DB...")
	// first, update accessed timestamp
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	tx, err := plData.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	timestampStmt, err := tx.Prepare(`UPDATE image SET accessed_timestamp = datetime() WHERE scene_id = ? AND file_type = ?`)
	if err != nil {
		return nil, err
	}
	defer timestampStmt.Close()
	_, err = timestampStmt.Exec(sceneID, fileType)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// now retrieve row
	ctx, cancelFunc = context.WithTimeout(context.Background(), time.Second*60)
	defer cancelFunc()
	tx, err = plData.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	stmt, err := tx.Prepare(`SELECT image_data FROM image WHERE scene_id = ? AND file_type = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var imgBytes []byte
	err = stmt.QueryRow(sceneID, fileType).Scan(&imgBytes)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	if fileType != "png" && fileType != "jpeg" {
		err = fmt.Errorf("invalid image file type found in db: %s", fileType)
		return nil, err
	}

	fmt.Println("Successfully retrieved encoded image in DB")
	return imgBytes, nil
}

func GetAll(plData *config.PhotolumData) ([]*Image, error) {
	fmt.Println("Retrieving all images row in DB...")
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*60)
	defer cancelFunc()
	tx, err := plData.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	stmt, err := tx.Prepare(`SELECT scene_id, file_type, created_timestamp, modified_timestamp, accessed_timestamp, image_data FROM image`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var sceneID string
	var fileType string
	var createdTimestamp time.Time
	var modifiedTimestamp time.Time
	var accessedTimestamp time.Time
	var imgBytes []byte

	totalImageRows := []*Image{}

	rows, err := stmt.Query()
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&sceneID, &fileType, &createdTimestamp, &modifiedTimestamp, &accessedTimestamp, &imgBytes)
		if err != nil {
			return nil, err
		}
		imgBytesReader := bytes.NewReader(imgBytes)
		var img image.Image
		switch fileType {
		case "png":
			img, err = png.Decode(imgBytesReader)
		case "jpeg":
			img, err = jpeg.Decode(imgBytesReader)
		default:
			img = nil
			err = fmt.Errorf("invalid image file type found in db: %s", fileType)
		}
		if err != nil {
			return nil, err
		}
		totalImageRows = append(totalImageRows, &Image{
			SceneID:           sceneID,
			FileType:          fileType,
			createdTimestamp:  createdTimestamp,
			modifiedTimestamp: modifiedTimestamp,
			accessedTimestamp: accessedTimestamp,
			ImageData:         img,
		})
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully retrieved all image rows in DB...")
	return totalImageRows, nil
}

func Update(plData *config.PhotolumData, image *Image) error {
	fmt.Println("Updating image row in DB...")
	// update row (and timestamp)
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*60)
	defer cancelFunc()
	tx, err := plData.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`UPDATE image SET modified_timestamp = datetime(), image_data = ? WHERE scene_id = ? AND file_type = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	var imgBytes []byte
	imgBytesBuffer := bytes.NewBuffer(imgBytes)
	err = png.Encode(imgBytesBuffer, image.ImageData)
	_, err = stmt.Exec(imgBytesBuffer.Bytes(), image.SceneID, image.FileType)
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
	stmt, err := tx.Prepare(`DELETE FROM image WHERE scene_id = ? AND file_type = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(image.SceneID, image.FileType)
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

func DoesExist(plData *config.PhotolumData, sceneID string, fileType string) (bool, error) {
	fmt.Println("Checking image existance in DB...")
	// first, update accessed timestamp
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	tx, err := plData.DB.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	stmt, err := tx.Prepare(`SELECT count(*) FROM image WHERE scene_id = ? AND file_type = ?`)
	if err != nil {
		return false, err
	}
	defer stmt.Close()
	var count int
	err = stmt.QueryRow(sceneID, fileType).Scan(&count)
	if err != nil {
		return false, err
	}
	err = tx.Commit()
	if err != nil {
		return false, err
	}
	return count != 0, nil
}
