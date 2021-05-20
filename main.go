package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/inkyblackness/imgui-go"

	"github.com/oyberntzen/opengl/res/imguiutil"
)

func init() {
	runtime.LockOSThread()
}

const (
	width  int32 = 700
	height int32 = 700
)

var (
	numberOfAgents int     = 50000
	trailWeight    float32 = 0.05
	decayRate      float32 = 0.005
	diffuseRate    float32 = 0.1
	moveSpeed      float32 = 0.001
	turnSpeed      float32 = 1.3
	sensorAngle    float32 = 0.5
	sensorDistance float32 = 0.005
	sensorSize     int32   = 2
)

type agents struct {
	xPositions []float32
	yPositions []float32
	angles     []float32
}

func initAgents(ssbo uint32) {
	data := agents{make([]float32, numberOfAgents), make([]float32, numberOfAgents), make([]float32, numberOfAgents)}
	for i := 0; i < numberOfAgents; i++ {
		data.angles[i] = rand.Float32() * math.Pi * 2
		data.xPositions[i] = 0.5 //rand.Float32()
		data.yPositions[i] = 0.5 //rand.Float32()
	}
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, ssbo)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, numberOfAgents*12, gl.Ptr(data.xPositions), gl.DYNAMIC_COPY)
}

func main() {
	checkError(glfw.Init())
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	window, err := glfw.CreateWindow(700, 700, "Hello World", nil, nil)
	checkError(err)
	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	checkError(gl.Init())

	gl.Enable(gl.DEBUG_OUTPUT)
	gl.DebugMessageCallback(errorCallback, nil)

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.BLEND)

	positions := []float32{
		-1, -1, 0, 0,
		1, -1, 1, 0,
		1, 1, 1, 1,
		-1, 1, 0, 1,
	}

	indicies := []uint32{
		0, 1, 2,
		2, 3, 0,
	}

	vb := CreateVertexBuffer(gl.Ptr(positions), len(positions)*4)
	defer vb.Delete()

	layout := CreateVertexBufferLayout()
	layout.PushFloat32(2)
	layout.PushFloat32(2)

	va := CreateVertexArray(vb, layout)
	defer va.Delete()

	ib := CreateIndexBuffer(indicies)
	defer ib.Delete()

	shader, err := CreateShader("./res/shaders/Basic.shader")
	checkError(err)

	texture := CreateEmptyTexture(width, height)
	defer texture.Delete()

	slot := 0
	texture.Bind(uint32(slot), true)
	shader.Bind()
	shader.SetUniform1i("u_Texture", int32(slot))

	var ssbo uint32
	gl.GenBuffers(1, &ssbo)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, ssbo)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 1, ssbo)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, 0)

	initAgents(ssbo)

	agentsShader, err := CreateComputeShader("./res/shaders/agents.shader")
	checkError(err)
	defer agentsShader.Delete()

	agentsShader.Bind()
	agentsShader.SetUniform1i("raw_texture", 0)
	agentsShader.SetUniform1i("width", width)
	agentsShader.SetUniform1i("height", height)

	textureShader, err := CreateComputeShader("./res/shaders/texture.shader")
	checkError(err)
	defer textureShader.Delete()

	textureShader.Bind()
	textureShader.SetUniform1i("raw_texture", 0)
	textureShader.SetUniform1i("width", width)
	textureShader.SetUniform1i("height", height)

	renderer := Renderer{}

	imguiContext := imgui.CreateContext(nil)
	defer imguiContext.Destroy()
	imguiIO := imgui.CurrentIO()

	imguiPlatform, err := imguiutil.NewGLFW(imguiIO, window)
	checkError(err)
	defer imguiPlatform.Dispose()

	imguiRenderer, err := imguiutil.NewOpenGL3(imguiIO)
	checkError(err)
	defer imguiRenderer.Dispose()

	time := uint32(0)
	for !window.ShouldClose() {
		agentsShader.Bind()

		agentsShader.SetUniform1ui("time", time)
		agentsShader.SetUniform1f("trailWeight", trailWeight)
		agentsShader.SetUniform1f("moveSpeed", moveSpeed)
		agentsShader.SetUniform1f("turnSpeed", turnSpeed)
		agentsShader.SetUniform1f("sensorAngle", sensorAngle)
		agentsShader.SetUniform1f("sensorDistance", sensorDistance)
		agentsShader.SetUniform1i("sensorSize", sensorSize)

		gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, ssbo)
		agentsShader.Run(uint32(numberOfAgents), 1, 1)

		textureShader.Bind()
		textureShader.SetUniform1f("decayRate", decayRate)
		textureShader.SetUniform1f("diffuseRate", diffuseRate)

		textureShader.Run(uint32(width), uint32(height), 1)
		renderer.Clear()

		imguiPlatform.NewFrame()

		imgui.NewFrame()

		{
			imgui.SliderFloat("Trail Weight", &trailWeight, 0, 3)
			imgui.SliderFloat("Decay Rate", &decayRate, 0, 0.1)
			imgui.SliderFloat("Diffuse Rate", &diffuseRate, 0, 1)
			imgui.SliderFloat("Move Speed", &moveSpeed, 0, 0.01)
			imgui.SliderFloat("Turn Speed", &turnSpeed, 0, 1.5)
			imgui.SliderFloat("Sensor Angle", &sensorAngle, 0, 1.5)
			imgui.SliderFloat("Sensor Distance", &sensorDistance, 0, 0.05)
			imgui.SliderInt("Sensor Size", &sensorSize, 1, 10)

			if imgui.Button("Reset") {
				initAgents(ssbo)
			}
		}

		imgui.Render()
		imguiRenderer.PreRender([3]float32{0.0, 0.0, 0.0})

		renderer.Draw(&va, &ib, &shader)

		imguiRenderer.Render(imguiPlatform.DisplaySize(), imguiPlatform.FramebufferSize(), imgui.RenderedDrawData())
		window.SwapBuffers()
		glfw.PollEvents()

		time++
	}
}
