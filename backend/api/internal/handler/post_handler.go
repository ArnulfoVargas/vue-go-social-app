package handler

import (
	"Server/internal/domain"
	"Server/internal/dto"
	"Server/internal/middleware"
	"Server/internal/model"
	"Server/internal/service"
	"Server/internal/validator"

	"github.com/gofiber/fiber/v3"
)

type PostHandler struct {
	validator    *validator.Validator
	postService  domain.PostService
	mediaService domain.MediaService
}

func NewPostHandler(validator *validator.Validator, postService domain.PostService, mediaService domain.MediaService) *PostHandler {
	return &PostHandler{
		validator:    validator,
		postService:  postService,
		mediaService: mediaService,
	}
}

func SetupPostRoutes(s fiber.Router, postHandler *PostHandler) {
	g := s.Group("/posts", middleware.Protected(service.ParseJWT))

	g.Post("/", postHandler.createPost)
}

// CreatePost handles the creation of a new post
// @Summary Create a new post
// @Description Create a new post with the provided content
// @Tags posts
// @Accept json
// @Produce json
// @Param post body dto.PostRequest true "Post content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Success 201 {object} dto.MessageResponse
func (p *PostHandler) createPost(c fiber.Ctx) error {
	id, ok := getUserIdFromLocals(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error:  "unauthorized",
			Status: fiber.StatusUnauthorized,
		})
	}

	var post dto.PostRequest
	if err := c.Bind().Form(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:  "invalid request body",
			Status: fiber.StatusBadRequest,
		})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:  "invalid request body",
			Status: fiber.StatusBadRequest,
		})
	}

	files, ok := form.File["images"]

	if !ok && post.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:  "content or image is required",
			Status: fiber.StatusBadRequest,
		})
	}

	if err := p.validator.Validate(post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	var media []*model.Media
	if ok {
		media, err = p.mediaService.UploadMany(files)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
				Error:  "failed to upload images",
				Status: fiber.StatusInternalServerError,
			})
		}
	}

	uploadedPost, err := p.postService.CreatePost(id, post)
	if err != nil {
		if ok {
			p.mediaService.DeleteManyByMedia(media)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:  "failed to create post",
			Status: fiber.StatusInternalServerError,
		})
	}

	if ok {
		if err := p.postService.AttachManyImages(uploadedPost.ID.Hex(), media); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
				Error:  "failed to attach images",
				Status: fiber.StatusInternalServerError,
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(dto.MessageResponse{
		Message: uploadedPost.ID.Hex(),
	})
}
