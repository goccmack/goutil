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
package stringset: Operations on a set of strings
*/
package stringset

type StringSet struct {
	set map[string]bool
}

func New() *StringSet {
	return &StringSet{make(map[string]bool)}
}

/*
Return ss to allow chained commands
*/
func (ss *StringSet) Add(s ...string) *StringSet {
	for _, s1 := range s {
		ss.set[s1] = true
	}
	return ss
}

func (ss *StringSet) Clone() *StringSet {
	return New().Add(ss.List()...)
}

func (ss *StringSet) Contain(s string) bool {
	_, exist := ss.set[s]
	return exist
}

func (this *StringSet) Equal(that *StringSet) bool {
	for s, _ := range this.set {
		if !that.Contain(s) {
			return false
		}
	}
	return true
}

func (ss *StringSet) Remove(s string) {
	delete(ss.set, s)
}

func (ss *StringSet) List() []string {
	sl := make([]string, 0, len(ss.set))
	for s, _ := range ss.set {
		sl = append(sl, s)
	}
	return sl
}
