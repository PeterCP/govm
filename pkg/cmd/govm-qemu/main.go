package main

import (
	"fmt"
	"os"

	cli "gopkg.in/urfave/cli.v2"
)

var cmd = cli.App{
	Name:  "govm-qemu",
	Usage: "QEMU wrapper to launch GoVM instances",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "name",
			Usage: "VM name",
		},
		&cli.StringFlag{
			Name:  "namespace",
			Usage: "VM namespace",
		},
		&cli.IntFlag{
			Name:  "vcpus",
			Usage: "Number of VCPUs for the VM",
		},
		&cli.IntFlag{
			Name:  "mem",
			Usage: "Amount of RAM for the VM (in MB)",
		},
		&cli.StringFlag{
			Name:  "image",
			Usage: "Image used to boot the VM",
		},
		&cli.BoolFlag{
			Name:  "efi",
			Usage: "Boot the VM using EFI",
		},
		&cli.StringSliceFlag{
			Name:    "share",
			Aliases: []string{"s", "v"},
			Usage:   "Share a directory with the VM (format: `LOCAL`:`VM`)",
		},
	},
	Action: runVM,
}

func main() {
	if err := cmd.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
