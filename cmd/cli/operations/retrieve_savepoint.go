package operations

import (
	"errors"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/afero"
)

var s3Schemes = map[string]bool{
	"s3":  true,
	"s3a": true,
	"s3p": true,
}

func (o RealOperator) retrieveLatestSavepoint(dir string) (string, error) {
	u, err := url.Parse(dir)

	if err == nil && s3Schemes[u.Scheme] {
		return o.retrieveLatestSavepointS3(u)
	}

	return o.retrieveLatestSavepointLocal(dir)
}

func (o RealOperator) retrieveLatestSavepointS3(dir *url.URL) (string, error) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		return "", errors.New("AWS_REGION env var must be specified for S3 savepoint directories")
	}

	config := &aws.Config{
		Region:      aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewEnvCredentials(),
	}

	sess, err := session.NewSession()
	if err != nil {
		return "", errors.New("creating S3 session: " + err.Error())
	}

	client := s3.New(sess, config)

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(dir.Host),
	}
	if dir.Path != "" {
		input.Prefix = aws.String(strings.TrimLeft(dir.Path, "/"))
	}

	var newestFile url.URL
	var newestTime time.Time
	err = client.ListObjectsV2Pages(input, func(output *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, object := range output.Contents {
			if strings.HasSuffix(*object.Key, "_metadata") && object.LastModified.After(newestTime) {
				newestTime = *object.LastModified
				newestFile = url.URL{Scheme: dir.Scheme, Host: dir.Host, Path: *object.Key}
			}
		}
		return true
	})
	if err != nil {
		return "", errors.New("listing S3 objects: " + err.Error())
	}

	return newestFile.String(), nil
}

func (o RealOperator) retrieveLatestSavepointLocal(dir string) (string, error) {
	if strings.HasSuffix(dir, "/") {
		dir = strings.TrimSuffix(dir, "/")
	}

	files, err := afero.ReadDir(o.Filesystem, dir)
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", errors.New("No savepoints present in directory: " + dir)
	}

	var newestFile string
	var newestTime int64
	for _, f := range files {
		filePath := dir + "/" + f.Name()
		fi, err := o.Filesystem.Stat(filePath)
		if err != nil {
			return "", err
		}
		currTime := fi.ModTime().Unix()
		if currTime > newestTime {
			newestTime = currTime
			newestFile = filePath
		}
	}

	return newestFile, nil
}
