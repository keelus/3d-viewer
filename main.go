package main

import (
	"fmt"
	"image/color"
	"slices"
	"time"

	math "github.com/chewxy/math32"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	SCREEN_WIDTH  float32 = 1280
	SCREEN_HEIGHT float32 = 720
	ASPECT_RATIO  float32 = SCREEN_HEIGHT / SCREEN_WIDTH

	DEFAULT_Z_OFFSET float32 = 8

	MODEL_FILENAME string = "mushroom.obj"

	NEAR_DISTANCE float32 = 0.1
	FAR_DISTANCE  float32 = 1000
	FOV_DEGREES   float32 = 90

	ROTATION_SPEED float32 = 1
	POSITION_SPEED float32 = 2
)

var (
	T_DELTA float32 = 0
)

func main() {
	rl.InitWindow(int32(SCREEN_WIDTH), int32(SCREEN_HEIGHT), "3d engine")

	lastFrame := time.Now()

	modelMesh := LoadMesh(MODEL_FILENAME)

	positionOffset := Vector4{0, 0, DEFAULT_Z_OFFSET, 0}
	rotationTheta := Vector4{0, 0, 0, 0}

	camera := Vector4{0, 0, 0, 1}

	matProj := projectionMatrix(ASPECT_RATIO, FOV_DEGREES, NEAR_DISTANCE, FAR_DISTANCE)

	for !rl.WindowShouldClose() {
		// Calculate T_DELTA for different framerate compatibility
		T_DELTA = float32(time.Now().Sub(lastFrame).Seconds())
		lastFrame = time.Now()

		// Handle keyboard input
		if rl.IsKeyDown(rl.KeyA) {
			rotationTheta.y += ROTATION_SPEED * T_DELTA
		}
		if rl.IsKeyDown(rl.KeyD) {
			rotationTheta.y -= ROTATION_SPEED * T_DELTA
		}

		if rl.IsKeyDown(rl.KeyW) {
			rotationTheta.x += ROTATION_SPEED * T_DELTA
		}
		if rl.IsKeyDown(rl.KeyS) {
			rotationTheta.x -= ROTATION_SPEED * T_DELTA
		}

		if rl.IsKeyDown(rl.KeyE) {
			rotationTheta.z += ROTATION_SPEED * T_DELTA
		}
		if rl.IsKeyDown(rl.KeyQ) {
			rotationTheta.z -= ROTATION_SPEED * T_DELTA
		}

		if rl.IsKeyDown(rl.KeyUp) {
			positionOffset.y -= POSITION_SPEED * T_DELTA
		}
		if rl.IsKeyDown(rl.KeyDown) {
			positionOffset.y += POSITION_SPEED * T_DELTA
		}

		if rl.IsKeyDown(rl.KeyRight) {
			positionOffset.x += POSITION_SPEED * T_DELTA
		}
		if rl.IsKeyDown(rl.KeyLeft) {
			positionOffset.x -= POSITION_SPEED * T_DELTA
		}

		if rl.IsKeyDown(rl.KeyX) {
			positionOffset.z += POSITION_SPEED * T_DELTA
		}
		if rl.IsKeyDown(rl.KeyC) {
			positionOffset.z -= POSITION_SPEED * T_DELTA
		}
		if rl.IsKeyDown(rl.KeyR) {
			positionOffset = Vector4{0, 0, DEFAULT_Z_OFFSET, 0}
			rotationTheta = Vector4{0, 0, 0, 0}
		}

		// Clear canvas
		rl.BeginDrawing()
		rl.ClearBackground(color.RGBA{0, 0, 0, 255})

		// Print FPS
		fpsColor := color.RGBA{255, 255, 255, 255}
		if rl.GetFPS() < 60 {
			fpsColor = color.RGBA{255, 0, 0, 255}
		}
		rl.DrawText(fmt.Sprintf("FPS: %d", rl.GetFPS()), int32(SCREEN_WIDTH)-120, 10, 22, fpsColor)

		rl.DrawText(fmt.Sprintf("Model triangles: %d", modelMesh.triangleAmount), int32(SCREEN_WIDTH)-180, 30, 18, rl.RayWhite)
		rl.DrawText(fmt.Sprintf("Model vertices: %d", modelMesh.vertexAmount), int32(SCREEN_WIDTH)-180, 50, 18, rl.RayWhite)
		// Rotation, translation & world matrix
		rotationMatrix := rotationMatrix(rotationTheta)
		translationMatrix := MakeTranslation(positionOffset.x, positionOffset.y, positionOffset.z)
		worldMatrix := rotationMatrix.multiplyMatrix(translationMatrix)

		triangles := []Triangle{}
		// Evaluate each triangle and save them if the normals face the camera. Handle clipping triangles.
		for _, tri := range modelMesh.tris {
			triTransformed := worldMatrix.multiplyTriangle(tri)

			// Calculate the normal of the triangle face
			line1 := triTransformed.vecs[1].Sub(triTransformed.vecs[0])
			line2 := triTransformed.vecs[2].Sub(triTransformed.vecs[0])
			normal := line1.CrossProduct(line2).Normalise()

			cameraRay := triTransformed.vecs[0].Sub(camera)

			if normal.Dot(cameraRay) < 0 {
				// Simple illumination via light direction
				lightDirection := Vector4{0, 1, -1, 1}.Normalise()
				ilumination := math.Max(0.1, lightDirection.Dot(normal))

				// Transform and project triangles
				clipped := ClipAgainstPlane(Vector4{0, 0, 0.1, 1}, Vector4{0, 0, 1, 1}, triTransformed)
				for n := 0; n < len(clipped); n++ {
					// Project triangles to 2D
					triProjected := matProj.multiplyTriangle(clipped[n])
					triProjected.ilum = ilumination

					// Apply depth
					triProjected.vecs[0] = triProjected.vecs[0].Div(triProjected.vecs[0].w)
					triProjected.vecs[1] = triProjected.vecs[1].Div(triProjected.vecs[1].w)
					triProjected.vecs[2] = triProjected.vecs[2].Div(triProjected.vecs[2].w)

					// Offset into view
					vOffsetView := Vector4{1, 1, 0, 0}
					triProjected.vecs[0] = triProjected.vecs[0].Add(vOffsetView)
					triProjected.vecs[1] = triProjected.vecs[1].Add(vOffsetView)
					triProjected.vecs[2] = triProjected.vecs[2].Add(vOffsetView)

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
		slices.SortFunc(triangles, func(t1, t2 Triangle) int {
			z1 := (t1.vecs[0].z + t1.vecs[1].z + t1.vecs[2].z) / 3
			z2 := (t2.vecs[0].z + t2.vecs[1].z + t2.vecs[2].z) / 3

			if z1 > z2 {
				return 1
			} else if z1 < z2 {
				return -1
			} else {
				return 0
			}
		})

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
						clipped = ClipAgainstPlane(Vector4{0, 0, 0, 1}, Vector4{0, 1, 0, 1}, test)
					case 1:
						clipped = ClipAgainstPlane(Vector4{0, SCREEN_HEIGHT, 0, 1}, Vector4{0, -1, 0, 1}, test)
					case 2:
						clipped = ClipAgainstPlane(Vector4{0, 0, 0, 1}, Vector4{1, 0, 0, 1}, test)
					case 3:
						clipped = ClipAgainstPlane(Vector4{SCREEN_WIDTH, 0, 0, 1}, Vector4{-1, 0, 0, 1}, test)
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

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
