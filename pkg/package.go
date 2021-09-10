package pkg

type Package struct {
	Header  PackageHeader
	Meta    PackageMeta
	Content PackageContent
}

type PackageHeader struct {
	MpkgVersion string
}

type PackageMeta struct {
	Name         string
	Version      string
	Dependencies []string
	//Maintainer   string
	Scripts []string
}

type PackageContent struct {
	Files []byte
}
