package currentuser

import (
	"context"
	"reflect" // reflect paketini import ediyoruz
	
	"github.com/gofiber/fiber/v2"
)

type contextKey string

const (
	ContextUserIDKey     contextKey = "user_id"
	ContextUserEmailKey  contextKey = "user_email"
	ContextUserTypeIDKey contextKey = "user_type_id"
)

// CurrentUser — context veya locals içinden alınan kullanıcı bilgilerini tutar
type CurrentUser struct {
	ID         uint
	Email      string
	UserTypeID uint
}

// FromFiber — Fiber locals içinden CurrentUser oluşturur
func FromFiber(c *fiber.Ctx) CurrentUser {
	authUserVal := c.Locals("authUser")
	if authUserVal == nil {
		return CurrentUser{}
	}
	
	// 1. Önce CurrentUser tipini kontrol et
	if au, ok := authUserVal.(CurrentUser); ok {
		return au
	}
	
	// 2. fiber.Map tipini kontrol et
	if au, ok := authUserVal.(fiber.Map); ok {
		return CurrentUser{
			ID:         convertToUint(au["ID"]),
			Email:      convertToString(au["Email"]),
			UserTypeID: convertToUint(au["UserTypeID"]),
		}
	}
	
	// 3. map[string]interface{} tipini kontrol et
	if au, ok := authUserVal.(map[string]interface{}); ok {
		return CurrentUser{
			ID:         convertToUint(au["ID"]),
			Email:      convertToString(au["Email"]),
			UserTypeID: convertToUint(au["UserTypeID"]),
		}
	}
	
	// 4. Reflection ile struct kontrolü
	rv := reflect.ValueOf(authUserVal)
	if rv.Kind() == reflect.Struct {
		// ID alanını al
		idField := rv.FieldByName("ID")
		emailField := rv.FieldByName("Email")
		userTypeIDField := rv.FieldByName("UserTypeID")
		
		if idField.IsValid() {
			var id uint
			var email string
			var userTypeID uint
			
			// ID değerini al
			switch idField.Kind() {
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				id = uint(idField.Uint())
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				id = uint(idField.Int())
			case reflect.Float32, reflect.Float64:
				id = uint(idField.Float())
			}
			
			// Email değerini al
			if emailField.IsValid() && emailField.Kind() == reflect.String {
				email = emailField.String()
			}
			
			// UserTypeID değerini al
			if userTypeIDField.IsValid() {
				switch userTypeIDField.Kind() {
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					userTypeID = uint(userTypeIDField.Uint())
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					userTypeID = uint(userTypeIDField.Int())
				case reflect.Float32, reflect.Float64:
					userTypeID = uint(userTypeIDField.Float())
				}
			}
			
			return CurrentUser{
				ID:         id,
				Email:      email,
				UserTypeID: userTypeID,
			}
		}
	}
	
	return CurrentUser{}
}

// SetToContext — CurrentUser bilgilerini context içine koyar
func SetToContext(ctx context.Context, user CurrentUser) context.Context {
	ctx = context.WithValue(ctx, ContextUserIDKey, user.ID)
	ctx = context.WithValue(ctx, ContextUserEmailKey, user.Email)
	ctx = context.WithValue(ctx, ContextUserTypeIDKey, user.UserTypeID)
	return ctx
}

// FromContext — context içinden CurrentUser oluşturur
func FromContext(ctx context.Context) CurrentUser {
	var cu CurrentUser
	if v := ctx.Value(ContextUserIDKey); v != nil {
		cu.ID = convertToUint(v)
	}
	if v := ctx.Value(ContextUserEmailKey); v != nil {
		if s, ok := v.(string); ok {
			cu.Email = s
		}
	}
	if v := ctx.Value(ContextUserTypeIDKey); v != nil {
		cu.UserTypeID = convertToUint(v)
	}
	return cu
}

// Yardımcı fonksiyonlar
func convertToUint(val interface{}) uint {
	if val == nil {
		return 0
	}
	
	switch v := val.(type) {
	case uint:
		return v
	case int:
		return uint(v)
	case float64:
		return uint(v)
	case float32:
		return uint(v)
	case int64:
		return uint(v)
	case uint64:
		return uint(v)
	default:
		return 0
	}
}

func convertToString(val interface{}) string {
	if val == nil {
		return ""
	}
	
	if s, ok := val.(string); ok {
		return s
	}
	
	return ""
}