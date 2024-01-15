package main

import (
	"fmt"
	"ne3d/ui"
	"path"
	"time"

	math "github.com/chewxy/math32"
	"github.com/ncruces/zenity"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	SCREEN_WIDTH  float32 = 1280
	SCREEN_HEIGHT float32 = 720
	ASPECT_RATIO  float32 = SCREEN_HEIGHT / SCREEN_WIDTH

	DEFAULT_Y_OFFSET   float32 = -1
	DEFAULT_Z_OFFSET   float32 = 2
	DEFAULT_Y_ROTATION float32 = math.Pi

	NEAR_DISTANCE float32 = 0.1
	FAR_DISTANCE  float32 = 1000
	FOV_DEGREES   float32 = 90

	ROTATION_SPEED float32 = 2
	POSITION_SPEED float32 = 10
)

var (
	modelMesh *Mesh

	surface     *sdl.Surface
	pixelBuffer []byte
	depthBuffer []float32

	flipNormals bool

	tDelta float32 = 0

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

	lblNoMeshLoaded ui.Label
)

func main() {
	// SDL and window setup
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Testing", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(SCREEN_WIDTH), int32(SCREEN_HEIGHT), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	// Window surface and depth buffer setup
	surface, err = window.GetSurface()
	if err != nil {
		panic(err)
	}

	pixelBuffer = surface.Pixels()
	depthBuffer = make([]float32, int(SCREEN_WIDTH)*int(SCREEN_HEIGHT))

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

	// Initialize 3D and misc things
	//LoadFile("Boat.obj")

	camera := Vector4{0, 0, 0, 1, -1}
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
		tDelta = float32(time.Now().Sub(lastFrame).Seconds())
		lblFps.SetText(fmt.Sprintf("FPS: %d", int(1/tDelta)))
		lastFrame = time.Now()

		// Update user mouse information
		newX, newY, _ := sdl.GetMouseState()
		diffX, diffY := curX-newX, curY-newY
		curX, curY = newX, newY

		// Handle keyboard and mouse input
		if CTRL_PRESSED && MOUSE_CLICK {
			var MOUSE_ROTATION_SPEED float32 = 5
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
					lightDirection := Vector4{0, 1, -1, 1, -1}.Normalise()
					ilumination := math.Max(0.1, lightDirection.Dot(normal))

					// Transform and project triangles
					clipped := ClipAgainstPlane(Vector4{0, 0, 0.1, 1, -1}, Vector4{0, 0, 1, 1, -1}, triTransformed)
					for n := 0; n < len(clipped); n++ {
						// Project triangles to 2D
						triProjected := matProj.multiplyTriangle(clipped[n])
						triProjected.ilum = ilumination

						// Apply depth
						triProjected.vecs[0] = triProjected.vecs[0].Div(triProjected.vecs[0].w)
						triProjected.vecs[1] = triProjected.vecs[1].Div(triProjected.vecs[1].w)
						triProjected.vecs[2] = triProjected.vecs[2].Div(triProjected.vecs[2].w)

						// Offset into view
						vOffsetView := Vector4{1, 1, 0, 0, -1}
						triProjected.vecs[0] = triProjected.vecs[0].Add(vOffsetView)
						triProjected.vecs[1] = triProjected.vecs[1].Add(vOffsetView)
						triProjected.vecs[2] = triProjected.vecs[2].Add(vOffsetView)

						triProjected.vecs[0].originalZ = triTransformed.vecs[0].originalZ
						triProjected.vecs[1].originalZ = triTransformed.vecs[1].originalZ
						triProjected.vecs[2].originalZ = triTransformed.vecs[2].originalZ

						// Expand to screen size
						triProjected.vecs[0].x *= 0.5 * SCREEN_WIDTH
						triProjected.vecs[0].y *= 0.5 * SCREEN_HEIGHT
						triProjected.vecs[1].x *= 0.5 * SCREEN_WIDTH
						triProjected.vecs[1].y *= 0.5 * SCREEN_HEIGHT
						triProjected.vecs[2].x *= 0.5 * SCREEN_WIDTH
						triProjected.vecs[2].y *= 0.5 * SCREEN_HEIGHT

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
						test := listTriangles[0]
						listTriangles = listTriangles[1:]
						nNewTriangles--

						// Clip against each plane (screen borders)
						switch p {
						case 0:
							clipped = ClipAgainstPlane(Vector4{0, 0, 0, 1, -1}, Vector4{0, 1, 0, 1, -1}, test)
						case 1:
							clipped = ClipAgainstPlane(Vector4{0, SCREEN_HEIGHT, 0, 1, -1}, Vector4{0, -1, 0, 1, -1}, test)
						case 2:
							clipped = ClipAgainstPlane(Vector4{0, 0, 0, 1, -1}, Vector4{1, 0, 0, 1, -1}, test)
						case 3:
							clipped = ClipAgainstPlane(Vector4{SCREEN_WIDTH, 0, 0, 1, -1}, Vector4{-1, 0, 0, 1, -1}, test)
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

		// Update and clean screen buffers
		window.UpdateSurface()

		for i := 0; i < len(pixelBuffer); i++ {
			pixelBuffer[i] = 0x00
		}

		for i := 0; i < len(depthBuffer); i++ {
			depthBuffer[i] = math.MaxFloat32
		}
	}
}

func LoadFile(filepath string) {
	modelMesh = LoadMesh(filepath)

	ResetCameraView()

	filename := path.Base(filepath)

	lblFileInfoName.SetText(filename)
	lblFileInfoTriangles.SetText(fmt.Sprintf("Triangles: %d", modelMesh.triangleAmount))
	lblFileInfoVertices.SetText(fmt.Sprintf("Vertices: %d", modelMesh.vertexAmount))

	width := math.Max(math.Max(float32(lblFileInfoName.GetRectWidth()), float32(lblFileInfoTriangles.GetRectWidth())), float32(lblFileInfoVertices.GetRectWidth()))

	cbFileInfo.UpdateRectToWidth(int32(width))
}

func ResetCameraView() {
	positionOffset = Vector4{0, DEFAULT_Y_OFFSET, DEFAULT_Z_OFFSET, 0, -1}
	rotationTheta = Vector4{0, DEFAULT_Y_ROTATION, 0, 0, -1}

	if modelMesh != nil {
		positionOffset.z = -modelMesh.lowestZ * 3
		positionOffset.y = -(modelMesh.lowestY + modelMesh.highestY) / 2
	}
}
