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

	for !window.ShouldClose() {

		window.SwapBuffers()
		glfw.PollEvents()
	}

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