package server

import (
	"bytes"
	"context"
	"fmt"
	"github.com/tokuhirom/blog4/server/sobs"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	backupBucketName = "your-backup-bucket-name" // Replace with your bucket name
)

func StartBackup(encryptionKey string, s3client *sobs.SobsClient) {
	slog.Info("Start taking backup")
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
	slog.Info("takeBackup")

	// Generate dump file path
	date := time.Now()
	dumpFileName := fmt.Sprintf("/tmp/blog3-backup-%s.sql", date.Format("2006-01-02T15-04-05"))
	encryptedFileName := dumpFileName + ".enc"
	slog.Info("mysqldump file name", slog.String("filename", dumpFileName))

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
		slog.Error("Error taking database dump", slog.Any("error", err))
		return
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			slog.Error("Error removing dump file", slog.String("file", name), slog.Any("error", err))
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
		slog.Error("Error encrypting dump file", slog.Any("error", err))
		return
	}
	slog.Info("Encrypted file created", slog.String("file", encryptedFileName))
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			slog.Error("Error removing encrypted file", slog.String("file", name), slog.Any("error", err))
		}
	}(encryptedFileName)

	// Upload file to S3
	slog.Info("Uploading file to S3", slog.String("file", encryptedFileName))
	fileContent, err := os.ReadFile(encryptedFileName)
	if err != nil {
		slog.Error("Error reading encrypted file", slog.Any("error", err))
		return
	}

	// TODO remove old back up files
	err = s3client.PutObjectToBackupBucket(
		context.Background(),
		filepath.Base(encryptedFileName),
		"application/octet-stream",
		int64(len(fileContent)),
		bytes.NewReader(fileContent))
	if err != nil {
		slog.Error("Error uploading file to S3", slog.Any("error", err))
		return
	}
	slog.Info("File uploaded to S3", slog.String("file", encryptedFileName))

	// Clean up temporary files
	slog.Info("Removing temporary files", slog.String("dump_file", dumpFileName), slog.String("encrypted_file", encryptedFileName))
}

func execCommand(command string) error {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command execution failed: %s, error: %v", string(output), err)
	}
	return nil
}
