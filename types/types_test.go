package types

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFileNode(t *testing.T) {
	file := fmt.Sprintf("%s%d.txt", os.TempDir(), time.Now().Unix())
	var fp FileNode = NewBaseFile(file)
	assert.Nil(t, fp.Create(os.O_CREATE|os.O_RDWR, 0644))
	defer fp.Close()
	fp.WriteString(time.Now().String())
	t.Log(fp.Abs())
	t.Log("name", fp.Name(), "extra", fp.ExtraName(), "dir", fp.BaseDir())
}

func TestBaseDir(t *testing.T) {
	//t.Run("by directory", func(t *testing.T) {
	//	dir := os.TempDir()
	//	var fp BaseNode = NewBaseDir(dir)
	//	assert.Equal(t, dir, fp.BaseDir())
	//})

	t.Run("by file", func(t *testing.T) {
		dir := os.TempDir() + "file"
		var fp BaseNode = NewBaseDir(dir)
		t.Log(fp.BaseDir())
		t.Log(dir)
	})
}
