package captcha

import (
	"errors"

	captchaLib "github.com/gbrayhan/microservices-go/src/infrastructure/lib/captcha"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
)

type ICaptchaService interface {
	Generate(id string) (*captchaLib.CaptchaResponse, error)
	Verify(id, answer string) (bool, error)
	Refresh(id string) (*captchaLib.CaptchaResponse, error)
}

type CaptchaUseCase struct {
	Logger         *logger.Logger
	captchaHandler *captchaLib.Captcha
}

func NewCaptchaUseCase(loggerInstance *logger.Logger, captchaHandler *captchaLib.Captcha) ICaptchaService {

	return &CaptchaUseCase{
		Logger:         loggerInstance,
		captchaHandler: captchaHandler,
	}
}

// Generate implements ICaptchaService.
func (c *CaptchaUseCase) Generate(id string) (captchaResp *captchaLib.CaptchaResponse, err error) {
	if id != "" {
		captchaResp = c.captchaHandler.Refresh(id)
		return
	}
	captchaResp = c.captchaHandler.Generate()
	if captchaResp == nil {
		err = errors.New("captcha generate failed")
	}
	return
}

// Verify implements ICaptchaService.
func (c *CaptchaUseCase) Verify(id, answer string) (bool, error) {
	isValid := c.captchaHandler.Verify(id, answer)
	return isValid, nil
}

// Refresh implements ICaptchaService.
func (c *CaptchaUseCase) Refresh(id string) (*captchaLib.CaptchaResponse, error) {
	captchaResp := c.captchaHandler.Refresh(id)
	return captchaResp, nil
}
