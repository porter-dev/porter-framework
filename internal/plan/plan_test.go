package plan

import (
	"testing"

	v "github.com/porterdev/ego/internal/value"
)

// Dequeues a Queue until empty to determine if it contains the same operations
// as specified in an array of operations. Since the order of some operations within
// a queue may change (due to hashmap randomness in Go's native "map" object), this
// function extracts the array operations into an auxiliary hashmap, and looks up
// the operations that are dequeued.
//
// This function is only meant to be used within a testing sequence, so it also logs
// operations that are not found using the (*testing.T).Errorf method.
func isOpSequenceEqual(q *OpQueue, ops []Operation, t *testing.T) bool {
	opMap := make(map[string]Operation)

	for _, op := range ops {
		opStr := op.ToString()
		opMap[opStr] = op
	}

	for i := 0; !q.IsEmpty(); i++ {
		op := q.Dequeue()
		opStr := op.ToString()

		// look up opStr to see if it exists in opMap
		if opFound, found := opMap[opStr]; found {
			if !op.IsEqual(&opFound) {
				t.Errorf("Operations not equal: expected %v, got %v", opFound, op)
			}
		} else {
			t.Errorf("Unmatched key %s in queue", opStr)
		}
	}

	return false
}

type planTest struct {
	old        v.Value
	new        v.Value
	operations []Operation
}

var simpleLiteralPlanTests = []planTest{
	planTest{
		old:        nil,
		new:        nil,
		operations: []Operation{},
	},
	planTest{
		old: v.Boolean(true),
		new: v.Boolean(false),
		operations: []Operation{
			Operation{
				Op:   UPDATE,
				Path: "",
				Old:  v.Boolean(true),
				New:  v.Boolean(false),
			},
		},
	},
	planTest{
		old: v.Boolean(true),
		new: v.Boolean(true),
		operations: []Operation{
			Operation{
				Op:   READ,
				Path: "",
				Old:  v.Boolean(true),
				New:  v.Boolean(true),
			},
		},
	},
	planTest{
		old: v.Boolean(true),
		new: nil,
		operations: []Operation{
			Operation{
				Op:   DELETE,
				Path: "",
				Old:  v.Boolean(true),
				New:  nil,
			},
		},
	},
	planTest{
		old: nil,
		new: v.Boolean(true),
		operations: []Operation{
			Operation{
				Op:   CREATE,
				Path: "",
				Old:  nil,
				New:  v.Boolean(true),
			},
		},
	},
	planTest{
		old: v.String("hello"),
		new: v.String("there"),
		operations: []Operation{
			Operation{
				Op:   UPDATE,
				Path: "",
				Old:  v.String("hello"),
				New:  v.String("there"),
			},
		},
	},
	planTest{
		old: v.String("hello"),
		new: v.String("hello"),
		operations: []Operation{
			Operation{
				Op:   READ,
				Path: "",
				Old:  v.String("hello"),
				New:  v.String("hello"),
			},
		},
	},
	planTest{
		old: v.String("hello"),
		new: nil,
		operations: []Operation{
			Operation{
				Op:   DELETE,
				Path: "",
				Old:  v.String("hello"),
				New:  nil,
			},
		},
	},
	planTest{
		old: nil,
		new: v.String("there"),
		operations: []Operation{
			Operation{
				Op:   CREATE,
				Path: "",
				Old:  nil,
				New:  v.String("there"),
			},
		},
	},
	planTest{
		old: v.Float(1.1),
		new: v.Float(1.2),
		operations: []Operation{
			Operation{
				Op:   UPDATE,
				Path: "",
				Old:  v.Float(1.1),
				New:  v.Float(1.2),
			},
		},
	},
	planTest{
		old: v.Float(1.1),
		new: v.Float(1.1),
		operations: []Operation{
			Operation{
				Op:   READ,
				Path: "",
				Old:  v.Float(1.1),
				New:  v.Float(1.1),
			},
		},
	},
	planTest{
		old: v.Float(1.1),
		new: nil,
		operations: []Operation{
			Operation{
				Op:   DELETE,
				Path: "",
				Old:  v.Float(1.1),
				New:  nil,
			},
		},
	},
	planTest{
		old: nil,
		new: v.Float(1.2),
		operations: []Operation{
			Operation{
				Op:   CREATE,
				Path: "",
				Old:  nil,
				New:  v.Float(1.2),
			},
		},
	},
	planTest{
		old: v.Integer(1),
		new: v.Integer(2),
		operations: []Operation{
			Operation{
				Op:   UPDATE,
				Path: "",
				Old:  v.Integer(1),
				New:  v.Integer(2),
			},
		},
	},
	planTest{
		old: v.Integer(1),
		new: v.Integer(1),
		operations: []Operation{
			Operation{
				Op:   READ,
				Path: "",
				Old:  v.Integer(1),
				New:  v.Integer(1),
			},
		},
	},
	planTest{
		old: v.Integer(1),
		new: nil,
		operations: []Operation{
			Operation{
				Op:   DELETE,
				Path: "",
				Old:  v.Integer(1),
				New:  nil,
			},
		},
	},
	planTest{
		old: nil,
		new: v.Integer(1),
		operations: []Operation{
			Operation{
				Op:   CREATE,
				Path: "",
				Old:  nil,
				New:  v.Integer(1),
			},
		},
	},
	planTest{
		old: v.Float(1.1),
		new: v.Integer(1),
		operations: []Operation{
			Operation{
				Op:   UPDATE,
				Path: "",
				Old:  v.Float(1.1),
				New:  v.Integer(1),
			},
		},
	},
}

func TestSimpleLiteralCreatePlan(t *testing.T) {
	for _, c := range simpleLiteralPlanTests {
		q := CreateOpQueue(c.old, c.new)

		isOpSequenceEqual(q, c.operations, t)
	}
}

var arrayPlanTests = []planTest{
	planTest{
		old: v.Array{},
		new: v.Array{
			v.Integer(0),
		},
		operations: []Operation{
			Operation{
				Op:   CREATE,
				Path: "[0]",
				Old:  nil,
				New:  v.Integer(0),
			},
		},
	},
	planTest{
		old: v.Array{
			v.Integer(0),
		},
		new: v.Array{},
		operations: []Operation{
			Operation{
				Op:   DELETE,
				Path: "[0]",
				Old:  v.Integer(0),
				New:  nil,
			},
		},
	},
	planTest{
		old: v.Array{
			v.Integer(0),
		},
		new: v.Array{},
		operations: []Operation{
			Operation{
				Op:   DELETE,
				Path: "[0]",
				Old:  v.Integer(0),
				New:  nil,
			},
		},
	},
	planTest{
		old: v.Array{
			v.Integer(0),
		},
		new: v.Array{
			v.Integer(1),
		},
		operations: []Operation{
			Operation{
				Op:   UPDATE,
				Path: "[0]",
				Old:  v.Integer(0),
				New:  v.Integer(1),
			},
		},
	},
	planTest{
		old: v.Array{
			v.Integer(1),
			v.Integer(2),
		},
		new: v.Array{
			v.Integer(1),
			v.Integer(2),
			v.Integer(3),
		},
		operations: []Operation{
			Operation{
				Op:   READ,
				Path: "[0]",
				Old:  v.Integer(1),
				New:  v.Integer(1),
			},
			Operation{
				Op:   READ,
				Path: "[1]",
				Old:  v.Integer(2),
				New:  v.Integer(2),
			},
			Operation{
				Op:   CREATE,
				Path: "[2]",
				Old:  nil,
				New:  v.Integer(3),
			},
		},
	},
	planTest{
		old: v.Array{
			v.Integer(0),
		},
		new: v.String("hello"),
		operations: []Operation{
			Operation{
				Op:   UPDATE,
				Path: "",
				Old: v.Array{
					v.Integer(0),
				},
				New: v.String("hello"),
			},
		},
	},
}

func TestArrayCreatePlan(t *testing.T) {
	for _, c := range arrayPlanTests {
		q := CreateOpQueue(c.old, c.new)

		isOpSequenceEqual(q, c.operations, t)
	}
}

var objectPlanTests = []planTest{
	planTest{
		old: v.Object{
			v.String("hello"): v.String("there1"),
		},
		new: v.Object{
			v.String("hello"): v.String("there2"),
		},
		operations: []Operation{
			Operation{
				Op:   UPDATE,
				Path: "[hello]",
				Old:  v.String("there1"),
				New:  v.String("there2"),
			},
		},
	},
	planTest{
		old: v.Object{
			v.String("hello"): v.String("there"),
			v.String("beep"):  v.String("boop"),
		},
		new: v.Object{
			v.String("hello"): v.String("there"),
		},
		operations: []Operation{
			Operation{
				Op:   READ,
				Path: "[hello]",
				Old:  v.String("there"),
				New:  v.String("there"),
			},
			Operation{
				Op:   DELETE,
				Path: "[beep]",
				Old:  v.String("boop"),
				New:  nil,
			},
		},
	},
	planTest{
		old: v.Object{
			v.String("hello"): v.String("there1"),
			v.String("beep1"): v.String("boop1"),
		},
		new: v.Object{
			v.String("hello"): v.String("there2"),
			v.String("beep2"): v.String("boop2"),
		},
		operations: []Operation{
			Operation{
				Op:   UPDATE,
				Path: "[hello]",
				Old:  v.String("there1"),
				New:  v.String("there2"),
			},
			Operation{
				Op:   DELETE,
				Path: "[beep1]",
				Old:  v.String("boop1"),
				New:  nil,
			},
			Operation{
				Op:   CREATE,
				Path: "[beep2]",
				Old:  nil,
				New:  v.String("boop2"),
			},
		},
	},
	planTest{
		old: v.Object{
			v.String("hello"): v.String("there1"),
		},
		new: v.String("hello"),
		operations: []Operation{
			Operation{
				Op:   UPDATE,
				Path: "",
				Old: v.Object{
					v.String("hello"): v.String("there1"),
				},
				New: v.String("hello"),
			},
		},
	},
}

func TestObjectCreatePlan(t *testing.T) {
	for _, c := range objectPlanTests {
		q := CreateOpQueue(c.old, c.new)

		isOpSequenceEqual(q, c.operations, t)
	}
}

// test taken from k8s deployment yaml
var k8sDeployment = planTest{
	old: v.Object{
		v.String("apiVersion"): v.String("apps/v1"),
		v.String("kind"):       v.String("Deployment"),
		v.String("metadata"): v.Object{
			v.String("name"): v.String("nginx-deployment"),
			v.String("labels"): v.Object{
				v.String("app"): v.String("nginx"),
			},
		},
		v.String("spec"): v.Object{
			v.String("replicas"): v.Integer(3),
			v.String("selector"): v.Object{
				v.String("matchLabels"): v.Object{
					v.String("app"): v.String("nginx"),
				},
			},
			v.String("template"): v.Object{
				v.String("metadata"): v.Object{
					v.String("labels"): v.Object{
						v.String("app"): v.String("nginx"),
					},
				},
				v.String("spec"): v.Object{
					v.String("containers"): v.Array{
						v.Object{
							v.String("name"):  v.String("nginx"),
							v.String("image"): v.String("nginx:1.14.2"),
							v.String("ports"): v.Array{
								v.Object{
									v.String("containerPort"): v.Integer(80),
								},
							},
						},
					},
				},
			},
		},
	},
	new: v.Object{
		v.String("apiVersion"): v.String("apps/v1beta1"),
		v.String("kind"):       v.String("Deployment"),
		v.String("metadata"): v.Object{
			v.String("name"):      v.String("nginx2-deployment"),
			v.String("namespace"): v.String("nginx"),
			v.String("labels"): v.Object{
				v.String("app"): v.String("nginx2"),
			},
		},
		v.String("spec"): v.Object{
			v.String("replicas"): v.Integer(5),
			v.String("selector"): v.Object{
				v.String("matchLabels"): v.Object{
					v.String("app"): v.String("nginx2"),
				},
			},
			v.String("template"): v.Object{
				v.String("metadata"): v.Object{
					v.String("labels"): v.Object{
						v.String("app"): v.String("nginx2"),
					},
				},
				v.String("spec"): v.Object{
					v.String("containers"): v.Array{
						v.Object{
							v.String("name"):  v.String("nginx2"),
							v.String("image"): v.String("nginx2:1.14.2"),
							v.String("ports"): v.Array{
								v.Object{
									v.String("containerPort"): v.Integer(3000),
								},
								v.Object{
									v.String("containerPort"): v.Integer(8000),
								},
							},
						},
					},
				},
			},
		},
	},
	operations: []Operation{
		Operation{
			Op:   UPDATE,
			Path: "[apiVersion]",
			Old:  v.String("apps/v1"),
			New:  v.String("apps/v1beta1"),
		},
		Operation{
			Op:   READ,
			Path: "[kind]",
			Old:  v.String("Deployment"),
			New:  v.String("Deployment"),
		},
		Operation{
			Op:   UPDATE,
			Path: "[metadata][name]",
			Old:  v.String("nginx-deployment"),
			New:  v.String("nginx2-deployment"),
		},
		Operation{
			Op:   UPDATE,
			Path: "[metadata][labels][app]",
			Old:  v.String("nginx"),
			New:  v.String("nginx2"),
		},
		Operation{
			Op:   CREATE,
			Path: "[metadata][namespace]",
			Old:  nil,
			New:  v.String("nginx"),
		},
		Operation{
			Op:   UPDATE,
			Path: "[spec][replicas]",
			Old:  v.Integer(3),
			New:  v.Integer(5),
		},
		Operation{
			Op:   UPDATE,
			Path: "[spec][selector][matchLabels][app]",
			Old:  v.String("nginx"),
			New:  v.String("nginx2"),
		},
		Operation{
			Op:   UPDATE,
			Path: "[spec][template][metadata][labels][app]",
			Old:  v.String("nginx"),
			New:  v.String("nginx2"),
		},
		Operation{
			Op:   UPDATE,
			Path: "[spec][template][spec][containers][0][name]",
			Old:  v.String("nginx"),
			New:  v.String("nginx2"),
		},
		Operation{
			Op:   UPDATE,
			Path: "[spec][template][spec][containers][0][image]",
			Old:  v.String("nginx:1.14.2"),
			New:  v.String("nginx2:1.14.2"),
		},
		Operation{
			Op:   UPDATE,
			Path: "[spec][template][spec][containers][0][ports][0][containerPort]",
			Old:  v.Integer(80),
			New:  v.Integer(3000),
		},
		Operation{
			Op:   CREATE,
			Path: "[spec][template][spec][containers][0][ports][1]",
			Old:  nil,
			New: v.Object{
				v.String("containerPort"): v.Integer(8000),
			},
		},
	},
}

func TestKubernetesDeploymentPlan(t *testing.T) {
	q := CreateOpQueue(k8sDeployment.old, k8sDeployment.new)

	isOpSequenceEqual(q, k8sDeployment.operations, t)
}

type resourcePathTest struct {
	paths      []string
	old        v.Value
	new        v.Value
	operations []Operation
}

var resourcePathTests = []resourcePathTest{
	resourcePathTest{
		paths: []string{
			"",
		},
		old: v.Object{
			v.String("hello"): v.Object{
				v.String("there"): v.String("general"),
			},
		},
		new: v.Object{
			v.String("hello"): v.Object{
				v.String("there"): v.String("general"),
			},
		},
		operations: []Operation{
			Operation{
				Op:   READ,
				Path: "",
				Old: v.Object{
					v.String("hello"): v.Object{
						v.String("there"): v.String("general"),
					},
				},
				New: v.Object{
					v.String("hello"): v.Object{
						v.String("there"): v.String("general"),
					},
				},
			},
		},
	},
	resourcePathTest{
		paths: []string{
			"",
		},
		old: v.Object{
			v.String("hello"): v.Object{
				v.String("there"): v.String("general1"),
			},
		},
		new: v.Object{
			v.String("hello"): v.Object{
				v.String("there"): v.String("general2"),
			},
		},
		operations: []Operation{
			Operation{
				Op:   UPDATE,
				Path: "",
				Old: v.Object{
					v.String("hello"): v.Object{
						v.String("there"): v.String("general1"),
					},
				},
				New: v.Object{
					v.String("hello"): v.Object{
						v.String("there"): v.String("general2"),
					},
				},
			},
		},
	},
	resourcePathTest{
		paths: []string{
			"[hello]",
		},
		old: v.Object{
			v.String("hello"): v.Object{
				v.String("there"): v.String("general1"),
			},
		},
		new: v.Object{
			v.String("hello"): v.Object{
				v.String("there"): v.String("general2"),
			},
		},
		operations: []Operation{
			Operation{
				Op:   UPDATE,
				Path: "[hello]",
				Old: v.Object{
					v.String("there"): v.String("general1"),
				},
				New: v.Object{
					v.String("there"): v.String("general2"),
				},
			},
		},
	},
}

func TestResourcePaths(t *testing.T) {
	for _, c := range resourcePathTests {
		q := CreateOpQueue(c.old, c.new, c.paths...)

		isOpSequenceEqual(q, c.operations, t)
	}
}
