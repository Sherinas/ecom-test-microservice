package handler

import (
	"net/http"
	"strconv"

	"github.com/Sherinas/ecommerce-microservices/APIGateway/client"

	"github.com/Sherinas/ecommerce-microservices/APIGateway/internal/middleware"
	"github.com/Sherinas/ecommerce-microservices/APIGateway/internal/util"
	"github.com/Sherinas/ecommerce-microservices/APIGateway/pb/admin"
	"github.com/Sherinas/ecommerce-microservices/APIGateway/pb/auth"
	"github.com/Sherinas/ecommerce-microservices/APIGateway/pb/product"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GatewayHandler struct {
	clients     *client.Clients
	logger      *zerolog.Logger
	validator   *util.AuthValidator
	adminSecret string
}

func NewGatewayHandler(clients *client.Clients, logger *zerolog.Logger, validator *util.AuthValidator, adminSecret string) *GatewayHandler {
	return &GatewayHandler{
		clients:     clients,
		logger:      logger,
		validator:   validator,
		adminSecret: adminSecret,
	}
}

func (h *GatewayHandler) RegisterRoutes(r *gin.Engine) {
	// Public routes (Auth)
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/signup", h.SignUp)
		authGroup.POST("/signin", h.SignIn)
		authGroup.POST("/logout", h.Logout)
	}

	// Admin routes (require JWT with role: admin)
	adminGroup := r.Group("/admin").Use(middleware.JWTAuthMiddleware(h.validator))
	{
		adminGroup.POST("/products", h.AddProduct)
		adminGroup.PUT("/products/:id", h.UpdateProduct)
		adminGroup.DELETE("/products/:id", h.DeleteProduct)
		adminGroup.GET("/products", h.ListAllProductsAdmin)
		adminGroup.GET("/products/:id", h.GetProductByIdAdmin)
		// adminGroup.GET("/users", h.ListUsers)
		// adminGroup.DELETE("/users/:id", h.DeleteUser)
		// adminGroup.PATCH("/users/:id/block", h.BlockUser)
	}

	// Public routes (Product)
	productGroup := r.Group("/products")
	{
		productGroup.GET("", h.ListAllProducts)
		productGroup.GET("/:id", h.GetProductById)
	}
}

func (h *GatewayHandler) SignUp(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Role     string `json:"role" binding:"required,oneof=user admin"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	// Restrict admin role creation
	ctx := c.Request.Context()
	if req.Role == "admin" {
		adminSecret := c.GetHeader("adminx--secret")
		if adminSecret != h.adminSecret {
			h.logger.Error().Msg("Invalid admin secret")

			c.JSON(http.StatusForbidden, gin.H{"error": "invalid admin secret"})
			return
		}
	}

	resp, err := h.clients.AuthClient.SignUp(ctx, &auth.SignUpRequest{
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	})
	if err != nil {
		h.logger.Error().Err(err).Str("email", req.Email).Msg("Failed to sign up")
		if status.Code(err) == codes.AlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign up"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": resp.Token})
}

func (h *GatewayHandler) SignIn(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	resp, err := h.clients.AuthClient.SignIn(c.Request.Context(), &auth.SignInRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		h.logger.Error().Err(err).Str("email", req.Email).Msg("Failed to sign in")
		if status.Code(err) == codes.NotFound || status.Code(err) == codes.Unauthenticated {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign in"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": resp.Token})
}

func (h *GatewayHandler) Logout(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	_, err := h.clients.AuthClient.Logout(c.Request.Context(), &auth.LogoutRequest{Token: req.Token})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to logout")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

func (h *GatewayHandler) AddProduct(c *gin.Context) {
	var req struct {
		Name     string  `json:"name" binding:"required"`
		Price    float32 `json:"price" binding:"required,gt=0"`
		Quantity int32   `json:"quantity" binding:"required,gte=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	token, _ := c.Get("jwt_token")
	ctx := h.validator.CreateGRPCMetadata(token.(string))

	resp, err := h.clients.AdminClient.AddProduct(ctx, &admin.AdminAddProductRequest{
		Name:     req.Name,
		Price:    req.Price,
		Quantity: req.Quantity,
	})
	if err != nil {
		h.logger.Error().Err(err).Str("name", req.Name).Msg("Failed to add product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": resp.Id})
}

func (h *GatewayHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	var req struct {
		Name     string  `json:"name" binding:"required"`
		Price    float32 `json:"price" binding:"required,gt=0"`
		Quantity int32   `json:"quantity" binding:"required,gte=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	token, _ := c.Get("jwt_token")
	ctx := h.validator.CreateGRPCMetadata(token.(string))

	_, err = h.clients.AdminClient.UpdateProduct(ctx, &admin.AdminUpdateProductRequest{
		Id:       id,
		Name:     req.Name,
		Price:    req.Price,
		Quantity: req.Quantity,
	})
	if err != nil {
		h.logger.Error().Err(err).Int64("id", id).Msg("Failed to update product")
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update product"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product updated successfully"})
}

func (h *GatewayHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	token, _ := c.Get("jwt_token")
	ctx := h.validator.CreateGRPCMetadata(token.(string))

	_, err = h.clients.AdminClient.DeleteProduct(ctx, &admin.AdminDeleteProductRequest{Id: id})
	if err != nil {
		h.logger.Error().Err(err).Int64("id", id).Msg("Failed to delete product")
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
}

func (h *GatewayHandler) ListAllProductsAdmin(c *gin.Context) {
	token, _ := c.Get("jwt_token")
	ctx := h.validator.CreateGRPCMetadata(token.(string))

	resp, err := h.clients.AdminClient.ListAllProducts(ctx, &admin.AdminListAllProductsRequest{})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to list products")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list products"})
		return
	}

	products := make([]map[string]interface{}, len(resp.Products))
	for i, p := range resp.Products {
		products[i] = map[string]interface{}{
			"id":       p.Id,
			"name":     p.Name,
			"price":    p.Price,
			"quantity": p.Quantity,
		}
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

func (h *GatewayHandler) GetProductByIdAdmin(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	token, _ := c.Get("jwt_token")
	ctx := h.validator.CreateGRPCMetadata(token.(string))

	resp, err := h.clients.AdminClient.GetProductById(ctx, &admin.AdminGetProductByIdRequest{Id: id})
	if err != nil {
		h.logger.Error().Err(err).Int64("id", id).Msg("Failed to get product")
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get product"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       resp.Id,
		"name":     resp.Name,
		"price":    resp.Price,
		"quantity": resp.Quantity,
	})
}

func (h *GatewayHandler) ListUsers(c *gin.Context) {
	token, _ := c.Get("jwt_token")
	ctx := h.validator.CreateGRPCMetadata(token.(string))

	resp, err := h.clients.AdminClient.ListUsers(ctx, &admin.AdminListUsersRequest{})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to list users")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users"})
		return
	}

	users := make([]map[string]interface{}, len(resp.Users))
	for i, u := range resp.Users {
		users[i] = map[string]interface{}{
			"id":         u.Id,
			"email":      u.Email,
			"is_blocked": u.IsBlocked,
		}
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *GatewayHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	token, _ := c.Get("jwt_token")
	ctx := h.validator.CreateGRPCMetadata(token.(string))

	_, err = h.clients.AdminClient.DeleteUser(ctx, &admin.AdminDeleteUserRequest{Id: id})
	if err != nil {
		h.logger.Error().Err(err).Int64("id", id).Msg("Failed to delete user")
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

func (h *GatewayHandler) BlockUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	token, _ := c.Get("jwt_token")
	ctx := h.validator.CreateGRPCMetadata(token.(string))

	_, err = h.clients.AdminClient.BlockUser(ctx, &admin.AdminBlockUserRequest{Id: id})
	if err != nil {
		h.logger.Error().Err(err).Int64("id", id).Msg("Failed to block user")
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to block user"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user blocked successfully"})
}

func (h *GatewayHandler) ListAllProducts(c *gin.Context) {
	resp, err := h.clients.ProductClient.ListAllProducts(c.Request.Context(), &product.ListAllProductsRequest{})
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to list products")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list products"})
		return
	}

	products := make([]map[string]interface{}, len(resp.Products))
	for i, p := range resp.Products {
		products[i] = map[string]interface{}{
			"id":       p.Id,
			"name":     p.Name,
			"price":    p.Price,
			"quantity": p.Quantity,
		}
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

func (h *GatewayHandler) GetProductById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	resp, err := h.clients.ProductClient.GetProductById(c.Request.Context(), &product.GetProductByIdRequest{Id: id})
	if err != nil {
		h.logger.Error().Err(err).Int64("id", id).Msg("Failed to get product")
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get product"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       resp.Id,
		"name":     resp.Name,
		"price":    resp.Price,
		"quantity": resp.Quantity,
	})
}
