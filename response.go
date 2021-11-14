package ksqldb

type CommandStatus struct {
	Message string
	Status  string
}

type Stream struct {
	Name   string
	Topic  string
	Format string
	Type   string
}

type Table struct {
	Name       string
	Topic      string
	Format     string
	Type       string
	IsWindowed bool
}

type Query struct {
	QueryString string
	Sinks       string
	ID          string // The query ID
}

// this is not complete yet
type Schema struct {
	Type string
	// Fields
}

type Field struct {
	Name   string
	Schema Schema
}

type QueryDescription struct {
	StatementText string
	Fields        []Field
	Sources       []string
	Sinks         []string
	ExecutionPlan string
	Topology      string
}

type KsqlResponseSlice []KsqlResponse
type StreamSlice []Stream
type TableSlice []Table
type QuerySlice []Query

type KsqlResponse struct {
	StatementText         string
	Warnings              []string
	Type                  string            `json:"@type"`
	CommandId             string            `json:"commandId,omitempty"`
	CommandSequenceNumber int64             `json:"commandSequenceNumber,omitempty"` // -1 if the operation was unsuccessful
	CommandStatus         CommandStatus     `json:"commandStatus,omitempty"`
	Stream                *StreamSlice      `json:"streams,omitempty"`
	Tables                *TableSlice       `json:"tables,omitempty"`
	Queries               *QuerySlice       `json:"queries,omitempty"`
	QueryDescription      *QueryDescription `json:"queryDescription,omitempty"`
}
