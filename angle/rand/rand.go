//  Copyright 2020 Marius Ackerman
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

/*
Package rand generates random angles.
*/
package rand

import (
	"math"
	"math/rand"

	"github.com/goccmack/goutil/angle"
)

/*
Return a random angle in radians sampled uniformly from [0, 2π).
*/
func Angle() float64 {
	return rand.Float64() * 2 * math.Pi
}

/*
Return a random angle in degrees sampled uniformly from [0, 360).
*/
func AngleDeg() float64 {
	return rand.Float64() * 360
}

/*
Return a random angle in [min, max) sampled from a Normal distribution.
All values are in degrees.
*/
func AngleNormDeg(mean, sdev float64) float64 {
	return angle.ToDeg(AngleNorm(angle.ToRad(mean), angle.ToRad(sdev)))
}

/*
Return a random angle in [min, max) sampled from a Normal distribution.
All values are in radians.
*/
func AngleNorm(mean, sdev float64) float64 {
	return math.Mod(rand.NormFloat64()*sdev+mean+2*math.Pi, 2*math.Pi)
}

/*
Return n random angles uniformly distributed over [0,2π)
*/
func Angles(n int) []float64 {
	angles := make([]float64, n)
	for i := 0; i < n; i++ {
		angles[i] = Angle()
	}
	return angles
}

/*
Return n random angles in radians normally distributed around mean with
standard deviation sdev.
*/
func AnglesNorm(n int, mean, sdev float64) []float64 {
	angles := make([]float64, n)
	for i := 0; i < n; i++ {
		angles[i] = AngleNorm(mean, sdev)
	}
	return angles
}

/*
Return n random angles in degrees normally distributed around mean with
standard deviation sdev.
*/
func AnglesNormDeg(n int, mean, sdev float64) []float64 {
	angles := make([]float64, n)
	for i := 0; i < n; i++ {
		angles[i] = AngleNormDeg(mean, sdev)
	}
	return angles
}
