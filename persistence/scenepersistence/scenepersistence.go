package scenepersistence

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/paulwrubel/photolum/config"
	"github.com/paulwrubel/photolum/config/renderstatus"
)

type Scene struct {
	SceneID           string
	RenderStatus      renderstatus.RenderStatus
	createdTimestamp  time.Time
	modifiedTimestamp time.Time
	accessedTimestamp time.Time
	ImageWidth        int
	ImageHeight       int
	ImageFileTypes    string
}

func Save(plData *config.PhotolumData, scene *Scene) (string, error) {
	fmt.Println("Creating scene row in DB...")
	if scene.SceneID == "" {
		newSceneID, err := uuid.NewRandom()
		if err != nil {
			return "", err
		}
		scene.SceneID = newSceneID.String()
	}
	scene.RenderStatus = renderstatus.Created

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	tx, err := plData.DB.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	stmt, err := tx.Prepare(`
		INSERT INTO scene (
			scene_id, 
			render_status, 
			image_width,
			image_height,
			image_file_types
		) VALUES (?, ?, ?, ?, ?)
		`)
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	_, err = stmt.Exec(scene.SceneID, scene.RenderStatus, scene.ImageWidth, scene.ImageHeight, scene.ImageFileTypes)
	if err != nil {
		return "", err
	}
	err = tx.Commit()
	if err != nil {
		return "", err
	}

	fmt.Println("Successfully created scene row in DB")
	return scene.SceneID, nil
}

func Get(plData *config.PhotolumData, sceneID string) (*Scene, error) {
	fmt.Println("Retrieving scene row in DB...")
	// first, update accessed timestamp
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	tx, err := plData.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	timestampStmt, err := tx.Prepare(`UPDATE scene SET accessed_timestamp = datetime() WHERE scene_id = ?`)
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
	stmt, err := tx.Prepare(`
		SELECT 
			render_status, 
			created_timestamp, 
			modified_timestamp, 
			accessed_timestamp, 
			image_width, 
			image_height, 
			image_file_types
		FROM scene 
		WHERE scene_id = ?
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var renderStatus string
	var createdTimestamp time.Time
	var modifiedTimestamp time.Time
	var accessedTimestamp time.Time
	var imageWidth int
	var imageHeight int
	var imageFileTypes string
	err = stmt.QueryRow(sceneID).Scan(
		&renderStatus,
		&createdTimestamp,
		&modifiedTimestamp,
		&accessedTimestamp,
		&imageWidth,
		&imageHeight,
		&imageFileTypes)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	retrievedSceneRow := &Scene{
		SceneID:           sceneID,
		RenderStatus:      renderstatus.RenderStatus(renderStatus),
		createdTimestamp:  createdTimestamp,
		modifiedTimestamp: modifiedTimestamp,
		accessedTimestamp: accessedTimestamp,
		ImageWidth:        imageWidth,
		ImageHeight:       imageHeight,
		ImageFileTypes:    imageFileTypes,
	}
	fmt.Println("Successfully retrieved scene row in DB")
	return retrievedSceneRow, nil
}

func GetAll(plData *config.PhotolumData) ([]*Scene, error) {
	fmt.Println("Retrieving all scene rows in DB...")
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()
	tx, err := plData.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	stmt, err := tx.Prepare(`
		SELECT 
			scene_id,
			render_status,
			created_timestamp, 
			modified_timestamp, 
			accessed_timestamp, 
			image_width,
			image_height,
			image_file_types
		FROM scene
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var sceneID string
	var renderStatus string
	var createdTimestamp time.Time
	var modifiedTimestamp time.Time
	var accessedTimestamp time.Time
	var imageWidth int
	var imageHeight int
	var imageFileTypes string

	totalSceneRows := []*Scene{}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(
			&sceneID,
			&renderStatus,
			&createdTimestamp,
			&modifiedTimestamp,
			&accessedTimestamp,
			&imageWidth,
			&imageHeight,
			&imageFileTypes)
		if err != nil {
			return nil, err
		}
		totalSceneRows = append(totalSceneRows, &Scene{
			SceneID:           sceneID,
			RenderStatus:      renderstatus.RenderStatus(renderStatus),
			createdTimestamp:  createdTimestamp,
			modifiedTimestamp: modifiedTimestamp,
			accessedTimestamp: accessedTimestamp,
			ImageWidth:        imageWidth,
			ImageHeight:       imageHeight,
			ImageFileTypes:    imageFileTypes,
		})
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully retrieved all image rows in DB...")
	return totalSceneRows, nil
}

func Update(plData *config.PhotolumData, scene *Scene) error {
	fmt.Println("Updating scene row in DB...")
	// update row (and timestamp)
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	tx, err := plData.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`
		UPDATE scene 
		SET 
			render_status = ?, 
			modified_timestamp = datetime(), 
			image_width = ?,
			image_height = ?,
			image_file_types = ? 
		WHERE scene_id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		scene.RenderStatus,
		scene.ImageWidth,
		scene.ImageHeight,
		scene.ImageFileTypes,
		scene.SceneID,
	)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	fmt.Println("Successfully updated scene row in DB")
	return nil
}

func Delete(plData *config.PhotolumData, scene *Scene) error {
	fmt.Println("Deleting scene row in DB...")
	// update row (and timestamp)
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	tx, err := plData.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`DELETE FROM scene WHERE scene_id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(scene.SceneID)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	fmt.Println("Successfully deleted scene row in DB...")
	return nil
}

func DoesExist(plData *config.PhotolumData, sceneID string) (bool, error) {
	fmt.Println("Checking scene existance in DB...")
	// first, update accessed timestamp
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	tx, err := plData.DB.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	stmt, err := tx.Prepare(`SELECT count(*) FROM scene WHERE scene_id = ?`)
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

func UpdateRenderStatus(plData *config.PhotolumData, sceneID string, status renderstatus.RenderStatus) error {
	scene, err := Get(plData, sceneID)
	if err != nil {
		return err
	}
	fmt.Printf("SceneID '%s' -  Updating Render Status: %s --> %s\n", sceneID, scene.RenderStatus, status)
	scene.RenderStatus = status
	Update(plData, scene)
	return nil
}
