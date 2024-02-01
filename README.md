<h1 align="center">3d-viewer</h1>

<p align="center">
  <a href="./LICENSE"><img src="https://img.shields.io/badge/⚖️ license-MIT-blue" alt="MIT License"></a>
  <img src="https://img.shields.io/github/stars/keelus/3d-viewer?color=red&logo=github" alt="stars">
</p>

## ℹ️ Description
A simple application to preview 3D models (currently .OBJ supported), where you can move the model around, flip normals, etc. It's implemented in Golang, and uses [SDL2](https://www.libsdl.org/) to show the window, handle user's input, and write pixels into the screen buffer.

I made it from scratch to learn 3D rendering concepts, mathematics and algorithms involved to display a 3D textured mesh into the screen.

## ✨ Features
- Fast 3D .obj file loading.
- Fast and smooth rendering.
- Simple camera system to move and rotate the object around.
- Buttons to change the resolution (to gain performance for more complex objects)
- Support for .obj 3D files and .mtl material files (with PNG and JPEG texture formats).

## 🐛 Known errors
- In very specific cases where the Z position of the camera is exactly 0, and rotation is default, some triangles might not be displayed correctly. This can be corrected by just moving or rotating the camera by a few pixels.
- In strange cases, when loading a 3D model, SDL2 could fail to render the mesh text information (unknown cause).

## 📸 Screenshots
<p align="center">
  <img src="https://github.com/keelus/go-3d-viewer/assets/86611436/3fd531f6-151e-4649-bc8c-a5341c047af2" style="width:75%;" /><br>
  <img src="https://github.com/keelus/go-3d-viewer/assets/86611436/344c17ce-c4ef-4510-8d41-7bb9954fdad4" style="width:75%;" />
</p>



## 🔨 Requirements
To use and/or compile the application, you will need to have [SDL2](https://github.com/libsdl-org/SDL/releases/latest) and [SDL2_TTF](https://github.com/libsdl-org/SDL_ttf/releases/latest) installed correctly in your system.

## ⬇️ Install & run it
The project is compatible with Windows, Linux and macOS, when requeriments are installed.

To use it, simply download the [latest release](https://github.com/keelus/3d-viewer/releases/latest) binary file and execute it (unzip and execute on Windows).

### 🐧 Linux & macOS
To make the downloaded binary executable, run:
```bash
chmod +x 3d_viewer-<rest of the file>
```
In newer versions of macOS, you might need to run `xattr -c 3d_viewer-<rest of the filename>` if you get an error message while opening the app.

Then, you can open it running:
```bash
./3d_viewer-<rest of the file>
```

## Compile
To compile the project, you will need SDL2 and SDL2_TTF properly installed in your system. Also, a C compiler could be needed (such as [GCC](https://gcc.gnu.org/)).
If you encounter any issues while compiling, please check [go-sdl2](https://github.com/veandco/go-sdl2) compiling guide.
### 🪟 Windows
To compile the app, just run:
```bash
go mod tidy
```
```bash
go build -o 3d_viewer.exe -ldflags "-s -w -H windowsgui"
```
Make sure to have `SDL2.dll` and `SDL2_ttf.dll` files in the same place of the `.exe`.
### 🐧 Linux or macOS
To compile the app, just run:
```bash
go mod tidy
```
```bash
go build -o 3d_viewer -ldflags "-s -w"
```

## 📰 References

## ⚠️ Disclaimer
The 3D models that are shown in the screenshots are for demonstration purposes only, I don't own the models.

## ⚖️ License
This project is open source under the terms of the [MIT License](./LICENSE)

<br />

Made by <a href="https://github.com/keelus">keelus</a> ✌️

](https://github.com/libsdl-org/SDL_ttf/releases/latest)https://github.com/libsdl-org/SDL_ttf/releases/latest