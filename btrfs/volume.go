package btrfs

import (
	"fmt"
	"path/filepath"
)

type Volume struct {
	Data     Segment
	Metadata Segment
	System   Segment

	uid     string
	baseDir string
}

func NewVolumeWithSysfsDir(uid string, sysfsDir string) Volume {
	baseDir := filepath.Join(sysfsDir, "fs/btrfs")

	return Volume{
		uid:      uid,
		baseDir:  baseDir,
		Data:     Segment{baseDir: filepath.Join(baseDir, uid, "allocation/data")},
		Metadata: Segment{baseDir: filepath.Join(baseDir, uid, "allocation/metadata")},
		System:   Segment{baseDir: filepath.Join(baseDir, uid, "allocation/system")},
	}
}

func NewVolume(uid string) Volume {
	return NewVolumeWithSysfsDir(uid, DEFAULT_SYSFS_DIR)
}

func GetVolumesFromDir(sysfsDir string) ([]Volume, error) {
	baseDir := filepath.Join(sysfsDir, "fs/btrfs")

	dirs, err := readDir(baseDir)
	if err != nil {
		return nil, err
	}

	result := []Volume{}
	for _, uid := range dirs {
		if uid != "features" {
			result = append(result, NewVolumeWithSysfsDir(uid, sysfsDir))
		}
	}

	return result, nil
}

func GetVolumes() ([]Volume, error) {
	return GetVolumesFromDir(DEFAULT_SYSFS_DIR)
}

func (v Volume) UID() string {
	return v.uid
}

func (v Volume) Label() (string, error) {
	return readString(filepath.Join(v.baseDir, v.uid, "label"))
}
func (v Volume) String() string {
	label, _ := v.Label()
	return fmt.Sprintf("%s (%s)", label, v.UID())
}
