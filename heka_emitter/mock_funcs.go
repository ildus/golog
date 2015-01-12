package heka_emitter

import "time"

var idGenerateBytes = func() []byte {
	// d1c7c768-b1be-4c70-93a6-9b52910d4baa.
	return []byte{0xd1, 0xc7, 0xc7, 0x68, 0xb1, 0xbe, 0x4c, 0x70, 0x93,
		0xa6, 0x9b, 0x52, 0x91, 0x0d, 0x4b, 0xaa}
}
var varosGetPid = func() int32 { return 1234 }
var timeNow = func() time.Time {
	// 2009-11-10 23:00:00 UTC; matches the Go Playground.
	return time.Unix(1257894000, 0).UTC()
}
