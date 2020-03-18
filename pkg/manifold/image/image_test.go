package image_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/manifold/tractor/pkg/manifold"
	"github.com/manifold/tractor/pkg/manifold/image"
	"github.com/manifold/tractor/pkg/manifold/object"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrite(t *testing.T) {
	fs := afero.NewMemMapFs()
	t.Run("writes and loads empty tree", func(t *testing.T) {
		img := image.NewWith(fs, "/image")
		root := object.New("root")
		require.Nil(t, img.Write(root))

		files := allfiles(t, fs)
		assert.Equal(t, 3, len(files))
		assertDir(t, files["/image"])
		assertDir(t, files["/image/obj"])
		assertFsObject(t, fs, "/image/obj/object.json", root)

		root2, err := img.Load()
		if assert.Nil(t, err) {
			assertObject(t, root, root2)
		}
	})

	// reset fs to build object tree with nested child objects.
	fs = afero.NewMemMapFs()
	t.Run("writes and loads full tree", func(t *testing.T) {
		img := image.NewWith(fs, "/image")
		root := object.New("root")
		obja := object.New("a")
		obja1 := object.New("1")
		obja11 := object.New("1")
		obja2 := object.New("2")
		objb := object.New("b")
		obja1.AppendChild(obja11)
		obja.AppendChild(obja1)
		obja.AppendChild(obja2)
		root.AppendChild(obja)
		root.AppendChild(objb)

		t.Logf("root children: %+v", root.Children())

		require.Nil(t, img.Write(root))

		files := allfiles(t, fs)
		assert.Equal(t, 13, len(files))
		assertDir(t, files["/image"])
		assertDir(t, files["/image/obj"])
		assertDir(t, files[objDir(obja)])
		assertDir(t, files[objDir(obja, obja1)])
		assertDir(t, files[objDir(obja, obja1, obja11)])
		assertDir(t, files[objDir(obja, obja2)])
		assertDir(t, files[objDir(objb)])
		assertFsObject(t, fs, "/image/obj/object.json", root)
		assertFsObject(t, fs, objPath(obja), obja)
		assertFsObject(t, fs, objPath(obja, obja1), obja1)
		assertFsObject(t, fs, objPath(obja, obja1, obja11), obja11)
		assertFsObject(t, fs, objPath(obja, obja2), obja2)
		assertFsObject(t, fs, objPath(objb), objb)

		root2, err := img.Load()
		if assert.Nil(t, err) {
			assertObject(t, root, root2)
		}
	})

	t.Run("deletes sub node", func(t *testing.T) {
		img := image.NewWith(fs, "/image")
		root, err := img.Load()
		require.Nil(t, err)

		obja := root.ChildAt(0)
		require.NotNil(t, obja)
		assert.Equal(t, "a", obja.Name())

		obja1 := obja.ChildAt(0)
		require.NotNil(t, obja)
		assert.Equal(t, "1", obja1.Name())

		objb := root.ChildAt(1)
		require.NotNil(t, objb)
		assert.Equal(t, "b", objb.Name())

		obja2 := obja.RemoveChildAt(1)
		require.NotNil(t, obja2)
		assert.Equal(t, "2", obja2.Name())

		require.Nil(t, img.Write(root))

		files := allfiles(t, fs)
		assert.Equal(t, 11, len(files))
		assertDir(t, files["/image"])
		assertDir(t, files["/image/obj"])
		assertDir(t, files[objDir(obja)])
		assertDir(t, files[objDir(obja, obja1)])
		assert.Nil(t, files[objDir(obja, obja2)])
		assertDir(t, files[objDir(objb)])
		assertFsObject(t, fs, "/image/obj/object.json", root)
		assertFsObject(t, fs, objPath(obja), obja)
		assertFsObject(t, fs, objPath(obja, obja1), obja1)
		assert.Nil(t, files[objPath(obja, obja2)])
		assertFsObject(t, fs, objPath(objb), objb)
	})
}

func allfiles(t *testing.T, fs afero.Fs) map[string]os.FileInfo {
	files := make(map[string]os.FileInfo)
	err := afero.Walk(fs, "/image", func(path string, info os.FileInfo, err error) error {
		files[path] = info
		return err
	})
	require.Nil(t, err)
	return files
}

var objNameRE = regexp.MustCompile("[^a-zA-Z0-9]+")

// builds slice of object path names from the given object hierarchy. No need
// to pass the root object, since it's in /image/obj already.
// ex: []string{"", "image", "obj", ...}
//
// path building ripped from image.pathNameFromImage()
func objParts(objects ...manifold.Object) []string {
	parts := make([]string, 3, len(objects)+3)
	parts[1] = "image"
	parts[2] = "obj"
	for _, o := range objects {
		id := o.ID()
		shortid := id[len(id)-8:]
		name := strings.ToLower(objNameRE.ReplaceAllString(o.Name(), ""))
		parts = append(parts, fmt.Sprintf("%s-%s", name[:min(8, len(name))], shortid))
	}
	return parts
}

// returns directory containing object.json.
func objDir(objects ...manifold.Object) string {
	parts := objParts(objects...)
	return strings.Join(parts, string(filepath.Separator))
}

// returns path to object.json.
func objPath(objects ...manifold.Object) string {
	parts := append(objParts(objects...), "object.json")
	return strings.Join(parts, string(filepath.Separator))
}

func assertDir(t *testing.T, fi os.FileInfo) {
	if assert.NotNil(t, fi) {
		assert.True(t, fi.IsDir())
	}
}

func assertObject(t *testing.T, expected, actual manifold.Object) {
	assert.Equal(t, expected.ID(), actual.ID())
	assert.Equal(t, expected.Name(), actual.Name())
	assert.Equal(t, expected.Path(), actual.Path())
	assert.Equal(t, len(expected.Children()), len(actual.Children()))
}

func assertSnapshot(t *testing.T, snap manifold.ObjectSnapshot, obj manifold.Object) {
	assert.Equal(t, snap.ID, obj.ID())
	assert.Equal(t, snap.Name, obj.Name())
	assert.Equal(t, len(snap.Children), len(obj.Children()))
}

func assertFsObject(t *testing.T, fs afero.Fs, path string, obj manifold.Object) manifold.ObjectSnapshot {
	snap := loadSnapshot(t, fs, path)
	assertSnapshot(t, snap, obj)
	return snap
}

func loadSnapshot(t *testing.T, fs afero.Fs, path string) manifold.ObjectSnapshot {
	by, err := afero.ReadFile(fs, path)
	require.Nil(t, err)
	var snap manifold.ObjectSnapshot
	err = json.Unmarshal(by, &snap)
	require.Nil(t, err)
	return snap
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
