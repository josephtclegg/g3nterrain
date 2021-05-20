package main

import (
	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	//"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/renderer/shaders"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/window"
	"time"
	//"math/rand"
	"image"
	"image/png"
	"os"
	"io"
	"fmt"
)

const turner_fragment = `precision highp float;

in vec3 Color;
out vec4 FragColor;

void main() {
  vec4 turnerOrange = vec4(1.0, 0.7, 0.5, 1.0);
  vec4 turnerBlue = vec4(0.149, 0.141, 0.912, 1.0);
  vec4 red = vec4(1.0, 0.0, 0.0, 1.0);
  vec4 green = vec4(0.0, 1.0, 0.0, 1.0);
  vec4 blue = vec4(0.0, 0.0, 1.0, 1.0);

  float interp = mod(gl_FragCoord.x, 3);
  interp = interp/3;

  float rgb1 = mod(gl_FragCoord.x, 6);
  float rgb2 = mod(gl_FragCoord.y, 6);

  if(rgb1 <= 2.0 && rgb2 <= 2.0){
    FragColor = red;
  } else if(rgb1 <= 4.0 && rgb2 <= 4.0) {
    FragColor = green;
  } else {
    FragColor = blue;
  }
}
`

const turner_vertex = `#include <attributes>

uniform mat4 MVP;

out vec3 Color;

void main() {
  Color = VertexColor;
  gl_Position = MVP * vec4(VertexPosition, 1.0);
}
`

func init() {
	shaders.AddShader("turner_vertex", turner_vertex)
	shaders.AddShader("turner_fragment", turner_fragment)
	shaders.AddProgram("turner", "turner_vertex", "turner_fragment")
}

func main() {
	// Create application and scene
	a := app.App()
	scene := core.NewNode()

	// Set the scene to be managed by the gui manager
	gui.Manager().Set(scene)

	// Create perspective camera
	cam := camera.New(1)
	cam.SetPosition(0, 0, 3)
	scene.Add(cam)

	// Set up the orbit control for the camera
	NewFpControl(cam)

	// Set up callback to update viewport and camera aspect ratio when
	// the window is resized
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := a.GetSize()
		a.Gls().Viewport(0, 0, int32(width*2), int32(height*2))
		// Update the camera's aspect ratio
		cam.SetAspect(float32(width) / float32(height))
	}
	a.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	file, err := os.Open("./hlm.png")
	if err != nil {
		fmt.Println("Error: Could not open height map!")
		os.Exit(1)
	}
	defer file.Close()
	pixels, err := getPixels(file)
	if err != nil {
		fmt.Println("Error: Could not load pixel vals!")
		os.Exit(1)
	}

	points := make([][]int, len(pixels))

	for i:=0; i < len(pixels); i++ {
		points[i] = make([]int, len(pixels[i]))
		for j:=0; j < len(points[i]); j++ {
			points[i][j] = pixels[i][j].R/3
		}
	}

	geom := NewSwathe(points)
	mat := material.NewStandard(math32.NewColor("LightBlue"))
	//mat2 := material.NewMaterial()
	//mat2.Init()
	//mat2.SetShader("turner")
	//mat2.SetShaderUnique(true)
	mesh := graphic.NewMesh(geom, mat)
	scene.Add(mesh)

	// Create and add a button to the scene
	scene.Add(light.NewAmbient(&math32.Color{1.0,1.0,1.0}, 0.5))
	//pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 5000.0)
	//pointLight.SetPosition(float32(len(points)/2), 100, float32(len(points)/2))

	//scene.Add(pointLight)

	lightColor := math32.NewColor("white");
	
	directionalLight := light.NewDirectional(lightColor, 0.7)
	directionalLight.SetPosition(float32(0.0), 100, float32(0.0))
	scene.Add(directionalLight);

	// Create and add an axis helper to the scene
	scene.Add(helper.NewAxes(0.5))

	// Set background color to gray
	a.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)

	// Run the application
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT |
			gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)
	})
}

func getPixels(file io.Reader) ([][]Pixel, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	var pixels [][]Pixel
	for y:=0; y<height; y++ {
		var row []Pixel
		for x:=0; x<width; x++ {
			row = append(row, rgbaToPixel(img.At(x,y).RGBA()))
		}
		pixels = append(pixels,row)
	}

	return pixels, nil
}

func rgbaToPixel(r, g, b, a uint32) Pixel {
	return Pixel{int(r/257), int(g/257), int(b/257), int(a/257)}
}

type Pixel struct {
	R int
	G int
	B int
	A int
}
