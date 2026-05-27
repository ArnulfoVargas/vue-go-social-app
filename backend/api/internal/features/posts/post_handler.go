package posts

import (
	"Server/internal/constants"
	"Server/internal/features/media"
	"Server/internal/helpers"
	"Server/internal/shared"
	"Server/internal/validator"
	"log"
	"mime/multipart"

	"github.com/gofiber/fiber/v3"
)

type PostHandler struct {
	validator    *validator.Validator
	postService  PostService
	mediaService media.MediaService
}

func NewPostHandler(validator *validator.Validator, postService PostService, mediaService media.MediaService) *PostHandler {
	return &PostHandler{
		validator:    validator,
		postService:  postService,
		mediaService: mediaService,
	}
}

func SetupPostRoutes(s fiber.Router, postHandler *PostHandler) {
	g := s.Group("/posts", shared.Protected(shared.ParseJWT))

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
// @Failure 400 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Success 201 {object} shared.MessageResponse
func (p *PostHandler) createPost(c fiber.Ctx) error {
	id, ok := helpers.GetUserIdFromLocals(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(shared.ErrorResponse{
			Error:  "unauthorized",
			Status: fiber.StatusUnauthorized,
		})
	}

	var post UpdatePostRequest
	if err := c.Bind().Form(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse{
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
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse{
			Error:  "content or image is required",
			Status: fiber.StatusBadRequest,
		})
	}

	if err := p.validator.Validate(post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	mediafiles := make([]*media.Media, 0)
	if len(files) > 0 {
		mediafiles, err = p.mediaService.UploadMany(files)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse{
				Error:  "failed to upload images",
				Status: fiber.StatusInternalServerError,
			})
		}
	}

	postModel := PostAdd{}

	images := media.ModelsFromPointer(mediafiles)

	if post.Content != "" {
		postModel.Content = post.Content
	}
	if len(images) > 0 {
		postModel.Media = images
	}

	uploadedPost, err := p.postService.CreatePost(id, postModel)
	if err != nil {
		if len(images) > 0 {
			if err := p.mediaService.DeleteManyByMedia(mediafiles); err != nil {
				log.Default().Println("ERROR: failed to rollback media")
			}
		}
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse{
			Error:  "failed to create post",
			Status: fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(shared.GenericResponse[PostResponse]{
		Data: PostResponse{
			ID:      uploadedPost.ID.Hex(),
			Content: uploadedPost.Content,
			Media:   media.MediasFromModels(uploadedPost.Media),
		},
	})
}

// @Summary Get a post
// @Description Get a post by its ID
// @Tags posts
// @Accept json
// @Produce json
// @Failure 400 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
func (p *PostHandler) getPost(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse{
			Error:  "id is required",
			Status: fiber.StatusBadRequest,
		})
	}

	post, err := p.postService.GetPost(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse{
			Error:  "failed to get post",
			Status: fiber.StatusInternalServerError,
		})
	}

	return c.Status(fiber.StatusOK).JSON(shared.GenericResponse[PostResponse]{
		Data: PostResponse{
			ID:      post.ID.Hex(),
			Content: post.Content,
			Media:   media.MediasFromModels(post.Media),
		},
	})
}

// @Summary Delete a post
// @Description Delete a post by its ID
// @Failure 400 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Success 200 {object} shared.MessageResponse
// @Router /api/v1/posts/{id} [delete]
func (p *PostHandler) deletePost(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse{
			Error:  "id is required",
			Status: fiber.StatusBadRequest,
		})
	}

	if err := p.postService.DeletePost(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse{
			Error:  "failed to delete post",
			Status: fiber.StatusInternalServerError,
		})
	}

	return c.SendStatus(200)
}

// @Summary Update a post
// @Description Update a post by its ID
// @Failure 400 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Success 200 {object} shared.MessageResponse
// @Router /api/v1/posts/{id} [put]
func (p *PostHandler) updatePost(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse{
			Error:  "id is required",
			Status: fiber.StatusBadRequest,
		})
	}

	var req UpdatePostRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse{
			Error:  "invalid request",
			Status: fiber.StatusBadRequest,
		})
	}

	post, err := p.postService.UpdatePost(id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse{
			Error:  "failed to update post",
			Status: fiber.StatusInternalServerError,
		})
	}

	return c.Status(200).JSON(PostResponse{
		ID:      post.ID.Hex(),
		Content: post.Content,
		Media:   media.MediasFromModels(post.Media),
	})
}

// @Summary Get posts by user
// @Description Get all posts by a user
// @Failure 400 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Success 200 {object} PostsResponse
// @Router /api/v1/posts/user/{id} [get]
func (p *PostHandler) getPostsByUser(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(shared.ErrorResponse{
			Error:  "id is required",
			Status: fiber.StatusBadRequest,
		})
	}

	userPosts, err := p.postService.GetPostsByUserId(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse{
			Error:  "failed to get posts by user",
			Status: fiber.StatusInternalServerError,
		})
	}

	return c.Status(200).JSON(PostsResponse{
		Posts: PostsFromModels(userPosts),
	})
}

// @Summary Get suggested posts
// @Description Get suggested posts for the authenticated user
// @Failure 400 {object} shared.ErrorResponse
// @Failure 500 {object} shared.ErrorResponse
// @Success 200 {object} PostsResponse
// @Router /api/v1/posts/suggested [get]
func (p *PostHandler) getSuggestedPosts(c fiber.Ctx) error {
	id, ok := helpers.GetUserIdFromLocals(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(shared.ErrorResponse{
			Error:  "unauthorized",
			Status: fiber.StatusUnauthorized,
		})
	}

	postsSug, err := p.postService.GetSuggestedPosts(id, constants.DEFAULT_PAGE_SIZE)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(shared.ErrorResponse{
			Error:  "failed to get suggested posts",
			Status: fiber.StatusInternalServerError,
		})
	}

	return c.Status(200).JSON(PostsResponse{
		Posts: PostsFromModels(postsSug),
	})
}
