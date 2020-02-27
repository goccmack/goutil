package files

import (
	"fmt"
	"testing"
)

const (
	fileSize, numFiles = 100, 3
)

func TestFiles1(t *testing.T) {
	rf := New("logs", "files_test", fileSize, numFiles)
	defer rf.Close()
	for i := 0; i < 25; i++ {
		str := fmt.Sprintf("%2d.......\n", i)
		if _, err := rf.Write([]byte(str)); err != nil {
			t.Fatal()
		}
	}
}
