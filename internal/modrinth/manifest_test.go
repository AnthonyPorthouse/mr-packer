package modrinth

import (
	"reflect"
	"testing"

	"github.com/MakeNowJust/heredoc/v2"
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
				data: []byte(heredoc.Doc(`
				{
					"dependencies": {
						"fabric-loader": "0.15.11",
						"minecraft": "1.21"
					},
					"files": [
						{
							"downloads": [
							    "https://cdn.modrinth.com/data/P7dR8mSH/versions/HXzEJYgV/fabric-api-0.100.1%2B1.21.jar"
							],
							"env": {
								"client": "required",
								"server": "required"
							},
							"fileSize": 2219441,
							"hashes": {
								"sha1": "63848b365c036a89e4aba51d14d6dc8301e32b27",
								"sha512": "9f24f8b4fd9ca14cede45e9b4b8232f007613edb9f108965db486fb2e3e52656c8630f16cd0fff219808822ad0a251172406bdd75ba1cb9fd68ae480e06b5b4b"
							},
							"path": "mods/fabric-api-0.100.1+1.21.jar"
						}
					],
					"formatVersion": 1,
					"game": "minecraft",
					"name": "Test Pack",
					"versionId": "1.0.0"
				}
				`)),
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
						FileSize: 2219441,
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
