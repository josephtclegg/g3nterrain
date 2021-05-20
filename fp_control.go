package main

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
	"math"
)

// FpEnabled specifies which control types are enabled.
type FpEnabled int

// The possible control types.
const (
	FpNone FpEnabled = 0x00
	FpRot  FpEnabled = 0x01
	FpZoom FpEnabled = 0x02
	FpPan  FpEnabled = 0x04
	FpKeys FpEnabled = 0x08
	FpAll  FpEnabled = 0xFF
)

// fpState bitmask
type fpState int

const (
	stateNone = fpState(iota)
	stateRotate
	stateZoom
	statePan
)

type FpControl struct {
	core.Dispatcher           // Embedded event dispatcher
	cam             *camera.Camera   // Controlled camera
	up              math32.Vector3 
	enabled         FpEnabled // Which controls are enabled
	state           fpState   // Current control state

	// Public properties
	MinPolarAngle     float32 // Minimum polar angle in radians (default is 0)
	MaxPolarAngle     float32 // Maximum polar angle in radians (default is Pi)
	MinAzimuthalAngle float32 // Minimum azimuthal angle in radians (default is negative infinity)
	MaxAzimuthalAngle float32 // Maximum azimuthal angle in radians (default is infinity)
	RotSpeed          float32 // Rotation speed factor (default is 1)
	ZoomSpeed         float32 // Zoom speed factor (default is 1)
	KeyRotSpeed       float32 // Rotation delta in radians used on each rotation key event (default is the equivalent of 15 degrees)
	KeyZoomSpeed      float32 // Zoom delta used on each zoom key event (default is 2)
	KeyPanSpeed       float32 // Pan delta used on each pan key event (default is 35)
	
	// Internal
	rotStart  math32.Vector2
	panStart  math32.Vector2
	zoomStart float32
	xrot      float32
	yrot      float32
}

// NewFpControl creates and returns a pointer to a new fp control for the specified camera.
func NewFpControl(cam *camera.Camera) *FpControl {

	fpc := new(FpControl)
	fpc.Dispatcher.Initialize()
	fpc.cam = cam
	fpc.up = *math32.NewVector3(0, 1, 0)
	fpc.enabled = FpAll

	gui.Manager().SetCursorFocus(fpc)
	fpc.state = stateRotate

	fpc.MinPolarAngle = 0
	fpc.MaxPolarAngle = math32.Pi // 180 degrees as radians
	fpc.MinAzimuthalAngle = float32(math.Inf(-1))
	fpc.MaxAzimuthalAngle = float32(math.Inf(1))
	fpc.RotSpeed = 1.0
	fpc.ZoomSpeed = 0.1
	fpc.KeyRotSpeed = 15 * math32.Pi / 180 // 15 degrees as radians
	fpc.KeyZoomSpeed = 2.0
	fpc.KeyPanSpeed = 10.0
	fpc.xrot = 0
	fpc.yrot = 0

	// Subscribe to events
	//gui.Manager().SubscribeID(window.OnMouseUp, &fpc, fpc.onMouse)
	//gui.Manager().SubscribeID(window.OnMouseDown, &fpc, fpc.onMouse)
	gui.Manager().SubscribeID(window.OnKeyDown, &fpc, fpc.onKey)
	gui.Manager().SubscribeID(window.OnKeyRepeat, &fpc, fpc.onKey)
	fpc.SubscribeID(window.OnCursor, &fpc, fpc.onCursor)

	return fpc
}

// Dispose unsubscribes from all events.
func (fpc *FpControl) Dispose() {

	gui.Manager().UnsubscribeID(window.OnMouseUp, &fpc)
	gui.Manager().UnsubscribeID(window.OnMouseDown, &fpc)
	fpc.UnsubscribeID(window.OnCursor, &fpc)
}

// Enabled returns the current FpEnabled bitmask.
func (fpc *FpControl) Enabled() FpEnabled {

	return fpc.enabled
}

// SetEnabled sets the current FpEnabled bitmask.
func (fpc *FpControl) SetEnabled(bitmask FpEnabled) {

	fpc.enabled = bitmask
}

// Rotate rotates the camera around the target by the specified angles.
func (fpc *FpControl) Rotate(thetaDelta, phiDelta float32) {

	//target := math32.NewVector3(0.0, phiDelta, thetaDelta)
	//fpc.cam.LookAt(target, math32.NewVector3(0.0, 1.0, 0.0))

	//fpc.cam.SetRotationZ(0.0)

	fpc.cam.SetRotationZ(2*math32.Pi)
	fpc.cam.SetRotationY(phiDelta)
	fpc.cam.SetRotationX(thetaDelta)
	//globrot := fpc.cam.Rotation()
	//xaxis := globrot.Component(1);
	//roll := globrot.Component(2);
	//angles := math32.NewVec3()
	//angles.X = thetaDelta
	//angles.Y = phiDelta
	//angles.Z = 0.0
	//angle := angles.Length()
	//axis := globrot.Multiply(angles.Normalize())
	//fpc.cam.RotateOnAxis(axis, angle)
//	xaxis := math32.NewVector3(1.0, 0.0, 0.0)
//	yaxis := math32.NewVector3(0.0, 1.0, 0.0)
	//zaxis := math32.NewVector3(0.0, 0.0, 1.0)

	// Why do we rotate on the zaxis?  To undo unwanted camera roll of course!
//	fpc.cam.RotateOnAxis(xaxis, phiDelta)
//	fpc.cam.RotateOnAxis(yaxis, thetaDelta)
	//fpc.cam.RotateOnAxis(zaxis, -1*roll);

}

// Pan pans the camera the specified amount on the plane perpendicular to the viewing direction.
func (fpc *FpControl) Pan(deltaX, deltaY, deltaZ float32) {

	fpc.cam.TranslateX(deltaX)
	fpc.cam.TranslateY(deltaY)
	fpc.cam.TranslateZ(deltaZ)
}

// onMouse is called when an OnMouseDown/OnMouseUp event is received.
func (fpc *FpControl) onMouse(evname string, ev interface{}) {

	// If nothing enabled ignore event
	if fpc.enabled == FpNone {
		return
	}

	switch evname {
	case window.OnMouseDown:
		gui.Manager().SetCursorFocus(fpc)
		mev := ev.(*window.MouseEvent)
		switch mev.Button {
		case window.MouseButtonLeft: // Rotate
			if fpc.enabled&FpRot != 0 {
				fpc.state = stateRotate
				fpc.rotStart.Set(mev.Xpos, mev.Ypos)
			}
		}
	case window.OnMouseUp:
		gui.Manager().SetCursorFocus(nil)
		fpc.state = stateNone
	}
	
}

// onCursor is called when an OnCursor event is received.
func (fpc *FpControl) onCursor(evname string, ev interface{}) {

	// If nothing enabled ignore event
	if fpc.enabled == FpNone || fpc.state == stateNone {
		return
	}

	mev := ev.(*window.CursorEvent)
	switch fpc.state {
	case stateRotate:
		c := -2 * math32.Pi * fpc.RotSpeed / fpc.winSize()
		difx := mev.Xpos-fpc.rotStart.X
		dify := mev.Ypos-fpc.rotStart.Y
		fpc.xrot += c*dify
		fpc.yrot += c*difx

		fpc.Rotate(fpc.xrot,fpc.yrot)
		//fpc.Rotate(c*(mev.Xpos-fpc.rotStart.X),
		//	c*(mev.Ypos-fpc.rotStart.Y))
		fpc.rotStart.Set(mev.Xpos, mev.Ypos)
	case statePan:
		//fpc.Pan(mev.Xpos-fpc.panStart.X,
		//	mev.Ypos0fpc.panStart.Y)
		//fpc.panStart.Set(mev.Xpos, mev.Ypos)
	}
}

// onKey is called when an OnKeyDown/OnKeyRepeat event is received.
func (fpc *FpControl) onKey(evname string, ev interface{}) {

	// If keyboard control is disabled ignore event
	if fpc.enabled == FpNone || fpc.state == stateNone {
		return
	}

	kev := ev.(*window.KeyEvent)

	switch kev.Key {
	case window.KeyW:
		fpc.Pan(0, 0, -fpc.KeyPanSpeed)
	case window.KeyS:
		fpc.Pan(0, 0, fpc.KeyPanSpeed)
	case window.KeyA:
		fpc.Pan(-fpc.KeyPanSpeed, 0, 0)
	case window.KeyD:
		fpc.Pan(fpc.KeyPanSpeed, 0, 0)
	}
}

// winSize returns the window height or width based on the camera reference axis
func (fpc *FpControl) winSize() float32 {

	width, size := window.Get().GetSize()
	if fpc.cam.Axis() == camera.Horizontal {
		size = width
	}
	return float32(size)
}
