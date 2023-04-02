package main

import (
	"github.com/dan-and-dna/gin-daemon/example/pb"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
)

func main() {
	msg := pb.Student{
		Name: "dan",
		Male: false,
	}

	buf, err := proto.Marshal(&msg)
	if err != nil {
		panic(err)
	}

	println(len(buf))

	//buf = []byte("dcdddddddddddd")
	if len(buf) > 0 {
		// Parse the tag (field number and wire type).
		num, wtyp, tagLen := protowire.ConsumeTag(buf)

		println(num)
		println(wtyp)
		println(tagLen)

	}

	err = proto.UnmarshalOptions{
		DiscardUnknown: true,
		AllowPartial:   true,
	}.Unmarshal(buf[:3], &msg)
	//err = proto.Unmarshal(buf[:3], &msg)
	if err != nil {
		panic(err)
	}

	println("ok")
}
