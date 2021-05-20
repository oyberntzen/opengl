package main

import (
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type VertexBuffer struct {
	rendererID uint32
}

func CreateVertexBuffer(data unsafe.Pointer, size int) VertexBuffer {
	buffer := VertexBuffer{}
	gl.GenBuffers(1, &buffer.rendererID)
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer.rendererID)
	gl.BufferData(gl.ARRAY_BUFFER, size, data, gl.STATIC_DRAW)
	return buffer
}

func (buffer *VertexBuffer) Bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer.rendererID)
}

func (buffer *VertexBuffer) Unbind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (buffer *VertexBuffer) Delete() {
	gl.DeleteBuffers(1, &buffer.rendererID)
}
