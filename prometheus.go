package awsom

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
	vpcExists, err := VpcExistsByName(prometheus.Vpc)
	if err != nil {
		return err
	}
	if !vpcExists {
		err = NewVpcBuilder(prometheus.Vpc).Create()
		if err != nil {
			return err
		}
	}

	clusterExists, err := EcsClusterExistsByName(prometheus.Vpc)
	if !clusterExists {
		err = NewEcsClusterBuilder(prometheus.Vpc).Create()
		if err != nil {
			return err
		}
	}

	err = NewEcsDeploymentBuilder(prometheus.Name, prometheus.Vpc, "prom/prometheus").Create()
	if err != nil {
		return err
	}

	return nil
}

func (prometheus *prometheusBuilder) WithName(name string) *prometheusBuilder {
	prometheus.Name = name
	return prometheus
}
