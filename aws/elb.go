package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	awsom_session "github.com/hekonsek/awsom-session"
)

type elasticLoadBalancerBuilder struct {
	Name string
}

func NewElasticLoadBalancer(name string) *elasticLoadBalancerBuilder {
	return &elasticLoadBalancerBuilder{
		Name: name,
	}
}

func (loadBalancer *elasticLoadBalancerBuilder) Create() error {
	sess, err := awsom_session.NewSession()
	if err != nil {
		return err
	}
	elbService := elbv2.New(sess)

	subnets, err := VpcSubnetsByName(loadBalancer.Name)
	if err != nil {
		return err
	}

	loadBalancerCreated, err := elbService.CreateLoadBalancer(&elbv2.CreateLoadBalancerInput{
		Name:    aws.String(loadBalancer.Name),
		Subnets: aws.StringSlice(subnets),
		Type:    aws.String(elbv2.LoadBalancerTypeEnumApplication),
	})
	if err != nil {
		return err
	}

	_, err = elbService.CreateListener(&elbv2.CreateListenerInput{
		Port:            aws.Int64(80),
		LoadBalancerArn: loadBalancerCreated.LoadBalancers[0].LoadBalancerArn,
		Protocol:        aws.String("HTTP"),
		DefaultActions: []*elbv2.Action{
			{
				Type: aws.String("fixed-response"),
				FixedResponseConfig: &elbv2.FixedResponseActionConfig{
					MessageBody: aws.String("Default backend."),
					StatusCode:  aws.String("200"),
				},
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func DeleteElasticLoadBalancer(name string) error {
	sess, err := awsom_session.NewSession()
	if err != nil {
		return err
	}
	elbService := elbv2.New(sess)

	arn, err := LoadBalancerArnByName(name)
	if err != nil {
		return err
	}
	_, err = elbService.DeleteLoadBalancer(&elbv2.DeleteLoadBalancerInput{
		LoadBalancerArn: aws.String(arn),
	})

	return err
}

func LoadBalancerByVpc(vpc string) (string, error) {
	sess, err := awsom_session.NewSession()
	if err != nil {
		return "", err
	}
	elbService := elbv2.New(sess)

	vpcId, err := VpcId(vpc)
	if err != nil {
		return "", err
	}

	loadBalancerInfo, err := elbService.DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{})
	if err != nil {
		return "", err
	}

	for _, loadBalancer := range loadBalancerInfo.LoadBalancers {
		if *loadBalancer.VpcId == vpcId {
			return *loadBalancer.LoadBalancerName, nil
		}
	}

	return "", nil
}

func LoadBalancerArnByVpcId(vpcId string) (string, error) {
	sess, err := awsom_session.NewSession()
	if err != nil {
		return "", err
	}
	elbService := elbv2.New(sess)

	loadBalancerInfo, err := elbService.DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{})
	if err != nil {
		return "", err
	}

	for _, loadBalancer := range loadBalancerInfo.LoadBalancers {
		if *loadBalancer.VpcId == vpcId {
			return *loadBalancer.LoadBalancerArn, nil
		}
	}

	return "", nil
}

func LoadBalancerArnByName(name string) (string, error) {
	sess, err := awsom_session.NewSession()
	if err != nil {
		return "", err
	}
	elbService := elbv2.New(sess)

	loadBalancerInfo, err := elbService.DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{
		Names: aws.StringSlice([]string{name}),
	})
	if err != nil {
		return "", err
	}

	return *loadBalancerInfo.LoadBalancers[0].LoadBalancerArn, nil
}

func LoadBalancerListenerArnByName(name string) (string, error) {
	sess, err := awsom_session.NewSession()
	if err != nil {
		return "", err
	}
	elbService := elbv2.New(sess)

	loadBalancerArn, err := LoadBalancerArnByName(name)
	if err != nil {
		return "", err
	}

	listeners, err := elbService.DescribeListeners(&elbv2.DescribeListenersInput{
		LoadBalancerArn: aws.String(loadBalancerArn),
	})
	if err != nil {
		return "", err
	}

	return *listeners.Listeners[0].ListenerArn, nil
}

type loadBalancerTargetGroupBuilderBuilder struct {
	Vpc  string
	Name string
	IPs  []string
}

func NewLoadBalancerTargetGroupBuilderBuilder(vpc string, name string) *loadBalancerTargetGroupBuilderBuilder {
	return &loadBalancerTargetGroupBuilderBuilder{
		Vpc:  vpc,
		Name: name,
	}
}

func (builder *loadBalancerTargetGroupBuilderBuilder) WithIPs(IPs []string) *loadBalancerTargetGroupBuilderBuilder {
	builder.IPs = IPs
	return builder
}

func (targetGroup *loadBalancerTargetGroupBuilderBuilder) Create() (arn string, err error) {
	sess, err := awsom_session.NewSession()
	if err != nil {
		return "", err
	}
	elbService := elbv2.New(sess)

	vpcId, err := VpcId(targetGroup.Vpc)
	if err != nil {
		return "", err
	}

	createdTargetGroup, err := elbService.CreateTargetGroup(&elbv2.CreateTargetGroupInput{
		Name:               aws.String(targetGroup.Name),
		Protocol:           aws.String("HTTP"),
		Port:               aws.Int64(80),
		HealthCheckEnabled: aws.Bool(true),
		HealthCheckPath:    aws.String("/"),
		VpcId:              aws.String(vpcId),
		TargetType:         aws.String("ip"),
	})
	if err != nil {
		return "", err
	}

	var targets []*elbv2.TargetDescription
	for _, ip := range targetGroup.IPs {
		targets = append(targets, &elbv2.TargetDescription{
			Port: aws.Int64(9090),
			Id:   aws.String(ip),
		})
	}
	_, err = elbService.RegisterTargets(&elbv2.RegisterTargetsInput{
		TargetGroupArn: createdTargetGroup.TargetGroups[0].TargetGroupArn,
		Targets:        targets,
	})
	if err != nil {
		return "", err
	}

	return *createdTargetGroup.TargetGroups[0].TargetGroupArn, nil
}

func LoadBalancerTargetGroupByVpc(vpc string) (string, error) {
	sess, err := awsom_session.NewSession()
	if err != nil {
		return "", err
	}
	elbService := elbv2.New(sess)

	vpcId, err := VpcId(vpc)
	if err != nil {
		return "", err
	}

	targetGroups, err := elbService.DescribeTargetGroups(&elbv2.DescribeTargetGroupsInput{})
	if err != nil {
		return "", err
	}

	for _, group := range targetGroups.TargetGroups {
		if *group.VpcId == vpcId {
			return *group.TargetGroupName, nil
		}
	}

	return "", nil
}

func LoadBalancerTargetGroupArnByName(name string) (string, error) {
	sess, err := awsom_session.NewSession()
	if err != nil {
		return "", err
	}
	elbService := elbv2.New(sess)

	targetGroups, err := elbService.DescribeTargetGroups(&elbv2.DescribeTargetGroupsInput{
		Names: aws.StringSlice([]string{name}),
	})
	if err != nil {
		return "", err
	}

	if len(targetGroups.TargetGroups) == 0 {
		return "", nil
	}
	return *targetGroups.TargetGroups[0].TargetGroupArn, nil
}

func LoadBalancerTargetGroupExistsByName(name string) (bool, error) {
	arn, err := LoadBalancerTargetGroupArnByName(name)
	return arn != "", err
}

func DeleteLoadBalancerTargetGroup(name string) error {
	sess, err := awsom_session.NewSession()
	if err != nil {
		return err
	}
	elbService := elbv2.New(sess)

	arn, err := LoadBalancerTargetGroupArnByName(name)
	if err != nil {
		return err
	}

	_, err = elbService.DeleteTargetGroup(&elbv2.DeleteTargetGroupInput{
		TargetGroupArn: aws.String(arn),
	})
	if err != nil {
		return err
	}

	return nil
}

func AssignLoadBalancerTargetGroup(loadBalancer string, targetGroup string, host string) error {
	sess, err := awsom_session.NewSession()
	if err != nil {
		return err
	}
	elbService := elbv2.New(sess)

	listenerArn, err := LoadBalancerListenerArnByName(loadBalancer)
	if err != nil {
		return err
	}

	targetGroupArn, err := LoadBalancerTargetGroupArnByName(targetGroup)
	if err != nil {
		return err
	}

	_, err = elbService.CreateRule(&elbv2.CreateRuleInput{
		ListenerArn: aws.String(listenerArn),
		Conditions: []*elbv2.RuleCondition{
			{
				Field: aws.String("host-header"),
				HostHeaderConfig: &elbv2.HostHeaderConditionConfig{
					Values: aws.StringSlice([]string{host}),
				},
			},
		},
		Actions: []*elbv2.Action{
			{
				TargetGroupArn: aws.String(targetGroupArn),
				Type:           aws.String("forward"),
			},
		},
		Priority: aws.Int64(100),
	})
	if err != nil {
		return err
	}

	return nil
}
