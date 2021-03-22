package dapr

type Configuration struct {
	APIVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Metadata   Metadata          `yaml:"metadata"`
	Spec       ConfigurationSpec `yaml:"spec"`
}

type ConfigurationSpec struct {
	Tracing ConfigurationSpecTracing `yaml:"tracing,omitempty"`
}

type ConfigurationSpecTracing struct {
	SamplingRate string                         `yaml:"samplingRate,omitempty"`
	Zipkin       ConfigurationSpecTracingZipkin `yaml:"zipkin,omitempty"`
}

type ConfigurationSpecTracingZipkin struct {
	EndpointAddress string `yaml:"endpointAddress,omitempty"`
}

func CreateConfiguration(name string, spec ConfigurationSpec) Configuration {
	return Configuration{
		APIVersion: "dapr.io/v1alpha1",
		Kind:       "Configuration",
		Metadata: Metadata{
			Name: name,
		},
		Spec: spec,
	}
}
