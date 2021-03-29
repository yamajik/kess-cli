package dapr

import (
	"strconv"
)

type Component struct {
	APIVersion string        `yaml:"apiVersion"`
	Kind       string        `yaml:"kind"`
	Metadata   Metadata      `yaml:"metadata"`
	Spec       ComponentSpec `yaml:"spec"`
}

type ComponentSpec struct {
	Type     string                      `yaml:"type"`
	Metadata []ComponentSpecMetadataItem `yaml:"metadata"`
}

type ComponentSpecMetadataItem struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

func CreateComponent(name string, spec ComponentSpec) Component {
	return Component{
		APIVersion: "dapr.io/v1alpha1",
		Kind:       "Component",
		Metadata: Metadata{
			Name: name,
		},
		Spec: spec,
	}
}

type RedisStateStoreComponentOptions struct {
	Host            string
	Password        string
	ActorStateStore bool
}

func CreateRedisStateStoreComponent(name string, options RedisStateStoreComponentOptions) Component {
	return CreateComponent(name, ComponentSpec{
		Type: "state.redis",
		Metadata: []ComponentSpecMetadataItem{
			{
				Name:  "redisHost",
				Value: options.Host,
			},
			{
				Name:  "redisPassword",
				Value: options.Password,
			},
			{
				Name:  "actorStateStore",
				Value: strconv.FormatBool(options.ActorStateStore),
			},
		},
	},
	)
}

type RedisPubsubComponentOptions struct {
	Host     string
	Password string
}

func CreateRedisPubsubComponent(name string, options RedisPubsubComponentOptions) Component {
	return CreateComponent(name, ComponentSpec{
		Type: "pubsub.redis",
		Metadata: []ComponentSpecMetadataItem{
			{
				Name:  "redisHost",
				Value: options.Host,
			},
			{
				Name:  "redisPassword",
				Value: options.Password,
			},
		},
	},
	)
}
