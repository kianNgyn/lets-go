package domain

import (
	"github.com/nkien0204/lets-go/internal/domain/entity/config"
	"github.com/nkien0204/lets-go/internal/domain/entity/generator"
)

type GeneratorRepository interface {
	GetRepoLatestVersion() (generator.RepoLatestVersionGetEntity, error)
	DownloadLatestAsset(generator.LatestAssetDownloadRequestEntity) error
}

type ConfigRepository interface {
	ReadConfigFile() (config.ConfigFileReadResponseEntity, error)
}
