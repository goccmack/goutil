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

package angle

import (
	"math"
	"testing"
)

func Test1(t *testing.T) {
	for deg := 0.; deg < 360; deg += 1 {
		if CosineSimilarityDeg(deg, deg) != 1 {
			t.Fail()
		}
	}
}

func Test2(t *testing.T) {
	for deg := 0.; deg < 360; deg += 1 {
		deg1 := AddDeg(deg, 180)
		if CosineSimilarityDeg(deg, deg1) != -1 {
			t.Errorf("%f, %f: %f", deg, deg1, CosineSimilarityDeg(deg, deg1))
		}
	}
}

func Test3(t *testing.T) {
	for deg1 := 0.; deg1 < 359; deg1 += 1 {
		deg2 := AddDeg(deg1, 90)
		if CosineSimilarityDeg(deg1, deg2) != 0 {
			t.Errorf("%f, %f: %.20f", deg1, deg2, CosineSimilarityDeg(deg1, deg2))
		}
	}
	for deg1 := 0.; deg1 < 359; deg1 += 1 {
		deg2 := AddDeg(deg1, 90+180)
		if CosineSimilarityDeg(deg1, deg2) != 0 {
			t.Errorf("%f, %f: %.20f", deg1, deg2, CosineSimilarityDeg(deg1, deg2))
		}
	}
}

func Test4(t *testing.T) {
	step := math.Pi / 4
	for th := 0.; th < 2*math.Pi; th += step {
		for i := 0.; i <= 4; i++ {
			diff := i * step
			if i > 0 {
				diff -= .1
			}
			sim := CosineSimilarity(th, Add(th, diff))
			if !Equal(sim, math.Cos(diff)) {
				t.Errorf("Sim(%f,%f)=%f(%f)", th, th+diff, sim, math.Cos(diff))
			}
		}
	}
}

func Test5(t *testing.T) {
	for d := 0.; d < 360; d += 1 {
		d1 := InvertDeg(d)
		d2 := InvertDeg(d1)
		if !Equal(d, d2) {
			t.Fail()
		}
		d3 := InvertDeg(d2)
		if !Equal(d2, d3) {
			t.Fail()
		}
	}
}

func Test6(t *testing.T) {
	for d := 0.; d < 2*math.Pi; d += math.Pi / 180 {
		d1 := InvertDeg(d)
		d2 := InvertDeg(d1)
		if !Equal(d, d2) {
			t.Fail()
		}
		d3 := InvertDeg(d2)
		if !Equal(d2, d3) {
			t.Fail()
		}
	}
}
