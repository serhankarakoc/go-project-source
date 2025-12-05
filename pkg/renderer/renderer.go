package renderer

import (
	"net/http"

	"zatrano/pkg/currentuser"
	"zatrano/pkg/flashmessages"
	"zatrano/pkg/formflash"

	"github.com/gofiber/fiber/v2"
)

const (
	CsrfTokenKey        = "CsrfToken"
	FlashSuccessKeyView = "Success"
	FlashErrorKeyView   = "Error"
	OldInputKey         = "Old"
	ValidationErrorsKey = "ValidationErrors"
)

func prepareRenderData(c *fiber.Ctx, data fiber.Map) fiber.Map {
	renderData := make(fiber.Map)

	// CSRF Token'ı ekle
	renderData[CsrfTokenKey] = c.Locals("csrf")

	// Flash mesajlarını al ve ekle
	flashData, _ := flashmessages.GetFlashMessages(c)
	if flashData.Success != "" {
		renderData[FlashSuccessKeyView] = flashData.Success
	}

	// FORMDAN GELEN VERİLERİ AL
	if formData, err := formflash.GetData(c); err == nil && len(formData) > 0 {
		renderData[OldInputKey] = formData
	}

	// VALİDASYON HATALARINI AL
	if validationErrors, err := formflash.GetValidationErrors(c); err == nil && len(validationErrors) > 0 {
		renderData[ValidationErrorsKey] = validationErrors
	}

	// CurrentUser bilgisini çek ve ekle
	currentUser := currentuser.FromFiber(c)
	if currentUser.ID != 0 {
		renderData["User"] = currentUser
	}

	var handlerError string
	if data == nil {
		data = fiber.Map{}
	}

	if errVal, ok := data[FlashErrorKeyView]; ok {
		if errStr, okStr := errVal.(string); okStr {
			handlerError = errStr
		}
	}

	if c != nil {
		data["Path"] = c.Path()
	}

	// Handler'dan gelen verileri renderData'ya birleştir
	for key, value := range data {
		renderData[key] = value
	}

	// Hata mesajlarını birleştirme
	combinedError := flashData.Error
	if handlerError != "" {
		if combinedError != "" {
			combinedError += " | " + handlerError
		} else {
			combinedError = handlerError
		}
	}

	if combinedError != "" {
		renderData[FlashErrorKeyView] = combinedError
	} else {
		delete(renderData, FlashErrorKeyView)
	}

	return renderData
}

func Render(c *fiber.Ctx, template string, layout string, data fiber.Map, statusCode ...int) error {
	status := http.StatusOK
	if len(statusCode) > 0 {
		status = statusCode[0]
	}

	finalData := prepareRenderData(c, data)

	if layout == "" {
		return c.Status(status).Render(template, finalData)
	} else {
		return c.Status(status).Render(template, finalData, layout)
	}
}
