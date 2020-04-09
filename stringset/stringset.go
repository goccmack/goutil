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
Package stringset: Operations on a set of strings
*/
package stringset

import "sort"

/*
StringSet implements a set of strings
*/
type StringSet struct {
	set map[string]bool
}

// New returns a new StringSet containing elements
func New(elements ...string) *StringSet {
	set := &StringSet{make(map[string]bool)}
	set.Add(elements...)
	return set
}

/*
Add elements to ss and return ss to allow chained commands
*/
func (ss *StringSet) Add(elements ...string) *StringSet {
	for _, e := range elements {
		ss.set[e] = true
	}
	return ss
}

/*
AddSet adds the elements of ss1 to ss and returns ss to allow chained commands
*/
func (ss *StringSet) AddSet(ss1 *StringSet) *StringSet {
	ss.Add(ss1.Elements()...)
	return ss
}

/*
Clone returns a deep copy of ss
*/
func (ss *StringSet) Clone() *StringSet {
	return New().Add(ss.Elements()...)
}

/*
Contain returns true iff ss contains s
*/
func (ss *StringSet) Contain(s string) bool {
	_, exist := ss.set[s]
	return exist
}

/*
Elements returns a slice containing the elements of ss
*/
func (ss *StringSet) Elements() []string {
	sl := make([]string, 0, len(ss.set))
	for s := range ss.set {
		sl = append(sl, s)
	}
	return sl
}

/*
ElementsSorted returns a slice containing the elements of ss sorted lexicographically
*/
func (ss *StringSet) ElementsSorted() []string {
	elements := ss.Elements()
	sort.Slice(elements, func(i, j int) bool { return elements[i] < elements[j] })
	return elements
}

/*
Equal returns true iff ss and ss1 have exactly the same elements
*/
func (ss *StringSet) Equal(ss1 *StringSet) bool {
	if ss.Len() != ss1.Len() {
		return false
	}
	for s := range ss.set {
		if !ss1.Contain(s) {
			return false
		}
	}
	return true
}

/*
Len returns the number of elements in ss
*/
func (ss *StringSet) Len() int {
	return len(ss.set)
}

/*
Remove element from ss and return ss to allow chained commands
*/
func (ss *StringSet) Remove(element string) *StringSet {
	delete(ss.set, element)
	return ss
}
