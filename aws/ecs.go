package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/giantswarm/retry-go"
	"github.com/go-errors/errors"
	awsom_session "github.com/hekonsek/awsom-session"
	"strings"
	"time"
)

func EcsService() (*ecs.ECS, error) {
	sess, err := awsom_session.NewSession()
	if err != nil {
		return nil, err
	}
	return ecs.New(sess), nil
}

// ECS task definition

type ecsTaskDefinitionBuilder struct {
	Name   string
	Image  string
	Memory int64
}

func (task *ecsTaskDefinitionBuilder) Create() error {
	ecsService, err := EcsService()
	if err != nil {
		return err
	}

	_, err = ecsService.RegisterTaskDefinition(&ecs.RegisterTaskDefinitionInput{
		RequiresCompatibilities: aws.StringSlice([]string{"FARGATE"}),
		NetworkMode:             aws.String("awsvpc"),
		Memory:                  aws.String(fmt.Sprintf("%d", task.Memory)),
		Cpu:                     aws.String("256"),
		Family:                  aws.String(task.Name),
		ContainerDefinitions: []*ecs.ContainerDefinition{
			{
				Name:   aws.String(task.Name),
				Image:  aws.String(task.Image),
				Memory: aws.Int64(task.Memory),
				PortMappings: []*ecs.PortMapping{
					{
						HostPort:      aws.Int64(8080),
						ContainerPort: aws.Int64(8080),
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func NewEcsTaskDefinitionBuilder(name string, image string) *ecsTaskDefinitionBuilder {
	return &ecsTaskDefinitionBuilder{
		Name:   name,
		Image:  image,
		Memory: 1024,
	}
}

func (task *ecsTaskDefinitionBuilder) WithMemory(memoryMegabytes int64) *ecsTaskDefinitionBuilder {
	task.Memory = memoryMegabytes
	return task
}

func DeleteEcsTaskDefinition(name string) error {
	ecsService, err := EcsService()
	if err != nil {
		return err
	}

	versions, err := ecsService.ListTaskDefinitions(&ecs.ListTaskDefinitionsInput{
		FamilyPrefix: aws.String(name),
	})
	if err != nil {
		return err
	}

	for _, version := range versions.TaskDefinitionArns {
		_, err = ecsService.DeregisterTaskDefinition(&ecs.DeregisterTaskDefinitionInput{
			TaskDefinition: version,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// Cluster

type ecsClusterBuilder struct {
	Name string
}

func NewEcsClusterBuilder(name string) *ecsClusterBuilder {
	return &ecsClusterBuilder{
		Name: name,
	}
}

func (cluster *ecsClusterBuilder) Create() error {
	ecsService, err := EcsService()
	if err != nil {
		return err
	}

	_, err = ecsService.CreateCluster(&ecs.CreateClusterInput{
		ClusterName: aws.String(cluster.Name),
	})
	if err != nil {
		return err
	}

	return nil
}

func EcsClusterArnByVpcId(vpcId string) (string, error) {
	ecsService, err := EcsService()
	if err != nil {
		return "", err
	}

	clusters, err := ecsService.ListClusters(&ecs.ListClustersInput{})
	if err != nil {
		return "", err
	}

	ec2Service, err := NewEc2Service()
	if err != nil {
		return "", err
	}

	for _, cluster := range clusters.ClusterArns {
		servicesArns, err := ecsService.ListServices(&ecs.ListServicesInput{
			Cluster: cluster,
		})
		if err != nil {
			return "", err
		}

		if len(servicesArns.ServiceArns) > 0 {
			services, err := ecsService.DescribeServices(&ecs.DescribeServicesInput{
				Cluster:  cluster,
				Services: servicesArns.ServiceArns,
			})
			if err != nil {
				return "", err
			}
			for _, service := range services.Services {
				subnets, err := ec2Service.DescribeSubnets(&ec2.DescribeSubnetsInput{
					SubnetIds: service.NetworkConfiguration.AwsvpcConfiguration.Subnets,
				})
				if err != nil {
					return "", err
				}
				for _, subnet := range subnets.Subnets {
					if *subnet.VpcId == vpcId {
						return *service.ClusterArn, err
					}
				}
			}
		}
	}

	return "", nil
}

func EcsClusterArnByName(name string) (string, error) {
	ecsService, err := EcsService()
	if err != nil {
		return "", err
	}

	clusters, err := ecsService.DescribeClusters(&ecs.DescribeClustersInput{
		Clusters: aws.StringSlice([]string{name}),
	})
	if len(clusters.Clusters) > 0 && *clusters.Clusters[0].Status != "INACTIVE" {
		return *clusters.Clusters[0].ClusterArn, nil
	}
	return "", nil
}

func EcsClusterExistsByName(name string) (bool, error) {
	clusterArn, err := EcsClusterArnByName(name)
	if err != nil {
		return false, err
	}
	return clusterArn != "", nil
}

func DeleteEcsCluster(clusterName string) error {
	ecsService, err := EcsService()
	if err != nil {
		return err
	}

	services, err := ecsService.ListServices(&ecs.ListServicesInput{
		Cluster: aws.String(clusterName),
	})
	if err != nil {
		return err
	}
	for _, serviceArn := range services.ServiceArns {
		_, err = ecsService.UpdateService(&ecs.UpdateServiceInput{
			Cluster:      aws.String(clusterName),
			Service:      serviceArn,
			DesiredCount: aws.Int64(0),
		})
		if err != nil {
			return err
		}
		_, err = ecsService.DeleteService(&ecs.DeleteServiceInput{
			Cluster: aws.String(clusterName),
			Service: serviceArn,
		})
		if err != nil {
			return err
		}
	}

	err = retry.Do(
		func() error {
			_, err = ecsService.DeleteCluster(&ecs.DeleteClusterInput{
				Cluster: aws.String(clusterName),
			})
			return err
		},
		retry.RetryChecker(func(err error) bool {
			return strings.Contains(err.Error(), "The Cluster cannot be deleted while Tasks are active.")
		}),
		retry.Timeout(2 * time.Minute),
		retry.Sleep(5*time.Second),
		retry.MaxTries(1000))
	if err != nil {
		return err
	}

	return nil
}

// ECS Application

type ecsDeploymentBuilder struct {
	Name    string
	Cluster string
	Image   string
}

func NewEcsDeploymentBuilder(name string, cluster string, image string) *ecsDeploymentBuilder {
	return &ecsDeploymentBuilder{
		Name:    name,
		Cluster: cluster,
		Image:   image,
	}
}

func (deployment *ecsDeploymentBuilder) Create() error {
	ecsService, err := EcsService()
	if err != nil {
		return err
	}

	err = NewEcsTaskDefinitionBuilder(deployment.Name, deployment.Image).Create()
	if err != nil {
		return err
	}

	subnets, err := VpcSubnetsByName(deployment.Cluster)
	if err != nil {
		return err
	}

	// Service with the same name is being currently deleted - we should wait for deletion process to stop
	for i := 0; i < 6; i++ {
		serviceState, err := ecsService.DescribeServices(&ecs.DescribeServicesInput{
			Services: aws.StringSlice([]string{deployment.Name}),
			Cluster:  aws.String(deployment.Cluster),
		})
		if err != nil {
			return err
		}
		if len(serviceState.Services) > 0 && *serviceState.Services[0].Status == "DRAINING" {
			if i == 5 {
				return errors.New(fmt.Sprintf("Service %s in cluster %s is draining for the past minute. Aborting startup.", deployment.Name, deployment.Cluster))
			}
			time.Sleep(time.Second * 10)
		} else {
			break
		}
	}

	_, err = ecsService.CreateService(&ecs.CreateServiceInput{
		ServiceName:    aws.String(deployment.Name),
		Cluster:        aws.String(deployment.Cluster),
		LaunchType:     aws.String("FARGATE"),
		TaskDefinition: aws.String(deployment.Name + ":1"),
		NetworkConfiguration: &ecs.NetworkConfiguration{
			AwsvpcConfiguration: &ecs.AwsVpcConfiguration{
				Subnets:        aws.StringSlice(subnets),
				AssignPublicIp: aws.String("ENABLED"),
			},
		},
		DesiredCount: aws.Int64(1),
		//LoadBalancers: []*ecs.LoadBalancer{
		//	{
		//		LoadBalancerName: aws.String(deployment.Cluster),
		//		ContainerPort: aws.Int64(9090),
		//		ContainerName: aws.String(deployment.Name),
		//		TargetGroupArn: aws.String(targetGroup),
		//	},
		//},
	})
	if err != nil {
		return err
	}

	var ips []string
	err = retry.Do(
		func() error {
			ips, err = EcsApplicationAddressesByName(deployment.Cluster, deployment.Name)
			if len(ips) < 1 {
				return errors.New("No tasks found.")
			}
			return err
		},
		retry.RetryChecker(func(err error) bool {
			return true
		}),
		retry.Timeout(time.Minute),
		retry.Sleep(5*time.Second),
		retry.MaxTries(1000))
	if err != nil {
		return err
	}

	_, err = NewLoadBalancerTargetGroupBuilderBuilder(deployment.Cluster).WithIPs(ips).Create()
	if err != nil {
		return err
	}

	err = AssignLoadBalancerTargetGroup(deployment.Cluster, deployment.Name, "/"+deployment.Name)
	if err != nil {
		return err
	}

	return nil
}

func DeleteEcsApplication(runtime string, name string) error {
	ecsService, err := EcsService()
	if err != nil {
		return err
	}

	_, err = ecsService.UpdateService(&ecs.UpdateServiceInput{
		Cluster:      aws.String(runtime),
		Service:      aws.String(name),
		DesiredCount: aws.Int64(0),
	})
	if err != nil {
		return err
	}
	_, err = ecsService.DeleteService(&ecs.DeleteServiceInput{
		Cluster: aws.String(runtime),
		Service: aws.String(name),
	})
	if err != nil {
		return err
	}

	return nil
}

func EcsApplicationAddressesByName(cluster string, name string) ([]string, error) {
	ecsService, err := EcsService()
	if err != nil {
		return nil, err
	}

	tasksArns, err := ecsService.ListTasks(&ecs.ListTasksInput{
		Cluster:     aws.String(cluster),
		ServiceName: aws.String(name),
	})
	if err != nil {
		return nil, err
	}

	tasks, err := ecsService.DescribeTasks(&ecs.DescribeTasksInput{
		Cluster: aws.String(cluster),
		Tasks:   tasksArns.TaskArns,
	})
	if err != nil {
		return nil, err
	}

	var addresses []string
	for _, task := range tasks.Tasks {
		if *task.Containers[0].LastStatus != "RUNNING" {
			continue
		}
		addresses = append(addresses, *task.Containers[0].NetworkInterfaces[0].PrivateIpv4Address)
	}

	return addresses, nil
}
