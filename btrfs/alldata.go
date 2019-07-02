package btrfs

type VolumeData struct {
	UID      string
	Label    string
	Data     SegmentData
	Metadata SegmentData
	System   SegmentData
	Devices  []DeviceData
}

func (v Volume) ReadAllData() (data VolumeData, err error) {
	data.UID = v.UID()

	if data.Label, err = v.Label(); err != nil {
		return VolumeData{}, err
	}
	if data.Data, err = v.Data.ReadAllData(); err != nil {
		return VolumeData{}, err
	}
	if data.Metadata, err = v.Metadata.ReadAllData(); err != nil {
		return VolumeData{}, err
	}
	if data.System, err = v.System.ReadAllData(); err != nil {
		return VolumeData{}, err
	}

	devs, err := v.Devices()
	if err != nil {
		return VolumeData{}, err
	}
	for _, dev := range devs {
		d, err := dev.ReadAllData()
		if err != nil {
			return VolumeData{}, err
		}
		data.Devices = append(data.Devices, d)
	}

	return
}

type SegmentData struct {
	TotalBytes int
	UsedBytes  int
}

func (s Segment) ReadAllData() (data SegmentData, err error) {
	if data.TotalBytes, err = s.TotalBytes(); err != nil {
		return SegmentData{}, err
	}
	if data.UsedBytes, err = s.UsedBytes(); err != nil {
		return SegmentData{}, err
	}

	return
}

type DeviceData struct {
	Name       string
	Model      string
	SectorsNum int
	SectorSize int
	Size       int
}

func (d Device) ReadAllData() (data DeviceData, err error) {
	data.Name = d.Name()

	if data.Model, err = d.Model(); err != nil {
		return DeviceData{}, err
	}
	if data.SectorsNum, err = d.SectorsNum(); err != nil {
		return DeviceData{}, err
	}
	if data.SectorSize, err = d.SectorSize(); err != nil {
		return DeviceData{}, err
	}
	if data.Size, err = d.Size(); err != nil {
		return DeviceData{}, err
	}

	return
}
