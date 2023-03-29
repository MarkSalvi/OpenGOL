package main

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"log"
	"runtime"
)

func init() {
	runtime.LockOSThread()
}

func main() {

	window := initGlfw()
	defer glfw.Terminate()
	program := initOpenGL()

	for !window.ShouldClose() {

		draw(window, program)

	}

}

func draw(window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(0.2, 0.3, 0.3, 1.0)

	gl.UseProgram(program)

	glfw.PollEvents()
	window.SwapBuffers()
}

func initGlfw() *glfw.Window {
	err := glfw.Init()
	if err != nil {
		log.Fatal("GLFW INIT ERROR: ", err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	window, err := glfw.CreateWindow(800, 600, "LearnOpenGl", nil, nil)
	defer glfw.Terminate()
	if err != nil {
		log.Fatal("GLFW WINDOW ERROR: ", err)
	}

	window.MakeContextCurrent()
	return window
}

func initOpenGL() uint32 {

	err := gl.Init()
	if err != nil {
		log.Fatal("OPENGL INIT ERROR: ", err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	prog := gl.CreateProgram()
	gl.LinkProgram(prog)
	return prog

}
