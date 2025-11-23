package generator

import (
	"github.com/azuyamat/pace/internal/config"
	"github.com/azuyamat/pace/internal/logger"
	"github.com/azuyamat/pace/internal/models"
)

type Generator interface {
	Generate() (config.Config, error)
}

type generatorImpl func() (config.Config, error)

type baseGenerator struct {
	projectType models.ProjectType
	impl        generatorImpl
}

func newGenerator(impl generatorImpl, projectType models.ProjectType) Generator {
	return &baseGenerator{
		projectType: projectType,
		impl:        impl,
	}
}

func (g *baseGenerator) Generate() (config.Config, error) {
	cfg := config.NewDefaultConfig()

	if err := g.beforeGenerate(cfg); err != nil {
		return config.Config{}, err
	}

	generatedCfg, err := g.impl()
	if err != nil {
		return config.Config{}, err
	}

	if err := g.afterGenerate(&generatedCfg); err != nil {
		return config.Config{}, err
	}

	return generatedCfg, nil
}

func (g *baseGenerator) beforeGenerate(cfg *config.Config) error {
	logger.Info("Generating template for %s project...", g.projectType)
	return nil
}

func (g *baseGenerator) afterGenerate(cfg *config.Config) error {
	logger.Info("Generated config.pace:")
	logger.Info("--------------------------------------------------")
	logger.Info("\n%s", cfg.String())
	logger.Info("--------------------------------------------------")
	logger.Info("Generated template for %s project", g.projectType)
	logger.Info("You may have to modify the generated config file according to your project structure.")
	return nil
}

func GetGeneratorByProjectType(projectType models.ProjectType) Generator {
	switch projectType {
	case models.ProjectTypeGo:
		return NewGoGenerator()
	case models.ProjectTypeNode:
		return NewNodeGenerator()
	case models.ProjectTypePython:
		return NewPythonGenerator()
	case models.ProjectTypeRust:
		return NewRustGenerator()
	}
	return nil
}
