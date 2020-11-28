package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/matsuyoshi30/localremo"
)

func main() {
	var get, post string
	var debug bool
	flag.StringVar(&get, "get", "", "What device you want to get signal")
	flag.StringVar(&post, "post", "", "Config file path to what you want to post signal")
	flag.BoolVar(&debug, "debug", false, "Debug mode")
	flag.Parse()

	entry, localRemoAddr, err := localremo.GetLocalRemoAddr()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't get local Remo address: %v\n", err)
		os.Exit(1)
	}
	if localRemoAddr == nil {
		fmt.Fprintln(os.Stderr, "Couldn't get local Remo address")
		os.Exit(1)
	}
	if debug {
		log.Println("ServiceRecord: ", entry.ServiceRecord)
		log.Println("Service HostNamet: ", entry.HostName)
		log.Println("Service Port: ", entry.Port)
		log.Println("Service Text: ", entry.Text)
		log.Println("Service TTL: ", entry.TTL)
		log.Println("Service AddrIPv4: ", entry.AddrIPv4)
		log.Println("Service AddrIPv6: ", entry.AddrIPv6)
	}

	client := localremo.NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if (get != "" && post != "") || (get == "" && post == "") {
		fmt.Fprintln(os.Stderr, "You have to select one option\n")
		os.Exit(1)
	}

	if get != "" {
		if out, err := client.Get(ctx, localRemoAddr); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to GET: %v\n", err)
			os.Exit(1)
		} else {
			b, err := json.Marshal(out)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to marshal: %v\n", err)
				os.Exit(1)
			}

			new, err := os.Create(get)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to create file: %v\n", err)
				os.Exit(1)
			}
			if _, err := new.Write(b); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write file: %v\n", err)
				os.Exit(1)
			}
		}
	} else if post != "" {
		ir, err := localremo.ReadJSON(post)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't read file: %v\n", err)
			os.Exit(1)
		}

		if err := client.Post(ctx, localRemoAddr, bytes.NewBuffer(ir)); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to GET: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("DONE!")
}
