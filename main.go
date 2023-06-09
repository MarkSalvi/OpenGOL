package main

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"image"
	"log"
	"math"
	"os"
	"runtime"
	"strings"
	"image/jpeg"
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
	uniform vec4 ourColor;
    out vec4 frag_colour;
    void main() {
        frag_colour = ourColor;
    }
` + "\x00"
)

func init() {
	runtime.LockOSThread()
}

type TextureID uint32


func main() {

	var triangle = []float32{
		0.5, 0.5, 0.0,
		0.5, -0.5, 0.0,
		-0.5, -0.5, 0.0,
		-0.5, 0.5, 0.0,
	}

	var indices = []uint32{
		0, 1, 3,
		1, 2, 3,
	}

	var trCoord = []float32{
		0.0, 0.0,  // lower-left corner
		1.0, 0.0,  // lower-right corner
		0.5, 1.0   // top-center corner
	}

	f, err := os.Open("inputimage.jpg")
	if err != nil {
		log.Fatal("No Such Image File Found")
	}
	defer f.Close()

	img, fmtName, err := image.Decode(f)
	if err != nil {
		// Handle error
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
	gl.UseProgram(prog)

	vao, vbo, ebo := makeVao(triangle, indices)

	for !window.ShouldClose() {

		draw(window, prog, vao)

	}
	gl.DeleteProgram(prog)
	gl.DeleteVertexArrays(1, &vao)
	gl.DeleteBuffers(1, &ebo)
	gl.DeleteBuffers(1, &vbo)
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

}

func draw(window *glfw.Window, program uint32, vao uint32) {

	processInput(window)

	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.UseProgram(program)

	timeValue := glfw.GetTime()
	greenValue := math.Sin(timeValue)/2 + 0.5
	vertexColorLocation := gl.GetUniformLocation(program, gl.Str("ourColor"+"\x00"))
	gl.Uniform4f(vertexColorLocation, 0, float32(greenValue), 0, 1)

	gl.BindVertexArray(vao)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.Ptr(nil))

	glfw.PollEvents()
	window.SwapBuffers()
}

func makeVao(points []float32, indices []uint32) (uint32, uint32, uint32) {
	var vao uint32
	var vbo uint32
	var ebo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &ebo)
	gl.GenBuffers(1, &vbo)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.Ptr(nil))
	gl.EnableVertexAttribArray(0)
	return vao, vbo, ebo
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

func processInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	} else if window.GetKey(glfw.KeyF1) == glfw.Press {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	} else if window.GetKey(glfw.KeyF2) == glfw.Press {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}
}

func GenTexture() TextureID{
	var texID uint32
	gl.GenTextures(1,&texID)
	gl.BindTexture(gl.TEXTURE_2D,texID)

	return TextureID(texID)
}

//Todo: std_image faster way
func LoadTextureAlpha(filename string) TextureID {

	infile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer infile.Close()

	img, err := jpeg.Decode(infile)
	if err != nil {
		panic(err)
	}

	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y

	pixels := make([]byte, w*h*4)
	bIndex := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[bIndex] = byte(r / 256)
			bIndex++
			pixels[bIndex] = byte(g / 256)
			bIndex++
			pixels[bIndex] = byte(b / 256)
			bIndex++
			pixels[bIndex] = byte(a / 256)
			bIndex++
		}
	}
	texture := GenTexture()
	gl.TexParameteri(gl.TEXTURE_2D,gl.TEXTURE_WRAP_S,gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D,gl.TEXTURE_WRAP_T,gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D,gl.TEXTURE_MIN_FILTER,gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D,gl.TEXTURE_MAG_FILTER,gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D,0,gl.RGBA,int32(w),int32(h),0,gl.RGBA,gl.UNSIGNED_BYTE,gl.Ptr(pixels))
	gl.GenerateMipmap(gl.TEXTURE_2D)
	return texture

}