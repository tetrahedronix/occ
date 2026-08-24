package main

import "fmt"

type Distro string

const (
	DistroDebian Distro = "debian"
	DistroUbuntu Distro = "ubuntu"
	DistroAlpine Distro = "alpine"
)

type PackageManager string

const (
	PMApt PackageManager = "apt"
	PMApk PackageManager = "apk"
)

// Un'unica mappa che serve sia da "lista dei distro supportati"
// Sia da "distro ⟶ package manager", invece di mantenerle separate
var distroPackageManager = map[Distro]PackageManager{
	DistroDebian: PMApt,
	DistroUbuntu: PMApt,
	DistroAlpine: PMApk,
}

func (d Distro) PackageManager() (PackageManager, error) {
	pm, ok := distroPackageManager[d]

	if !ok {
		return "", fmt.Errorf("unsupported distro: %q", d)
	}

	return pm, nil
}

type Runtime string

const (
	RuntimePodman Runtime = "podman"
	RuntimeDocker Runtime = "docker"
)

var validRuntimes = map[Runtime]bool{
	RuntimeDocker: true,
	RuntimePodman: true,
}

func (r Runtime) Validate() error {
	if !validRuntimes[r] {
		return fmt.Errorf("unsupported runtime: %q", r)
	}

	return nil
}
