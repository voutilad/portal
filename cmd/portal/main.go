package main

import (
	"fmt"
	"github.com/voutilad/portal"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

const blocksize = (1 << 20) // 1 MiB

func doit(config configItem, wg *sync.WaitGroup) {
	defer wg.Done()

	path := config.pipePath
	portal := config.portalConfig

	err := syscall.Mkfifo(path, 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr, "syscall.Mkfifo('%s'): %s\n", path, err)
		return
	}
	defer syscall.Unlink(path)

again:
	fmt.Printf("Opening portal to %s @ %s\n", portal, path)
	pipe, err := os.OpenFile(path, os.O_WRONLY, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "os.OpenFile('%s'): %s\n", path, err)
		return
	}

	fmt.Printf("Portal %s now transmitting...\n", portal)
	cnt, _ := io.Copy(pipe, portal)
	pipe.Close()

	fmt.Printf("Wrote %d bytes from portal %s\n", cnt, portal)
	portal, err = portal.Clone()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed clone!\n")
		return
	}
	goto again
}

func checkForCredentials() {
	// First, check for an explicit token
	_, ok := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")
	if ok {
		fmt.Fprintf(os.Stdout, "Using auth token from environmet\n")
		return
	}

	if !ok {
		fmt.Fprintln(os.Stderr, "Could not find GOOGLE_APPLICATION_CREDENTIALS environment variable!")
		usage()
	}
}

// Simple container for a portal and a pipe file path
type configItem struct {
	portalConfig portal.Portal
	pipePath     string
}

func usage() {
	fmt.Println("usage: portal [--gcs=<bucketName>,<objectName>:<fifo path>] [--gsm=<projectId>,<secret>,<version>:<fifo path>]")
	os.Exit(1)
}

// Parse an arg like: --gcs=bucket,object:/path/to/fifo
func parseArgs(args []string) ([]configItem, error) {
	var config []configItem

	for _, arg := range args {
		val := strings.Split(arg, "=")
		if len(val) < 2 {
			return nil, fmt.Errorf("bad argument: %s", arg)
		}

		parts := strings.Split(val[1], ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("bad argument: %s", arg)
		}

		portalParts := strings.Split(parts[0], ",")
		path := parts[1]

		switch val[0] {
		case "--gcs":
			p, err := portal.NewGcsPortal(portalParts[0], portalParts[1])
			if err != nil {
				return config, err
			}
			config = append(config, configItem{p, path})

		case "--gsm":
			version := portal.LatestVersion
			if len(portalParts) > 2 {
				v, err := strconv.Atoi(portalParts[2])
				if err != nil {
					return config, err
				}
				version = v
			}
			p, err := portal.NewGsmPortal(portalParts[0], portalParts[1], version)
			if err != nil {
				return config, err
			}
			config = append(config, configItem{p, path})

		default:
			continue
		}
	}

	return config, nil
}

func main() {
	var wg sync.WaitGroup

	config, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		usage()
	}
	if len(config) < 1 {
		usage()
	}

	for _, cfg := range config {
		wg.Add(1)
		go doit(cfg, &wg)

		// TODO: this should be handled better!
		defer cfg.portalConfig.Close()
	}

	wg.Wait()
}
