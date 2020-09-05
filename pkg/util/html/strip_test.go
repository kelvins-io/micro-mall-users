package html

import (
	"testing"
)

func TestStripTags(t *testing.T) {
	h := `<html>

<head>
    <meta charset="utf-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <title>API设计心得</title>

    <meta name="HandheldFriendly" content="True" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
`
	if StripTags(h) != "" {
		t.Error("html标签过滤不全")
	}
}
