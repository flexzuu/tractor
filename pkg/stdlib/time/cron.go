package time

import "github.com/manifold/tractor/pkg/workspace/view"

type CronManager struct {
	Hello string
}

func (c *CronManager) InspectorUI() view.Element {
	return view.El("div", view.Attrs{"class": "mx-4 flex"}, []view.Element{
		view.El("atom.Knob", nil, nil),
		view.El("atom.Slider", nil, nil),
	})
}
