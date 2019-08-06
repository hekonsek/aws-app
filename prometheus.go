package awsom

import "github.com/hekonsek/awsom/aws"

type prometheusBuilder struct {
	Name string
	Vpc  string
}

func NewPrometheusBuilder() *prometheusBuilder {
	return &prometheusBuilder{
		Name: "prometheus",
		Vpc:  "monitoring",
	}
}

func (prometheus *prometheusBuilder) Create() error {
	vpcExists, err := aws.VpcExistsByName(prometheus.Vpc)
	if err != nil {
		return err
	}
	if !vpcExists {
		err = aws.NewVpcBuilder(prometheus.Vpc).Create()
		if err != nil {
			return err
		}
	}

	err = aws.NewElasticLoadBalancer(prometheus.Vpc).Create()
	if err != nil {
		return err
	}

	clusterExists, err := aws.EcsClusterExistsByName(prometheus.Vpc)
	if !clusterExists {
		err = aws.NewEcsClusterBuilder(prometheus.Vpc).Create()
		if err != nil {
			return err
		}
	}

	err = aws.NewEcsDeploymentBuilder(prometheus.Name, prometheus.Vpc, "prom/prometheus").Create()
	if err != nil {
		return err
	}

	return nil
}

func (prometheus *prometheusBuilder) WithName(name string) *prometheusBuilder {
	prometheus.Name = name
	return prometheus
}

func (prometheus *prometheusBuilder) WithVPc(vpc string) *prometheusBuilder {
	prometheus.Vpc = vpc
	return prometheus
}
