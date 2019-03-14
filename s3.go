package awsom

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Bucket struct {
	Name string
}

func (bucket *S3Bucket) CreateOrUpdate() error {
	sess, err := CreateSession()
	if err != nil {
		return err
	}
	s3Service := s3.New(sess)

	bucketExists, err := S3BucketExists(bucket.Name)
	if err != nil {
		return err
	}
	if !bucketExists {
		_, err = s3Service.CreateBucket(&s3.CreateBucketInput{
			Bucket: aws.String(bucket.Name),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func S3BucketExists(name string) (bool, error) {
	sess, err := CreateSession()
	if err != nil {
		return false, err
	}
	s3Service := s3.New(sess)

	buckets, err := s3Service.ListBuckets(&s3.ListBucketsInput{})
	for _, bucket := range buckets.Buckets {
		if *bucket.Name == name {
			return true, nil
		}
	}

	return false, nil
}

func DeleteS3Bucket(name string) error {
	sess, err := CreateSession()
	if err != nil {
		return err
	}
	s3Service := s3.New(sess)

	_, err = s3Service.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(name),
	})
	if err != nil {
		return err
	}

	return nil
}
