package btrfs

import (
	"path/filepath"
)

type Segment struct {
	baseDir string
}

func (s Segment) TotalBytes() (int, error) {
	return readInt(filepath.Join(s.baseDir, "total_bytes"))
}

func (s Segment) UsedBytes() (int, error) {
	return readInt(filepath.Join(s.baseDir, "bytes_used"))
}
