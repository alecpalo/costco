package s3

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
	"io"
	"k8s.io/klog/v2"
)

type S3 struct {
	client *s3.Client
	bucket string
}

func Init() *S3 {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"), // MinIO doesn't require a region, but AWS SDK does
	)
	if err != nil {
		klog.Fatal(err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true // Required for MinIO
	})

	s := S3{
		client: client,
	}

	return &s
}

// StartMultiPartUpload starts a multipart layer upload, returning the uuid of the upload
// and an error, if any occurred throughout the process.
func (s *S3) StartMultiPartUpload(key string) (uuid.UUID, error) {
	input := &s3.CreateMultipartUploadInput{
		Bucket:            aws.String(s.bucket),
		Key:               aws.String(key),
		ChecksumAlgorithm: types.ChecksumAlgorithmSha256,
	}

	output, err := s.client.CreateMultipartUpload(context.Background(), input)
	if err != nil {
		return uuid.Nil, err
	}

	id, err := uuid.Parse(*output.UploadId)

	return id, err
}

// PutChunk does a thing I promise
func (s *S3) PutChunk(key string, id uuid.UUID, partNumber int32, data io.Reader) error {
	rawData, err := io.ReadAll(data)
	if err != nil {
		klog.Error(err)
		return err
	}

	checkSum := fmt.Sprintf("%x", sha256.Sum256(rawData))

	input := s3.UploadPartInput{
		Bucket:            aws.String(s.bucket),
		Key:               aws.String(key),
		UploadId:          aws.String(id.String()),
		Body:              data,
		PartNumber:        aws.Int32(partNumber),
		ChecksumAlgorithm: types.ChecksumAlgorithmSha256,
		ChecksumSHA256:    aws.String(checkSum),
	}

	output, err := s.client.UploadPart(context.Background(), &input)
	if err != nil {
		return err
	}

	if *output.ChecksumSHA256 != checkSum {
		return fmt.Errorf("checksum mismatch")
	}

	return nil
}

// CompleteLayer completes a chunked layer upload.
func (s *S3) CompleteLayer(key string, id uuid.UUID) error {
	input := s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(s.bucket),
		Key:      aws.String(key),
		UploadId: aws.String(id.String()),
	}
	_, err := s.client.CompleteMultipartUpload(context.Background(), &input)

	return err
}

func (s *S3) GetLayer(key string) (io.ReadCloser, error) {
	return nil, nil
}

// PutLayer puts an entire layer in with a single request (a monolithic upload)
func (s *S3) PutLayer(key string, data io.Reader) error {
	rawData, err := io.ReadAll(data)
	if err != nil {
		klog.Error(err)
		return err
	}

	checkSum := fmt.Sprintf("%x", sha256.Sum256(rawData))

	input := s3.PutObjectInput{
		Bucket:            aws.String(s.bucket),
		Key:               aws.String(key),
		Body:              data,
		ChecksumSHA256:    aws.String(checkSum),
		ChecksumAlgorithm: types.ChecksumAlgorithmSha256,
	}

	output, err := s.client.PutObject(context.Background(), &input)
	if err != nil {
		return err
	}

	if *output.ChecksumSHA256 != checkSum {
		return fmt.Errorf("checksum mismatch")
	}

	return nil
}

func (s *S3) DeleteLayer(key string) error {
	return nil
}

func (s *S3) ListLayers(key string) ([]string, error) {
	return nil, nil
}

func (s *S3) FindLayer(key string) (bool, error) {
	return false, nil
}

func (s *S3) CreateImage() error {
	return nil
}

func (s *S3) DeleteImage() error {
	return nil
}

func (s *S3) GetImage() error {
	return nil
}

func (s *S3) FindImage() error {
	return nil
}

func (s *S3) FindRepo(repo string) (bool, error) {
	return false, nil
}
