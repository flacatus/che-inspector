package version

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestVersion(t *testing.T) {
	cmd := AddVersionCommand()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.Execute()
	_, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
}
