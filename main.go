package main

import (
	"log"
	"time"

	math "github.com/chewxy/math32"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	SCREEN_WIDTH  float32 = 1280
	SCREEN_HEIGHT float32 = 720
	ASPECT_RATIO  float32 = SCREEN_HEIGHT / SCREEN_WIDTH

	DEFAULT_Z_OFFSET float32 = 8

	MODEL_FILENAME string = "sunflower.obj"

	NEAR_DISTANCE float32 = 0.1
	FAR_DISTANCE  float32 = 1000
	FOV_DEGREES   float32 = 90

	ROTATION_SPEED float32 = 2
	POSITION_SPEED float32 = 10
)

var sur *sdl.Surface
var pixelBuffer []byte
var zBuffer []float32

var (
	T_DELTA float32 = 0
)

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Testing", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(SCREEN_WIDTH), int32(SCREEN_HEIGHT), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	defer window.Destroy()

	modelMesh := LoadMesh("mushroom.obj")
	// modelMesh := Mesh{
	// 	tris: []Triangle{
	// 		{
	// 			vecs: [3]Vector4{
	// 				{0, 0, 10, 1, -1},
	// 				{0, 0, 4, 1, -1},
	// 				{0, 2, 4, 1, -1},
	// 			},
	// 		},
	// 	},
	// }

	positionOffset := Vector4{0, 0, DEFAULT_Z_OFFSET + 4, 0, -1}
	rotationTheta := Vector4{0, 180, 0, 0, -1}

	camera := Vector4{0, 0, 0, 1, -1}

	matProj := projectionMatrix(ASPECT_RATIO, FOV_DEGREES, NEAR_DISTANCE, FAR_DISTANCE)

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	sur = surface

	pixelBuffer = surface.Pixels()
	zBuffer = make([]float32, int(SCREEN_WIDTH)*int(SCREEN_HEIGHT))

	var theta float32 = 0
	lastFrame := time.Now()

	curX, curY, _ := sdl.GetMouseState()

	lastFpsRecord := time.Now()
	frames := 0

	var CTRL_PRESSED = false
	var UP_PRESSED = false
	var DOWN_PRESSED = false
	var R_PRESSED = false
	var MOUSE_CLICK = false
	var MOUSE_WHEEL_UP = false
	var MOUSE_WHEEL_DOWN = false
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
				if e.Keysym.Sym == sdl.K_r {
					R_PRESSED = e.State == sdl.PRESSED
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

		frames++
		if time.Now().Sub(lastFpsRecord).Milliseconds() > 1000 {
			log.Printf("FPS: %d", frames)
			frames = 0
			lastFpsRecord = time.Now()
		}

		T_DELTA = float32(time.Now().Sub(lastFrame).Seconds())
		lastFrame = time.Now()

		newX, newY, _ := sdl.GetMouseState()
		diffX, diffY := curX-newX, curY-newY
		curX, curY = newX, newY

		if CTRL_PRESSED && MOUSE_CLICK {
			var MOUSE_ROTATION_SPEED float32 = 5
			if diffX != 0 {
				if diffX > 1 {
					rotationTheta.y += MOUSE_ROTATION_SPEED * T_DELTA
				} else if diffX < -1 {
					rotationTheta.y -= MOUSE_ROTATION_SPEED * T_DELTA
				}
			}
			if diffY != 0 {
				if diffY < -1 {
					rotationTheta.x += MOUSE_ROTATION_SPEED * T_DELTA
				} else if diffY > 1 {
					rotationTheta.x -= MOUSE_ROTATION_SPEED * T_DELTA
				}
			}
		}

		if MOUSE_WHEEL_UP {
			positionOffset.z -= 5 * POSITION_SPEED * T_DELTA
		} else if MOUSE_WHEEL_DOWN {
			positionOffset.z += 5 * POSITION_SPEED * T_DELTA
		}

		// Handle keyboard input
		if UP_PRESSED {
			positionOffset.y -= POSITION_SPEED * T_DELTA
		}
		if DOWN_PRESSED {
			positionOffset.y += POSITION_SPEED * T_DELTA
		}

		if R_PRESSED {
			positionOffset = Vector4{0, 0, DEFAULT_Z_OFFSET, 0, -1}
			rotationTheta = Vector4{0, 0, 0, 0, -1}
		}

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

			if normal.Dot(cameraRay) < 0 {
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

		// Sort by Z depth. Far triangles are drawn first.
		// slices.SortFunc(triangles, func(t1, t2 Triangle) int {
		// 	z1 := (t1.vecs[0].z + t1.vecs[1].z + t1.vecs[2].z) / 3
		// 	z2 := (t2.vecs[0].z + t2.vecs[1].z + t2.vecs[2].z) / 3

		// 	if z1 > z2 {
		// 		return 1
		// 	} else if z1 < z2 {
		// 		return -1
		// 	} else {
		// 		return 0
		// 	}
		// })

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

		window.UpdateSurface()

		for i := 0; i < len(pixelBuffer); i++ {
			pixelBuffer[i] = 0x00
		}

		for i := 0; i < len(zBuffer); i++ {
			zBuffer[i] = math.MaxFloat32
		}

		theta += 1 * T_DELTA
		_ = theta
	}
}
