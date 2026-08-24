package main

import "errors"

var ErrUnsupportedDistro = errors.New("unsupported distro")

type DistroInfo struct {
	PackageManager string
	BaseImage      string
}

var DistroMap = map[string]DistroInfo{
	"debian":   {PackageManager: "apt", BaseImage: "debian"},
	"ubuntu":   {PackageManager: "apt", BaseImage: "ubuntu"},
	"alpine":   {PackageManager: "apk", BaseImage: "alpine"},
}

func GetDistroInfo(distro string) (DistroInfo, error) {
	info, ok := DistroMap[distro]
	if !ok {
		return DistroInfo{}, ErrUnsupportedDistro
	}
	return info, nil
}

func IsSupportedDistro(distro string) bool {
	_, ok := DistroMap[distro]
	return ok
}