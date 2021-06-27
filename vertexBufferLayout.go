package main

import "github.com/go-gl/gl/v4.6-core/gl"

type LayoutElement struct {
	count      int32
	gltype     uint32
	typeSize   int
	normalized bool
}

type VertexBufferLayout struct {
	elements []LayoutElement
	stride   int32
}

func CreateVertexBufferLayout() VertexBufferLayout {
	layout := VertexBufferLayout{}
	layout.stride = 0
	return layout
}

func (layout *VertexBufferLayout) PushFloat32(count int32) {
	layout.elements = append(layout.elements, LayoutElement{count, gl.FLOAT, 4, false})
	layout.stride += 4 * count
}
func (layout *VertexBufferLayout) PushUint32(count int32) {
	layout.elements = append(layout.elements, LayoutElement{count, gl.UNSIGNED_INT, 4, false})
	layout.stride += 4 * count
}
func (layout *VertexBufferLayout) PushUint8(count int32) {
	layout.elements = append(layout.elements, LayoutElement{count, gl.UNSIGNED_BYTE, 1, true})
	layout.stride += 1 * count
}

func (layout *VertexBufferLayout) GetStride() int32 {
	return layout.stride
}

func (layout *VertexBufferLayout) GetElements() []LayoutElement {
	return layout.elements
}
