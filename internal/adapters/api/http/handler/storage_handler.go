package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"io"
	"starter/internal/adapters/api/http"
	"starter/internal/core/storage"
)

type StorageHandler struct {
	Storage storage.Storage
}

func NewStorageHandler(storage storage.Storage) *StorageHandler {
	return &StorageHandler{
		Storage: storage,
	}
}

func (handler *StorageHandler) Upload(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		log.Error().Err(err).Msg("failed to upload file")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(
			http.ErrorResponse("Failed to upload file"),
		)
	}

	f, err := file.Open()
	defer f.Close()

	if err != nil {
		log.Error().Err(err).Msg("failed to open file")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(
			http.ErrorResponse("Failed to open file"),
		)
	}

	data, err := io.ReadAll(f)
	if err != nil {
		log.Error().Err(err).Msg("failed to read file")
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(
			http.ErrorResponse("Failed to open file"),
		)
	}

	uploadRequest := storage.UploadRequest{
		Bucket: "bucket",
		Name:   file.Filename,
		Data:   data,
	}

	upload := handler.Storage.Upload(ctx.Context(), uploadRequest)
	if upload == nil {
		log.Error().Err(err).Msg("failed to save file")
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			http.ErrorResponse("Failed to save file"),
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		http.SuccessResponse(upload, "Upload successful"))
}
