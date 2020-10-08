package tracing

import (
	"context"
	"image"
	"math"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/paulwrubel/photolum/config"
	"github.com/paulwrubel/photolum/config/geometry"
	"github.com/paulwrubel/photolum/config/shading"
	"github.com/paulwrubel/photolum/enumeration/renderstatus"
	"github.com/paulwrubel/photolum/persistence/renderpersistence.go"
	"github.com/sirupsen/logrus"
	"golang.org/x/image/draw"
	"golang.org/x/sync/semaphore"
)

func RunWorker(plData *config.PhotolumData,
	log *logrus.Entry,
	parameters *config.Parameters,
	renderName string,
	encodingChan chan<- *config.TracingPayload) {
	log.Debug("running tracing worker")

	// create new image
	img := image.NewRGBA64(image.Rect(0, 0, parameters.ImageWidth, parameters.ImageHeight))

	tiles := getTiles(parameters, img)

	// shuffle tiles
	// rand.Shuffle(len(tiles), func(i, j int) {
	// 	tiles[i], tiles[j] = tiles[j], tiles[i]
	// })

	roundChan := make(chan bool)
	tileChan := make(chan bool)
	doneChan := make(chan bool)
	go runProgressWorker(plData, log, renderName, parameters.RoundCount, len(tiles), roundChan, tileChan, doneChan)

	for round := 1; round <= parameters.RoundCount; round++ {
		log.Debugf("beginning round %d", round)
		traceRound(parameters, log, img, tiles, round, tileChan)
		log.Debugf("round %d finished, copying image", round)
		bounds := img.Bounds()
		imgCopy := image.NewRGBA64(bounds)
		draw.Draw(imgCopy, bounds, img, bounds.Min, draw.Src)
		payload := &config.TracingPayload{
			FileType: parameters.FileType,
			Image:    imgCopy,
		}
		log.Debugf("image copied, sending to encoder")
		encodingChan <- payload
		roundChan <- true
	}
	doneChan <- true
	close(encodingChan)

	err := renderpersistence.UpdateRenderStatus(plData, log, renderName, renderstatus.Completed)
	if err != nil {
		log.WithError(err).Error("error setting render to completed")
		renderpersistence.UpdateRenderStatus(plData, log, renderName, renderstatus.Error)
	}

	log.Debug("closing tracing worker")
}

func runProgressWorker(plData *config.PhotolumData,
	log *logrus.Entry,
	renderName string,
	totalRounds int,
	totalTiles int,
	roundChan <-chan bool,
	tileChan <-chan bool,
	doneChan <-chan bool) {
	completedRounds := 0
	completedTiles := 0
	roundPercentage := 1.0 / float64(totalRounds)
	for {
		select {
		case <-roundChan:
			completedRounds++
			completedTiles = 0
			_ = renderpersistence.UpdateCompletedRounds(plData, log, renderName, uint32(completedRounds))
		case <-tileChan:
			completedTiles++
			progress := (float64(completedRounds) / float64(totalRounds)) + roundPercentage*(float64(completedTiles)/float64(totalTiles))
			_ = renderpersistence.UpdateRenderProgress(plData, log, renderName, progress)
		case <-doneChan:
			return
		}
	}
}

func traceRound(params *config.Parameters,
	log *logrus.Entry,
	img *image.RGBA64,
	tiles []config.Tile,
	roundNum int,
	tileChan chan<- bool) {

	sem := semaphore.NewWeighted(int64(runtime.NumCPU()))

	wg := sync.WaitGroup{}
	for _, tile := range tiles {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		sem.Acquire(context.Background(), 1)
		wg.Add(1)
		go traceTile(params, log, r, img, sem, &wg, tile, roundNum, tileChan)
	}
	wg.Wait()
}

// traceTile iterates over the pixels in a tile and writes the received colors to the image
func traceTile(p *config.Parameters,
	log *logrus.Entry,
	rng *rand.Rand,
	img *image.RGBA64,
	sem *semaphore.Weighted,
	wg *sync.WaitGroup,
	t config.Tile,
	roundNum int,
	tileChan chan<- bool) {

	defer sem.Release(1)
	defer wg.Done()
	//log.Tracef("tracing tile id: %s", t.ID)
	for y := t.Origin.Y; y < t.Origin.Y+t.Span.Y; y++ {
		for x := t.Origin.X; x < t.Origin.X+t.Span.X; x++ {
			pixelColor := tracePixel(p, int(x), int(y), rng)

			// we need to weight the color value by what round we just finished
			imgColor := shading.MakeColor(img.RGBA64At(int(x), p.ImageHeight-int(y)-1))
			weightedColor := imgColor.MultScalar(float64(roundNum - 1)).Add(pixelColor).DivScalar(float64(roundNum))
			img.SetRGBA64(int(x), p.ImageHeight-int(y)-1, weightedColor.ToRGBA64())
		}
	}
	tileChan <- true
	// dc <- 1
}

// tracePixel gets the color for a pixel
func tracePixel(p *config.Parameters, x, y int, rng *rand.Rand) shading.Color {
	pixelColor := shading.Color{}
	for s := 0; s < p.SamplesPerRound; s++ {
		// pick a random spot on the pixel to shoot a ray into
		// this is purely random, NOT stratified
		u := (float64(x) + rng.Float64()) / float64(p.ImageWidth)
		v := (float64(y) + rng.Float64()) / float64(p.ImageHeight)

		ray := p.Scene.Camera.GetRay(u, v, rng)

		tempColor := traceRay(p, ray, rng, 0)
		pixelColor = pixelColor.Add(tempColor)
	}
	if p.UseScalingTruncation {
		return pixelColor.DivScalar(float64(p.SamplesPerRound)).ScaleDown(1.0).Pow(1.0 / p.GammaCorrection)
	}
	return pixelColor.DivScalar(float64(p.SamplesPerRound)).Clamp(0, 1).Pow(1.0 / p.GammaCorrection)

}

// traceRay casts in individual ray into the scene
func traceRay(parameters *config.Parameters, r geometry.Ray, rng *rand.Rand, depth int) shading.Color {

	// if we've gone too deep...
	if depth > parameters.MaxBounces {
		// ...just return BLACK
		return shading.ColorBlack
	}
	// check if we've hit something
	rayHit, hitSomething := parameters.Scene.Objects.Intersection(r, parameters.TMin, parameters.TMax)
	// if we did not hit something...
	if !hitSomething {
		// ...return the background color
		// TODO: add support for HDR skymaps
		return parameters.BackgroundColor
	}

	mat := rayHit.Material

	// if the surface is BLACK, it's not going to let any incoming light contribute to the outgoing color
	// so we can safely say no light is reflected and simply return the emittance of the material
	if mat.Reflectance(rayHit.U, rayHit.V) == shading.ColorBlack {
		return mat.Emittance(rayHit.U, rayHit.V)
	}

	// get the reflection incoming ray
	scatteredRay, wasScattered := rayHit.Material.Scatter(*rayHit, rng)
	// if no ray could have reflected to us, we just return BLACK
	if !wasScattered {
		return shading.ColorBlack
	}
	// get the color that came to this point and gave us the outgoing ray
	incomingColor := traceRay(parameters, scatteredRay, rng, depth+1)
	// return the (very-roughly approximated) value of the rendering equation
	return mat.Emittance(rayHit.U, rayHit.V).Add(mat.Reflectance(rayHit.U, rayHit.V).MultColor(incomingColor))
}

// getTiles creates and return a grid of tiles on the image
func getTiles(p *config.Parameters, i *image.RGBA64) []config.Tile {
	tiles := []config.Tile{}
	idNum := 0
	for y := 0; y < p.ImageHeight; y += p.TileHeight {
		for x := 0; x < p.ImageWidth; x += p.TileWidth {
			idNum++
			width := math.Min(float64(p.TileWidth), float64(p.ImageWidth-x))
			height := math.Min(float64(p.TileHeight), float64(p.ImageHeight-y))
			tiles = append(tiles, config.Tile{
				ID: strconv.Itoa(idNum),
				Origin: geometry.Point{
					X: float64(x),
					Y: float64(y),
				},
				Span: geometry.Vector{
					X: width,
					Y: height,
				},
			})
		}
	}
	return tiles
}
