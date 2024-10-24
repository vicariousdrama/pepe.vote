package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	nostr "github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

var nsec string

func main() {
	flag.StringVar(&nsec, "nsec", "", "")
	flag.Parse()
	sk := ""
	if _, s, e := nip19.Decode(nsec); e == nil {
		sk = s.(string)
	} else {
		fmt.Printf("Error decoding nsec: %s\n", e.Error())
		return
	}
	pub, _ := nostr.GetPublicKey(sk)

	readRelays := []string{"wss://relay.damus.io", "wss://relay.primal.net"}
	imageUrls := []string{}
	ctx := context.Background()
	relay1, err := nostr.RelayConnect(ctx, readRelays[0])
	if err != nil {
		panic(err)
	}
	relay2, err := nostr.RelayConnect(ctx, readRelays[1])
	if err != nil {
		panic(err)
	}
	for imageNum, imageurl := range imageUrls {
		fmt.Println("making event...")
		ev := nostr.Event{
			PubKey:    pub,
			CreatedAt: nostr.Now(),
			Kind:      3939,
			Tags:      []nostr.Tag{[]string{"L", "vote.pepe"}, []string{"l", "voteoptions"}, []string{"t", "pepevote"}, []string{"r", imageurl}},
		}
		err = ev.Sign(sk)
		if err != nil {
			fmt.Printf("error signing: %s\n", err.Error())
		}
		if err != nil {
		}
		if err := relay1.Publish(ctx, ev); err != nil {
			fmt.Printf("relay 1 err: %s\n", err.Error())
			return
			// continue
		}
		time.Sleep(time.Second * 5)
		fmt.Printf("image num to realy 1: %d\n", imageNum)

		if err := relay2.Publish(ctx, ev); err != nil {
			fmt.Printf("relay 2 err: %s\n", err.Error())
			return
		}
		fmt.Printf("image num to realy 2: %d\n", imageNum)
		time.Sleep(time.Second * 5)

	}
}
