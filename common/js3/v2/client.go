package js3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	awstrace "github.com/DataDog/dd-trace-go/contrib/aws/aws-sdk-go-v2/v2/aws"
	"github.com/Freelance-launchpad/backend-interview/common/jerror/v2"
	"github.com/Freelance-launchpad/backend-interview/common/jlog/v2"
	"github.com/Freelance-launchpad/backend-interview/common/jutils"
	"github.com/aws/aws-sdk-go-v2/aws"
	signer "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	transfermanagertypes "github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type s3API interface {
	DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, opts ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	HeadObject(ctx context.Context, params *s3.HeadObjectInput, opts ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
	GetObject(ctx context.Context, params *s3.GetObjectInput, opts ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

type s3UploaderAPI interface {
	UploadObject(ctx context.Context, input *transfermanager.UploadObjectInput, opts ...func(*transfermanager.Options)) (*transfermanager.UploadObjectOutput, error)
}

type s3PresignerAPI interface {
	PresignGetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.PresignOptions)) (*signer.PresignedHTTPRequest, error)
}

type s3Client struct {
	bucket    string
	client    s3API
	uploader  s3UploaderAPI
	presigner s3PresignerAPI
}

type Configuration struct {
	Bucket string `yaml:"bucket"`
	Region string `yaml:"region"`
}

func New(appName string, conf Configuration) FileProvider {
	awsConfig, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(conf.Region))
	if err != nil {
		slog.Error("Failed to load AWS config", slog.Any("error", err))
		os.Exit(1)
	}

	awstrace.AppendMiddleware(&awsConfig, awstrace.WithService(appName))

	client := s3.NewFromConfig(awsConfig, func(o *s3.Options) {
		o.UsePathStyle = true
		o.Logger = &jlog.AWSLogger{}
	})

	return &s3Client{
		bucket:    conf.Bucket,
		client:    client,
		uploader:  transfermanager.New(client),
		presigner: s3.NewPresignClient(client),
	}
}

func (c *s3Client) Store(ctx context.Context, identifier string, readSeeker io.ReadSeeker) error {
	return c.StoreWithMetadata(ctx, identifier, nil, readSeeker)
}

func (c *s3Client) StoreWithMetadata(ctx context.Context, identifier string, metadata map[string]string, readSeeker io.ReadSeeker) error {
	const errorMsg = "S3Client.StoreWithMetadata has failed"

	contentType, err := jutils.GetMimeType(readSeeker)
	if err != nil {
		return fmt.Errorf("%s: %w", errorMsg, err)
	}

	err = c.StoreWithMetadataAndContentType(ctx, identifier, contentType, metadata, readSeeker)
	if err != nil {
		return fmt.Errorf("%s to store the object on s3: %w", errorMsg, err)
	}

	return nil
}

func (c *s3Client) StoreWithMetadataAndContentType(ctx context.Context, identifier, contentType string, metadata map[string]string, readSeeker io.ReadSeeker) error {
	if cursor, err := readSeeker.Seek(0, io.SeekStart); err != nil || cursor != 0 {
		return fmt.Errorf("S3Client.StoreWithMetadataAndContentType: Seek err (before mime): %w", err)
	}

	_, err := c.uploader.UploadObject(ctx, &transfermanager.UploadObjectInput{
		Bucket:       &c.bucket,
		Key:          &identifier,
		Body:         readSeeker,
		StorageClass: transfermanagertypes.StorageClassIntelligentTiering,
		ContentType:  &contentType,
		Metadata:     metadata,
	})
	if err != nil {
		return fmt.Errorf("S3Client.StoreWithMetadataAndContentType: UploadObject err: %w", err)
	}
	return nil
}

func (c *s3Client) Delete(ctx context.Context, identifier string) error {
	const errorMsg = "S3Client.Delete has failed"

	_, err := c.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(identifier),
	})
	if err != nil {
		if _, ok := errors.AsType[*types.NoSuchKey](err); ok {
			return jerror.NewNotFoundError("FILE_NOT_FOUND")
		}
		return fmt.Errorf("%s to delete object from s3: %w", errorMsg, err)
	}

	return nil
}

func (c *s3Client) GetURL(ctx context.Context, identifier string) (string, error) {
	url, _, err := c.GetURLWithMetadata(ctx, identifier)
	return url, err
}

func (c *s3Client) GetURLWithMetadata(ctx context.Context, identifier string) (string, map[string]string, error) {
	errorMsg := "S3Client.GetURLWithMetadata has failed with identifier: " + identifier

	output, err := c.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(identifier),
	})
	if err != nil {
		return "", nil, fmt.Errorf("%s to HeadObject: %w", errorMsg, err)
	}

	presignedRequest, err := c.presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket:               aws.String(c.bucket),
		Key:                  aws.String(identifier),
		ResponseCacheControl: aws.String("no-store, no-cache, must-revalidate"),
	}, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		return "", nil, fmt.Errorf("%s to presign request: %w", errorMsg, err)
	}

	return presignedRequest.URL, output.Metadata, nil
}

func (c *s3Client) Download(ctx context.Context, identifier string) (FileOutput, error) {
	const errorMsg = "S3Client.Download has failed"

	result, err := c.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(identifier),
	})
	if err != nil {
		return FileOutput{}, fmt.Errorf("%s to download the object from S3: %w", errorMsg, err)
	}

	fileOutput := FileOutput{
		Reader: result.Body,
	}

	path := strings.Split(identifier, "/")
	fileOutput.FileName = path[len(path)-1]

	if result.ContentLength != nil {
		fileOutput.ContentLength = *result.ContentLength
	}

	if result.ContentType != nil {
		fileOutput.ContentType = *result.ContentType
	}

	return fileOutput, nil
}
