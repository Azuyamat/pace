package detector

import (
	"os"

	"github.com/azuyamat/pace/internal/models"
)

type detector interface {
	isDetected(projectPath string) bool
	ProjectType() models.ProjectType
}

type baseDetector struct {
	indicators  []indicatorFunc
	projectType models.ProjectType
}

func (d *baseDetector) isDetected(projectPath string) bool {
	for _, indicator := range d.indicators {
		if !indicator(projectPath) {
			return false
		}
	}
	return true
}

func (d *baseDetector) ProjectType() models.ProjectType {
	return d.projectType
}

func newBaseDetector(projectType models.ProjectType, indicators ...indicatorFunc) *baseDetector {
	return &baseDetector{
		indicators:  indicators,
		projectType: projectType,
	}
}

var detectors = []detector{
	newBaseDetector(
		models.ProjectTypeGo,
		anyFileExistsIndicator(
			"go.mod",
			"go.sum",
			"main.go",
		),
		anyDirExistsIndicator(
			"cmd",
			"pkg",
		),
	),
	newBaseDetector(
		models.ProjectTypeNode,
		fileExistsIndicator("package.json"),
		dirExistsIndicator("node_modules"),
	),
	newBaseDetector(
		models.ProjectTypePython,
		anyFileExistsIndicator(
			"*.py",
			"main.py",
		),
		anyDirExistsIndicator(
			"src",
			"tests",
		),
	),
	newBaseDetector(
		models.ProjectTypeRust,
		fileExistsIndicator("Cargo.toml"),
		dirExistsIndicator("src"),
	),
}

func DetectProjectType(projectPath string) models.ProjectType {
	for _, detector := range detectors {
		if detector.isDetected(projectPath) {
			return detector.ProjectType()
		}
	}
	return models.ProjectTypeUnknown
}

func DetectCurrentProjectType() models.ProjectType {
	cwd, err := os.Getwd()
	if err != nil {
		return models.ProjectTypeUnknown
	}
	return DetectProjectType(cwd)
}
