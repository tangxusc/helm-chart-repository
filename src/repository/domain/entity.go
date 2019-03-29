package domain

import "time"

const APIVersionV1 = "v1"

type IndexFile struct {
	APIVersion string                   `json:"apiVersion" yaml:"apiVersion"`
	Generated  time.Time                `json:"generated" yaml:"generated"`
	Entries    map[string]ChartVersions `json:"entries" yaml:"entries"`
	PublicKeys []string                 `json:"publicKeys,omitempty" yaml:"publicKeys,omitempty"`
}

type ChartVersions []*ChartVersion

// ChartVersion represents a chart entry in the IndexFile
type ChartVersion struct {
	Metadata `json:"" yaml:",inline"`
	URLs     []string  `json:"urls" yaml:"urls"`
	Created  time.Time `json:"created,omitempty" yaml:"created,omitempty"`
	Removed  bool      `json:"removed,omitempty" yaml:"removed,omitempty"`
	Digest   string    `json:"digest,omitempty" yaml:"digest,omitempty"`
}

type Metadata struct {
	// The name of the chart
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// The URL to a relevant project page, git repo, or contact person
	Home string `json:"home,omitempty" yaml:"home,omitempty"`
	// Source is the URL to the source code of this chart
	Sources []string `json:"sources,omitempty" yaml:"sources,omitempty"`
	// A SemVer 2 conformant version string of the chart
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
	// A one-sentence description of the chart
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// A list of string keywords
	Keywords []string `json:"keywords,omitempty" yaml:"keywords,omitempty"`
	// A list of name and URL/email address combinations for the maintainer(s)
	Maintainers []*Maintainer `json:"maintainers,omitempty" yaml:"maintainers,omitempty"`
	// The name of the template engine to use. Defaults to 'gotpl'.
	Engine string `json:"engine,omitempty" yaml:"engine,omitempty"`
	// The URL to an icon file.
	Icon string `json:"icon,omitempty" yaml:"icon,omitempty"`
	// The API Version of this chart.
	ApiVersion string `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	// The condition to check to enable chart
	Condition string `json:"condition,omitempty" yaml:"condition,omitempty"`
	// The tags to check to enable chart
	Tags string `json:"tags,omitempty" yaml:"tags,omitempty"`
	// The version of the application enclosed inside of this chart.
	AppVersion string `json:"appVersion,omitempty" yaml:"appVersion,omitempty"`
	// Whether or not this chart is deprecated
	Deprecated bool `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	// TillerVersion is a SemVer constraints on what version of Tiller is required.
	// See SemVer ranges here: https://github.com/Masterminds/semver#basic-comparisons
	TillerVersion string `json:"tillerVersion,omitempty" yaml:"tillerVersion,omitempty"`
	// Annotations are additional mappings uninterpreted by Tiller,
	// made available for inspection by other applications.
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	// KubeVersion is a SemVer constraint specifying the version of Kubernetes required.
	KubeVersion string `json:"kubeVersion,omitempty" yaml:"kubeVersion,omitempty"`
}

// Maintainer describes a Chart maintainer.
type Maintainer struct {
	// Name is a user name or organization name
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Email is an optional email address to contact the named maintainer
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
	// Url is an optional URL to an address for the named maintainer
	Url string `json:"url,omitempty" yaml:"url,omitempty"`
}
