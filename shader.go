package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"gopkg.in/gookit/color.v1"
)

type Shader struct {
	rendererID           uint32
	filepath             string
	uniformLocationCache map[string]int32
}

func CreateShader(filepath string) (Shader, error) {
	shader := Shader{}
	shader.filepath = filepath
	shader.uniformLocationCache = make(map[string]int32)

	vertexShader, fragmentShader, err := shader.parseShader()
	if err != nil {
		return Shader{}, err
	}

	shader.rendererID, err = shader.createShader(vertexShader, fragmentShader)
	if err != nil {
		return Shader{}, err
	}

	return shader, nil
}

func (shader *Shader) Bind() {
	gl.UseProgram(shader.rendererID)
}

func (shader *Shader) Unbind() {
	gl.UseProgram(0)
}

func (shader *Shader) Delete() {
	gl.DeleteProgram(shader.rendererID)
}

func (shader *Shader) SetUniform4f(name string, v0 float32, v1 float32, v2 float32, v3 float32) {
	location := shader.getUniformLocation(name)
	if location != -1 {
		gl.Uniform4f(location, v0, v1, v2, v3)
	}
}

func (shader *Shader) SetUniform1i(name string, v0 int32) {
	location := shader.getUniformLocation(name)
	if location != -1 {
		gl.Uniform1i(location, v0)
	}
}

func (shader *Shader) SetUniformMat4f(name string, matrix mgl32.Mat4) {
	location := shader.getUniformLocation(name)
	if location != -1 {
		gl.UniformMatrix4fv(location, 1, false, &matrix[0])
	}
}

func (shader *Shader) getUniformLocation(name string) int32 {
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

func (shader *Shader) compileShader(shaderType uint32, source string) (uint32, error) {
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
func (shader *Shader) parseShader() (string, string, error) {
	file, err := os.Open(shader.filepath)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	type shaderType int
	const NONE shaderType = -1
	const VERTEX shaderType = 0
	const FRAGMENT shaderType = 1

	var shaderStrings [2]string
	currentShaderType := NONE

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "#shader") {
			if strings.Contains(line, "vertex") {
				currentShaderType = VERTEX
			} else if strings.Contains(line, "fragment") {
				currentShaderType = FRAGMENT
			}
		} else {
			shaderStrings[currentShaderType] = shaderStrings[currentShaderType] + line + "\n"
		}
	}

	if err := scanner.Err(); err != nil {
		return "", "", err
	}

	return shaderStrings[0] + "\x00", shaderStrings[1] + "\x00", nil
}

func (shader *Shader) createShader(vertexShader, fragmentShader string) (uint32, error) {
	program := gl.CreateProgram()

	vs, err := shader.compileShader(gl.VERTEX_SHADER, vertexShader)
	if err != nil {
		return 0, err
	}
	fs, err := shader.compileShader(gl.FRAGMENT_SHADER, fragmentShader)
	if err != nil {
		return 0, err
	}

	gl.AttachShader(program, vs)
	gl.AttachShader(program, fs)
	gl.LinkProgram(program)
	gl.ValidateProgram(program)

	gl.DeleteShader(vs)
	gl.DeleteShader(fs)

	return program, nil
}
