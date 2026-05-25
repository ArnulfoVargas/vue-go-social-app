package handler

import (
	"Server/internal/constants"
	"Server/internal/domain"
	"Server/internal/dto"
	"Server/internal/middleware"
	"Server/internal/model"
	"Server/internal/service"
	"Server/internal/validator"
	"log"
	"mime/multipart"

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
	g.Get("/suggested", postHandler.getSuggestedPosts)
	g.Get("/:id", postHandler.getPost)
	g.Put("/:id", postHandler.updatePost)
	g.Delete("/:id", postHandler.deletePost)
	g.Get("/user/:id", postHandler.getPostsByUser)
}

// CreatePost handles the creation of a new post
// @Router /api/v1/posts [post]
// @Summary Create a new post
// @Description Create a new post with the provided content
// @Tags posts
// @Security BearerAuth
// @Accept multipart/form-data
// @Param content formData string false "Post content"
// @Param files formData file false "Attached files"
// @Produce json
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

	var post dto.UpdatePostRequest
	if err := c.Bind().Form(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:  "invalid request body",
			Status: fiber.StatusBadRequest,
		})
	}

	var files []*multipart.FileHeader
	form, err := c.MultipartForm()
	if err == nil {
		files = form.File["images"]
	}

	if len(files) == 0 && post.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:  "content or image is required",
			Status: fiber.StatusBadRequest,
		})
	}

	if err := p.validator.Validate(post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	media := make([]*model.Media, 0)
	if len(files) > 0 {
		media, err = p.mediaService.UploadMany(files)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
				Error:  "failed to upload images",
				Status: fiber.StatusInternalServerError,
			})
		}
	}

	postModel := dto.PostAdd{}

	images := dto.ModelsFromPointer(media)

	if post.Content != "" {
		postModel.Content = post.Content
	}
	if len(images) > 0 {
		postModel.Media = images
	}

	uploadedPost, err := p.postService.CreatePost(id, postModel)
	if err != nil {
		if len(images) > 0 {
			if err := p.mediaService.DeleteManyByMedia(media); err != nil {
				log.Default().Println("ERROR: failed to rollback media")
			}
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:  "failed to create post",
			Status: fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.GenericResponse[dto.PostResponse]{
		Data: dto.PostResponse{
			ID:      uploadedPost.ID.Hex(),
			Content: uploadedPost.Content,
			Media:   dto.MediasFromModels(uploadedPost.Media),
		},
	})
}

func (p *PostHandler) getPost(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:  "id is required",
			Status: fiber.StatusBadRequest,
		})
	}

	post, err := p.postService.GetPost(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:  "failed to get post",
			Status: fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.GenericResponse[dto.PostResponse]{
		Data: dto.PostResponse{
			ID:      post.ID.Hex(),
			Content: post.Content,
			Media:   dto.MediasFromModels(post.Media),
		},
	})
}

func (p *PostHandler) deletePost(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:  "id is required",
			Status: fiber.StatusBadRequest,
		})
	}

	if err := p.postService.DeletePost(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:  "failed to delete post",
			Status: fiber.StatusInternalServerError,
		})
	}

	return c.SendStatus(200)
}

func (p *PostHandler) updatePost(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:  "id is required",
			Status: fiber.StatusBadRequest,
		})
	}

	var req dto.UpdatePostRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:  "invalid request",
			Status: fiber.StatusBadRequest,
		})
	}

	post, err := p.postService.UpdatePost(id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:  "failed to update post",
			Status: fiber.StatusInternalServerError,
		})
	}

	return c.Status(200).JSON(dto.PostResponse{
		ID:      post.ID.Hex(),
		Content: post.Content,
		Media:   dto.MediasFromModels(post.Media),
	})
}

func (p *PostHandler) getPostsByUser(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:  "id is required",
			Status: fiber.StatusBadRequest,
		})
	}

	posts, err := p.postService.GetPostsByUserId(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:  "failed to get posts by user",
			Status: fiber.StatusInternalServerError,
		})
	}

	return c.Status(200).JSON(dto.PostsResponse{
		Posts: dto.PostsFromModels(posts),
	})
}

func (p *PostHandler) getSuggestedPosts(c fiber.Ctx) error {
	id, ok := getUserIdFromLocals(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
			Error:  "unauthorized",
			Status: fiber.StatusUnauthorized,
		})
	}

	posts, err := p.postService.GetSuggestedPosts(id, constants.DEFAULT_PAGE_SIZE)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:  "failed to get suggested posts",
			Status: fiber.StatusInternalServerError,
		})
	}

	return c.Status(200).JSON(dto.PostsResponse{
		Posts: dto.PostsFromModels(posts),
	})
}
