package services

import (
	"context"

	"github.com/FranciscoMendes10866/queues/config"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

func UploadImageFromURL(URL string, MangaName string) (string, error) {
	result, err := config.Bucket.Upload.Upload(context.Background(), URL, uploader.UploadParams{
		Folder: MangaName,
	})

	if err != nil {
		return "", err
	}

	return result.SecureURL, nil
}
