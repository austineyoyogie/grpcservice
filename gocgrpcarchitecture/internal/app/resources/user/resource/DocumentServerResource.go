package resource

import (
	"bytes"
	"context"
	"errors"
	"gocgrpcarchitecture/generate/proto/documentProto"
	"gocgrpcarchitecture/internal/app/resources/user/models"
	"gocgrpcarchitecture/internal/app/resources/user/service"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

var (
	ErrInvalidType = errors.New("file type or extension not allowed")
	ErrEmptyFile   = errors.New("file content cannot be empty")
)

// Allowed MIME types mapped to their valid extensions
var allowedTypes = map[string]string{
	"application/pdf": ".pdf",
	"image/jpeg":      ".jpg",
	"image/png":       ".png",
}

type DocumentServerResource struct {
	documentProto.UnimplementedDocumentServiceServer
	DocumentService service.DocumentService
	DB              *gorm.DB
}

func (D *DocumentServerResource) ValidateAndSave(ctx context.Context, Filename, Title string, buffer *bytes.Buffer) (*models.Document, error) {
	fileData := buffer.Bytes()
	if len(fileData) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "File content is empty")
	}

	// Check extension
	ext := strings.ToLower(filepath.Ext(Filename))
	if ext != ".docx" && ext != ".pdf" {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid file extension")
	}

	// Sneak peek the magic bytes (First 512 bytes)
	sniffSize := min(len(fileData), 512)

	detectedMime := http.DetectContentType(fileData[:sniffSize])

	// Match MIME type structure and verify alignment with the provided extension
	expectedExt, allowed := allowedTypes[detectedMime]
	if !allowed || (ext != expectedExt && !(detectedMime == "application/pdf" || detectedMime == "image/jpeg" || detectedMime == "image/png" && ext == ".jpeg")) {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid file type")
	}
	document := &models.Document{
		Id:       uuid.New().String(),
		Filename: Filename,
		Title:    Title,
		MimeType: detectedMime,
		Size:     uint64(len(fileData)),
		Data:     fileData,
	}
	if err := D.DocumentService.Create(ctx, document); err != nil {
		return nil, err
	}
	return document, nil

}

func (D *DocumentServerResource) UploadDocument(stream documentProto.DocumentService_UploadDocumentServer) error {
	//var Filename string
	//var Title string
	var buffer bytes.Buffer

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.Internal, "Failed to read stream. %v", err)
		}

		switch data := req.Data.(type) {
		case *documentProto.UploadRequest_Metadata:
			//Filename = data.Metadata.Filename
			//Title = data.Metadata.Title
		case *documentProto.UploadRequest_ChunkData:
			if _, err := buffer.Write(data.ChunkData); err != nil {
				return status.Errorf(codes.Internal, "Failed to write chunk. %v", err)
			}
		}
	}

	//document, err := D.ValidateAndSave(stream.Context(), Filename, Title, &buffer)
	//if err != nil {
	//	if errors.Is(err, ErrInvalidType) {
	//		return status.Errorf(codes.InvalidArgument, "Invalid document. %v", err)
	//	}
	//	return status.Errorf(codes.Internal, "database storage failed: %v", err)
	//}
	return stream.SendAndClose(&documentProto.UploadResponse{
		//Id:       document.Id,
		//Filename: document.Filename,
		//Size:     document.Size,
	})
}

func (D *DocumentServerResource) DownloadDocument(req *documentProto.DownloadRequest, stream documentProto.DocumentService_DownloadDocumentServer) error {
	doc, err := D.DocumentService.GetByID(stream.Context(), req.GetMetadata().Id)
	if err != nil {
		return status.Errorf(codes.NotFound, "Failed to get document. %v", err)
	}
	reader := bytes.NewReader(doc.Data)
	chunkSize := 32 * 1024 // 32KB Chunk
	buffer := make([]byte, chunkSize)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.Internal, "Failed to chunk document. %v", err)
		}
		if err := stream.Send(&documentProto.DownloadResponse{ChunkData: buffer[:n]}); err != nil {
			return status.Errorf(codes.Internal, "Failed to send chunk. %v", err)
		}
	}
	return nil
}

func (D *DocumentServerResource) GetDocument(ctx context.Context, request *documentProto.DownloadRequest) (*documentProto.DownloadResponse, error) {
	document, err := D.DocumentService.GetByID(ctx, request.GetMetadata().Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Failed to get document. %v", err)
	}
	return &documentProto.DownloadResponse{ChunkData: document.Data}, nil
}
func (D *DocumentServerResource) UpdateDocument(ctx context.Context, request *documentProto.UpdateRequest) (*documentProto.UpdateResponse, error) {
	//document, err := D.DocumentService.GetByID(ctx, request.GetMetadata().Id, request.Metadata().Title, request.Metadata().Description)
	if err, _ := D.DocumentService.GetByID(ctx, request.GetMetadata().Id); err != nil {
		return nil, status.Errorf(codes.NotFound, "Failed to get document. %v", err)
	}

	return &documentProto.UpdateResponse{Success: true}, nil
}

/*
func (s *DocumentServerResource) UploadDocument(stream documentProto.DocumentService_UploadDocumentServer) error {
	var fileName string
	var fileType string
	var fileBuffer bytes.Buffer

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.Internal, "Failed to read stream. %v", err)
		}

		switch data := req.Data.(type) {
		case *documentProto.UploadRequest_Metadata:
			fileName = data.Metadata.Name
			fileType = data.Metadata.FileType
		case *documentProto.UploadRequest_Chunk:
			fileBuffer.Write(data.Chunk)
		}
	}

	docID := uuid.New().String()
	uploadDir := "./storage/"
	_ = os.MkdirAll(uploadDir, os.ModePerm)
	filePath := filepath.Join(uploadDir, fmt.Sprintf("%s_%s", docID, fileName))

	if err := os.WriteFile(filePath, fileBuffer.Bytes(), 0644); err != nil {
		return status.Errorf(codes.Internal, "Failed to write file to disk. %v", err)
	}
	docEntity := &models.Document{
		Id:       docID,
		Name:     fileName,
		FileType: fileType,
		FilePath: filePath,
	}
	if err := s.DocumentService.Save(docEntity); err != nil {
		return status.Errorf(codes.Internal, "Failed to database persistence. %v", err)
	}

	return stream.SendAndClose(&documentProto.UploadResponse{
		Id:      docID,
		Message: "Document uploaded successfully.",
		Success: true,
	})
}

func (s *DocumentServerResource) GetDocument(ctx context.Context, req *documentProto.GetDocumentRequest, stream documentProto.DocumentService_GetDocumentServer) error {
	doc, err := s.DocumentService.FindById(ctx, req.GetId())
	if err != nil {
		return status.Errorf(codes.NotFound, "Document not found. %v", err)
	}
	file, err := os.Open(doc.FilePath)
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to access local file asset: %v", err)
	}
	defer file.Close()

	// Step A: Stream metadata first
	err = stream.Send(&documentProto.GetDocumentResponse{
		Data: &documentProto.GetDocumentResponse_Metadata{
			Metadata: &documentProto.Metadata{
				Name:     doc.Name,
				FileType: doc.FileType,
			},
		},
	})
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to send metadata. %v", err)
	}

	// Step B: Stream byte content chunks
	buffer := make([]byte, 1024) // 1KB buffer
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.Internal, "Failed to read file. %v", err)
		}
		err = stream.Send(&documentProto.GetDocumentResponse{
			Data: &documentProto.GetDocumentResponse_Chuck{
				Chuck: buffer[:n],
			},
		})
		if err != nil {
			return status.Errorf(codes.Internal, "Failed to send chunk. %v", err)
		}
	}
	return nil
}

func (s *DocumentServerResource) UpdateDocument(ctx context.Context, req *documentProto.UpdateDocumentRequest) (*documentProto.UpdateDocumentResponse, error) {
	doc, err := s.DocumentService.FindById(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Document not found. %v", err)
	}
	// Adjust the document entity with the new values
	doc.Name = req.GetName()
	doc.FileType = req.GetFileType()

	if err := s.DocumentService.Update(ctx, doc); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update document. %v", err)
	}
	return &documentProto.UpdateDocumentResponse{
		Message: "Document fields adjusted accurately.",
		Success: true,
	}, nil
}

func (s *DocumentServerResource) DeleteDocument(ctx context.Context, req *documentProto.DeleteDocumentRequest) (*documentProto.DeleteDocumentResponse, error) {
	doc, err := s.DocumentService.FindById(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Document not found. %v", err)
	}

	// Remove matching physical file from storage
	if err := os.Remove(doc.FilePath); err != nil && !os.IsNotExist(err) {
		return nil, status.Errorf(codes.Internal, "Failed erasing local physical document. %v", err)

	}
	if err := s.DocumentService.Delete(ctx, req.GetId()); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete document. %v", err)
	}
	return &documentProto.DeleteDocumentResponse{
		Message: "Document fully discarded from storage successfully.",
		Success: true,
	}, nil
}
*/
//=======================================
