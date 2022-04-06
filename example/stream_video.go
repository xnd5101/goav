package main

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/giorgisio/goav/avcodec"
	"github.com/giorgisio/goav/avformat"
	"github.com/giorgisio/goav/avutil"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a src file")
		os.Exit(1)
	}

	// Register all formats and codecs
	avformat.AvRegisterAll()
	avcodec.AvcodecRegisterAll()

	// Open video file
	srcUrl := os.Args[1]
	outputUrl := "rtmp://127.0.0.1:1935/live/stream"
	pFormatContext := avformat.AvformatAllocContext()
	defer pFormatContext.AvformatCloseInput()

	if avformat.AvformatOpenInput(&pFormatContext, srcUrl, nil, nil) != 0 {
		fmt.Printf("Unable to open file %s\n", os.Args[1])
		os.Exit(1)
	}

	// Retrieve stream information
	if pFormatContext.AvformatFindStreamInfo(nil) < 0 {
		fmt.Println("Couldn't find stream information")
		os.Exit(1)
	}

	// Dump information about file onto standard error
	pFormatContext.AvDumpFormat(0, srcUrl, 0)

	// Find the first video stream
	videoindex := -1
	for i := 0; i < int(pFormatContext.NbStreams()); i++ {
		if pFormatContext.Streams()[i].CodecParameters().AvCodecGetType() == avformat.AVMEDIA_TYPE_VIDEO {
			videoindex = i
			break
		}
	}
	fmt.Println("videoindex=", videoindex)

	if videoindex == -1 {
		fmt.Println("Could not find video stream")
		os.Exit(1)
	}

	pOctxOut := avformat.AvformatAllocContext()
	defer pOctxOut.AvformatCloseInput()
	avformat.AvformatAllocOutputContext2(&pOctxOut, nil, "flv", outputUrl)
	for i := 0; i < int(pFormatContext.NbStreams()); i++ {
		pCodecCtxOrig := pFormatContext.Streams()[i].Codec()
		defer (*avcodec.Context)(unsafe.Pointer(pCodecCtxOrig)).AvcodecClose()
		out_stream := pOctxOut.AvformatNewStream((*avformat.AvCodec)(pCodecCtxOrig.GetCodec()))
		// out_stream := pOctxOut.AvformatNewStream(nil)
		if out_stream == nil {
			fmt.Println("could not alloc new stream")
			os.Exit(1)
		}

		ret := avcodec.AvcodecParametersCopy(out_stream.CodecParameters(), pFormatContext.Streams()[i].CodecParameters())
		if ret != 0 {
			fmt.Println("parameters copy fail")
			os.Exit(1)
		}

		out_stream.Codec().SetCodecTag(0)
		out_stream.CodecParameters().SetParamsCodecTag(0)
	}

	// pCodecCtxOrig := pFormatContext.Streams()[videoindex].Codec()
	// defer (*avcodec.Context)(unsafe.Pointer(pCodecCtxOrig)).AvcodecClose()
	// out_stream := pOctxOut.AvformatNewStream((*avformat.AvCodec)(pCodecCtxOrig.GetCodec()))
	// if out_stream == nil {
	// 	fmt.Println("could not alloc new stream")
	// 	os.Exit(1)
	// }

	// ret := avcodec.AvcodecParametersCopy(out_stream.CodecParameters(), pFormatContext.Streams()[videoindex].CodecParameters())
	// if ret != 0 {
	// 	fmt.Println("parameters copy fail")
	// 	os.Exit(1)
	// }
	// out_stream.Codec().SetCodecTag(0)
	// out_stream.CodecParameters().SetParamsCodecTag(0)

	// err := avformat.AvIOOpen1(outputUrl, avformat.AVIO_FLAG_WRITE, pOctxOut.Pb())
	// if err != nil {
	// 	fmt.Println("could not open io, code=", err)
	// }

	pIOCtx, err := avformat.AvIOOpen(outputUrl, avformat.AVIO_FLAG_WRITE)
	if err != nil {
		fmt.Println("could not open io, code=", err)
	}
	pOctxOut.SetPb(pIOCtx)
	// time.Sleep(time.Second * 5)

	// d := &avutil.Dictionary{}
	// d.AvDictSet("key", "value", 0)
	var opt *avutil.Dictionary = nil
	opt.AvDictSet("key", "value", 0)

	ret1 := pOctxOut.AvformatWriteHeader(nil)
	fmt.Println("ret=", ret1)
	if ret1 < 0 {
		fmt.Println("write head fail")
		os.Exit(1)
	}

	packet := avcodec.AvPacketAlloc()
	defer packet.AvFreePacket()
	for {

		if pFormatContext.AvReadFrame(packet) < 0 {
			fmt.Println("read frame fail")
			os.Exit(1)
		}

		// if packet.StreamIndex() != videoindex {
		// 	fmt.Println("not video frame")
		// 	continue
		// }

		if (packet.Flags() & avformat.AV_PKT_FLAG_KEY) == 1 {
			// fmt.Println("key frame")
		}

		if pOctxOut.AvWriteFrame(packet) < 0 {
			fmt.Println("write frame fail")
			break
		}
		// packet.AvFreePacket()
	}

}
