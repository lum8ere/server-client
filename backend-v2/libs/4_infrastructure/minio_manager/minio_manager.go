package minio_manager

// ЗАКЛАД НА БУДУЮЩЕЕ

// import (
// 	"backed-api-v2/libs/4_common/smart_context"
// 	"fmt"
// 	"io"
// 	"os"

// 	"github.com/minio/minio-go/v7"
// 	"github.com/minio/minio-go/v7/pkg/credentials"
// )

// type MinioManager struct {
// 	client *minio.Client
// }

// // NewMinioManager initializes the MinioManager by connecting to the MinIO server using environment variables.
// func NewMinioManager() (*MinioManager, error) {
// 	endpoint := os.Getenv("MINIO_ENDPOINT")
// 	if endpoint == "" {
// 		return nil, fmt.Errorf("MINIO_ENDPOINT is not set")
// 	}

// 	accessKey := os.Getenv("MINIO_ACCESS_KEY")
// 	if accessKey == "" {
// 		return nil, fmt.Errorf("MINIO_ACCESS_KEY is not set")
// 	}

// 	secretKey := os.Getenv("MINIO_SECRET_KEY")
// 	if secretKey == "" {
// 		return nil, fmt.Errorf("MINIO_SECRET_KEY is not set")
// 	}

// 	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

// 	client, err := minio.New(endpoint, &minio.Options{
// 		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
// 		Secure: useSSL,
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
// 	}

// 	return &MinioManager{
// 		client: client,
// 	}, nil
// }

// func (m *MinioManager) UploadFile(sctx smart_context.ISmartContext, bucketName, objectName string, file io.Reader, fileSize int64) (string, error) {
// 	// Загружаем объект
// 	_, err := m.client.PutObject(sctx.GetContext(), bucketName, objectName, file, fileSize, minio.PutObjectOptions{
// 		ContentType: "application/octet-stream", // Можно определить тип контента динамически
// 	})
// 	if err != nil {
// 		return "", fmt.Errorf("не удалось загрузить объект: %v", err)
// 	}

// 	// Генерация URL (можно использовать presigned URL или другой метод)
// 	fileURL := fmt.Sprintf("http://%s/%s/%s", m.client.EndpointURL().Host, bucketName, objectName)
// 	return fileURL, nil
// }
