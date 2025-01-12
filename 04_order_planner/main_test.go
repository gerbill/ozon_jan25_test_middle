package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	t.Parallel()

	for i := 1; ; i++ {
		file, err := os.Open(fmt.Sprintf("tests/%d", i))
		if i < 13 {
			continue
		}
		if i > 14 {
			t.Fail()
			break
		}
		if err != nil {
			if i > 30 {
				break
			}
			continue
		}
		t.Run(file.Name(), func(t *testing.T) {
			in := bufio.NewReader(file)
			fileName := fmt.Sprintf("tests/%d.a", i)
			t.Log(">>>>>>", fileName)
			expected, err := os.ReadFile(fileName)
			require.Nil(t, err)

			var buffer bytes.Buffer
			out := bufio.NewWriter(&buffer)

			Run(in, out)

			out.Flush()

			result, err := io.ReadAll(bufio.NewReader(&buffer))
			require.Nil(t, err)

			require.Equal(t, string(expected), string(result))
		})
	}
}
