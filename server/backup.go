package server

import (
	"bytes"
	"context"
	"fmt"
	"github.com/tokuhirom/blog4/server/sobs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	backupBucketName = "your-backup-bucket-name" // Replace with your bucket name
)

func StartBackup(encryptionKey string, s3client *sobs.SobsClient) {
	log.Printf("Start taking backup")
	time.AfterFunc(1*time.Second, func() {
		takeBackup(encryptionKey, s3client)
	})
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			takeBackup(encryptionKey, s3client)
		}
	}
}

func takeBackup(encryptionKey string, s3client *sobs.SobsClient) {
	log.Print("takeBackup")

	// Generate dump file path
	date := time.Now()
	dumpFileName := fmt.Sprintf("/tmp/blog3-backup-%s.sql", date.Format("2006-01-02T15-04-05"))
	encryptedFileName := dumpFileName + ".enc"
	fmt.Println("mysqldump file name:", dumpFileName)

	// Execute mysqldump command
	err := execCommand(fmt.Sprintf(
		"mysqldump --host=%s --port=%s --user=%s --password=%s %s > %s",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
		dumpFileName,
	))
	if err != nil {
		log.Fatalf("Error taking database dump: %v", err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Printf("Error removing dump file: %s, %v", name, err)
		}
	}(dumpFileName)

	// Encrypt the dump file
	err = execCommand(fmt.Sprintf(
		"openssl enc -aes-256-cbc -salt -in %s -out %s -pass pass:%s",
		dumpFileName,
		encryptedFileName,
		encryptionKey,
	))
	if err != nil {
		log.Fatalf("Error encrypting dump file: %v", err)
	}
	fmt.Printf("Encrypted file created: %s\n", encryptedFileName)
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Printf("Error removing encrypted file: %s, %v", name, err)
		}
	}(encryptedFileName)

	// Upload file to S3
	fmt.Printf("Uploading file to S3: %s\n", encryptedFileName)
	fileContent, err := os.ReadFile(encryptedFileName)
	if err != nil {
		log.Fatalf("Error reading encrypted file: %v", err)
	}

	// TODO remove old back up files
	err = s3client.PutObjectToBackupBucket(
		context.Background(),
		filepath.Base(encryptedFileName),
		"application/octet-stream",
		int64(len(fileContent)),
		bytes.NewReader(fileContent))
	if err != nil {
		log.Printf("Error uploading file to S3: %v", err)
	}
	fmt.Printf("File uploaded to S3: %s\n", encryptedFileName)

	// Clean up temporary files
	fmt.Printf("Removing temporary files: %s, %s\n", dumpFileName, encryptedFileName)
}

func execCommand(command string) error {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command execution failed: %s, error: %v", string(output), err)
	}
	return nil
}
