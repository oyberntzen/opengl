package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"gopkg.in/gookit/color.v1"
)

type ComputeShader struct {
	rendererID           uint32
	filepath             string
	uniformLocationCache map[string]int32
}

func CreateComputeShader(filepath string) (ComputeShader, error) {
	shader := ComputeShader{}
	shader.filepath = filepath
	shader.uniformLocationCache = make(map[string]int32)

	content, err := ioutil.ReadFile(filepath)
	computeShader, err := shader.compileShader(gl.COMPUTE_SHADER, string(content)+"\x00")
	if err != nil {
		return ComputeShader{}, err
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, computeShader)
	gl.LinkProgram(program)
	gl.ValidateProgram(program)
	shader.rendererID = program

	return shader, nil
}

func (shader *ComputeShader) Bind() {
	gl.UseProgram(shader.rendererID)
}

func (shader *ComputeShader) Unbind() {
	gl.UseProgram(0)
}

func (shader *ComputeShader) Delete() {
	gl.DeleteProgram(shader.rendererID)
}

func (shader *ComputeShader) Run(groupsX, groupsY, groupsZ uint32) {
	shader.Bind()
	gl.DispatchCompute(groupsX, groupsY, groupsZ)
	gl.MemoryBarrier(gl.SHADER_IMAGE_ACCESS_BARRIER_BIT)
}

func (shader *ComputeShader) SetUniform4f(name string, v0 float32, v1 float32, v2 float32, v3 float32) {
	location := shader.getUniformLocation(name)
	if location != -1 {
		gl.Uniform4f(location, v0, v1, v2, v3)
	}
}

func (shader *ComputeShader) SetUniform1i(name string, v0 int32) {
	location := shader.getUniformLocation(name)
	if location != -1 {
		gl.Uniform1i(location, v0)
	}
}

func (shader *ComputeShader) SetUniformMat4f(name string, matrix mgl32.Mat4) {
	location := shader.getUniformLocation(name)
	if location != -1 {
		gl.UniformMatrix4fv(location, 1, false, &matrix[0])
	}
}

func (shader *ComputeShader) SetUniform1f(name string, v0 float32) {
	location := shader.getUniformLocation(name)
	if location != -1 {
		gl.Uniform1f(location, v0)
	}
}

func (shader *ComputeShader) SetUniform1ui(name string, v0 uint32) {
	location := shader.getUniformLocation(name)
	if location != -1 {
		gl.Uniform1ui(location, v0)
	}
}

func (shader *ComputeShader) getUniformLocation(name string) int32 {
	if location, ok := shader.uniformLocationCache[name]; ok {
		return location
	}

	location := gl.GetUniformLocation(shader.rendererID, gl.Str(name+"\x00"))
	if location == -1 {
		color.Warn.Prompt(fmt.Sprintf("Did not find %v uniform", name))
	}
	shader.uniformLocationCache[name] = location
	return location
}

func (shader *ComputeShader) compileShader(shaderType uint32, source string) (uint32, error) {
	id := gl.CreateShader(shaderType)
	csources, free := gl.Strs(source)
	gl.ShaderSource(id, 1, csources, nil)
	free()
	gl.CompileShader(id)

	var result int32
	gl.GetShaderiv(id, gl.COMPILE_STATUS, &result)
	if result == gl.FALSE {
		var length int32
		gl.GetShaderiv(id, gl.INFO_LOG_LENGTH, &length)

		message := strings.Repeat("\x00", int(length+1))
		gl.GetShaderInfoLog(id, length, &length, gl.Str(message))
		gl.DeleteShader(id)
		if shaderType == gl.VERTEX_SHADER {
			return 0, fmt.Errorf("Failed to compile vertex shader %v\n", message)
		} else if shaderType == gl.FRAGMENT_SHADER {
			return 0, fmt.Errorf("Failed to compile fragment shader %v\n", message)
		} else {
			return 0, fmt.Errorf("Failed to compile unknown shader %v\n", message)
		}

	}

	return id, nil
}
