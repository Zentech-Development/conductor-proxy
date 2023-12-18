package domain

import "context"

type Property struct {
	Name         string `json:"name" binding:"required"`
	FriendlyName string `json:"friendlyName" binding:"required"`
	DataType     string `json:"dataType" binding:"required"`
	Required     bool   `json:"required" binding:"required"`
	DefaultValue any    `json:"defaultValue" binding:"required"`
	HasDefault   bool   `json:"hasDefault" binding:"required"`
}

type Parameter struct {
	Name         string `json:"name" binding:"required"`
	FriendlyName string `json:"friendlyName" binding:"required"`
	DataType     string `json:"dataType" binding:"required"`
	Required     bool   `json:"required" binding:"required"`
	DefaultValue any    `json:"defaultValue" binding:"required"`
	HasDefault   bool   `json:"hasDefault" binding:"required"`
	Type         string `json:"type" binding:"required"`
}

type Endpoint struct {
	Name       string      `json:"name" binding:"required"`
	Path       string      `json:"path" binding:"required"`
	Method     string      `json:"method" binding:"required"`
	Parameters []Parameter `json:"parameters" binding:"required"`
}

type Resource struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	FriendlyName string     `json:"friendlyName"`
	ServiceID    string     `json:"serviceId"`
	Properties   []Property `json:"properties"`
	Endpoints    []Endpoint `json:"endpoints"`
}

type ResourceInput struct {
	Name         string     `json:"name" binding:"required"`
	FriendlyName string     `json:"friendlyName" binding:"required"`
	ServiceID    string     `json:"serviceId" binding:"required"`
	Properties   []Property `json:"properties" binding:"required"`
	Endpoints    []Endpoint `json:"endpoints" binding:"required"`
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
