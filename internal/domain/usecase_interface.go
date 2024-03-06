package domain

import (
	"github.com/nkien0204/lets-go/internal/domain/entity/config"
	"github.com/nkien0204/lets-go/internal/domain/entity/generator"
)

type GeneratorUsecase interface {
	Generate(generator.OnlineGeneratorInputEntity) error
}

type ConfigUsecase interface {
	LoadConfig() *config.Cfg
}
