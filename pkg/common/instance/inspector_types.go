package instance

import corev1 "k8s.io/api/core/v1"

// The CheInspectorSpec defines all specs of Che suites
type CheInspectorSpec struct {
	Deployment CheDeploymentSpec `yaml:"deployment"`
	Tests      []CheTestsSpec    `yaml:"tests"`
}

// The CheDeploymentSpec defines the type of deployment to deploy a Che instance. Supported deployments: chectl/crwctl
type CheDeploymentSpec struct {
	Cli CliSpec `yaml:"cli,omitempty"`
}

// The CliSpec defines the flags used by Che cli
type CliSpec struct {
	Flags string `yaml:"flags"`
}

// The CheTestsSpec define the information about the suites to execute against Che instance.
// Supported suites are: happy-path, test-harness
type CheTestsSpec struct {
	Name      string           `yaml:"name"`
	Namespace string           `yaml:"namespace,omitempty"`
	Image     string           `yaml:"image"`
	Args      []string         `yaml:"args,omitempty"`
	Env       []corev1.EnvVar  `yaml:"env"`
	Artifacts CheArtifactsSpec `yaml:"artifacts,omitempty"`
}

// The CheArtifactsSpec define the information where to store tests artifacts.
type CheArtifactsSpec struct {
	FromContainerPath string `yaml:"fromContainerPath"`
	To                string `yaml:"to"`
}

// The CheReporterSpec define a basic reporter to send suites results. Options supported: slack
type CheReporterSpec struct {
	CI       string `yaml:"ci"`
	Provider string `yaml:"provider"`
	Token    string `yaml:"token"`
	URL      string `yaml:"url"`
}

// The CheInspector allows defining and managing Che suites
type CheInspector struct {
	Name            string            `yaml:"name"`
	Version         string            `yaml:"version"`
	Ide             string            `yaml:"ide"`
	Spec            CheInspectorSpec  `yaml:"spec"`
	CleanAfterTests bool              `yaml:"cleanAfterTests,omitempty"`
	Report          []CheReporterSpec `yaml:"suites,omitempty"`
}
