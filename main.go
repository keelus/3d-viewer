package main

import (
	"3d-viewer/ui"
	"fmt"
	"log"
	"path/filepath"
	"sync"
	"time"

	"math"

	"github.com/ncruces/zenity"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	SCREEN_WIDTH  int = 1280
	SCREEN_HEIGHT int = 720

	ASPECT_RATIO float64 = float64(SCREEN_HEIGHT) / float64(SCREEN_WIDTH)

	DEFAULT_Y_OFFSET   float64 = -1
	DEFAULT_Z_OFFSET   float64 = 10
	DEFAULT_Y_ROTATION float64 = math.Pi

	NEAR_DISTANCE float64 = 0.1
	FAR_DISTANCE  float64 = 1000
	FOV_DEGREES   float64 = 90

	ROTATION_SPEED float64 = 2
	POSITION_SPEED float64 = 10

	BG_COLOR uint32 = 0xff202020 // 0xAABBGGRR
)

var (
	SCALE_FACTOR                            int // To scale down the render resolution. For 1280x720, can be: x1, x2, x4, x8, x16
	RENDER_WIDTH, RENDER_HEIGHT             int
	RENDER_WIDTH_FLOAT, RENDER_HEIGHT_FLOAT float64
	RENDER_WIDTH_HALF, RENDER_HEIGHT_HALF   float64
)

var (
	modelMesh *Mesh

	surface *sdl.Surface

	renderBuffer []byte
	screenBuffer []byte

	depthBuffer       []float64
	depthBufferLength int

	flipNormals bool

	tDelta float64 = 0

	positionOffset, rotationTheta Vector4
)

var (
	CTRL_PRESSED     bool = false
	UP_PRESSED       bool = false
	DOWN_PRESSED     bool = false
	MOUSE_CLICK      bool = false
	MOUSE_WHEEL_UP   bool = false
	MOUSE_WHEEL_DOWN bool = false
)

var (
	fontBig     *ttf.Font
	fontRegular *ttf.Font
	fontSmall   *ttf.Font

	btnLoadMesh ui.Button
	lblLoadMesh ui.Label

	cbFileInfo           ui.ContentBlock
	lblFileInfoName      ui.Label
	lblFileInfoTriangles ui.Label
	lblFileInfoVertices  ui.Label

	cbFps  ui.ContentBlock
	lblFps ui.Label

	cbVisualTools             ui.ContentBlock
	lblVisualToolsTitle       ui.Label
	btnVisualToolsFlipNormals ui.Button
	lblVisualToolsFlipNormals ui.Label
	btnVisualToolsResetView   ui.Button
	lblVisualToolsResetView   ui.Label

	cbResolution       ui.ContentBlock
	lblResolutionTitle ui.Label
	btnResolution1     ui.Button
	lblResolution1     ui.Label
	btnResolution2     ui.Button
	lblResolution2     ui.Label
	btnResolution4     ui.Button
	lblResolution4     ui.Label
	btnResolution8     ui.Button
	lblResolution8     ui.Label
	btnResolution16    ui.Button
	lblResolution16    ui.Label

	lblNoMeshLoaded ui.Label
)

func main() {
	// SDL and window setup
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("3D viewer", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(SCREEN_WIDTH), int32(SCREEN_HEIGHT), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	// Window surface and depth buffer setup
	surface, err = window.GetSurface()
	if err != nil {
		panic(err)
	}

	setScale(1)

	// Load TTF module and load fonts
	if err = ttf.Init(); err != nil {
		return
	}
	defer ttf.Quit()

	fontBig = ui.LoadFont("font.ttf", 18)
	defer fontBig.Close()

	fontRegular = ui.LoadFont("font.ttf", 17)
	defer fontRegular.Close()

	fontSmall = ui.LoadFont("font.ttf", 14)
	defer fontSmall.Close()

	// Load custom UI elements
	btnLoadMesh := ui.NewButton(110/2+20, 25/2+10, 110, 25, ui.NewMargin(10, 10), ui.CENTER_CENTER, 0xffffffff, 0xdddddddd, 0xbbbbbbbb)
	lblLoadMesh := ui.NewLabel(110/2+20, 25/2+10+3, "Load file", ui.NewMargin(10, 10), ui.CENTER_CENTER, sdl.Color{R: 0, G: 0, B: 0, A: 255}, fontSmall)

	cbFileInfo = ui.NewContentBlock(1280, 0, 150, 50, ui.NewMargin(10, 10), ui.NewPadding(10, 13), ui.TOP_RIGHT, 0x001a1a1a)
	lblFileInfoName = ui.NewLabel(1280, 5, " ", ui.NewMargin(20, 10), ui.TOP_RIGHT, sdl.Color{R: 255, G: 255, B: 255, A: 255}, fontRegular)
	lblFileInfoTriangles = ui.NewLabel(1280, 30, " ", ui.NewMargin(20, 10), ui.TOP_RIGHT, sdl.Color{R: 255, G: 255, B: 255, A: 255}, fontSmall)
	lblFileInfoVertices = ui.NewLabel(1280, 50, " ", ui.NewMargin(20, 10), ui.TOP_RIGHT, sdl.Color{R: 255, G: 255, B: 255, A: 255}, fontSmall)

	cbFps = ui.NewContentBlock(int32(SCREEN_WIDTH)/2, 0, 85, 20, ui.NewMargin(0, 10), ui.NewPadding(0, 0), ui.TOP_CENTER, 0x00000000)
	lblFps = ui.NewLabel(int32(SCREEN_WIDTH)/2, 0, " ", ui.NewMargin(0, 10), ui.TOP_CENTER, sdl.Color{R: 255, G: 255, B: 255, A: 255}, fontSmall)

	cbVisualTools = ui.NewContentBlock(1280, int32(SCREEN_HEIGHT), 110, 80, ui.NewMargin(10, 10), ui.NewPadding(10, 10), ui.BOTTOM_RIGHT, 0x001a1a1a)
	lblVisualToolsTitle = ui.NewLabel(1280-65, int32(SCREEN_HEIGHT)-75, "Visual tools", ui.NewMargin(20, 10), ui.BOTTOM_CENTER, sdl.Color{R: 255, G: 255, B: 255, A: 255}, fontSmall)
	btnVisualToolsFlipNormals = ui.NewButton(1280-110/2, int32(SCREEN_HEIGHT)-65, 110, 25, ui.NewMargin(20, 10), ui.CENTER_CENTER, 0xffffffff, 0xdddddddd, 0xbbbbbbbb)
	btnVisualToolsResetView = ui.NewButton(1280-110/2, int32(SCREEN_HEIGHT)-35, 110, 25, ui.NewMargin(20, 10), ui.CENTER_CENTER, 0xffffffff, 0xdddddddd, 0xbbbbbbbb)
	lblVisualToolsFlipNormals = ui.NewLabel(1280-110/2, int32(SCREEN_HEIGHT)-60-2, "Flip normals", ui.NewMargin(20, 10), ui.CENTER_CENTER, sdl.Color{R: 0, G: 0, B: 0, A: 255}, fontSmall)
	lblVisualToolsResetView = ui.NewLabel(1280-110/2, int32(SCREEN_HEIGHT)-30-2, "Reset view", ui.NewMargin(20, 10), ui.CENTER_CENTER, sdl.Color{R: 0, G: 0, B: 0, A: 255}, fontSmall)

	lblNoMeshLoaded = ui.NewLabel(int32(SCREEN_WIDTH)/2, int32(SCREEN_HEIGHT)/2, "Load a 3D file to preview it (.obj supported)", ui.NewMargin(0, 0), ui.CENTER_CENTER, sdl.Color{R: 127, G: 127, B: 127, A: 255}, fontBig)

	cbResolution = ui.NewContentBlock(0, int32(SCREEN_HEIGHT), 140, 50, ui.NewMargin(10, 10), ui.NewPadding(10, 10), ui.BOTTOM_LEFT, 0x001a1a1a)
	lblResolutionTitle = ui.NewLabel(100, int32(SCREEN_HEIGHT)-45, "Resolution", ui.NewMargin(20, 10), ui.BOTTOM_CENTER, sdl.Color{R: 255, G: 255, B: 255, A: 255}, fontSmall)

	btnResolution1 = ui.NewButton(50, int32(SCREEN_HEIGHT)-35, 25, 25, ui.NewMargin(20, 10), ui.CENTER_CENTER, 0xffffffff, 0xdddddddd, 0xbbbbbbbb)
	lblResolution1 = ui.NewLabel(100/2, int32(SCREEN_HEIGHT)-32, "x1", ui.NewMargin(20, 10), ui.CENTER_CENTER, sdl.Color{R: 0, G: 0, B: 0, A: 255}, fontSmall)
	btnResolution2 = ui.NewButton(80, int32(SCREEN_HEIGHT)-35, 25, 25, ui.NewMargin(20, 10), ui.CENTER_CENTER, 0xffffffff, 0xdddddddd, 0xbbbbbbbb)
	lblResolution2 = ui.NewLabel(100/2+30, int32(SCREEN_HEIGHT)-32, "/2", ui.NewMargin(20, 10), ui.CENTER_CENTER, sdl.Color{R: 0, G: 0, B: 0, A: 255}, fontSmall)
	btnResolution4 = ui.NewButton(110, int32(SCREEN_HEIGHT)-35, 25, 25, ui.NewMargin(20, 10), ui.CENTER_CENTER, 0xffffffff, 0xdddddddd, 0xbbbbbbbb)
	lblResolution4 = ui.NewLabel(100/2+60, int32(SCREEN_HEIGHT)-32, "/4", ui.NewMargin(20, 10), ui.CENTER_CENTER, sdl.Color{R: 0, G: 0, B: 0, A: 255}, fontSmall)
	btnResolution8 = ui.NewButton(140, int32(SCREEN_HEIGHT)-35, 25, 25, ui.NewMargin(20, 10), ui.CENTER_CENTER, 0xffffffff, 0xdddddddd, 0xbbbbbbbb)
	lblResolution8 = ui.NewLabel(100/2+90, int32(SCREEN_HEIGHT)-32, "/8", ui.NewMargin(20, 10), ui.CENTER_CENTER, sdl.Color{R: 0, G: 0, B: 0, A: 255}, fontSmall)
	btnResolution16 = ui.NewButton(170, int32(SCREEN_HEIGHT)-35, 25, 25, ui.NewMargin(20, 10), ui.CENTER_CENTER, 0xffffffff, 0xdddddddd, 0xbbbbbbbb)
	lblResolution16 = ui.NewLabel(100/2+120, int32(SCREEN_HEIGHT)-32, "/16", ui.NewMargin(20, 10), ui.CENTER_CENTER, sdl.Color{R: 0, G: 0, B: 0, A: 255}, fontSmall)

	// Initialize 3D and misc things
	camera := Vector4{0, 0, 0, 1, -1, NewTexVector(0, 0, 0)}
	flipNormals = false
	matProj := projectionMatrix(ASPECT_RATIO, FOV_DEGREES, NEAR_DISTANCE, FAR_DISTANCE)

	lastFrame := time.Now()
	curX, curY, _ := sdl.GetMouseState()

	// Main loop
	running := true
	for running {
		MOUSE_WHEEL_UP = false
		MOUSE_WHEEL_DOWN = false
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			case *sdl.KeyboardEvent:
				e := event.(*sdl.KeyboardEvent)

				if e.Keysym.Sym == sdl.K_LCTRL {
					CTRL_PRESSED = e.State == sdl.PRESSED
				}
				if e.Keysym.Sym == sdl.K_UP {
					UP_PRESSED = e.State == sdl.PRESSED
				}
				if e.Keysym.Sym == sdl.K_DOWN {
					DOWN_PRESSED = e.State == sdl.PRESSED
				}

				break
			case *sdl.MouseButtonEvent:
				e := event.(*sdl.MouseButtonEvent)

				if e.Button == sdl.BUTTON_LEFT {
					MOUSE_CLICK = e.State == sdl.PRESSED
				}
				break
			case *sdl.MouseWheelEvent:
				e := event.(*sdl.MouseWheelEvent)

				if e.Y > 0 {
					MOUSE_WHEEL_UP = true
				} else if e.Y < 0 {
					MOUSE_WHEEL_DOWN = true
				}

			}
		}

		// Get tDelta
		tDelta = float64(time.Now().Sub(lastFrame).Seconds())
		lblFps.SetText(fmt.Sprintf("FPS: %d", int(1/tDelta)))
		lastFrame = time.Now()

		// Update user mouse information
		newX, newY, _ := sdl.GetMouseState()
		diffX, diffY := curX-newX, curY-newY
		curX, curY = newX, newY

		// Handle keyboard and mouse input
		if CTRL_PRESSED && MOUSE_CLICK {
			var MOUSE_ROTATION_SPEED float64 = 5
			if diffX != 0 {
				if diffX > 1 {
					rotationTheta.y += MOUSE_ROTATION_SPEED * tDelta
				} else if diffX < -1 {
					rotationTheta.y -= MOUSE_ROTATION_SPEED * tDelta
				}
			}
			if diffY != 0 {
				if diffY < -1 {
					rotationTheta.x += MOUSE_ROTATION_SPEED * tDelta
				} else if diffY > 1 {
					rotationTheta.x -= MOUSE_ROTATION_SPEED * tDelta
				}
			}
		}

		if MOUSE_WHEEL_UP {
			positionOffset.z -= 5 * POSITION_SPEED * tDelta
		} else if MOUSE_WHEEL_DOWN {
			positionOffset.z += 5 * POSITION_SPEED * tDelta
		}

		if UP_PRESSED {
			positionOffset.y -= POSITION_SPEED * tDelta
		}
		if DOWN_PRESSED {
			positionOffset.y += POSITION_SPEED * tDelta
		}

		if pressed := btnLoadMesh.UpdateAndGetStatus(curX, curY, MOUSE_CLICK); pressed {
			selected, _ := zenity.SelectFile(
				zenity.Filename("/"),
				zenity.FileFilters{
					{
						Name:     "OBJ files",
						Patterns: []string{"*.obj"},
						CaseFold: false,
					},
				})
			if selected != "" {
				log.Print(selected)
				LoadFile(selected)
				continue
			}
		}

		if pressed := btnVisualToolsFlipNormals.UpdateAndGetStatus(curX, curY, MOUSE_CLICK); pressed {
			flipNormals = !flipNormals
		}

		if pressed := btnVisualToolsResetView.UpdateAndGetStatus(curX, curY, MOUSE_CLICK); pressed {
			ResetCameraView()
		}

		if pressed := btnResolution1.UpdateAndGetStatus(curX, curY, MOUSE_CLICK); pressed {
			setScale(1)
		}
		if pressed := btnResolution2.UpdateAndGetStatus(curX, curY, MOUSE_CLICK); pressed {
			setScale(2)
		}
		if pressed := btnResolution4.UpdateAndGetStatus(curX, curY, MOUSE_CLICK); pressed {
			setScale(4)
		}
		if pressed := btnResolution8.UpdateAndGetStatus(curX, curY, MOUSE_CLICK); pressed {
			setScale(8)
		}
		if pressed := btnResolution16.UpdateAndGetStatus(curX, curY, MOUSE_CLICK); pressed {
			setScale(16)
		}

		// Main 3D code
		if modelMesh != nil {
			// Rotation, translation & world matrix
			rotationMatrix := rotationMatrix(rotationTheta)
			translationMatrix := MakeTranslation(positionOffset.x, positionOffset.y, positionOffset.z)
			worldMatrix := rotationMatrix.multiplyMatrix(translationMatrix)

			triangles := []Triangle{}
			// Evaluate each triangle and save them if the normals face the camera. Handle clipping triangles.
			for _, tri := range modelMesh.tris {
				triTransformed := worldMatrix.multiplyTriangle(tri)
				triTransformed.vecs[0].originalZ = triTransformed.vecs[0].z
				triTransformed.vecs[1].originalZ = triTransformed.vecs[1].z
				triTransformed.vecs[2].originalZ = triTransformed.vecs[2].z

				// Calculate the normal of the triangle face
				line1 := triTransformed.vecs[1].Sub(triTransformed.vecs[0])
				line2 := triTransformed.vecs[2].Sub(triTransformed.vecs[0])
				normal := line1.CrossProduct(line2).Normalise()

				cameraRay := triTransformed.vecs[0].Sub(camera)

				if (normal.Dot(cameraRay) < 0 && !flipNormals) || (normal.Dot(cameraRay) > 0 && flipNormals) {
					// Simple illumination via light direction
					lightDirection := Vector4{0, 1, -1, 1, -1, NewTexVector(0, 0, 0)}.Normalise()
					ilumination := math.Max(0.1, lightDirection.Dot(normal))

					// Transform and project triangles
					clipped := ClipAgainstPlane(Vector4{0, 0, 0.1, 1, -1, NewTexVector(0, 0, 0)}, Vector4{0, 0, 1, 1, -1, NewTexVector(0, 0, 0)}, triTransformed)
					for n := 0; n < len(clipped); n++ {
						// Project triangles to 2D
						triProjected := matProj.multiplyTriangle(clipped[n])
						triProjected.ilum = ilumination
						triProjected.tex = tri.tex

						// Apply depth
						triProjected.vecs[0].texVec.u /= triProjected.vecs[0].w
						triProjected.vecs[1].texVec.u /= triProjected.vecs[1].w
						triProjected.vecs[2].texVec.u /= triProjected.vecs[2].w
						triProjected.vecs[0].texVec.v /= triProjected.vecs[0].w
						triProjected.vecs[1].texVec.v /= triProjected.vecs[1].w
						triProjected.vecs[2].texVec.v /= triProjected.vecs[2].w
						triProjected.vecs[0].texVec.w = 1 / triProjected.vecs[0].w
						triProjected.vecs[1].texVec.w = 1 / triProjected.vecs[1].w
						triProjected.vecs[2].texVec.w = 1 / triProjected.vecs[2].w

						triProjected.vecs[0] = triProjected.vecs[0].Div(triProjected.vecs[0].w)
						triProjected.vecs[1] = triProjected.vecs[1].Div(triProjected.vecs[1].w)
						triProjected.vecs[2] = triProjected.vecs[2].Div(triProjected.vecs[2].w)

						// Offset into view
						vOffsetView := Vector4{1, 1, 0, 0, -1, NewTexVector(0, 0, 0)}
						triProjected.vecs[0] = triProjected.vecs[0].Add(vOffsetView)
						triProjected.vecs[1] = triProjected.vecs[1].Add(vOffsetView)
						triProjected.vecs[2] = triProjected.vecs[2].Add(vOffsetView)

						triProjected.vecs[0].originalZ = triTransformed.vecs[0].originalZ
						triProjected.vecs[1].originalZ = triTransformed.vecs[1].originalZ
						triProjected.vecs[2].originalZ = triTransformed.vecs[2].originalZ

						// Expand to screen size
						triProjected.vecs[0].x *= RENDER_WIDTH_HALF
						triProjected.vecs[0].y *= RENDER_HEIGHT_HALF
						triProjected.vecs[1].x *= RENDER_WIDTH_HALF
						triProjected.vecs[1].y *= RENDER_HEIGHT_HALF
						triProjected.vecs[2].x *= RENDER_WIDTH_HALF
						triProjected.vecs[2].y *= RENDER_HEIGHT_HALF

						triangles = append(triangles, triProjected)
					}
				}
			}

			// Check for clipping and draw triangles to screen.
			for _, triToRaster := range triangles {
				clipped := []Triangle{}
				listTriangles := []Triangle{}
				listTriangles = append(listTriangles, triToRaster)
				nNewTriangles := 1

				for p := 0; p < 4; p++ {
					nTrisToAdd := 0
					for nNewTriangles > 0 {
						curTri := listTriangles[0]
						listTriangles = listTriangles[1:]
						nNewTriangles--

						// Clip against each plane (screen borders)
						switch p {
						case 0:
							clipped = ClipAgainstPlane(Vector4{0, 0, 0, 1, -1, NewTexVector(0, 0, 0)}, Vector4{0, 1, 0, 1, -1, NewTexVector(0, 0, 0)}, curTri)
						case 1:
							clipped = ClipAgainstPlane(Vector4{0, RENDER_HEIGHT_FLOAT, 0, 1, -1, NewTexVector(0, 0, 0)}, Vector4{0, -1, 0, 1, -1, NewTexVector(0, 0, 0)}, curTri)
						case 2:
							clipped = ClipAgainstPlane(Vector4{0, 0, 0, 1, -1, NewTexVector(0, 0, 0)}, Vector4{1, 0, 0, 1, -1, NewTexVector(0, 0, 0)}, curTri)
						case 3:
							clipped = ClipAgainstPlane(Vector4{RENDER_WIDTH_FLOAT, 0, 0, 1, -1, NewTexVector(0, 0, 0)}, Vector4{-1, 0, 0, 1, -1, NewTexVector(0, 0, 0)}, curTri)
						}

						nTrisToAdd = len(clipped)
						for w := 0; w < nTrisToAdd; w++ {
							listTriangles = append(listTriangles, clipped[w])
						}
					}

					nNewTriangles = len(listTriangles)
				}

				for _, tri := range listTriangles {
					tri.Draw()
				}
			}
		}

		if SCALE_FACTOR > 1 {
			var wg sync.WaitGroup

			for y := 0; y < int(RENDER_HEIGHT); y++ {
				wg.Add(1)
				go func(y int) {
					defer wg.Done()
					for x := 0; x < RENDER_WIDTH; x++ {

						screenIdx := y*int(SCREEN_WIDTH)*SCALE_FACTOR*4 + x*SCALE_FACTOR*4
						renderIdx := (y*RENDER_WIDTH + x) * 4

						for b := 0; b < 4; b++ {
							for dy := 0; dy < SCALE_FACTOR; dy++ {
								for dx := 0; dx < SCALE_FACTOR; dx++ {
									screenBuffer[screenIdx+int(SCREEN_WIDTH)*4*dy+4*dx+b] = renderBuffer[renderIdx+b]
								}
							}
						}
					}
				}(y)
			}

			wg.Wait()
		}

		// Draw UI elements
		btnLoadMesh.Draw(surface)
		lblLoadMesh.Draw(surface)

		if modelMesh != nil {
			cbFileInfo.Draw(surface)
			lblFileInfoName.Draw(surface)
			lblFileInfoTriangles.Draw(surface)
			lblFileInfoVertices.Draw(surface)
		} else {
			lblNoMeshLoaded.Draw(surface)
		}

		cbFps.Draw(surface)
		lblFps.Draw(surface)

		cbVisualTools.Draw(surface)
		lblVisualToolsTitle.Draw(surface)
		btnVisualToolsFlipNormals.Draw(surface)
		btnVisualToolsResetView.Draw(surface)
		lblVisualToolsFlipNormals.Draw(surface)
		lblVisualToolsResetView.Draw(surface)

		cbResolution.Draw(surface)
		lblResolutionTitle.Draw(surface)
		btnResolution1.Draw(surface)
		lblResolution1.Draw(surface)
		btnResolution2.Draw(surface)
		lblResolution2.Draw(surface)
		btnResolution4.Draw(surface)
		lblResolution4.Draw(surface)
		btnResolution8.Draw(surface)
		lblResolution8.Draw(surface)
		btnResolution16.Draw(surface)
		lblResolution16.Draw(surface)

		//lblResolution1.Draw(surface)
		//btnResolution2.Draw(surface)
		//lblResolution2.Draw(surface)
		//btnResolution4.Draw(surface)
		//lblResolution4.Draw(surface)
		//btnResolution8.Draw(surface)
		//lblResolution8.Draw(surface)
		//btnResolution16.Draw(surface)
		//lblResolution16.Draw(surface)

		// Update and clean screen buffers
		window.UpdateSurface()

		for i := 0; i < len(renderBuffer); i += 4 {
			renderBuffer[i] = byte((BG_COLOR >> 16) & 0xff)
			renderBuffer[i+1] = byte((BG_COLOR >> 8) & 0xff)
			renderBuffer[i+2] = byte(BG_COLOR & 0xff)
			renderBuffer[i+3] = byte((BG_COLOR >> 24) & 0xff)
		}

		for i := 0; i < len(depthBuffer); i++ {
			depthBuffer[i] = math.MaxFloat64
		}
	}

}

func LoadFile(modelFilePath string) {
	modelMesh = ParseObj(modelFilePath)

	ResetCameraView()

	filename := filepath.Base(modelFilePath)

	lblFileInfoName.SetText(filename)
	lblFileInfoTriangles.SetText(fmt.Sprintf("Triangles: %d", modelMesh.triangleAmount))
	lblFileInfoVertices.SetText(fmt.Sprintf("Vertices: %d", modelMesh.vertexAmount))

	width := math.Max(math.Max(float64(lblFileInfoName.GetRectWidth()), float64(lblFileInfoTriangles.GetRectWidth())), float64(lblFileInfoVertices.GetRectWidth()))

	cbFileInfo.UpdateRectToWidth(int32(width))
}

func ResetCameraView() {
	positionOffset = Vector4{0, DEFAULT_Y_OFFSET, DEFAULT_Z_OFFSET, 0, -1, NewTexVector(0, 0, 0)}
	rotationTheta = Vector4{0, DEFAULT_Y_ROTATION, 0, 0, -1, NewTexVector(0, 0, 0)}

	if modelMesh != nil {
		positionOffset.z = -modelMesh.lowestZ * 3
		positionOffset.y = -(modelMesh.lowestY + modelMesh.highestY) / 2
	}
}

func setScale(scale int) {
	if scale != 1 && scale != 2 && scale != 4 && scale != 8 && scale != 16 {
		log.Fatalf("Unexpected resolution scale '%d'", scale)
	}

	SCALE_FACTOR = scale

	RENDER_WIDTH = SCREEN_WIDTH / SCALE_FACTOR
	RENDER_HEIGHT = SCREEN_HEIGHT / SCALE_FACTOR

	RENDER_WIDTH_FLOAT = float64(RENDER_WIDTH)
	RENDER_HEIGHT_FLOAT = float64(RENDER_HEIGHT)

	RENDER_WIDTH_HALF = RENDER_WIDTH_FLOAT * 0.5
	RENDER_HEIGHT_HALF = RENDER_HEIGHT_FLOAT * 0.5

	if SCALE_FACTOR > 1 {
		renderBuffer = make([]byte, RENDER_WIDTH*RENDER_HEIGHT*4)
		screenBuffer = surface.Pixels()
	} else {
		renderBuffer = surface.Pixels()
	}

	depthBuffer = make([]float64, RENDER_WIDTH*RENDER_HEIGHT)
	depthBufferLength = len(depthBuffer)
}
