package utilities

import (
	"github.com/sampson-golang/utilities/boolable"
	"github.com/sampson-golang/utilities/container"
	"github.com/sampson-golang/utilities/container/merge"
	"github.com/sampson-golang/utilities/env"
	"github.com/sampson-golang/utilities/output"
	"github.com/sampson-golang/utilities/strutil"
)

var (
	// Type conversion utilities
	Boolable = boolable.From

	// Container utilities
	Contains  = container.Contains
	Dig       = container.Dig
	DigAssign = container.DigAssign

	// Environment variable utilities
	EnvExists        = env.Exists
	GetEnv           = env.Get
	GetPresentEnv    = env.GetPresent
	LookupEnv        = env.Lookup
	LookupPresentEnv = env.LookupPresent

	// Merging utilities
	MergeParams  = merge.Params
	MergeStructs = merge.Structs

	// Output utilities
	Prettify    = output.Prettify
	PrettyPrint = output.PrettyPrint

	// String utilities
	Squish = strutil.Squish
)

type Set container.Set
