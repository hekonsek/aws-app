package awsom

import "github.com/hekonsek/awsom/aws"

type envBuilder struct {
	Name string
}

func NewEnvBuilder(name string) *envBuilder {
	return &envBuilder{
		Name: name,
	}
}

func (env *envBuilder) Create() error {
	return aws.NewVpcBuilder(env.Name).Create()
}