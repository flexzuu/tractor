package state

import (
	"sync"

	"github.com/manifold/tractor/pkg/manifold"
	"github.com/manifold/tractor/pkg/manifold/library"
	"github.com/manifold/tractor/pkg/manifold/prefab"
	"github.com/manifold/tractor/pkg/workspace/ui"
	"github.com/manifold/tractor/pkg/workspace/ui/field"
)

type UIProvider interface {
	InspectorUI() []ui.Element
}

type State struct {
	Projects        []ui.Project       `msgpack:"projects"`
	CurrentProject  string             `msgpack:"currentProject"`
	Components      []ui.ComponentType `msgpack:"components"`
	Prefabs         []ui.Prefab        `msgpack:"prefabs"`
	Hierarchy       []string           `msgpack:"hierarchy"`
	Nodes           map[string]ui.Node `msgpack:"nodes"`
	NodePaths       map[string]string  `msgpack:"nodePaths"`
	SelectedNode    string             `msgpack:"selectedNode"`
	EditorsEndpoint string             `msgpack:"editorsEndpoint"`

	mu sync.Mutex
}

func New(root manifold.Object) *State {
	state := &State{
		Projects:       []ui.Project{},
		CurrentProject: "dev",
		Nodes:          make(map[string]ui.Node),
		NodePaths:      make(map[string]string),
	}
	for _, com := range library.Registered() {
		state.Components = append(state.Components, ui.ComponentType{
			Name:     com.Name,
			Filepath: com.Filepath,
		})
	}
	for _, pf := range prefab.Registered() {
		state.Prefabs = append(state.Prefabs, ui.Prefab{
			Name: pf.Name,
			ID:   pf.ID,
		})
	}
	state.Update(root)
	return state
}

func (s *State) Update(root manifold.Object) {
	// reset/clear nodes
	s.Hierarchy = []string{}
	s.Nodes = make(map[string]ui.Node)
	// walk every object in the tree
	manifold.Walk(root, func(n manifold.Object) {
		// all the nodes paths
		s.Hierarchy = append(s.Hierarchy, n.Path())
		// start a node struct based on node passed in
		node := ui.Node{
			Name:   n.Name(),
			Active: true,
			// Dir:        n.Dir,
			Path:       n.Path(),
			Index:      n.SiblingIndex(),
			ID:         n.ID(),
			Components: []ui.Component{},
		}
		for _, com := range n.Components() {
			// get all the fields of the component
			fields := field.FromComponent(com)

			// see if component provides custom ui
			var customUI []ui.Element
			uip, ok := com.Pointer().(UIProvider)
			if ok {
				customUI = uip.InspectorUI()
			}

			// look up the filepath for this component
			var filepath string
			if com.ID() != "" {
				rc := library.LookupID(com.ID())
				if rc == nil {
					panic("component ID not in library: " + com.ID())
				}
				filepath = rc.Filepath
			} else {
				rc := library.Lookup(com.Name())
				if rc == nil {
					panic("component name not in library: " + com.Name())
				}
				filepath = rc.Filepath
			}

			// look up related components to this component
			var related []string
			for _, rc := range library.Related(library.Lookup(com.Name())) {
				related = append(related, rc.Name)
			}

			// add component to frontend node's components
			node.Components = append(node.Components, ui.Component{
				Name:     com.Name(),
				Filepath: filepath,
				Fields:   fields,
				Related:  related,
				CustomUI: customUI,
			})
		}
		// add the node to state
		s.mu.Lock()
		s.Nodes[n.ID()] = node
		s.NodePaths[n.Path()] = n.ID()
		s.mu.Unlock()
	})
}
