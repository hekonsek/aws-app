package awsom

import (
	"github.com/hekonsek/awsom/aws"
	log "github.com/sirupsen/logrus"
	"strings"
)

type envBuilder struct {
	Name   string
	Domain string
}

func NewEnvBuilder(name string, domain string) *envBuilder {
	return &envBuilder{
		Name:   name,
		Domain: domain,
	}
}

func (env *envBuilder) Create() error {
	domainSegments := strings.Split(env.Domain, ".")
	rootDomainSegments := domainSegments[len(domainSegments)-2:]
	rootDomain := strings.Join(rootDomainSegments, ".")

	zoneExists, err := aws.HostedZoneExists(rootDomain)
	if err != nil {
		return err
	}
	if !zoneExists {
		err := aws.NewHostedZoneBuilder(rootDomain).Create()
		if err != nil {
			return err
		}
	}
	err = aws.TagHostedZone(rootDomain, "env:"+env.Name, env.Domain)
	if err != nil {
		return err
	}

	return aws.NewVpcBuilder(env.Name).Create()
}

func DeleteEnv(name string) error {
	domain, _, err := aws.FirstHostedZoneTag("env:" + name)
	if err != nil {
		return err
	}

	if domain != "" {
		err = aws.DeleteHostedZoneTag(domain, "env:"+name)
		if err != nil {
			return err
		}
	} else {
		log.Debugf("Domain is not associated with environment %s. Skipping deletion.", name)
	}

	ecsCluster, err := aws.EcsClusterArnByName(name)
	if err != nil {
		return err
	}
	if ecsCluster != "" {
		err = aws.DeleteEcsCluster(ecsCluster)
		if err != nil {
			return err
		}
	}

	loadBalancerName, err := aws.LoadBalancerByVpc(name)
	if err != nil {
		return err
	}
	if loadBalancerName != "" {
		err = aws.DeleteElasticLoadBalancer(loadBalancerName)
		if err != nil {
			return err
		}
	}

	targetGroup, err := aws.LoadBalancerTargetGroupByVpc(name)
	if err != nil {
		return err
	}
	if targetGroup != "" {
		err = aws.DeleteLoadBalancerTargetGroup(targetGroup)
		if err != nil {
			return err
		}
	}

	err = aws.DeleteVpc(name)
	if err != nil {
		return err
	}

	return nil
}
