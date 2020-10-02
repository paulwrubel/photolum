package tracing

// import (
// 	"fmt"
// 	"image"
// 	"image/color"
// 	"math"
// 	"strings"

// 	"github.com/paulwrubel/photolum/config"
// 	"github.com/paulwrubel/photolum/config/renderstatus"
// 	"github.com/paulwrubel/photolum/persistence/imagepersistence"
// 	"github.com/paulwrubel/photolum/persistence/scenepersistence"
// 	"github.com/paulwrubel/photolum/service/imageservice"
// 	"github.com/paulwrubel/photolum/service/sceneservice"
// )

// func TraceImage(plData *config.PhotolumData, scene *scenepersistence.Scene) {
// 	newImage := image.NewRGBA64(image.Rect(0, 0, scene.ImageWidth, scene.ImageHeight))
// 	for y := 0; y < scene.ImageHeight; y++ {
// 		for x := 0; x < scene.ImageWidth; x++ {
// 			col := color.RGBA64{
// 				R: uint16(0.0 * float64(math.MaxUint16)),
// 				G: uint16((float64(x) / float64(scene.ImageWidth)) * float64(math.MaxUint16)),
// 				B: uint16((float64(y) / float64(scene.ImageHeight)) * float64(math.MaxUint16)),
// 				A: uint16(1.0 * float64(math.MaxUint16))}
// 			newImage.SetRGBA64(x, y, col)
// 		}
// 	}
// 	for _, fileType := range strings.Split(scene.ImageFileTypes, ",") {
// 		err := imageservice.SaveOrUpdate(plData, &imagepersistence.Image{
// 			SceneID:   scene.SceneID,
// 			FileType:  fileType,
// 			ImageData: newImage,
// 		})
// 		if err != nil {
// 			fmt.Printf("Error in tracer.go: %s\n", err.Error())
// 			sceneservice.UpdateRenderStatus(plData, scene.SceneID, renderstatus.Error)
// 			return
// 		}
// 	}
// 	err := sceneservice.UpdateRenderStatus(plData, scene.SceneID, renderstatus.Completed)
// 	if err != nil {
// 		fmt.Printf("Error in tracer.go: %s\n", err.Error())
// 		scenepersistence.UpdateRenderStatus(plData, scene.SceneID, renderstatus.Error)
// 		return
// 	}
// }
