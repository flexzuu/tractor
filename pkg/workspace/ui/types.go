package ui

type Field struct {
	Type   string      `msgpack:"type"`
	Name   string      `msgpack:"name"`
	Path   string      `msgpack:"path"`
	Value  interface{} `msgpack:"value"`
	Enum   []string    `msgpack:"enum"`
	Fields []Field     `msgpack:"fields"`
}

type Button struct {
	Name    string `msgpack:"name"`
	Path    string `msgpack:"path"`
	OnClick string `msgpack:"onclick"`
}

type Component struct {
	Name     string    `msgpack:"name"`
	Filepath string    `msgpack:"filepath"`
	Fields   []Field   `msgpack:"fields"`
	Buttons  []Button  `msgpack:"buttons"`
	Related  []string  `msgpack:"related"`
	CustomUI []Element `msgpack:"customUI"`
}

type Node struct {
	Name       string      `msgpack:"name"`
	Path       string      `msgpack:"path"`
	Dir        string      `msgpack:"dir"`
	ID         string      `msgpack:"id"`
	Index      int         `msgpack:"index"`
	Active     bool        `msgpack:"active"`
	Components []Component `msgpack:"components"`
}

type Project struct {
	Name string `msgpack:"name"`
	Path string `msgpack:"path"`
}

type Prefab struct {
	Name string `msgpack:"name"`
	ID   string `msgpack:"id"`
}

type ComponentType struct {
	Filepath string `msgpack:"filepath"`
	Name     string `msgpack:"name"`
}
