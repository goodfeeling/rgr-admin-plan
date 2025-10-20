package user

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainRole "github.com/gbrayhan/microservices-go/src/domain/sys/role"
	domainUser "github.com/gbrayhan/microservices-go/src/domain/user"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/user"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Structures
type NewUserRequest struct {
	ID        int    `json:"id"`
	HeaderImg string `json:"header_img"`
	UserName  string `json:"user_name" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Phone     string `json:"phone"`
	Status    int16  `json:"status"`
	NickName  string `json:"nick_name"`
}

type ResponseUser struct {
	ID        int64             `json:"id"`
	UUID      string            `json:"uuid"`
	UserName  string            `json:"user_name"`
	Email     string            `json:"email"`
	NickName  string            `json:"nick_name"`
	Status    int16             `json:"status"`
	Phone     string            `json:"phone"`
	HeaderImg string            `json:"header_img"`
	Roles     []domainRole.Role `json:"roles"`
	CreatedAt domain.CustomTime `json:"created_at,omitempty"`
	UpdatedAt domain.CustomTime `json:"updated_at,omitempty"`
}

type IUserController interface {
	NewUser(ctx *gin.Context)
	GetAllUsers(ctx *gin.Context)
	GetUsersByID(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
	SearchPaginated(ctx *gin.Context)
	SearchByProperty(ctx *gin.Context)
	UserBindRoles(ctx *gin.Context)
	ResetPassword(ctx *gin.Context)
	EditPassword(ctx *gin.Context)
	ChangePassword(ctx *gin.Context)
}

type UserController struct {
	userService domainUser.IUserService
	Logger      *logger.Logger
}

func NewUserController(userService domainUser.IUserService, loggerInstance *logger.Logger) IUserController {
	return &UserController{userService: userService, Logger: loggerInstance}
}

// CreateUser
// @Summary create user
// @Description create user
// @Tags user create
// @Accept json
// @Produce json
// @Param book body NewUserRequest true  "JSON Data"
// @Success 200 {object} controllers.CommonResponseBuilder
// @Router /v1/user [post]
func (c *UserController) NewUser(ctx *gin.Context) {
	c.Logger.Info("Creating new user")
	var request NewUserRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new user", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	userModel, err := c.userService.Create(toUsecaseMapper(&request))
	if err != nil {
		c.Logger.Error("Error creating user", zap.Error(err), zap.String("email", request.Email))
		_ = ctx.Error(err)
		return
	}
	userResponse := controllers.NewCommonResponseBuilder[*ResponseUser]().
		Data(domainToResponseMapper(userModel)).
		Message("success").
		Status(0).
		Build()
	c.Logger.Info("User created successfully", zap.String("email", request.Email), zap.Int64("id", userModel.ID))
	ctx.JSON(http.StatusOK, userResponse)
}

// GetAllUsers
// @Summary get all users by
// @Description get  all users by where
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} domain.CommonResponse[[]ResponseUser]
// @Router /v1/user [get]
func (c *UserController) GetAllUsers(ctx *gin.Context) {
	c.Logger.Info("Getting all users")
	users, err := c.userService.GetAll()
	if err != nil {
		c.Logger.Error("Error getting all users", zap.Error(err))
		appError := domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Successfully retrieved all users", zap.Int("count", len(*users)))
	ctx.JSON(http.StatusOK, domain.CommonResponse[*[]*ResponseUser]{
		Data: arrayDomainToResponseMapper(users),
	})
}

// GetUserByID
// @Summary get users
// @Description get users by id
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} ResponseUser
// @Router /v1/user/{id} [get]
func (c *UserController) GetUsersByID(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid user ID parameter", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("user id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Getting user by ID", zap.Int("id", userID))
	user, err := c.userService.GetByID(userID)
	if err != nil {
		c.Logger.Error("Error getting user by ID", zap.Error(err), zap.Int("id", userID))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("Successfully retrieved user by ID", zap.Int("id", userID))
	ctx.JSON(http.StatusOK, domainToResponseMapper(user))
}

// UpdateUserInfo
// @Summary update userinfo
// @Description update userinfo
// @Tags userinfo
// @Accept json
// @Produce json
// @Param book body map[string]any  true  "JSON Data"
// @Success 200 {array} controllers.CommonResponseBuilder[ResponseUser]
// @Router /v1/user [put]
func (c *UserController) UpdateUser(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid user ID parameter for update", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Updating user", zap.Int("id", userID))
	var requestMap map[string]any
	err = controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		c.Logger.Error("Error binding JSON for user update", zap.Error(err), zap.Int("id", userID))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	err = updateValidation(requestMap)
	if err != nil {
		c.Logger.Error("Validation error for user update", zap.Error(err), zap.Int("id", userID))
		_ = ctx.Error(err)
		return
	}
	userUpdated, err := c.userService.Update(int64(userID), requestMap)
	if err != nil {
		c.Logger.Error("Error updating user", zap.Error(err), zap.Int("id", userID))
		_ = ctx.Error(err)
		return
	}
	response := controllers.NewCommonResponseBuilder[*ResponseUser]().
		Data(domainToResponseMapper(userUpdated)).
		Message("success").
		Status(0).
		Build()
	c.Logger.Info("User updated successfully", zap.Int("id", userID))
	ctx.JSON(http.StatusOK, response)
}

// DeleteUser
// @Summary delete user
// @Description delete user by id
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} domain.CommonResponse[int]
// @Router /v1/user/{id} [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid user ID parameter for deletion", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Deleting user", zap.Int("id", userID))
	err = c.userService.Delete(userID)
	if err != nil {
		c.Logger.Error("Error deleting user", zap.Error(err), zap.Int("id", userID))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("User deleted successfully", zap.Int("id", userID))
	ctx.JSON(http.StatusOK, domain.CommonResponse[int]{
		Data:    userID,
		Message: "success",
		Status:  0,
	})
}

// SearchUsersPageList
// @Summary search users
// @Description search users by query
// @Tags search users
// @Accept json
// @Produce json
// @Param page query string false "page"
// @Param pageSize query string false "PageSize"
// @Param sortBy query string false "sortBy"
// @Param sortDirection query string false "sortDirection"
// @Param status_match query string  false "status"
// @Param user_name_like query string false "userName"
// @Success 200 {object} domain.PageList[[]ResponseUser]
// @Router /v1/user/search [get]
func (c *UserController) SearchPaginated(ctx *gin.Context) {
	c.Logger.Info("Searching users with pagination")

	// Parse query parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	if pageSize < 1 {
		pageSize = 10
	}

	// Build filters
	filters := domain.DataFilters{
		Page:     page,
		PageSize: pageSize,
	}

	// Parse like filters
	likeFilters := make(map[string][]string)
	for field := range user.ColumnsUserMapping {
		if values := ctx.QueryArray(field + "_like"); len(values) > 0 {
			likeFilters[field] = values
		}
	}
	filters.LikeFilters = likeFilters

	// Parse exact matches
	matches := make(map[string][]string)
	for field := range user.ColumnsUserMapping {
		if values := ctx.QueryArray(field + "_match"); len(values) > 0 {
			matches[field] = values
		}
	}
	filters.Matches = matches

	// Parse date range filters
	var dateRanges []domain.DateRangeFilter
	for field := range user.ColumnsUserMapping {
		startStr := ctx.Query(field + "_start")
		endStr := ctx.Query(field + "_end")

		if startStr != "" || endStr != "" {
			dateRange := domain.DateRangeFilter{Field: field}

			if startStr != "" {
				if startTime, err := time.Parse(time.RFC3339, startStr); err == nil {
					dateRange.Start = &startTime
				}
			}

			if endStr != "" {
				if endTime, err := time.Parse(time.RFC3339, endStr); err == nil {
					dateRange.End = &endTime
				}
			}

			dateRanges = append(dateRanges, dateRange)
		}
	}
	filters.DateRangeFilters = dateRanges

	// Parse sorting
	sortBy := ctx.QueryArray("sortBy")
	if len(sortBy) > 0 {
		filters.SortBy = sortBy
	} else {
		filters.SortBy = []string{}
	}

	sortDirection := domain.SortDirection(ctx.DefaultQuery("sortDirection", "asc"))
	if sortDirection.IsValid() {
		filters.SortDirection = sortDirection
	}

	result, err := c.userService.SearchPaginated(filters)
	if err != nil {
		c.Logger.Error("Error searching users", zap.Error(err))
		_ = ctx.Error(err)
		return
	}
	type PageResult = domain.PageList[*[]*ResponseUser]
	response := controllers.NewCommonResponseBuilder[PageResult]().
		Data(PageResult{
			List:       arrayDomainToResponseMapper(result.Data),
			Total:      result.Total,
			Page:       result.Page,
			PageSize:   result.PageSize,
			TotalPages: result.TotalPages,
			Filters:    filters,
		}).
		Message("success").
		Status(0).
		Build()

	c.Logger.Info("Successfully searched users",
		zap.Int64("total", result.Total),
		zap.Int("page", result.Page))
	ctx.JSON(http.StatusOK, response)
}

// SearchByProperty
// @Summary  search by property
// @Description search by property
// @Tags search
// @Accept json
// @Produce json
// @Success 200 {array} []string
// @Router /v1/user/search-property [get]
func (c *UserController) SearchByProperty(ctx *gin.Context) {
	property := ctx.Query("property")
	searchText := ctx.Query("searchText")

	if property == "" || searchText == "" {
		c.Logger.Error("Missing property or searchText parameter")
		appError := domainErrors.NewAppError(errors.New("missing property or searchText parameter"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	// Validate property
	allowed := map[string]bool{
		"userName":  true,
		"email":     true,
		"firstName": true,
		"lastName":  true,
		"status":    true,
	}
	if !allowed[property] {
		c.Logger.Error("Invalid property for search", zap.String("property", property))
		appError := domainErrors.NewAppError(errors.New("invalid property"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	coincidences, err := c.userService.SearchByProperty(property, searchText)
	if err != nil {
		c.Logger.Error("Error searching by property", zap.Error(err), zap.String("property", property))
		_ = ctx.Error(err)
		return
	}

	c.Logger.Info("Successfully searched by property",
		zap.String("property", property),
		zap.Int("results", len(*coincidences)))
	ctx.JSON(http.StatusOK, coincidences)
}

// UserBindRoles
// @Summary user bind role
// @Description user bind role multiple
// @Tags user role
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Router /v1/user/{id}/role [post]
func (c *UserController) UserBindRoles(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid user ID parameter ", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	var requestMap map[string]any
	err = controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		c.Logger.Error("Error binding JSON for user bind role update", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	err = c.userService.UserBindRoles(int64(userId), requestMap)
	if err != nil {
		c.Logger.Error("Error updating  user bind role ", zap.Error(err), zap.Int("id", userId))
		_ = ctx.Error(err)
		return
	}
	response := controllers.NewCommonResponseBuilder[bool]().
		Data(true).
		Message("success").
		Status(0).
		Build()
	c.Logger.Info("Role updated successfully", zap.Int("id", userId))
	ctx.JSON(http.StatusOK, response)
}

// ResetPassword
// @Summary reset password
// @Description reset password
// @Tags password
// @Accept json
// @Produce json
// @Success 200 {object} domain.CommonResponse
// @Router /v1/user/{id}/reset-password [post]
func (c *UserController) ResetPassword(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid user ID parameter ", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	userModal, err := c.userService.ResetPassword(int64(userId))
	if err != nil {
		c.Logger.Error("Error updating  user bind role ", zap.Error(err), zap.Int("id", userId))
		_ = ctx.Error(err)
		return
	}
	response := controllers.NewCommonResponseBuilder[*domainUser.User]().
		Data(userModal).
		Message("success").
		Status(0).
		Build()
	c.Logger.Info("Role updated successfully", zap.Int("id", userId))
	ctx.JSON(http.StatusOK, response)
}

// EditPassword
// @Summary edit password
// @Description edit password
// @Tags password edit
// @Accept json
// @Produce json
// @Success 200 {object} domain.CommonResponse
// @Router /v1/user/{id}/edit-password [post]
func (c *UserController) EditPassword(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid user ID parameter ", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Creating new user")
	var request domainUser.PasswordEditRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for password data", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	userModal, err := c.userService.EditPassword(int64(userId), request)
	if err != nil {
		c.Logger.Error("Error updating  user bind role ", zap.Error(err), zap.Int("id", userId))
		_ = ctx.Error(err)
		return
	}
	response := controllers.NewCommonResponseBuilder[*domainUser.User]().
		Data(userModal).
		Message("success").
		Status(0).
		Build()
	c.Logger.Info("Role updated successfully", zap.Int("id", userId))
	ctx.JSON(http.StatusOK, response)
}

// ChangePassword implements IUserController.
// @Summary change password
// @Description change password by email
// @Tags password email
// @Accept json
// @Produce json
// @Success 200 {object} domain.CommonResponse
// @Router /v1/user/change-password [POST]
func (c *UserController) ChangePassword(ctx *gin.Context) {
	c.Logger.Info("change user password")
	var request domainUser.ChangePasswordRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for password data", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	userID, _ := controllers.NewAppUtils(ctx).GetUserID()
	if userID == 0 {
		c.Logger.Error("Error getting user ID")
		appError := domainErrors.NewAppError(errors.New("user id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	userModal, err := c.userService.ChangePasswordById(int64(userID), request.NewPasswd, ctx.Query("token"))
	if err != nil {
		c.Logger.Error("Error updating  user bind role ", zap.Error(err), zap.Int("id", userID))
		_ = ctx.Error(err)
		return
	}
	response := controllers.NewCommonResponseBuilder[*domainUser.User]().
		Data(userModal).
		Message("success").
		Status(0).
		Build()
	c.Logger.Info("Role updated successfully", zap.Int("id", userID))
	ctx.JSON(http.StatusOK, response)
}

// Mappers
func domainToResponseMapper(domainUser *domainUser.User) *ResponseUser {
	return &ResponseUser{
		ID:        domainUser.ID,
		UserName:  domainUser.UserName,
		Email:     domainUser.Email,
		NickName:  domainUser.NickName,
		UUID:      domainUser.UUID,
		Phone:     domainUser.Phone,
		HeaderImg: domainUser.HeaderImg,
		Status:    domainUser.Status,
		Roles:     domainUser.Roles,
		CreatedAt: domain.CustomTime{Time: domainUser.CreatedAt},
		UpdatedAt: domain.CustomTime{Time: domainUser.UpdatedAt},
	}
}

func arrayDomainToResponseMapper(users *[]domainUser.User) *[]*ResponseUser {
	res := make([]*ResponseUser, len(*users))
	for i, u := range *users {
		res[i] = domainToResponseMapper(&u)
	}
	return &res
}

func toUsecaseMapper(req *NewUserRequest) *domainUser.User {
	return &domainUser.User{
		UserName:  req.UserName,
		NickName:  req.NickName,
		Email:     req.Email,
		HeaderImg: req.HeaderImg,
		Phone:     req.Phone,
		Status:    req.Status,
	}
}
