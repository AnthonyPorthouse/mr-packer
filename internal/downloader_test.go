package internal

import "testing"
import "github.com/spf13/afero"

func TestDownloadFiles(t *testing.T) {
	type args struct {
		manifest    *Manifest
		environment Environment
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				manifest: &Manifest{
					FormatVersion: 1,
					Game:          "minecraft",
					Name:          "",
					VersionId:     "1.0.0",
					Dependencies: &ManifestDependencies{
						Minecraft:    "1.21",
						FabricLoader: "0.15.11",
					},
					Files: []*ManifestFile{
						{
							Path:     "mods/fabric-api-0.100.1+1.21.jar",
							FileSize: 2219441,
							Env: &ManifestFileEnv{
								Client: Required,
								Server: Required,
							},
							Hashes: &ManifestFileHashes{
								Sha1:   "63848b365c036a89e4aba51d14d6dc8301e32b27",
								Sha512: "9f24f8b4fd9ca14cede45e9b4b8232f007613edb9f108965db486fb2e3e52656c8630f16cd0fff219808822ad0a251172406bdd75ba1cb9fd68ae480e06b5b4b",
							},
							Downloads: []string{
								"https://cdn.modrinth.com/data/P7dR8mSH/versions/HXzEJYgV/fabric-api-0.100.1%2B1.21.jar",
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DownloadFiles(tt.args.manifest, tt.args.environment, afero.NewMemMapFs()); (err != nil) != tt.wantErr {
				t.Errorf("DownloadFiles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
