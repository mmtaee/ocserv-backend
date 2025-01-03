package ocGroup

import (
	"api/internal/repository"
	"api/pkg/ocserv"
	"api/pkg/utils"
	"api/pkg/validator"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Controller struct {
	validator       validator.CustomValidatorInterface
	ocservGroupRepo repository.OcservGroupRepositoryInterface
}

func New() *Controller {
	return &Controller{
		validator:       validator.NewCustomValidator(),
		ocservGroupRepo: repository.NewOcservGroupRepository(),
	}
}

// UpdateDefaultOcservGroup Create Superuser account
//
// @Summary      Update Ocserv Defaults Group
// @Description  Update Ocserv Defaults Group initializing step
// @Tags         Ocserv Group
// @Accept       json
// @Produce      json
// @Failure      401 {object} middlewares.Unauthorized
// @Param        request body  models.OcGroupConfig true "oc group default config"
// @Success      200  {object}  nil
// @Failure      400 {object} utils.ErrorResponse
// @Router       /api/v1/ocserv/group [post]
func (ctrl *Controller) UpdateDefaultOcservGroup(c echo.Context) error {
	var data ocserv.OcGroupConfig
	if err := ctrl.validator.Validate(c, &data); err != nil {
		return utils.BadRequest(c, err.(error))
	}
	err := ctrl.ocservGroupRepo.UpdateDefaultGroup(c.Request().Context(), data)
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusAccepted, nil)
}

func (ctrl *Controller) Groups(c echo.Context) error {
	groups, err := ctrl.ocservGroupRepo.Groups(c.Request().Context())
	if err != nil {
		return utils.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, groups)
}

func (ctrl *Controller) CreateGroup(c echo.Context) error {
	return nil
}

func (ctrl *Controller) UpdateGroup(c echo.Context) error {
	return nil
}

func (ctrl *Controller) DeleteGroup(c echo.Context) error {
	return nil
}
