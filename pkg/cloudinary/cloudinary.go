package cloudinarypkg

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/iqbal2604/dear-talk-api.git/pkg/config"
)

type CloudinaryClient struct {
	cld *cloudinary.Cloudinary
}

func boolPtr(b bool) *bool { return &b }

func NewCloudinaryClient(cfg *config.CloudinaryConfig) (*CloudinaryClient, error) {
	cld, err := cloudinary.NewFromParams(
		cfg.CloudName,
		cfg.APIKey,
		cfg.APISecret,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to init cloudinary: %w", err)
	}
	return &CloudinaryClient{cld: cld}, nil
}

func (c *CloudinaryClient) UploadAvatar(ctx context.Context, file multipart.File, userID uint) (string, error) {
	uploadResult, err := c.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID:  fmt.Sprintf("deartalk/avatars/user_%d", userID),
		Folder:    "deartalk/avatars",
		Overwrite: boolPtr(true),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload avatar: %w", err)
	}
	return uploadResult.SecureURL, nil
}

func (c *CloudinaryClient) DeleteAvatar(ctx context.Context, userID uint) error {
	publicID := fmt.Sprintf("deartalk/avatars/user_%d", userID)
	_, err := c.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}
