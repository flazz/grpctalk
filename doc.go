//go:generate mkdir -p point
//go:generate protoc point.proto --go_out=plugins=grpc:point

package main

import (
	"encoding/json"
	"fmt"

	"github.com/flazz/atxgophers/grpctalk/point"
	"github.com/gogo/protobuf/proto"
)

func main() {
	p := point.Point{X: 100, Y: 100}

	tryBytes(p)
	tryText(p)
	tryJSON(p)
}

func tryBytes(p point.Point) {
	b, err := proto.Marshal(&p)
	if err != nil {
		panic(err)
	}
	fmt.Printf("as bytes: %#v\n", b)

	var p2 point.Point
	if err = proto.Unmarshal(b, &p2); err != nil {
		panic(err)
	}
	fmt.Println("byte round trip?", p == p2)
}

func tryText(p point.Point) {
	fmt.Println("as text:", p.String())
	var p2 point.Point
	if err := proto.UnmarshalText(p.String(), &p2); err != nil {
		panic(err)
	}
	fmt.Println("text round trip?", p == p2)

}

func tryJSON(p point.Point) {
	b, err := json.Marshal(&p)
	if err != nil {
		panic(err)
	}
	fmt.Println("as JSON:", string(b))

	var p2 point.Point
	if err := json.Unmarshal(b, &p2); err != nil {
		panic(err)
	}
	fmt.Println("json round trip?", p == p2)
}
