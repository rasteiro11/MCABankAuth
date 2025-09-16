package http

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rasteiro11/MCABankAuth/src/auth/service"
	"github.com/rasteiro11/PogCore/pkg/server"
	"github.com/rasteiro11/PogCore/pkg/transport/rest"
	"github.com/rasteiro11/PogCore/pkg/validator"
)

var AuthGroupPath = "/user/auth"

type (
	HandlerOpt func(*handler)
	handler    struct {
		authService service.AuthService
	}
)

func WithAuthService(authService service.AuthService) HandlerOpt {
	return func(u *handler) {
		u.authService = authService
	}
}

func NewHandler(server server.Server, opts ...HandlerOpt) {
	h := &handler{}

	for _, opt := range opts {
		opt(h)
	}

	server.AddHandler("/signin", AuthGroupPath, http.MethodPost, h.Login)
	server.AddHandler("/register", AuthGroupPath, http.MethodPost, h.Register)
}

var ErrNotAuthorized = errors.New("not authorized")

var _ Handler = (*handler)(nil)

// Login godoc
// @Summary Sign in with email and password
// @Description Authenticate a user and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body loginRequest true "Login credentials"
// @Success 200 {object} authSession
// @Failure 400 {object} any
// @Failure 401 {object} any
// @Router /user/auth/signin [post]
func (h *handler) Login(c *fiber.Ctx) error {
	req := &loginRequest{}

	if err := c.BodyParser(req); err != nil {
		return rest.NewStatusBadRequest(c, err)
	}

	if _, err := validator.IsRequestValid(req); err != nil {
		return rest.NewResponse(c, http.StatusBadRequest, rest.WithBody(err)).JSON(c)
	}

	creds, err := h.authService.Login(c.Context(), MapLoginRequestToUser(req))
	if err != nil {
		return rest.NewStatusUnauthorized(c, err)
	}

	return rest.NewStatusOk(c, rest.WithBody(MapUserLoginResponseToHTTP(creds)))
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body registerRequest true "Registration data"
// @Success 201 {object} authSession
// @Failure 400 {object} any
// @Failure 422 {object} any
// @Router /user/auth/register [post]
func (h *handler) Register(c *fiber.Ctx) error {
	req := &registerRequest{}

	if err := c.BodyParser(req); err != nil {
		return rest.NewStatusBadRequest(c, err)
	}

	if _, err := validator.IsRequestValid(req); err != nil {
		return rest.NewResponse(c, http.StatusBadRequest, rest.WithBody(err)).JSON(c)
	}

	creds, err := h.authService.Register(c.Context(), MapRegisterRequestToDTO(req))
	if err != nil {
		return rest.NewStatusUnprocessableEntity(c, err)
	}

	return rest.NewStatusCreated(c, rest.WithBody(MapUserRegisterResponseToHTTP(creds)))
}
