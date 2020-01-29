package gen

import (
	"bytes"
	"fmt"
	"path"

	. "github.com/dave/jennifer/jen"
)

func qualObj(id string) string {
	return fmt.Sprintf("workspace/pkg/obj/%s", id)
}

func ObjectPackageIndex(objIDs []string) ([]byte, error) {
	f := NewFile("obj")
	f.ImportName("github.com/manifold/tractor/pkg/manifold/library", "library")
	for _, objId := range objIDs {
		f.ImportName(qualObj(objId), objId)
		f.ImportAlias(qualObj(objId), objId)
	}

	f.Func().Id("relPath").Params(Id("subpath").Id("string")).Id("string").Block(
		List(Id("_"), Id("filename"), Id("_"), Id("_")).Op(":=").Qual("runtime", "Caller").Call(Lit(1)),
		Return(Qual("path", "Join").Call(Qual("path", "Dir").Call(Id("filename")), Id("subpath"))),
	)
	f.Line()
	var registrations []Code
	for _, objId := range objIDs {
		registrations = append(registrations,
			Qual("github.com/manifold/tractor/pkg/manifold/library", "Register").Call(
				Op("&").Qual(qualObj(objId), "Main").Values(),
				Lit(objId),
				Id("relPath").Call(Lit(path.Join(objId, "component.go"))),
			),
		)
	}
	f.Func().Id("init").Params().Block(registrations...)
	buf := &bytes.Buffer{}
	err := f.Render(buf)
	return buf.Bytes(), err
}
