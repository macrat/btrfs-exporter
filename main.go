package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/macrat/btrfs-exporter/btrfs"
)

const (
	NAMESPACE = "btrfs"
)

func newDesc(name string, labels ...string) *prometheus.Desc {
	return prometheus.NewDesc(fmt.Sprintf("%s_%s", NAMESPACE, name), "", labels, prometheus.Labels{})
}

var (
	totalBytes       = newDesc("total_bytes", "uid", "label", "type")
	usedBytes        = newDesc("used_bytes", "uid", "label", "type")
	deviceSectorsNum = newDesc("device_sectors_num", "uid", "label", "name", "device", "model")
	deviceSectorSize = newDesc("device_sector_bytes", "uid", "label", "name", "device", "model")
	deviceSize       = newDesc("device_total_bytes", "uid", "label", "name", "device", "model")
)

type Collector struct {
	sysfsDir string
}

func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- totalBytes
	ch <- usedBytes
	ch <- deviceSectorsNum
	ch <- deviceSectorSize
	ch <- deviceSize
}

func (c Collector) Collect(ch chan<- prometheus.Metric) {
	vs, err := btrfs.GetVolumesFromDir(c.sysfsDir)
	if err != nil {
		panic(err.Error())
	}

	for _, v := range vs {
		data, err := v.ReadAllData()
		if err != nil {
			panic(err.Error())
		}

		ch <- prometheus.MustNewConstMetric(totalBytes, prometheus.GaugeValue, float64(data.Data.TotalBytes), data.UID, data.Label, "data")
		ch <- prometheus.MustNewConstMetric(usedBytes, prometheus.GaugeValue, float64(data.Data.UsedBytes), data.UID, data.Label, "data")
		ch <- prometheus.MustNewConstMetric(totalBytes, prometheus.GaugeValue, float64(data.Metadata.TotalBytes), data.UID, data.Label, "metadata")
		ch <- prometheus.MustNewConstMetric(usedBytes, prometheus.GaugeValue, float64(data.Metadata.UsedBytes), data.UID, data.Label, "metadata")
		ch <- prometheus.MustNewConstMetric(totalBytes, prometheus.GaugeValue, float64(data.System.TotalBytes), data.UID, data.Label, "system")
		ch <- prometheus.MustNewConstMetric(usedBytes, prometheus.GaugeValue, float64(data.System.UsedBytes), data.UID, data.Label, "system")

		for _, dev := range data.Devices {
			device := fmt.Sprintf("/dev/%s", dev.Name)
			ch <- prometheus.MustNewConstMetric(deviceSectorsNum, prometheus.GaugeValue, float64(dev.SectorsNum), data.UID, data.Label, dev.Name, device, dev.Model)
			ch <- prometheus.MustNewConstMetric(deviceSectorSize, prometheus.GaugeValue, float64(dev.SectorSize), data.UID, data.Label, dev.Name, device, dev.Model)
			ch <- prometheus.MustNewConstMetric(deviceSize, prometheus.GaugeValue, float64(dev.Size), data.UID, data.Label, dev.Name, device, dev.Model)
		}
	}
}

var (
	addr     = flag.String("listen", "localhost:9999", "The address to listen")
	sysfsDir = flag.String("sysfs", "/sys", "The directory of sysfs")
)

func main() {
	flag.Parse()

	prometheus.MustRegister(Collector{sysfsDir: *sysfsDir})

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<a href=\"/metrics\">metrics</a>")
	})

	log.Printf("listen on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
