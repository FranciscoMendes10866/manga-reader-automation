package config

import "github.com/cloudinary/cloudinary-go"

var Bucket *cloudinary.Cloudinary

func ConnectBucket() error {
	var err error

	Bucket, err = cloudinary.NewFromParams("dvjbsarnb", "634899784977361", "gu8VaPsnHXXvNc3Jh4qLmca-WII")
	if err != nil {
		return err
	}

	return nil
}
