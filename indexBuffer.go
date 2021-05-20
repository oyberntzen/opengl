package main

import (
	"github.com/go-gl/gl/v4.6-core/gl"
)

type IndexBuffer struct {
	rendererID uint32
	count      int32
}

func CreateIndexBuffer(indicies []uint32) IndexBuffer {
	buffer := IndexBuffer{}
	buffer.count = int32(len(indicies))
	gl.GenBuffers(1, &buffer.rendererID)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, buffer.rendererID)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indicies)*4, gl.Ptr(indicies), gl.STATIC_DRAW)
	return buffer
}

func (buffer *IndexBuffer) Bind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, buffer.rendererID)
}

func (buffer *IndexBuffer) Unbind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}

func (buffer *IndexBuffer) Delete() {
	gl.DeleteBuffers(1, &buffer.rendererID)
}

func (buffer *IndexBuffer) GetCount() int32 {
	return buffer.count
}
