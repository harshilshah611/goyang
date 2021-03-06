// Copyright 2019 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package yang

import (
	"errors"
	"fmt"
	"testing"
)

// TestNode provides a framework for processing tests that can check particular
// nodes being added to the grammar. It can be used to ensure that particular
// statement combinations are supported, especially where they are opaque to
// the YANG library.
func TestNode(t *testing.T) {
	tests := []struct {
		desc      string
		inFn      func(*Modules) (Node, error)
		inModules map[string]string
		wantNode  func(Node) error
	}{{
		desc: "import reference statement",
		inFn: func(ms *Modules) (Node, error) {

			m, err := ms.FindModuleByPrefix("t")
			if err != nil {
				return nil, fmt.Errorf("can't find module in %v", ms)
			}

			if len(m.Import) == 0 {
				return nil, fmt.Errorf("node %v is missing imports", m)
			}

			return m.Import[0], nil
		},
		inModules: map[string]string{
			"test": `
				module test {
					prefix "t";
					namespace "urn:t";

					import foo {
						prefix "f";
						reference "bar";
					}
				}
			`,
			"foo": `
				module foo {
					prefix "f";
					namespace "urn:f";
				}
			`,
		},
		wantNode: func(n Node) error {
			is, ok := n.(*Import)
			if !ok {
				return fmt.Errorf("got node: %v, want type: import", n)
			}

			switch {
			case is.Reference == nil:
				return errors.New("did not get expected reference, got: nil, want: *yang.Statement")
			case is.Reference.Statement().Argument != "bar":
				return fmt.Errorf("did not get expected reference, got: %v, want: 'bar'", is.Reference.Statement())
			}

			return nil
		},
	}, {
		desc: "import description statement",
		inFn: func(ms *Modules) (Node, error) {

			m, err := ms.FindModuleByPrefix("t")
			if err != nil {
				return nil, fmt.Errorf("can't find module in %v", ms)
			}

			if len(m.Import) == 0 {
				return nil, fmt.Errorf("node %v is missing imports", m)
			}

			return m.Import[0], nil
		},
		inModules: map[string]string{
			"test": `
				module test {
					prefix "t";
					namespace "urn:t";

					import foo {
						prefix "f";
						description "foo module";
					}
				}
			`,
			"foo": `
				module foo {
					prefix "f";
					namespace "urn:f";
				}
			`,
		},
		wantNode: func(n Node) error {
			is, ok := n.(*Import)
			if !ok {
				return fmt.Errorf("got node: %v, want type: import", n)
			}

			switch {
			case is.Description == nil:
				return errors.New("did not get expected reference, got: nil, want: *yang.Statement")
			case is.Description.Statement().Argument != "foo module":
				return fmt.Errorf("did not get expected reference, got: '%v', want: 'foo module'", is.Description.Statement().Argument)
			}

			return nil
		},
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			ms := NewModules()

			for n, m := range tt.inModules {
				if err := ms.Parse(m, n); err != nil {
					t.Errorf("error parsing module %s, got: %v, want: nil", n, err)
				}
			}

			node, err := tt.inFn(ms)
			if err != nil {
				t.Fatalf("cannot run in function, %v", err)
			}

			if err := tt.wantNode(node); err != nil {
				t.Fatalf("failed check function, %v", err)
			}
		})
	}
}
