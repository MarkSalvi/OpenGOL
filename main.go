package main

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"log"
	"runtime"
	"strings"
)

const (
	vertexShaderSource = `
    #version 330 core
    layout (location = 0) in vec3 vp;
    void main() {
        gl_Position = vec4(vp, 1.0);
    }
` + "\x00"

	fragmentShaderSource = `
    #version 330 core
    out vec4 frag_colour;
    void main() {
        frag_colour = vec4(1, 0.5, 0.2, 1);
    }
` + "\x00"
)

func init() {
	runtime.LockOSThread()
}

func main() {

	var triangle = []float32{
		0.5, 0.5, 0.0, // top right
		0.5, -0.5, 0.0, // bottom right
		-0.5, -0.5, 0.0, // bottom left
		-0.5, 0.5, 0.0,
	}

	var indices = []uint{
		0, 1, 3,
		1, 2, 3,
	}

	err := glfw.Init()
	if err != nil {
		log.Fatal("GLFW INIT ERROR: ", err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	window, err := glfw.CreateWindow(800, 600, "Sium", nil, nil)
	defer glfw.Terminate()
	if err != nil {
		log.Fatal("GLFW WINDOW ERROR: ", err)
	}

	window.MakeContextCurrent()

	err1 := gl.Init()
	if err1 != nil {
		log.Fatal("OPENGL INIT ERROR: ", err1)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	prog := gl.CreateProgram()
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		log.Fatal("vertexShader compile Error: ")
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		log.Fatal("fragmentShader compile Error: ")
	}
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)

	vao := makeVao(triangle, indices)

	for !window.ShouldClose() {

		draw(window, prog, vao)

	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

}

func draw(window *glfw.Window, program uint32, vao uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(program)
	gl.BindVertexArray(vao)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.Ptr(nil))

	glfw.PollEvents()
	window.SwapBuffers()
}

func makeVao(points []float32, indices []uint) uint32 {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)
	return vao
}

func compileShader(source string, fragmentType uint32) (uint32, error) {
	shader := gl.CreateShader(fragmentType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		loge := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(loge))

		return 0, fmt.Errorf("failed to compile %v: %v", source, loge)
	}

	return shader, nil
}
