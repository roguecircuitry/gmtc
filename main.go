package main

import (
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

	mb := NewMeshBuilder()

	//create a chunk
	chunk := NewChunk()

	//loop over each block
	chunk.ForEachBlock(func(block *BlockData) {
		//do some math to figure out what the block type should be
		if block.y < 5 {
			block.data.id = 1
		} else if block.y < 6 && rand.Float32() < 0.9 {
			block.data.id = 1
		} else {
			block.data.id = 0
		}
		//replace the block
		chunk.SetBlockData(block)
	})

	//loop over each block and cube info (calculates neighbors automatically)
	chunk.ForEachCubeInfo(func(cubeInfo *CubeInfo, block *BlockData) {
		//if block is renderable
		if block.data.id != 0 {
			mb.Cube(cubeInfo) //draw a cuboid, taking neighbor occlusion into account
		}
	})

	//recalculate normals (currently not working properly)
	mb.AutoNorms()

	//write mesh to geometry
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
	pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 50.0)
	pointLight.SetPosition(5, 5, 5)
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
