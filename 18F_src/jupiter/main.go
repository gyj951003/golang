package main

import (
	"image"
	"math"
	"fmt"
	"os"
	"strconv"
)

//G is the gravitational constant in the gravitational force equation.  It is declared as a "global" constant that can be accessed by all functions.
const G = 6.67408e-11

//data setup.
type Body struct {
	name                             string
	mass, radius                     float64
	position, velocity, acceleration OrderedPair
	red, green, blue                 uint8
}

type OrderedPair struct {
	x, y float64
}

type Universe struct {
	bodies []Body
	width  float64
}

//Angle computes the angle formed by the vectors to two given ordered pairs.
//I give this function to you for free because the math is tedious.
func Angle(p1, p2 OrderedPair) float64 {
	dx := p2.x - p1.x
	dy := p2.y - p1.y
	if dx > 0 && dy > 0 {
		return math.Atan(dy / dx)
	} else if dx > 0 && dy < 0 {
		return 2.0*math.Pi + math.Atan(dy/dx)
	} else if dx < 0 && dy > 0 {
		return math.Pi + math.Atan(dy/dx)
	} else if dx < 0 && dy < 0 {
		return math.Pi + math.Atan(dy/dx)
	} else if dx == 0 && dy < 0 {
		return 3.0 * math.Pi / 2.0
	} else if dx == 0 && dy > 0 {
		return math.Pi / 2.0
	} else if dx > 0 && dy == 0 {
		return 0.0
	} else if dx < 0 && dy == 0 {
		return math.Pi
	} else {
		return 0.0
	}
}

func Dist(b, b2 Body) float64 {
  deltaX := b.position.x - b2.position.x
  deltaY := b.position.y - b2.position.y
  return math.Sqrt(math.Pow(deltaX, 2.0) + math.Pow(deltaY, 2.0))
}

func UpdatePosition(b Body, t float64) OrderedPair {
	var newPos OrderedPair
  newPos.x = b.position.x + b.velocity.x * t
  newPos.y = b.position.y + b.velocity.y * t
  return newPos
}

func UpdateVelocity(b Body, t float64) OrderedPair {
	var newVelocity OrderedPair
	newVelocity.x = b.velocity.x + b.acceleration.x * t + 0.5 * b.acceleration.x * t * t
	newVelocity.y = b.velocity.y + b.acceleration.y * t + + 0.5 * b.acceleration.y * t * t
	return newVelocity
}

func ComputeForce(b, b2 Body) OrderedPair {
	var force OrderedPair
	dist := Dist(b, b2)
	mag := G * b.mass * b2.mass / (dist * dist)
	theta := Angle(b.position, b2.position)
	force.x = mag * math.Cos(theta)
	force.y = mag * math.Sin(theta)
	return force
}

func ComputeNetForce(currentUniverse Universe, b Body) OrderedPair {
	var netForce OrderedPair
	netForce.x = 0
	netForce.y = 0
	for _, body := range currentUniverse.bodies {
		if body.name == b.name {
			continue
		}
		force := ComputeForce(b, body)
		netForce.x += force.x
		netForce.y += force.y
	}
	return netForce
}

func UpdateAcceleration(currentUniverse Universe, b Body) OrderedPair {
	var accel OrderedPair
	force := ComputeNetForce(currentUniverse, b)
	accel.x = force.x / b.mass
	accel.y = force.y / b.mass
	return accel
}


//Update Universe based on the force of gravity
func UpdateUniverse(currentUniverse Universe, time float64) Universe {
	//deep cooy the currentUniverse to the newUniverse 
	var newUniverse Universe
	newUniverse.bodies = make([]Body, 0)
	newUniverse.width = currentUniverse.width
	for _, body := range currentUniverse.bodies {
		newUniverse.bodies = append(newUniverse.bodies, body)
	}

	for i := range currentUniverse.bodies {
		newUniverse.bodies[i].acceleration = UpdateAcceleration(currentUniverse, currentUniverse.bodies[i])
		newUniverse.bodies[i].velocity = UpdateVelocity(currentUniverse.bodies[i], time)
		newUniverse.bodies[i].position = UpdatePosition(currentUniverse.bodies[i], time)
	}
	return newUniverse
}

// Simulate the force within the universe, update bodies's motions and status
func SimulateGravity(initialUniverse Universe, numSteps int, time float64) []Universe {
	timePoints := make([]Universe, numSteps + 1)
	timePoints[0] = initialUniverse
	for i := 1; i <= numSteps; i++ {
		timePoints[i] = UpdateUniverse(timePoints[i-1], time)
	}
	return timePoints
}


//AnimateSystem takes a collection of Universe objects along with a canvas width
//parameter and generates a slice of images corresponding to drawing each Universe
//on a canvasWidth x canvasWidth canvas
func AnimateSystem(timePoints []Universe, canvasWidth int, k int) []image.Image {
	images := make([]image.Image, 0)
	for i, u := range timePoints {
		if i % k == 0 {
			images = append(images, DrawToCanvas(u, canvasWidth))
		}
	}
	return images
}

//DrawToCanvas generates the image corresponding to a canvas after drawing a Universe
//object's bodies on a square canvas that is canvasWidth pixels x canvasWidth pixels
func DrawToCanvas(u Universe, canvasWidth int) image.Image {
	c := CreateNewCanvas(canvasWidth, canvasWidth)

	// first, make a black background
	c.SetFillColor(MakeColor(0, 0, 0))
	c.ClearRect(0, 0, canvasWidth, canvasWidth)
	c.Fill()

	for i, b := range u.bodies {
		// make a white circle at b's position with the appropriate width
		c.SetFillColor(MakeColor(b.red, b.green, b.blue))
		if i == 0 {
			c.Circle((b.position.x/u.width)*float64(canvasWidth), (b.position.y/u.width)*float64(canvasWidth), (b.radius/u.width)*float64(canvasWidth))
		} else {
			c.Circle((b.position.x/u.width)*float64(canvasWidth), (b.position.y/u.width)*float64(canvasWidth), (b.radius/u.width)*float64(canvasWidth)*10)
		}
		c.Fill()
	}
	return c.img
}

func main() {
	// declaring objects
	var jupiter, io, europa, ganymede, callisto Body

	jupiter.name = "Jupiter"
	io.name = "Io"
	europa.name = "Europa"
	ganymede.name = "Ganymede"
	callisto.name = "Callisto"

	jupiter.red, jupiter.green, jupiter.blue = 223, 227, 202
	io.red, io.green, io.blue = 249, 249, 165
	europa.red, europa.green, europa.blue = 132, 83, 52
	ganymede.red, ganymede.green, ganymede.blue = 76, 0, 153
	callisto.red, callisto.green, callisto.blue = 0, 153, 76

	jupiter.mass = 1.898 * math.Pow(10, 27)
	io.mass = 8.9319 * math.Pow(10, 22)
	europa.mass = 4.7998 * math.Pow(10, 22)
	ganymede.mass = 1.4819 * math.Pow(10, 23)
	callisto.mass = 1.0759 * math.Pow(10, 23)

	jupiter.radius = 71000000
	io.radius = 1821000
	europa.radius = 1569000
	ganymede.radius = 2631000
	callisto.radius = 2410000

	jupiter.position.x, jupiter.position.y = 2000000000, 2000000000
	io.position.x, io.position.y = 2000000000-421600000, 2000000000
	europa.position.x, europa.position.y = 2000000000, 2000000000+670900000
	ganymede.position.x, ganymede.position.y = 2000000000+1070400000, 2000000000
	callisto.position.x, callisto.position.y = 2000000000, 2000000000-1882700000

	jupiter.velocity.x, jupiter.velocity.y = 0, 0
	io.velocity.x, io.velocity.y = 0, -17320
	europa.velocity.x, europa.velocity.y = -13740, 0
	ganymede.velocity.x, ganymede.velocity.y = 0, 10870
	callisto.velocity.x, callisto.velocity.y = 8200, 0

	var jupiterSystem Universe
	jupiterSystem.width = 4000000000
	jupiterSystem.bodies = []Body{jupiter, io, europa, ganymede, callisto}

	numSteps, err1 := strconv.Atoi(os.Args[1])
	if err1 != nil {
		fmt.Println("Error parsing numSteps!")
		os.Exit(1)
	}

	time, err2 := strconv.ParseFloat(os.Args[2], 64)
	if err2 != nil {
		fmt.Println("Error parsing time!")
		os.Exit(1)
	}

	canvasWidth, err3 := strconv.Atoi(os.Args[3])
	if err3 != nil {
		fmt.Println("Error parsing canvasWidth!")
		os.Exit(1)
	}

	k, err4 := strconv.Atoi(os.Args[4])
	if err4 != nil {
		fmt.Println("Error parsing k!")
		os.Exit(1)
	}

	timePoints := SimulateGravity(jupiterSystem, numSteps, time)

	fmt.Println("Simulation done! Now draw animations!")


	image := AnimateSystem(timePoints, canvasWidth, k)

	Process(image, "jupiter")
}
