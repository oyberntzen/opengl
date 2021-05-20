package main

import (
	"github.com/go-gl/gl/v4.6-core/gl"
)

type VertexArray struct {
	rendererID uint32
}

func CreateVertexArray(buffer VertexBuffer, layout VertexBufferLayout) VertexArray {
	array := VertexArray{}
	gl.GenVertexArrays(1, &array.rendererID)
	array.Bind()

	buffer.Bind()
	var i uint32
	offset := uintptr(0)
	for i = 0; int(i) < len(layout.elements); i++ {
		element := layout.elements[i]
		gl.EnableVertexAttribArray(i)
		gl.VertexAttribPointerWithOffset(i, element.count, element.gltype, element.normalized, layout.GetStride(), offset)
		offset += uintptr(element.count) * uintptr(element.typeSize)
	}

	return array
}

func (array *VertexArray) Bind() {
	gl.BindVertexArray(array.rendererID)
}

func (array *VertexArray) Unbind() {
	gl.BindVertexArray(0)
}

func (array *VertexArray) Delete() {
	gl.DeleteVertexArrays(1, &array.rendererID)
}
