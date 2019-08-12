package aws

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	awsom_session "github.com/hekonsek/awsom-session"
	"strings"
)

type hostedZoneBuilder struct {
	Domain string
}

func NewHostedZoneBuilder(domain string) *hostedZoneBuilder {
	return &hostedZoneBuilder{
		Domain: domain,
	}
}

func (hostedZone *hostedZoneBuilder) Create() error {
	sess, err := awsom_session.NewSession()
	if err != nil {
		return err
	}
	route53Service := route53.New(sess)

	_, err = route53Service.CreateHostedZone(&route53.CreateHostedZoneInput{
		Name:            aws.String(hostedZone.Domain),
		CallerReference: aws.String(hostedZone.Domain),
	})
	if err != nil {
		return err
	}

	return nil
}

func HostedZoneExists(domain string) (bool, error) {
	id, err := HostedZoneIdByDomain(domain)
	if err != nil {
		return false, err
	}
	return id != "", nil
}

func HostedZoneIdByDomain(domain string) (string, error) {
	sess, err := awsom_session.NewSession()
	if err != nil {
		return "", err
	}
	route53Service := route53.New(sess)

	zones, err := route53Service.ListHostedZones(&route53.ListHostedZonesInput{})
	if err != nil {
		return "", err
	}

	for _, zone := range zones.HostedZones {
		if *zone.Name == domain+"." {
			return *zone.Id, nil
		}
	}

	return "", nil
}

func HostedZoneDomainById(id string) (string, error) {
	sess, err := awsom_session.NewSession()
	if err != nil {
		return "", err
	}
	route53Service := route53.New(sess)

	zones, err := route53Service.GetHostedZone(&route53.GetHostedZoneInput{
		Id: aws.String(id),
	})
	if err != nil {
		return "", err
	}

	if zones.HostedZone != nil {
		return *zones.HostedZone.Name, nil
	}

	return "", nil
}

func DeleteHostedZone(domain string) error {
	sess, err := awsom_session.NewSession()
	if err != nil {
		return err
	}
	route53Service := route53.New(sess)

	zoneId, err := HostedZoneIdByDomain(domain)
	if err != nil {
		return err
	}

	_, err = route53Service.DeleteHostedZone(&route53.DeleteHostedZoneInput{
		Id: aws.String(zoneId),
	})
	if err != nil {
		return err
	}

	return nil
}

func TagHostedZone(domain string, key string, value string) error {
	sess, err := awsom_session.NewSession()
	if err != nil {
		return err
	}
	route53Service := route53.New(sess)

	zoneId, err := HostedZoneIdByDomain(domain)
	if err != nil {
		return err
	}

	_, err = route53Service.ChangeTagsForResource(&route53.ChangeTagsForResourceInput{
		ResourceType: aws.String("hostedzone"),
		ResourceId:   aws.String(zoneId),
		AddTags: []*route53.Tag{
			{
				Key:   aws.String(key),
				Value: aws.String(value),
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func DeleteHostedZoneTag(domain string, key string) error {
	sess, err := awsom_session.NewSession()
	if err != nil {
		return err
	}
	route53Service := route53.New(sess)

	zoneId, err := HostedZoneIdByDomain(domain)
	if err != nil {
		return err
	}

	_, err = route53Service.ChangeTagsForResource(&route53.ChangeTagsForResourceInput{
		ResourceType:  aws.String("hostedzone"),
		ResourceId:    aws.String(zoneId),
		RemoveTagKeys: []*string{aws.String(key)},
	},
	)
	if err != nil {
		return err
	}

	return nil
}

func HostedZoneTags(domain string) (map[string]string, error) {
	sess, err := awsom_session.NewSession()
	if err != nil {
		return nil, err
	}
	route53Service := route53.New(sess)

	zoneId, err := HostedZoneIdByDomain(domain)
	if err != nil {
		return nil, err
	}
	zone, err := route53Service.ListTagsForResource(&route53.ListTagsForResourceInput{
		ResourceType: aws.String("hostedzone"),
		ResourceId:   aws.String(zoneId),
	})
	if err != nil {
		return nil, err
	}
	tags := map[string]string{}
	for _, tag := range zone.ResourceTagSet.Tags {
		tags[*tag.Key] = *tag.Value
	}
	return tags, nil
}

func FirstHostedZoneTag(key string) (string, string, error) {
	sess, err := awsom_session.NewSession()
	if err != nil {
		return "", "", err
	}
	route53Service := route53.New(sess)

	zones, err := route53Service.ListHostedZones(&route53.ListHostedZonesInput{})
	if err != nil {
		return "", "", err
	}
	zonesIds := []*string{}
	for _, zone := range zones.HostedZones {
		zonesIds = append(zonesIds, aws.String(strings.Replace(*zone.Id, "/hostedzone/", "", 1)))
	}
	if len(zonesIds) > 10 {
		return "", "", errors.New("cannot list more than 10 hosted zones")
	}
	tags, err := route53Service.ListTagsForResources(&route53.ListTagsForResourcesInput{
		ResourceType: aws.String("hostedzone"),
		ResourceIds:  zonesIds,
	})
	for _, tagSet := range tags.ResourceTagSets {
		for _, tag := range tagSet.Tags {
			if *tag.Key == key {
				domain, err := HostedZoneDomainById(*tagSet.ResourceId)
				if err != nil {
					return "", "", err
				}
				return domain, *tag.Value, nil
			}
		}
	}

	return "", "", nil
}
