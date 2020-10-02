package tracingservice

// import (
// 	"github.com/paulwrubel/photolum/config"
// 	"github.com/paulwrubel/photolum/config/renderstatus"
// 	"github.com/paulwrubel/photolum/service/sceneservice"
// 	"github.com/paulwrubel/photolum/tracing"
// )

// func StartRender(plData *config.PhotolumData, sceneID string) error {
// 	scene, err := sceneservice.Get(plData, sceneID)
// 	if err != nil {
// 		return err
// 	}
// 	err = sceneservice.UpdateRenderStatus(plData, sceneID, renderstatus.Running)
// 	if err != nil {
// 		return err
// 	}
// 	go tracing.TraceImage(plData, scene)
// 	return nil
// }

// func StopRender(plData *config.PhotolumData, sceneID string) error {
// 	err := sceneservice.UpdateRenderStatus(plData, sceneID, renderstatus.Stopping)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
