package modrinth

import (
	"encoding/json"
	"fmt"
)

type Manifest struct {
	FormatVersion uint                  `json:"formatVersion"`
	Game          string                `json:"game"`
	VersionId     string                `json:"versionId"`
	Name          string                `json:"name"`
	Summary       string                `json:"summary,omitempty"`
	Dependencies  *ManifestDependencies `json:"dependencies,omitempty"`
	Files         []*ManifestFile       `json:"files"`
}

type ManifestDependencies struct {
	Minecraft    string `json:"minecraft,omitempty"`
	Forge        string `json:"forge,omitempty"`
	NeoForge     string `json:"neoforge,omitempty"`
	FabricLoader string `json:"fabric-loader,omitempty"`
	QuiltLoader  string `json:"quilt-loader,omitempty"`
}

type ManifestFile struct {
	Path      string              `json:"path"`
	FileSize  int64               `json:"fileSize"`
	Downloads []string            `json:"downloads"`
	Hashes    *ManifestFileHashes `json:"hashes"`
	Env       *ManifestFileEnv    `json:"env,omitempty"`
}

type ManifestFileHashes struct {
	Sha1   string `json:"sha1"`
	Sha512 string `json:"sha512"`
}

type ManifestFileEnv struct {
	Client EnvironmentSupport `json:"client"`
	Server EnvironmentSupport `json:"server"`
}

type EnvironmentSupport string

const (
	Required    EnvironmentSupport = "required"
	Optional    EnvironmentSupport = "optional"
	Unsupported EnvironmentSupport = "unsupported"
)

type Environment string

const (
	Client Environment = "client"
	Server Environment = "server"
)

func ParseManifest(data []byte) (*Manifest, error) {
	var manifest *Manifest

	json.Unmarshal(data, &manifest)

	fmt.Println(manifest)

	return manifest, nil
}
