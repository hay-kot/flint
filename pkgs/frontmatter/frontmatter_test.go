package frontmatter_test

import (
	"strings"
	"testing"
	"time"

	"github.com/hay-kot/flint/pkgs/frontmatter"
	"github.com/stretchr/testify/assert"
)

var YAMLFrontMatter = `---
title: "Hello World"
date: 2012-12-12
tags:
  - foo
  - bar
categories:
  - foo
  - bar
---`

func Test_Read_YAML(t *testing.T) {
	assert := assert.New(t)

	reader := strings.NewReader(YAMLFrontMatter)

	fm, err := frontmatter.Read(reader)
	assert.NoError(err)

	data := fm.Data()

	assert.Equal("Hello World", data["title"])
	assert.Equal([]interface{}{"foo", "bar"}, data["tags"])
	assert.Equal([]interface{}{"foo", "bar"}, data["categories"])

	dt, _ := time.Parse("2006-01-02", "2012-12-12")
	assert.Equal(dt, data["date"])
}

var YAMLOrderedFrontMatter = `---
tags:
    - foo
    - bar
categories:
    - foo
    - bar
title: Hello World
date: 2012-12-12
---`

// func Test_Write_YAML_Ordered(t *testing.T) {
// 	reader := strings.NewReader(YAMLFrontMatter)

// 	fm, err := frontmatter.Read(reader)
// 	assert.NoError(t, err)

// 	order := frontmatter.WithOrder("tags", "categories", "title", "date")

// 	writer := &strings.Builder{}

// 	_, err = frontmatter.Write(writer, fm, order)
// 	assert.NoError(t, err)

// 	println(writer.String())

// 	assert.Equal(t, YAMLOrderedFrontMatter, writer.String())

// }
// func Test_Write_YAML(t *testing.T) {
// 	reader := strings.NewReader(YAMLFrontMatter)

// 	fm, err := frontmatter.Read(reader)

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	writer := &strings.Builder{}

// 	_, err = frontmatter.Write(writer, fm)

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	assert.Equal(t, YAMLFrontMatter, writer.String())
// }
