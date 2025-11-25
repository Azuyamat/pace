package models

type ProjectType string

const (
	ProjectTypeUnknown ProjectType = "unknown"
	ProjectTypeGo      ProjectType = "go"
	ProjectTypeNode    ProjectType = "node"
	ProjectTypePython  ProjectType = "python"
	ProjectTypeRust    ProjectType = "rust"
)

func ParseProjectType(s string) ProjectType {
	switch s {
	case "go":
		return ProjectTypeGo
	case "node":
		return ProjectTypeNode
	case "python":
		return ProjectTypePython
	case "rust":
		return ProjectTypeRust
	case "unknown":
		return ProjectTypeUnknown
	case "default":
		return ProjectTypeUnknown
	default:
		return ProjectTypeUnknown
	}
}
