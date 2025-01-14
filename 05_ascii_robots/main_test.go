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
	// testIteration(t, 14)
	for i := 1; i <= 18; i++ {
		testIteration(t, i)
	}
}

func testIteration(t *testing.T, i int) {
	file, err := os.Open(fmt.Sprintf("tests/%d", i))
	if err != nil {
		if i > 30 {
			return
		}
		return
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
		// t.Fail()
	})
}
