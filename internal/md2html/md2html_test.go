package md2html

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func readFile(fp string) ([]byte, error) {
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func TestConvert(t *testing.T) {
	a := assert.New(t)

	md, err := readFile("./testdata/test.md")
	a.NoError(err)
	a.NotEmpty(md)

	res, err := Convert(md)
	a.NoError(err)
	a.NotNil(res)

	html, err := readFile("./testdata/test.html")
	a.NoError(err)
	a.NotEmpty(html)

	a.Equal(string(html), res.Content)
	a.Equal(map[string]interface{}{
		"title":  "Blog Title",
		"labels": []interface{}{"Go", "Test Post"},
	}, res.Meta)
}
