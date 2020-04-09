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

// Package stringslice contains functions on slices of strings
package stringslice

import (
	"regexp"
)

// Clone returns a clone of s1
func Clone(s1 []string) (s2 []string) {
	s2 = make([]string, len(s1))

	copy(s2, s1)
	return
}

/*
Equal returns true iff s1 contains exactly the same strings as s2. The order of strings
may be different in s1 and s2.
*/
func Equal(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for _, e := range s1 {
		if !Contains(s2, e) {
			return false
		}
	}
	return true
}

/*
Contains returns true iff s contains at least one instance of e
*/
func Contains(s []string, e string) bool {
	for _, se := range s {
		if se == e {
			return true
		}
	}
	return false
}

/*
Find returns a list of indices in ss of strings equal to s.
Find returns a nil slice if ss does not contain s.
*/
func Find(ss []string, s string) (indices []int) {
	for i, s1 := range ss {
		if s1 == s {
			indices = append(indices, i)
		}
	}
	return
}

/*
MatchRegex returns true iff at least one of the strins in ss matches re.
*/
func MatchRegex(ss []string, re *regexp.Regexp) bool {
	for _, s := range ss {
		if re.MatchString(s) {
			return true
		}
	}
	return false
}

/*
Diff returns a minus all elements of b
*/
func Diff(a, b []string) (diff []string) {
	for _, e := range a {
		if !Contains(b, e) {
			diff = append(diff, e)
		}
	}
	return
}

/*
RemoveDuplicates returns a slice containing one instance of every string in in.
The order of strings returned is random.
*/
func RemoveDuplicates(in []string) []string {
	out := []string{}
	smap := make(map[string]bool)
	for _, s := range in {
		if _, exist := smap[s]; !exist {
			smap[s] = true
			out = append(out, s)
		}
	}
	return out
}

/*
Reverse returns a slice of string in reverse order of ss
*/
func Reverse(ss []string) []string {
	rev := make([]string, len(ss))
	for i, j := 0, len(ss)-1; i < len(ss); i, j = i+1, j-1 {
		rev[j] = ss[i]
	}
	return rev
}
