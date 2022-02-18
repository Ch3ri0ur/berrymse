package bmff

import (
	// Bytes to transform values to a bitstream
	"bytes"

	// IO Writer to write Bitstream
	"io"
)

//TODO Maybe upgrade to https://pkg.go.dev/github.com/edgeware/mp4ff to create bmff

// References:
// ISO/IEC 14496 Part 12
// ISO/IEC 14496 Part 15

// ISO/IEC 14496 Part 14 is not used.

//TODO maybe changing it to https://go.dev/src/encoding/binary/binary.go
//a := binary.LittleEndian.Uint16(sampleA)
//a := binary.BigEndian.Uint16(sampleA)

// Writes Int as Byte Array of n Length in BIG-ENDIAN style into given Bytebuffer(io.Writer)
// Lowest Byte of the given Int is at [n-1] Pos
// v = value to write, n = number of bytes to write
func writeInt(w io.Writer, v int, n int) {
	b := make([]byte, n)
	// Write lowest Byte of Value in Highest Byte that isn't set
	// Shifting value by one Byte and do it again (next Byte of value)
	for i := 0; i < n; i++ {
		b[n-i-1] = byte(v & 0xff)
		v >>= 8
	}
	w.Write(b)
}

// Writes string in given Bytebuffer(io.Writer)
func writeString(w io.Writer, s string) {
	w.Write([]byte(s))
}

// Writes Box into given Bytebuffer(io.Writer)
// 1. Length of Box (4Byte)
// 2. Name of Box (4byte)
// 3. Content via Callback
func writeTag(w io.Writer, tag string, cb func(w io.Writer)) {
	var b bytes.Buffer
	cb(&b)                    // callback
	writeInt(w, b.Len()+8, 4) // box size
	writeString(w, tag)       // box type
	w.Write(b.Bytes())        // box content
}

// File Type Box
// Defines Type of File, Version and Compatible ISO Files.
func WriteFTYP(w io.Writer) {
	writeTag(w, "ftyp", func(w io.Writer) {
		writeString(w, "isom")                 // major brand
		writeInt(w, 0x200, 4)                  // minor version
		writeString(w, "isomiso2iso5avc1mp41") // compatible brands
	})
}

// Movie Box
// Metadata Container for Presentation.
func WriteMOOV(w io.Writer, width, height uint16, sps, pps []byte) {
	writeTag(w, "moov", func(w io.Writer) {
		writeMVHD(w)
		writeTRAK(w, width, height, sps, pps)
		writeMVEX(w)
	})
}

// Movie Header Box
// Generic Info about the movie
func writeMVHD(w io.Writer) {
	writeTag(w, "mvhd", func(w io.Writer) {
		writeInt(w, 0, 4)          // version and flags
		writeInt(w, 0, 4)          // creation time
		writeInt(w, 0, 4)          // modification time
		writeInt(w, 1000, 4)       // timescale
		writeInt(w, 0, 4)          // duration (all 1s == unknown)
		writeInt(w, 0x00010000, 4) // rate (1.0 == normal)
		writeInt(w, 0x0100, 2)     // volume (1.0 == normal)
		writeInt(w, 0, 2)          // reserved
		writeInt(w, 0, 4)          // reserved
		writeInt(w, 0, 4)          // reserved
		writeInt(w, 0x00010000, 4) // matrix
		writeInt(w, 0x0, 4)        // matrix
		writeInt(w, 0x0, 4)        // matrix
		writeInt(w, 0x0, 4)        // matrix
		writeInt(w, 0x00010000, 4) // matrix
		writeInt(w, 0x0, 4)        // matrix
		writeInt(w, 0x0, 4)        // matrix
		writeInt(w, 0x0, 4)        // matrix
		writeInt(w, 0x40000000, 4) // matrix
		writeInt(w, 0, 4)          // pre-defined
		writeInt(w, 0, 4)          // pre-defined
		writeInt(w, 0, 4)          // pre-defined
		writeInt(w, 0, 4)          // pre-defined
		writeInt(w, 0, 4)          // pre-defined
		writeInt(w, 0, 4)          // pre-defined
		writeInt(w, -1, 4)         // next track id
	})
}

// Track Box
// Metadata to one stream
func writeTRAK(w io.Writer, width, height uint16, sps, pps []byte) {
	writeTag(w, "trak", func(w io.Writer) {
		writeTKHD(w, width, height)
		writeMDIA(w, width, height, sps, pps)
	})
}

// Track Header Box
func writeTKHD(w io.Writer, width, height uint16) {
	writeTag(w, "tkhd", func(w io.Writer) {
		writeInt(w, 7, 4)               // version and flags (track enabled)
		writeInt(w, 0, 4)               // creation time
		writeInt(w, 0, 4)               // modification time
		writeInt(w, 1, 4)               // track id
		writeInt(w, 0, 4)               // reserved
		writeInt(w, 0, 4)               // duration
		writeInt(w, 0, 4)               // reserved
		writeInt(w, 0, 4)               // reserved
		writeInt(w, 0, 2)               // layer
		writeInt(w, 0, 2)               // alternate group
		writeInt(w, 0, 2)               // volume (ignored for video tracks)
		writeInt(w, 0, 2)               // reserved
		writeInt(w, 0x00010000, 4)      // matrix
		writeInt(w, 0x0, 4)             // matrix
		writeInt(w, 0x0, 4)             // matrix
		writeInt(w, 0x0, 4)             // matrix
		writeInt(w, 0x00010000, 4)      // matrix
		writeInt(w, 0x0, 4)             // matrix
		writeInt(w, 0x0, 4)             // matrix
		writeInt(w, 0x0, 4)             // matrix
		writeInt(w, 0x40000000, 4)      // matrix
		writeInt(w, int(width)<<16, 4)  // width (fixed-point 16.16 format)
		writeInt(w, int(height)<<16, 4) // height (fixed-point 16.16 format)
	})
}

// Media Box
func writeMDIA(w io.Writer, width, height uint16, sps, pps []byte) {
	writeTag(w, "mdia", func(w io.Writer) {
		writeMDHD(w)
		writeHDLR(w)
		writeMINF(w, width, height, sps, pps)
	})
}

// Media Header Box
func writeMDHD(w io.Writer) {
	writeTag(w, "mdhd", func(w io.Writer) {
		writeInt(w, 0, 4)      // version and flags
		writeInt(w, 0, 4)      // creation time
		writeInt(w, 0, 4)      // modification time
		writeInt(w, 10000, 4)  // timescale
		writeInt(w, 0, 4)      // duration
		writeInt(w, 0x55c4, 2) // language ('und' == undefined)
		writeInt(w, 0, 2)      // pre-defined
	})
}

// Handler Box
func writeHDLR(w io.Writer) {
	writeTag(w, "hdlr", func(w io.Writer) {
		writeInt(w, 0, 4)                        // version and flags
		writeInt(w, 0, 4)                        // pre-defined
		writeString(w, "vide")                   // handler type
		writeInt(w, 0, 4)                        // reserved
		writeInt(w, 0, 4)                        // reserved
		writeInt(w, 0, 4)                        // reserved
		writeString(w, "MicroMSE Video Handler") // name
		writeInt(w, 0, 1)                        // null-terminator
	})
}

// Media Information Box
func writeMINF(w io.Writer, width, height uint16, sps, pps []byte) {
	writeTag(w, "minf", func(w io.Writer) {
		writeVMHD(w)
		writeDINF(w)
		writeSTBL(w, width, height, sps, pps)
	})
}

// Video Media Header Box
func writeVMHD(w io.Writer) {
	writeTag(w, "vmhd", func(w io.Writer) {
		writeInt(w, 1, 4) // version and flags
		writeInt(w, 0, 2) // graphics mode
		writeInt(w, 0, 2) // opcolor
		writeInt(w, 0, 2) // opcolor
		writeInt(w, 0, 2) // opcolor
	})
}

// Data Information Box
func writeDINF(w io.Writer) {
	writeTag(w, "dinf", func(w io.Writer) {
		writeDREF(w)
	})
}

// Data Reference Box
func writeDREF(w io.Writer) {
	writeTag(w, "dref", func(w io.Writer) {
		writeInt(w, 0, 4) // version and flags
		writeInt(w, 1, 4) // entry count
		writeURL(w)
	})
}

// URL Box
func writeURL(w io.Writer) {
	writeTag(w, "url ", func(w io.Writer) {
		writeInt(w, 1, 4) // version and flags
	})
}

// Sample Table Box
func writeSTBL(w io.Writer, width, height uint16, sps, pps []byte) {
	writeTag(w, "stbl", func(w io.Writer) {
		writeSTSD(w, width, height, sps, pps)
		writeSTSZ(w)
		writeSTSC(w)
		writeSTTS(w)
		writeSTCO(w)
	})
}

// Sample Description Box
func writeSTSD(w io.Writer, width, height uint16, sps, pps []byte) {
	writeTag(w, "stsd", func(w io.Writer) {
		writeInt(w, 0, 6) // reserved
		writeInt(w, 1, 2) // data reference index
		writeAVC1(w, width, height, sps, pps)
	})
}

// Advanced Codec (H264) Box
func writeAVC1(w io.Writer, width, height uint16, sps, pps []byte) {
	writeTag(w, "avc1", func(w io.Writer) {
		writeInt(w, 0, 6)           // reserved
		writeInt(w, 1, 2)           // data reference index
		writeInt(w, 0, 2)           // pre-defined
		writeInt(w, 0, 2)           // reserved
		writeInt(w, 0, 4)           // pre-defined
		writeInt(w, 0, 4)           // pre-defined
		writeInt(w, 0, 4)           // pre-defined
		writeInt(w, int(width), 2)  // width
		writeInt(w, int(height), 2) // height
		writeInt(w, 0x00480000, 4)  // horizontal resolution: 72 dpi
		writeInt(w, 0x00480000, 4)  // vertical resolution: 72 dpi
		writeInt(w, 0, 4)           // data size: 0
		writeInt(w, 1, 2)           // frame count: 1
		w.Write(make([]byte, 32))   // compressor name
		writeInt(w, 0x18, 2)        // depth
		writeInt(w, 0xffff, 2)      // pre-defined
		writeAVCC(w, sps, pps)
	})
}

// AVC Configuration Box
// MPEG-4 Part 15 extension
// See ISO/IEC 14496-15:2004 5.3.4.1.2
func writeAVCC(w io.Writer, sps, pps []byte) {
	writeTag(w, "avcC", func(w io.Writer) {
		writeInt(w, 1, 1)           // configuration version
		writeInt(w, int(sps[1]), 1) // H.264 profile (0x64 == high)
		writeInt(w, int(sps[2]), 1) // H.264 profile compatibility
		writeInt(w, int(sps[3]), 1) // H.264 level (0x28 == 4.0)
		writeInt(w, 0xff, 1)        // nal unit length - 1 (upper 6 bits == 1)
		writeInt(w, 0xe1, 1)        // number of sps (upper 3 bits == 1)
		writeInt(w, len(sps), 2)    //len of sps
		w.Write(sps)                //sps
		writeInt(w, 1, 1)           // number of pps
		writeInt(w, len(pps), 2)    //len pps
		w.Write(pps)                //pps
	})
}

// Sample Size Box
func writeSTSZ(w io.Writer) {
	writeTag(w, "stsz", func(w io.Writer) {
		writeInt(w, 0, 4) // version and flags
		writeInt(w, 0, 4) // sample size
		writeInt(w, 0, 4) // sample count
	})
}

// Sample to Chunk Box
func writeSTSC(w io.Writer) {
	writeTag(w, "stsc", func(w io.Writer) {
		writeInt(w, 0, 4) // version and flags
		writeInt(w, 0, 4) // entry count
	})
}

// Time to Sample Box
func writeSTTS(w io.Writer) {
	writeTag(w, "stts", func(w io.Writer) {
		writeInt(w, 0, 4) // version and flags
		writeInt(w, 0, 4) // entry count
	})
}

// Chunk Offset Box
func writeSTCO(w io.Writer) {
	writeTag(w, "stco", func(w io.Writer) {
		writeInt(w, 0, 4) // version and flags
		writeInt(w, 0, 4) // entry count
	})
}

// Movie Extends Box
func writeMVEX(w io.Writer) {
	writeTag(w, "mvex", func(w io.Writer) {
		writeMEHD(w)
		writeTREX(w)
	})
}

// Movie Extends Header Box
func writeMEHD(w io.Writer) {
	writeTag(w, "mehd", func(w io.Writer) {
		writeInt(w, 0, 4) // version and flags
		writeInt(w, 0, 4) // fragment duration
	})
}

// Track Extends Box
func writeTREX(w io.Writer) {
	writeTag(w, "trex", func(w io.Writer) {
		writeInt(w, 0, 4)          // version and flags
		writeInt(w, 1, 4)          // track id
		writeInt(w, 1, 4)          // default sample description index
		writeInt(w, 0, 4)          // default sample duration
		writeInt(w, 0, 4)          // default sample size
		writeInt(w, 0x00010000, 4) // default sample flags
	})
}

// Movie Fragment Box
// The MOOF Box holds the information to a Movie Fragment, which is in the following MDAT Object.
func WriteMOOF(w io.Writer, seq int, data []byte) {
	writeTag(w, "moof", func(w io.Writer) {
		writeMFHD(w, seq)
		writeTRAF(w, seq, data)
	})
}

// Movie Fragment Header Box
func writeMFHD(w io.Writer, seq int) {
	writeTag(w, "mfhd", func(w io.Writer) {
		writeInt(w, 0, 4)   // version and flags
		writeInt(w, seq, 4) // sequence number
	})
}

// Track Fragment Box
func writeTRAF(w io.Writer, seq int, data []byte) {
	writeTag(w, "traf", func(w io.Writer) {
		writeTFHD(w)
		writeTFDT(w, seq)
		writeTRUN(w, data)
	})
}

// Track Fragment Header Box
func writeTFHD(w io.Writer) {
	writeTag(w, "tfhd", func(w io.Writer) {
		writeInt(w, 0x020020, 4)   // version and flags
		writeInt(w, 1, 4)          // track ID
		writeInt(w, 0x01010000, 4) // default sample flags
	})
}

// Track Fragment Base Media Decode Time Box
func writeTFDT(w io.Writer, seq int) {
	writeTag(w, "tfdt", func(w io.Writer) {
		writeInt(w, 0x01000000, 4) // version and flags
		writeInt(w, 330*seq, 8)    // base media decode time
	})
}

// Track Fragment Run Box
func writeTRUN(w io.Writer, data []byte) {
	writeTag(w, "trun", func(w io.Writer) {
		writeInt(w, 0x00000305, 4) // version and flags
		writeInt(w, 1, 4)          // sample count
		writeInt(w, 0x70, 4)       // data offset
		if (len(data) > 0) && (data[0]&0x1f == 0x5) {
			writeInt(w, 0x02000000, 4) // first sample flags (i-frame)
		} else {
			writeInt(w, 0x01010000, 4) // first sample flags (not i-frame)
		}
		writeInt(w, 330, 4)         // sample duration
		writeInt(w, 4+len(data), 4) // sample size
	})
}

// Media Data Box
// The Media Box holds an Media Sample. In this Project this is a NAL Unit.
func WriteMDAT(w io.Writer, data []byte) {
	writeTag(w, "mdat", func(w io.Writer) {
		writeInt(w, len(data), 4)
		w.Write(data)
	})
}
