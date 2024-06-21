package internal

import (
	"reflect"
	"testing"
)

func TestParseManifest(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *Manifest
		wantErr bool
	}{
		{
			name: "Happy path",
			args: args{
				data: []byte("{\n    \"dependencies\": {\n        \"fabric-loader\": \"0.15.11\",\n        \"minecraft\": \"1.21\"\n    },\n    \"files\": [\n        {\n            \"downloads\": [\n                \"https://cdn.modrinth.com/data/P7dR8mSH/versions/HXzEJYgV/fabric-api-0.100.1%2B1.21.jar\"\n            ],\n            \"env\": {\n                \"client\": \"required\",\n                \"server\": \"required\"\n            },\n            \"fileSize\": 2219441,\n            \"hashes\": {\n                \"sha1\": \"63848b365c036a89e4aba51d14d6dc8301e32b27\",\n                \"sha512\": \"9f24f8b4fd9ca14cede45e9b4b8232f007613edb9f108965db486fb2e3e52656c8630f16cd0fff219808822ad0a251172406bdd75ba1cb9fd68ae480e06b5b4b\"\n            },\n            \"path\": \"mods/fabric-api-0.100.1+1.21.jar\"\n        }\n    ],\n    \"formatVersion\": 1,\n    \"game\": \"minecraft\",\n    \"name\": \"Test Pack\",\n    \"versionId\": \"1.0.0\"\n}"),
			},
			want: &Manifest{
				FormatVersion: 1,
				Name:          "Test Pack",
				VersionId:     "1.0.0",
				Game:          "minecraft",
				Dependencies: &ManifestDependencies{
					Minecraft:    "1.21",
					FabricLoader: "0.15.11",
				},
				Files: []*ManifestFile{
					{
						Path:      "mods/fabric-api-0.100.1+1.21.jar",
						Downloads: []string{"https://cdn.modrinth.com/data/P7dR8mSH/versions/HXzEJYgV/fabric-api-0.100.1%2B1.21.jar"},
						Env: &ManifestFileEnv{
							Client: Required,
							Server: Required,
						},
						Hashes: &ManifestFileHashes{
							Sha1:   "63848b365c036a89e4aba51d14d6dc8301e32b27",
							Sha512: "9f24f8b4fd9ca14cede45e9b4b8232f007613edb9f108965db486fb2e3e52656c8630f16cd0fff219808822ad0a251172406bdd75ba1cb9fd68ae480e06b5b4b",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseManifest(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseManifest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseManifest() got = %v, want %v", got, tt.want)
			}
		})
	}
}
