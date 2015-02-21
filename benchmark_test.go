package flatjson

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"
)

func BenchmarkFlatJSON(b *testing.B) {
	lines := loadObjects(b, "dump.json")

	b.ResetTimer()
	for i, line := range lines {
		b.SetBytes(int64(len(line)))
		_, _, err := scanObject(line, func(num Number) {
			// dont care
		}, func(str String) {
			// dont care
		}, func(bl Bool) {
			// dont care
		}, func(null Null) {
			// dont care
		})
		if err != nil {
			b.Errorf("line %d: %v", i, err)
		}
	}
}

func BenchmarkEncodingJSON(b *testing.B) {
	lines := loadObjects(b, "dump.json")
	q := struct{}{}
	b.ResetTimer()
	for i, line := range lines {
		b.SetBytes(int64(len(line)))

		err := json.Unmarshal(line, &q)
		if err != nil {
			b.Errorf("line %d: %v", i, err)
		}
	}
}

func loadObjects(b *testing.B, filename string) [][]byte {
	var objects [][]byte

	f, err := os.Open(filename)
	if err != nil {
		b.Error(err)
	}
	defer f.Close()

	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanLines)

	for scan.Scan() {
		if len(objects) == b.N {
			return objects
		}
		// Text() to bytes, to force a copy of the memory,
		// otherwise Bytes() will recycle the bytes
		objects = append(objects, []byte(scan.Text()))
	}

	if scan.Err() != nil {
		b.Error(scan.Err())
	}

	for i := 0; len(objects) < b.N; i++ {
		objects = append(objects, []byte(string(objects[i])))
	}

	if len(objects) < b.N {
		b.Errorf("only %d lines in %q", len(objects), filename)
	}

	return objects
}
