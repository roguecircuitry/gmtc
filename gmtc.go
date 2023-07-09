package main

import (
	"fmt"
	"math"
	"math/rand"
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

func NewMeshBuilder() *MeshBuilder {
	triCountSoftMax := 64
	floatsPerVertex := 3
	floatsPerUV := 2
	vertsPerTri := 3

	result := MeshBuilder{
		verts: math32.NewArrayF32(0, floatsPerVertex*vertsPerTri*triCountSoftMax), //3 floats per vertex * 3 vertices per triangle * 64 triangles
		norms: math32.NewArrayF32(0, floatsPerVertex*triCountSoftMax),             //3 floats per vertex * 1 vertex per triangle * 64 triangles
		uvs:   math32.NewArrayF32(0, floatsPerUV*vertsPerTri*triCountSoftMax),
		inds:  math32.NewArrayU32(0, vertsPerTri*triCountSoftMax),
	}
	return &result
}

func (mb *MeshBuilder) MeshAppend(
	verts []float32,
	norms []float32,
	uvs []float32,
	inds []uint32,
) *MeshBuilder {
	offset := uint32(mb.verts.Size() / 3)

	for i, ind := range inds {
		inds[i] = ind + offset
	}

	mb.verts.Append(verts...)
	mb.norms.Append(norms...)
	mb.uvs.Append(uvs...)
	mb.inds.Append(inds...)
	return mb
}

func (mb *MeshBuilder) MeshWrite(
	result *geometry.Geometry,
) *geometry.Geometry {
	result.SetIndices(mb.inds)
	result.AddVBO(gls.NewVBO(mb.verts).AddAttrib(gls.VertexPosition))
	result.AddVBO(gls.NewVBO(mb.norms).AddAttrib(gls.VertexNormal))
	result.AddVBO(gls.NewVBO(mb.uvs).AddAttrib(gls.VertexTexcoord))
	return result
}

type CubeInfo struct {
	top, bottom, north, south, east, west bool
	min, max                              math32.Vector3
}

func calcTriNormal(a, b, c, out *math32.Vector3) {
	A := b.Clone().Sub(a)
	B := c.Clone().Sub(a)

	out.Set(
		A.Y*B.Z-A.Z*B.Y,
		A.Z*B.X-A.X*B.Z,
		A.X*B.Y-A.Y*B.X,
	)
}

func (mb *MeshBuilder) Cube(c *CubeInfo) *MeshBuilder {
	minx := c.min.X
	miny := c.min.Y
	minz := c.min.Z

	maxx := c.max.X
	maxy := c.max.Y
	maxz := c.max.Z

	//MIN-Z plane
	if c.south {
		mb.MeshAppend(
			[]float32{
				minx, miny, minz,
				minx, maxy, minz,
				maxx, maxy, minz,
				maxx, miny, minz,
			},
			[]float32{
				0, 0, 0,
				0, 0, 0,
			},
			[]float32{
				minx, miny,
				minx, maxy,
				maxx, maxy,
				maxx, miny,
			},
			[]uint32{
				0, 1, 2,
				0, 2, 3,
			},
		)
	}

	//MAX-Z plane
	if c.north {
		mb.MeshAppend(
			[]float32{
				minx, miny, maxz,
				minx, maxy, maxz,
				maxx, maxy, maxz,
				maxx, miny, maxz,
			},
			[]float32{
				0, 0, 0,
				0, 0, 0,
			},
			[]float32{
				minx, miny,
				minx, maxy,
				maxx, maxy,
				maxx, miny,
			},
			[]uint32{
				0, 2, 1,
				0, 3, 2,
			},
		)
	}

	//MAX-X plane
	if c.west {
		mb.MeshAppend(
			[]float32{
				maxx, miny, minz,
				maxx, miny, maxz,
				maxx, maxy, maxz,
				maxx, maxy, minz,
			},
			[]float32{
				0, 0, 0,
				0, 0, 0,
			},
			[]float32{
				minx, miny,
				minx, maxy,
				maxx, maxy,
				maxx, miny,
			},
			[]uint32{
				0, 2, 1,
				0, 3, 2,
			},
		)
	}

	//MIN-X plane
	if c.east {
		mb.MeshAppend(
			[]float32{
				minx, miny, minz,
				minx, miny, maxz,
				minx, maxy, maxz,
				minx, maxy, minz,
			},
			[]float32{
				0, 0, 0,
				0, 0, 0,
			},
			[]float32{
				minx, miny,
				minx, maxy,
				maxx, maxy,
				maxx, miny,
			},
			[]uint32{
				0, 1, 2,
				0, 2, 3,
			},
		)
	}

	//MAX-Y plane
	if c.top {
		mb.MeshAppend(
			[]float32{
				minx, maxy, minz,
				maxx, maxy, minz,
				maxx, maxy, maxz,
				minx, maxy, maxz,
			},
			[]float32{
				0, 0, 0,
				0, 0, 0,
			},
			[]float32{
				minx, miny,
				minx, maxy,
				maxx, maxy,
				maxx, miny,
			},
			[]uint32{
				0, 2, 1,
				0, 3, 2,
			},
		)
	}

	//MIN-Y plane
	if c.bottom {
		mb.MeshAppend(
			[]float32{
				minx, miny, minz,
				maxx, miny, minz,
				maxx, miny, maxz,
				minx, miny, maxz,
			},
			[]float32{
				0, 0, 0,
				0, 0, 0,
			},
			[]float32{
				minx, miny,
				minx, maxy,
				maxx, maxy,
				maxx, miny,
			},
			[]uint32{
				0, 1, 2,
				0, 2, 3,
			},
		)
	}

	return mb
}

func Vector3FromArray(out *math32.Vector3, from math32.ArrayF32, offset uint32) {
	out.Set(
		from[offset],
		from[offset+1],
		from[offset+2],
	)
}

type TriInfo struct {
	a, b, c math32.Vector3
}

func (mb *MeshBuilder) GetTriFromInd(ind uint32, out *TriInfo) {
	ai := mb.inds[ind]
	bi := mb.inds[ind+1]
	ci := mb.inds[ind+2]

	Vector3FromArray(&out.a, mb.verts, ai)
	Vector3FromArray(&out.b, mb.verts, bi)
	Vector3FromArray(&out.b, mb.verts, ci)
}

func (mb *MeshBuilder) GetTri(triIndex int, out *TriInfo) {
	ind := uint32(triIndex * 3)
	mb.GetTriFromInd(ind, out)
}

func TriIndexToVertexIndex(triIndex int) int {
	return triIndex * 3
}

func VertIndexToTriIndex(ind int) int {
	return ind / 3
}

func (mb *MeshBuilder) ForEachTri(cb func(ind uint32, triIndex int, tri *TriInfo)) {
	tri := new(TriInfo)

	for i := 0; i < mb.inds.Len(); i += 3 {
		ind := mb.inds[i]
		mb.GetTriFromInd(ind, tri)
		triIndex := VertIndexToTriIndex(i)

		cb(ind, triIndex, tri)
	}
}

func (mb *MeshBuilder) AutoNorms() *MeshBuilder {
	normal := math32.NewVec3()

	mb.ForEachTri(func(ind uint32, triIndex int, tri *TriInfo) {

		calcTriNormal(&tri.a, &tri.b, &tri.c, normal)

		fmt.Printf("Triangle #%v: %v, normal: %v\n", triIndex, tri, normal)

		normalIndex := triIndex * 3

		mb.norms[normalIndex] = normal.X
		mb.norms[normalIndex+1] = normal.Y
		mb.norms[normalIndex+2] = normal.Z
	})
	return mb
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

	cubeInfo := &CubeInfo{
		top:    true,
		bottom: true,
		north:  true,
		south:  true,
		east:   true,
		west:   true,
		min:    *math32.NewVector3(0, 0, 0),
		max:    *math32.NewVector3(1, 1, 1),
	}

	mb := NewMeshBuilder()

	for i := 0; i < 10; i++ {
		mb.Cube(cubeInfo)
		cubeInfo.min.Set(
			float32(math.Floor(float64(rand.Float32()*10))),
			float32(math.Floor(float64(rand.Float32()*10))),
			float32(math.Floor(float64(rand.Float32()*10))),
		)

		cubeInfo.max.Copy(&cubeInfo.min).AddScalar(1)
	}

	mb.AutoNorms()

	fmt.Printf("mb.norms: %v\n", mb.norms)

	mb.MeshWrite(geom)

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
