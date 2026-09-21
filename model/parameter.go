package model

import "encoding/json"

// Parameter defines a workflow or task parameter.
type Parameter struct {
	Name        string            `json:"name"`
	Type        string            `json:"type,omitempty"` // string, int, float, bool, json, list
	Value       json.RawMessage   `json:"value,omitempty"`
	Default     json.RawMessage   `json:"default,omitempty"`
	Description string            `json:"description,omitempty"`
	Enum        []json.RawMessage `json:"enum,omitempty"`
	ValueFrom   *ValueFrom        `json:"valueFrom,omitempty"`
}

// ValueFrom specifies a source for a parameter value.
//
// Exactly one source should be set. Resolution tries them in the fixed order
// Path → Parameter → Expression → SecretKeyRef and uses the first non-empty
// field, silently ignoring the rest (mutual exclusivity is not yet enforced by
// validate.go). Resolution failures do not fail the task: they are reported to
// the ErrorSink and the parameter falls back to its Default, or is left empty.
type ValueFrom struct {
	// Path is a key looked up directly in the evaluation environment (EvalVars),
	// e.g. "tasks.fetch.outputs.parameters.body". The value found there is
	// marshalled to JSON as-is.
	//
	// Deprecated alias of Parameter. Unlike Parameter it is NOT normalised, so a
	// legacy "workflow.arguments.parameters.x" key will not resolve. Prefer
	// Parameter. The name comes from the original schema wording ("read value
	// from a file path"); reading from a file is not implemented — in practice
	// this branch is an env lookup identical to Parameter minus normalisation.
	Path string `json:"path,omitempty"`

	// Parameter references a value from another parameter or task output by env
	// key. Recognised prefixes:
	//
	//	workflow.parameters.<name>                workflow-level argument
	//	workflow.arguments.parameters.<name>      legacy alias of the above
	//	inputs.parameters.<name>                  current template's bound input
	//	tasks.<name>.outputs.parameters.<param>   sibling task output
	//	tasks.<name>.phase | .code | .msg         sibling task state
	//	loop_iter.*                               loop iteration context
	//	<custom>.<name>                           user-registered vars.Source
	//
	// The legacy "workflow.arguments.parameters.<name>" form is normalised to
	// "workflow.parameters.<name>" before lookup.
	Parameter string `json:"parameter,omitempty"`

	// Expression is interpolated ({{...}}) and then evaluated by the configured
	// expr.Evaluator. Requires an evaluator to be injected; without one this
	// source fails and falls back like any other resolution error.
	Expression string `json:"expression,omitempty"`

	// SecretKeyRef reads the value from the configured secret.Provider. The value
	// is never placed into the evaluation environment, keeping secrets out of
	// expression scope.
	SecretKeyRef *SecretKeyRef `json:"secretKeyRef,omitempty"`
}

// SecretKeyRef references a key in a secret store.
type SecretKeyRef struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

// Artifact defines an input/output artifact.
type Artifact struct {
	Name    string          `json:"name"`
	From    string          `json:"from,omitempty"`
	Source  *ArtifactSource `json:"source,omitempty"`
	Archive *Archive        `json:"archive,omitempty"`
}

// ArtifactSource specifies the storage backend for an artifact.
type ArtifactSource struct {
	Type   string          `json:"type"`   // oss, local, http
	Config json.RawMessage `json:"config"` // ossSourceConfig | localSourceConfig | httpSourceConfig
}

// Archive defines the archive format for an artifact.
// TODO: reserved for artifact upload/download feature; not yet used in the engine.
type Archive struct {
	Type string `json:"type,omitempty"` // none, tar, tar.gz, zip
}

// OSSSourceConfig is the configuration for OSS artifact storage.
// TODO: reserved for artifact upload/download feature; not yet used in the engine.
type OSSSourceConfig struct {
	Bucket   string `json:"bucket"`
	Key      string `json:"key"`
	Endpoint string `json:"endpoint"`
}

// LocalSourceConfig is the configuration for local artifact storage.
// TODO: reserved for artifact upload/download feature; not yet used in the engine.
type LocalSourceConfig struct {
	Path string `json:"path"`
}

// HTTPSourceConfig is the configuration for HTTP artifact storage.
// TODO: reserved for artifact upload/download feature; not yet used in the engine.
type HTTPSourceConfig struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
}
