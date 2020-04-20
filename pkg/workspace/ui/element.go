package ui

type Element struct {
	Name     string
	Attrs    Attrs
	Children []Element
}

func E(name string, attrs Attrs, children []Element) Element {
	return Element{name, attrs, children}
}

type Attrs map[string]string
