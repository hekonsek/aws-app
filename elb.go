package awsom

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

type ApplicationLoadBalancer struct {
	Name string
}

func (loadBalancer *ApplicationLoadBalancer) CreateOrUpdate() error {
	sess, err := CreateSession()
	if err != nil {
		return err
	}
	elbService := elbv2.New(sess)

	subnets, err := VpcSubnetsByName(loadBalancer.Name)
	if err != nil {
		return err
	}

	_, err = elbService.CreateLoadBalancer(&elbv2.CreateLoadBalancerInput{
		Name:    aws.String(loadBalancer.Name),
		Subnets: aws.StringSlice(subnets),
		Type:    aws.String(elbv2.LoadBalancerTypeEnumApplication),
	})
	if err != nil {
		return err
	}

	return nil
}

func DeleteLoadBalancer(name string) error {
	sess, err := CreateSession()
	if err != nil {
		return err
	}
	elbService := elb.New(sess)

	_, err = elbService.DeleteLoadBalancer(&elb.DeleteLoadBalancerInput{
		LoadBalancerName: aws.String(name),
	})

	return err
}
