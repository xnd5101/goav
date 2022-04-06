package main

import (
	"log"

	"github.com/xnd5101/goav/avcodec"
	"github.com/xnd5101/goav/avdevice"
	"github.com/xnd5101/goav/avfilter"
	"github.com/xnd5101/goav/avformat"
	"github.com/xnd5101/goav/avutil"
	"github.com/xnd5101/goav/swresample"
	"github.com/xnd5101/goav/swscale"
)

func main() {

	// Register all formats and codecs
	avformat.AvRegisterAll()
	avcodec.AvcodecRegisterAll()

	log.Printf("AvFilter Version:\t%v", avfilter.AvfilterVersion())
	log.Printf("AvDevice Version:\t%v", avdevice.AvdeviceVersion())
	log.Printf("SWScale Version:\t%v", swscale.SwscaleVersion())
	log.Printf("AvUtil Version:\t%v", avutil.AvutilVersion())
	log.Printf("AvCodec Version:\t%v", avcodec.AvcodecVersion())
	log.Printf("Resample Version:\t%v", swresample.SwresampleLicense())

}
