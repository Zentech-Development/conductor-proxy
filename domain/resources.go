package domain

import "context"

type Property struct {
	Name         string `json:"name"`
	FriendlyName string `json:"friendlyName"`
	DataType     string `json:"dataType"`
	Required     bool   `json:"required"`
	DefaultValue any    `json:"defaultValue"`
	HasDefault   bool   `json:"hasDefault"`
}

type Parameter struct {
	Name         string `json:"name"`
	FriendlyName string `json:"friendlyName"`
	DataType     string `json:"dataType"`
	Required     bool   `json:"required"`
	DefaultValue any    `json:"defaultValue"`
	HasDefault   bool   `json:"hasDefault"`
	Type         string `json:"type"`
}

type Endpoint struct {
	Name       string      `json:"name"`
	Path       string      `json:"path"`
	Method     string      `json:"method"`
	Parameters []Parameter `json:"parameters"`
}

type Resource struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	FriendlyName string     `json:"friendlyName"`
	AppID        string     `json:"appId"`
	Properties   []Property `json:"properties"`
	Endpoints    []Endpoint `json:"endpoints"`
}

type ResourceInput struct {
	Name         string     `json:"name"`
	FriendlyName string     `json:"friendlyName"`
	AppID        string     `json:"appId"`
	Properties   []Property `json:"properties"`
	Endpoints    []Endpoint `json:"endpoints"`
}

type ResourceRepo interface {
	GetByID(ctx context.Context, id string) (Resource, error)
	Add(ctx context.Context, resource Resource) (Resource, error)
}

type ResourceHandlers interface {
	GetByID(id string, userGroups []string) (Resource, error)
	Add(resource ResourceInput, userGroups []string) (Resource, error)
}

const (
	DataTypeString = "string"
	DataTypeBool   = "bool"
	DataTypeInt    = "int"
	DataTypeNum    = "num"
	DataTypeDate   = "date"
	DataTypeObject = "object"
	DataTypeArray  = "array"
)

const (
	ParameterTypePath     = "path"
	ParameterTypeQuery    = "query"
	ParameterTypeHeader   = "header"
	ParameterTypeBody     = "body"
	ParameterTypeBodyFlat = "bodyFlat"
)
