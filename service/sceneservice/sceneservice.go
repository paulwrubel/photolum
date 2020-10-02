package sceneservice

// import (
// 	"github.com/paulwrubel/photolum/config"
// 	"github.com/paulwrubel/photolum/config/renderstatus"
// 	"github.com/paulwrubel/photolum/persistence/scenepersistence"
// )

// func Save(plData *config.PhotolumData, scene *scenepersistence.Scene) (string, error) {
// 	return scenepersistence.Save(plData, scene)
// }

// func Get(plData *config.PhotolumData, sceneID string) (*scenepersistence.Scene, error) {
// 	return scenepersistence.Get(plData, sceneID)
// }

// func GetAll(plData *config.PhotolumData) ([]*scenepersistence.Scene, error) {
// 	return scenepersistence.GetAll(plData)
// }

// func Update(plData *config.PhotolumData, scene *scenepersistence.Scene) error {
// 	return scenepersistence.Update(plData, scene)
// }

// func Delete(plData *config.PhotolumData, sceneID string) error {
// 	return scenepersistence.Delete(plData, &scenepersistence.Scene{
// 		SceneID: sceneID,
// 	})
// }

// func UpdateRenderStatus(plData *config.PhotolumData, sceneID string, renderStatus renderstatus.RenderStatus) error {
// 	return scenepersistence.UpdateRenderStatus(plData, sceneID, renderStatus)
// }

// func DoesExist(plData *config.PhotolumData, sceneID string) (bool, error) {
// 	return scenepersistence.DoesExist(plData, sceneID)
// }
