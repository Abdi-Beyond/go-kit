package s3

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Presigner struct {
	client *s3.PresignClient
}

func NewPresigner(s3Client *s3.Client) *Presigner {
	return &Presigner{
		client: s3.NewPresignClient(s3Client),
	}
}
