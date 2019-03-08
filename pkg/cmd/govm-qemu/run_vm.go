package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/govm-project/govm/pkg/optparser"
	uuid "github.com/nu7hatch/gouuid"
	cli "gopkg.in/urfave/cli.v2"
)

func runVM(c *cli.Context) error {
	exe := getExecutable()
	rootDir := c.String("root")
	if _, err := os.Stat(rootDir); err != nil {
		return err
	}
	shares, err := parseShares(c.StringSlice("share"))
	if err != nil {
		return err
	}
	cloud := c.Bool("cloud")
	efi := c.Bool("efi")
	vcpus := c.Int("vcpus")
	mem := c.Int("mem")

	cmd := exec.Command("qemu-img", "info", filepath.Join(rootDir, "base.img"))
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	cow := false
	if strings.Contains(string(out), "qcow2") {
		cow = true
	}

	var args []string

	// Convert base image to qcow2 format
	if !cow {
		cmd = exec.Command("qemu-img", "convert", "-O", "qcow2",
			filepath.Join(rootDir, "base.img"),
			filepath.Join(rootDir, "base.img"))
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	// Create root image if it does not exist
	_, err = os.Stat(filepath.Join(rootDir, "root.img"))
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		cmd = exec.Command("qemu-img", "create", "-f", "qcow2", "-b",
			filepath.Join(rootDir, "base.img"), filepath.Join(rootDir, "root.img"))
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	// Root image mount args
	args = append(args, "-drive", fmt.Sprintf(
		"if=virtio,file=%s,format=qcow2,id=data",
		filepath.Join(rootDir, "root.img")))

	// Shared directories mount args
	for i := range shares {
		args = append(args, "-virtfs", fmt.Sprintf(
			"local,id=%s,path=%s,security_model=passthrough,mount_tag=%s",
			shares[i].id, shares[i].host, shares[i].guest))
	}

	// KVM args
	args = append(args,
		"-nodefaults",
		"-device", "virtio-balloon-pci,id=balloon1",
		"-realtime", "mlock=off",
		"-msg", "timestamp=on",
		"-chardev", "pty,id=charserial0",
		"-device", "isa-serial,chardev=charserial0,id=serial0",
		"-serial", "stdio",
		"-object", "rng-random,filename=/dev/urandom,id=rng0",
		"-device", "virtio-rng-pci,rng=rng0",
	)

	// Size args
	args = append(args,
		"-m", strconv.Itoa(mem),
		"-smp", strconv.Itoa(vcpus),
	)

	// Enable KVM if supported by host
	out, err = ioutil.ReadFile("/proc/cpuinfo")
	if err != nil {
		return err
	}
	outStr := string(out)
	if strings.Contains(outStr, "vmx") && strings.Contains(outStr, "svm") {
		args = append(args, "-enable-kvm", "-machine", "accel=kvm,usb=off")
	} else {
		args = append(args, "-machine", "usb=off")
	}

	args = append(args, "-vga", "qxl", "-display", "none")

	mac, err := genMACAddr()
	if err != nil {
		return err
	}

	return nil
}

func getExecutable() string {
	switch runtime.GOARCH {
	case "amd64":
		return "qemu-system-x86_64"
	default:
		return "qemu-system-" + runtime.GOARCH
	}
}

type share struct {
	host, guest, id string
}

func parseShares(strs []string) ([]share, error) {
	rand.Seed(0)
	var shares []share
	for _, str := range strs {
		var sh share
		if strings.Contains(str, ":") {
			parts := strings.Split(str, ":")
			if len(parts) != 2 {
				return shares, fmt.Errorf("invalid share format: %s", str)
			}
			sh.host = parts[0]
			sh.guest = parts[1]
		} else {
			pos, kv := optparser.ParseOpts(str, true)
			if len(pos) != 0 {
				return shares, fmt.Errorf("invalid share format: %s", str)
			}
			if kv["host"] == "" || kv["guest"] == "" {
				return shares, fmt.Errorf("invalid share format: %s", str)
			}
			sh.host = kv["host"]
			sh.guest = kv["guest"]
			sh.id = kv["id"]
		}
		if sh.id == "" {
			id, err := uuid.NewV4()
			if err != nil {
				return shares, err
			}
			sh.id = id.String()
		}
		shares = append(shares, sh)
	}
	return shares, nil
}

func genMACAddr() (net.HardwareAddr, error) {
	mac := net.HardwareAddr(make([]byte, 6))
	mac[0] = 0xFE
	mac[1] = 0x05
	_, err := rand.Read(mac[2:])
	if err != nil {
		return nil, err
	}
	return mac, nil
}

func bash(cmd string) error {
	return exec.Command("bash", "-c", cmd).Run()
}
