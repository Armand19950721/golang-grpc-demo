package s3Utils

import (
	"errors"
	"io"
	"service/protos/ArContent"
	"service/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	sess, _ = session.NewSession(&aws.Config{
		Region:      aws.String("ap-northeast-1"),
		Credentials: credentials.NewStaticCredentials(utils.GetEnv("AWS_S3_IP"), utils.GetEnv("AWS_S3_SECRET"), ""),
	})

	downloader = s3manager.NewDownloader(sess)
	svc        = s3.New(sess)
	bucket     = getDefaultBucket()
)

func DownloadFile(name string) ([]byte, error) {
	utils.PrintObj(name, "DownloadFile")
	buf := aws.NewWriteAtBuffer([]byte{})

	numBytes, err := downloader.Download(buf,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(name),
		})
	if err != nil {
		utils.PrintObj("Unable to download name: " + name + " , err:" + err.Error())
		return []byte{}, err
	}

	utils.PrintObj("Downloaded: " + name + ", numBytes:" + utils.ToString(int(numBytes)) + " bytes")

	return buf.Bytes(), nil
}

func getDefaultBucket() string {
	return utils.GetEnv("AWS_S3_BUCKET")
}

func MoveFileWithBucket(fileNameFrom string, imageTypeFrom ArContent.ArContentImageType, imageTypeTo ArContent.ArContentImageType) error {
	err := CopyFileWithBucket(fileNameFrom, imageTypeFrom, fileNameFrom, imageTypeTo)
	if err != nil {
		return err
	}

	err = DeleteFile(fileNameFrom, imageTypeFrom)
	if err != nil {
		return err
	}

	return nil
}

func CopyFileWithBucket(fileNameFrom string, imageTypeFrom ArContent.ArContentImageType, fileNameTo string, imageTypeTo ArContent.ArContentImageType) error {
	utils.PrintObj([]string{fileNameFrom, imageTypeFrom.String(), imageTypeTo.String()}, "CopyFileWithBucket")

	// check file source exist
	exist, err := FileExists(fileNameFrom, imageTypeFrom)
	if err != nil {
		return err
	}

	if !exist {
		// return errors.New("file source not exist: " + fileNameFrom)
		utils.PrintObj("file source not exist: " + fileNameFrom + ", stop CopyFileWithBucket")
		return nil
	}

	// get folder path
	fromFolder, err := utils.GetFolderPath(imageTypeFrom)
	if err != nil {
		return err
	}

	fromTo, err := utils.GetFolderPath(imageTypeTo)
	if err != nil {
		return err
	}

	pathFrom := bucket + fromFolder
	pathTo := bucket + fromTo

	utils.PrintObj([]string{pathFrom, pathTo})

	// Copy the item
	res, err := svc.CopyObject(&s3.CopyObjectInput{
		CopySource: aws.String(pathFrom + fileNameFrom),
		Bucket:     aws.String(pathTo),
		Key:        aws.String(fileNameTo),
	})

	if err != nil {
		return err
	}

	utils.PrintObj(res)

	// Wait to see if the item got copied
	err = svc.WaitUntilObjectExists(&s3.HeadObjectInput{Bucket: aws.String(pathFrom), Key: aws.String(fileNameFrom)})
	if err != nil {
		return err
	}

	return nil
}

func DeleteFile(fileName string, imageType ArContent.ArContentImageType) error {
	utils.PrintObj([]string{fileName, imageType.String()}, "DeleteFile")
	if fileName == "" {
		return nil
	}

	// get folder path
	path, err := utils.GetFolderPath(imageType)
	if err != nil {
		return err
	}

	// Delete
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucket + path), Key: aws.String(fileName)})
	if err != nil {
		return err
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	})

	if err != nil {
		return err
	}

	return nil
}

func FileExists(fileName string, imageType ArContent.ArContentImageType) (bool, error) {
	// get folder path
	path, err := utils.GetFolderPath(imageType)
	if err != nil {
		return false, err
	}

	// find
	_, err = svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket + path),
		Key:    aws.String(fileName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NotFound":
				return false, nil
			default:
				return false, err
			}
		}
		return false, err
	}
	return true, nil
}

func UploadImageFile(file io.Reader, fileName string, imageType ArContent.ArContentImageType) error {
	utils.PrintObj([]string{fileName, imageType.String()}, "UploadFile")
	// get folder path
	path, err := utils.GetFolderPath(imageType)
	if err != nil {
		return err
	}

	// upload
	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket + path),
		Key:    aws.String(fileName),
		Body:   file,
	})

	if err != nil {
		utils.PrintObj(err.Error())
		return err
	}

	return nil
}

func UploadFile(file io.Reader, fileName string, path string) error {
	// check
	if path == "" {
		return errors.New("path is empty")
	}

	// upload
	uploader := s3manager.NewUploader(sess)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket + path),
		Key:    aws.String(fileName),
		Body:   file,
	})

	if err != nil {
		utils.PrintObj(err.Error())
		return err
	}

	return nil
}
