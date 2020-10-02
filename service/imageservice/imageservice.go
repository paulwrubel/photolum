package imageservice

// import (
// 	"github.com/paulwrubel/photolum/config"
// 	"github.com/paulwrubel/photolum/persistence/imagepersistence"
// )

// func Save(plData *config.PhotolumData, image *imagepersistence.Image) (string, error) {
// 	return imagepersistence.Save(plData, image)
// }

// func Get(plData *config.PhotolumData, sceneID string, fileType string) (*imagepersistence.Image, error) {
// 	return imagepersistence.Get(plData, sceneID, fileType)
// }

// func GetEncoded(plData *config.PhotolumData, sceneID string, fileType string) ([]byte, error) {
// 	return imagepersistence.GetEncoded(plData, sceneID, fileType)
// }

// func GetAll(plData *config.PhotolumData) ([]*imagepersistence.Image, error) {
// 	return imagepersistence.GetAll(plData)
// }

// func Update(plData *config.PhotolumData, image *imagepersistence.Image) error {
// 	return imagepersistence.Update(plData, image)
// }

// func Delete(plData *config.PhotolumData, sceneID string, fileType string) error {
// 	return imagepersistence.Delete(plData, &imagepersistence.Image{
// 		SceneID: sceneID,
// 	})
// }

// func SaveOrUpdate(plData *config.PhotolumData, image *imagepersistence.Image) error {
// 	imageDoesExist, err := imagepersistence.DoesExist(plData, image.SceneID, image.FileType)
// 	if err != nil {
// 		return err
// 	}
// 	if imageDoesExist {
// 		err = imagepersistence.Update(plData, image)
// 	} else {
// 		_, err = imagepersistence.Save(plData, image)
// 	}
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func DoesExist(plData *config.PhotolumData, sceneID string, fileType string) (bool, error) {
// 	return imagepersistence.DoesExist(plData, sceneID, fileType)
// }
