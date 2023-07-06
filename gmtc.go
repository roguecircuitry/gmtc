package main

import (
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/window"
)

type MeshBuilder struct {
	verts, norms, uvs math32.ArrayF32
	inds              math32.ArrayU32
}

func NewMeshBuilder() MeshBuilder {
	result := MeshBuilder{
		verts: math32.NewArrayF32(0, 3*3*64),
		norms: math32.NewArrayF32(0, 3*3*64),
		uvs:   math32.NewArrayF32(0, 2*3*64),
		inds:  math32.NewArrayU32(0, 3*64),
	}
	return result
}

func (mb MeshBuilder) MeshAppend(
	verts []float32,
	norms []float32,
	uvs []float32,
	inds []uint32,
) MeshBuilder {
	mb.verts.Append(verts...)
	mb.norms.Append(norms...)
	mb.uvs.Append(uvs...)
	mb.inds.Append(inds...)
	return mb
}

func (mb MeshBuilder) MeshWrite(
	result *geometry.Geometry,
) *geometry.Geometry {
	result.SetIndices(mb.inds)
	result.AddVBO(gls.NewVBO(mb.verts).AddAttrib(gls.VertexPosition))
	result.AddVBO(gls.NewVBO(mb.norms).AddAttrib(gls.VertexNormal))
	result.AddVBO(gls.NewVBO(mb.uvs).AddAttrib(gls.VertexTexcoord))
	return result
}

func NewChunkMesh() *geometry.Geometry {
	result := geometry.NewGeometry()

	// w, h, d := 8,8,8

	vertexCount := 3

	verts := math32.NewArrayF32(0, vertexCount*3)
	norms := math32.NewArrayF32(0, vertexCount*3)
	uvs := math32.NewArrayF32(0, vertexCount*2)
	indices := math32.NewArrayU32(0, vertexCount)

	verts.Append(
		0, 0, 0,
		1, 0, 0,
		1, 1, 0,
	)
	norms.Append(
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	)
	uvs.Append(
		0, 0,
		1, 0,
		1, 1,
	)
	indices.Append(0, 1, 2)

	result.SetIndices(indices)
	result.AddVBO(gls.NewVBO(verts).AddAttrib(gls.VertexPosition))
	result.AddVBO(gls.NewVBO(norms).AddAttrib(gls.VertexNormal))
	result.AddVBO(gls.NewVBO(uvs).AddAttrib(gls.VertexTexcoord))

	return result
}

func main() {

	// Create application and scene
	a := app.App()
	a.IWindow.(*window.GlfwWindow).SetTitle("go-minetest-client")
	scene := core.NewNode()

	// Set the scene to be managed by the gui manager
	gui.Manager().Set(scene)

	// Create perspective camera
	cam := camera.New(1)
	cam.SetPosition(0, 0, 3)
	scene.Add(cam)

	// Set up orbit control for the camera
	camera.NewOrbitControl(cam)

	// Set up callback to update viewport and camera aspect ratio when the window is resized
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := a.GetSize()
		a.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		cam.SetAspect(float32(width) / float32(height))
	}
	a.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	geom := geometry.NewGeometry()

	NewMeshBuilder().MeshAppend(
		[]float32{0, 0, 0, 1, 0, 0, 1, 1, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 1, 0, 1, 1},
		[]uint32{0, 1, 2},
	).MeshAppend(
		[]float32{0, 0, 0, 1, 1, 0, 0, 1, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 1, 0, 1, 1},
		[]uint32{3, 4, 5},
	).MeshWrite(geom)

	// Create a blue torus and add it to the scene
	// geom := NewChunkMesh() //geometry.NewTorus(1, .4, 12, 32, math32.Pi*2)
	mat := material.NewStandard(math32.NewColor("DarkBlue"))
	mesh := graphic.NewMesh(geom, mat)
	scene.Add(mesh)

	// Create and add a button to the scene
	btn := gui.NewButton("Make Red")
	btn.SetPosition(100, 40)
	btn.SetSize(40, 40)
	btn.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		mat.SetColor(math32.NewColor("DarkRed"))
	})
	scene.Add(btn)

	// Create and add lights to the scene
	scene.Add(light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8))
	pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	pointLight.SetPosition(1, 0, 2)
	scene.Add(pointLight)

	// Create and add an axis helper to the scene
	scene.Add(helper.NewAxes(0.5))

	// Set background color to gray
	a.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)

	// Run the application
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)
	})
}
