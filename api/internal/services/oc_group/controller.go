package ocGroup

import (
	"api/internal/repository"
	_ "api/internal/routes/middlewares"
	"api/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/mmtaee/go-oc-utils/handler/ocgroup"
	"net/http"
)

type Controller struct {
	validator       utils.CustomValidatorInterface
	ocservGroupRepo repository.OcservGroupRepositoryInterface
}

func New() *Controller {
	return &Controller{
		validator:       utils.NewCustomValidator(),
		ocservGroupRepo: repository.NewOcservGroupRepository(),
	}
}

// Groups 		 List Of Groups
//
// @Summary      List Of Groups
// @Description  List Of Groups Sort By Name
// @Tags         Ocserv Group
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Success      200 {array}  ocgroup.OcservGroupConfigInfo
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Router       /api/v1/ocserv/groups [get]
func (ctrl *Controller) Groups(c echo.Context) error {
	var (
		groups *[]ocgroup.OcservGroupConfigInfo
		err    error
	)
	groups, err = ctrl.ocservGroupRepo.Groups(c.Request().Context())
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, groups)
}

// GroupNames 		 List Of Group Names
//
// @Summary      List Of Group Names
// @Description  List Of Group Names Sort By Name
// @Tags         Ocserv Group
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Success      200  {array} string
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Router       /api/v1/ocserv/groups/names [get]
func (ctrl *Controller) GroupNames(c echo.Context) error {
	names, err := ctrl.ocservGroupRepo.GroupNames(c.Request().Context())
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, names)
}

// UpdateDefaultOcservGroup  Update Ocserv Defaults Group
//
// @Summary      Update Ocserv Defaults Group
// @Description  Update Ocserv Defaults Group initializing step
// @Tags         Ocserv Group
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param        request body  ocgroup.OcservGroupConfig true "oc group default config"
// @Success      200  {object}  nil
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Router       /api/v1/ocserv/groups/defaults [post]
func (ctrl *Controller) UpdateDefaultOcservGroup(c echo.Context) error {
	var data ocgroup.OcservGroupConfig
	if err := ctrl.validator.Validate(c, &data); err != nil {
		return utils.BadRequest(c, err)
	}
	err := ctrl.ocservGroupRepo.UpdateDefaultGroup(c.Request().Context(), &data)
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusAccepted, nil)
}

// CreateGroup   Create Ocserv Group
//
// @Summary      Create Ocserv Group
// @Description  Create Ocserv Group by given name
// @Tags         Ocserv Group
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param        request body  CreateGroupRequest true "oc group config"
// @Success      200  {object}  nil
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Router       /api/v1/ocserv/groups [post]
func (ctrl *Controller) CreateGroup(c echo.Context) error {
	var data CreateGroupRequest
	if err := ctrl.validator.Validate(c, &data); err != nil {
		return utils.BadRequest(c, err)
	}
	err := ctrl.ocservGroupRepo.CreateOrUpdateGroup(c.Request().Context(), data.Name, data.Config)
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, nil)
}

// UpdateGroup   Update Ocserv Group
//
// @Summary      Update Ocserv Group
// @Description  Update Ocserv Group
// @Tags         Ocserv Group
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param 		 name path string true "Group Name"
// @Param        request body  ocgroup.OcservGroupConfig true "oc group config"
// @Success      200  {object}  nil
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Router       /api/v1/ocserv/groups/:name [post]
func (ctrl *Controller) UpdateGroup(c echo.Context) error {
	var data ocgroup.OcservGroupConfig
	if err := ctrl.validator.Validate(c, &data); err != nil {
		return utils.BadRequest(c, err)
	}
	err := ctrl.ocservGroupRepo.CreateOrUpdateGroup(c.Request().Context(), c.Param("name"), &data)
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, nil)
}

// DeleteGroup  Delete Ocserv Group
//
// @Summary      Delete Ocserv Group
// @Description  Delete Ocserv Group by given name
// @Tags         Ocserv Group
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param 		 name path string true "Group Name"
// @Success      200  {object}  nil
// @Failure      400 {object} utils.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Router       /api/v1/ocserv/groups/:name [delete]
func (ctrl *Controller) DeleteGroup(c echo.Context) error {
	err := ctrl.ocservGroupRepo.DeleteGroup(c.Request().Context(), c.Param("name"))
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusNoContent, nil)
}
