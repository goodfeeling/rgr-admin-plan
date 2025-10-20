package captcha

import (
	captchaLib "github.com/gbrayhan/microservices-go/src/infrastructure/lib/captcha"
)

type ICaptchaService interface {
	Generate(id string) (*captchaLib.CaptchaResponse, error)
	Verify(id, answer string) (bool, error)
	Refresh(id string) (*captchaLib.CaptchaResponse, error)
}
