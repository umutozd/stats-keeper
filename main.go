package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/umutozd/stats-keeper/protos/statspb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	msg := &statspb.ComponentDate{
		Timestamps: []*timestamppb.Timestamp{
			timestamppb.New(time.Now()),
		},
	}
	d, err := json.MarshalIndent(msg, "", " ")
	if err != nil {
		log.Fatalf("error marshaling message: %v", err)
	}

	fmt.Printf("marshaled message: %s", string(d))
}
