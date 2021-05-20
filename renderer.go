package main

import (
	"fmt"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
	"gopkg.in/gookit/color.v1"
)

type Renderer struct {
}

func (renderer *Renderer) Clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func (renderer *Renderer) Draw(va *VertexArray, ib *IndexBuffer, shader *Shader) {
	va.Bind()
	ib.Bind()
	shader.Bind()

	gl.DrawElements(gl.TRIANGLES, ib.GetCount(), gl.UNSIGNED_INT, nil)
}

func clearGLError() {
	for gl.GetError() != gl.NO_ERROR {

	}
}

func checkGLError() {
	for err := gl.GetError(); err != gl.NO_ERROR; {
		color.Error.Prompt(fmt.Sprintf("[OpengGL Error] (%v)", err))
		panic(nil)
	}
}

func errorCallback(source, gltype, id, severity uint32, length int32, message string, userParam unsafe.Pointer) {
	if gltype == gl.DEBUG_TYPE_ERROR {
		color.Error.Prompt(fmt.Sprintf("[OpengGL Error] type: %v, severity: %v, message: '%v'", gltype, severity, message))
		panic(nil)
	} else {
		//color.Info.Prompt(fmt.Sprintf("[OpengGL Message] type: %v, severity: %v, message: '%v'", gltype, severity, message))
	}
}

func checkError(err error) {
	if err != nil {
		color.Error.Prompt(err.Error())
		panic(nil)
	}
}
