package btrfs

import (
	"path/filepath"
)

func (v Volume) Devices() ([]Device, error) {
	dirs, err := readDir(filepath.Join(v.baseDir, v.uid, "devices"))
	if err != nil {
		return nil, err
	}

	result := []Device{}
	for _, name := range dirs {
		result = append(result, Device{
			baseDir: filepath.Join(v.baseDir, v.uid, "devices", name),
			name:    name,
		})
	}

	return result, nil
}

type Device struct {
	baseDir string
	name    string
}

func (d Device) Name() string {
	return d.name
}

func (d Device) Model() (string, error) {
	return readString(filepath.Join(d.baseDir, "device/model"))
}

func (d Device) SectorsNum() (int, error) {
	return readInt(filepath.Join(d.baseDir, "size"))
}

func (d Device) SectorSize() (int, error) {
	return readInt(filepath.Join(d.baseDir, "queue/hw_sector_size"))
}

func (d Device) Size() (int, error) {
	n, err := d.SectorsNum()
	if err != nil {
		return 0, err
	}
	s, err := d.SectorSize()
	if err != nil {
		return 0, err
	}
	return n * s, nil
}
