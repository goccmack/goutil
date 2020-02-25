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
Package angle contains utility routines for computing angles.
*/
package angle

import (
	"math"
)

const (
	FP_IGNORE = 1.e-9
)

/*
Return (θ1+θ2) mod 2π
*/
func Add(θ1, θ2 float64) float64 {
	return math.Mod(θ1+θ2, 2*math.Pi)
}

/*
Return (θ1+θ2) mod 360
*/
func AddDeg(θ1, θ2 float64) float64 {
	return math.Mod(θ1+θ2, 360)
}

/*
Returns cosine similarity between unit direction vectors. θ1 and θ2 are in [0,2π) or (-π,π).
The result is in [0,1], where 0 means the vectors are orthogonal and 1 that θ1 = θ2 or θ1 = -θ2.
*/
func CosineSimilarity(θ1, θ2 float64) float64 {
	θ := Diff(θ1, θ2)
	i, f := math.Modf(math.Cos(θ))
	if math.Abs(f) < FP_IGNORE {
		return i
	}
	return i + f
}

/*
θ1 and θ2 are in degrees.
See CosineSimilarity for details.
*/
func CosineSimilarityDeg(θ1, θ2 float64) float64 {
	return CosineSimilarity(ToRad(θ1), ToRad(θ2))
}

/*
Return the smaller angle between θ1 and θ2 in radians
*/
func Diff(θ1, θ2 float64) float64 {
	θ := math.Abs(θ1 - θ2)
	if θ > math.Pi {
		θ = 2*math.Pi - θ
	}
	return θ
}

func Equal(θ1, θ2 float64) bool {
	d := Diff(θ1, θ2)
	return d < FP_IGNORE
}

/*
Invert angle θ by turning it by π radians.
*/
func Invert(θ float64) float64 {
	return math.Mod(θ+math.Pi, 2*math.Pi)
}

/*
Invert angle θ by turning it by 180 degrees.
*/
func InvertDeg(θ float64) float64 {
	return math.Mod(θ+180, 360)
}

/*
Convert degrees to radians.
*/
func ToRad(deg float64) float64 {
	return deg * math.Pi / 180
}

/*
Convert radians to degrees.
*/
func ToDeg(rad float64) float64 {
	return rad * 180 / math.Pi
}
