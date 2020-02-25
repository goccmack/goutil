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

package stringset

import (
	"testing"

	"github.com/goccmack/goutil/stringslice"
)

var (
	s1 = []string{"a", "b", "c"}
)

func Test1(t *testing.T) {
	ss := New().Add(s1...)
	if !stringslice.Equal(s1, ss.List()) {
		t.Fail()
	}
}

func Test2(t *testing.T) {
	ss := New().Add(s1...)
	for _, s := range s1 {
		if !ss.Contain(s) {
			t.Errorf("Expected ss to contain %s", s)
		}
	}
}

func Test3(t *testing.T) {
	ss := New().Add(s1...)
	for i, s := range s1 {
		ss.Remove(s)
		if ss.Contain(s) {
			t.Fail()
		}
		for j := i + 1; j < len(s1); j++ {
			if !ss.Contain(s1[j]) {
				t.Fail()
			}
		}
	}
}
