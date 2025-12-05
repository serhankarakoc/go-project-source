package services

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"strings"

	"zatrano/configs/logconfig" // NOT: Kendi log paketinize göre bu satırı düzenleyin

	"go.uber.org/zap"
)

// IMailService, e-posta işlemleri için genel arayüzü tanımlar.
type IMailService interface {
	SendMail(to, subject, body string) error
}

// MailService, IMailService arayüzünü uygular ve e-posta gönderme işlevselliğini barındırır.
// Tüm yapılandırmasını ortam değişkenlerinden alır.
type MailService struct {
	host        string
	port        string
	username    string
	password    string
	fromAddress string
	fromName    string
	encryption  string
}

// NewMailService, ortam değişkenlerini okuyarak yeni ve yapılandırılmış bir MailService örneği oluşturur.
// Port belirtilmemişse, şifreleme türüne göre mantıklı bir varsayılan atar (tls->587, ssl->465).
func NewMailService() IMailService {
	encryption := strings.ToLower(getEnvWithDefault("MAIL_ENCRYPTION", "tls"))
	port := getEnvWithDefault("MAIL_PORT", "")

	if port == "" {
		switch encryption {
		case "ssl":
			port = "465"
		case "tls":
			port = "587"
		default:
			// Şifresiz bağlantılar için varsayılan port
			port = "25"
		}
	}

	return &MailService{
		host:        getEnvWithDefault("MAIL_HOST", "localhost"),
		port:        port,
		username:    getEnvWithDefault("MAIL_USERNAME", ""),
		password:    getEnvWithDefault("MAIL_PASSWORD", ""),
		fromAddress: getEnvWithDefault("MAIL_FROM_ADDRESS", ""),
		fromName:    getEnvWithDefault("MAIL_FROM_NAME", ""),
		encryption:  encryption,
	}
}

// getEnvWithDefault, bir ortam değişkenini okur veya anahtar bulunamazsa varsayılan bir değer döndürür.
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// SendMail, servisin dışarıya açılan ana e-posta gönderme metodudur.
func (m *MailService) SendMail(to, subject, body string) error {
	if to == "" {
		return fmt.Errorf("alıcı e-posta adresi (to) boş olamaz")
	}
	if m.fromAddress == "" {
		return fmt.Errorf("gönderen e-posta adresi (MAIL_FROM_ADDRESS) tanımlanmamış")
	}

	message, err := m.buildMessage(to, subject, body)
	if err != nil {
		return fmt.Errorf("e-posta mesajı oluşturulamadı: %w", err)
	}

	// Asıl gönderme işlemini yapan iç metoda yönlendir.
	err = m.send(to, message)
	if err != nil {
		logconfig.Log.Error("E-posta gönderimi başarısız oldu", zap.Error(err), zap.String("to", to))
		return err // Hatanın detayını üst katmana döndür
	}

	logconfig.Log.Info("E-posta başarıyla gönderildi", zap.String("to", to))
	return nil
}

// buildMessage, e-posta başlıklarını (From, To, Subject, Content-Type) ve gövdesini oluşturur.
func (m *MailService) buildMessage(to, subject, body string) ([]byte, error) {
	if subject == "" {
		subject = "(Konu Belirtilmemiş)"
	}

	// MAIL_FROM_NAME tanımlıysa "Ad <eposta>" formatını kullan, değilse sadece e-postayı kullan.
	fromHeader := m.fromAddress
	if m.fromName != "" {
		fromHeader = fmt.Sprintf("\"%s\" <%s>", m.fromName, m.fromAddress)
	}

	header := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n\r\n",
		fromHeader, to, subject)

	return []byte(header + body), nil
}

// send, MAIL_ENCRYPTION değerine göre doğru bağlantı yöntemini seçer ve e-postayı gönderir.
func (m *MailService) send(to string, message []byte) error {
	address := fmt.Sprintf("%s:%s", m.host, m.port)
	auth := smtp.PlainAuth("", m.username, m.password, m.host)

	switch m.encryption {
	case "tls":
		// STARTTLS Yöntemi (Port 587 için)
		client, err := smtp.Dial(address)
		if err != nil {
			return fmt.Errorf("SMTP sunucusuna bağlanılamadı: %w", err)
		}
		defer client.Quit()

		if err := client.StartTLS(&tls.Config{ServerName: m.host, InsecureSkipVerify: false}); err != nil {
			return fmt.Errorf("STARTTLS başlatılamadı: %w", err)
		}

		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("kimlik doğrulama başarısız: %w", err)
		}

		return sendMailWithClient(client, m.fromAddress, to, message)

	case "ssl":
		// SMTPS/SSL Yöntemi (Port 465 için)
		tlsConfig := &tls.Config{ServerName: m.host, InsecureSkipVerify: false}
		conn, err := tls.Dial("tcp", address, tlsConfig)
		if err != nil {
			return fmt.Errorf("SSL ile TLS bağlantısı kurulamadı: %w", err)
		}

		client, err := smtp.NewClient(conn, m.host)
		if err != nil {
			return fmt.Errorf("SSL bağlantısı üzerinden SMTP istemcisi oluşturulamadı: %w", err)
		}
		defer client.Quit()

		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("kimlik doğrulama başarısız: %w", err)
		}

		return sendMailWithClient(client, m.fromAddress, to, message)

	default:
		// Şifresiz Yöntem (Port 25 veya 587 - Tavsiye Edilmez)
		return smtp.SendMail(address, auth, m.fromAddress, []string{to}, message)
	}
}

// sendMailWithClient, var olan bir SMTP istemcisi üzerinden e-posta gönderme adımlarını çalıştırır.
// Bu yardımcı fonksiyon, TLS ve SSL durumlarındaki kod tekrarını önler.
func sendMailWithClient(client *smtp.Client, from, to string, message []byte) error {
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("SMTP göndericisi (%s) ayarlanamadı: %w", from, err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("SMTP alıcısı (%s) ayarlanamadı: %w", to, err)
	}
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("SMTP DATA komutu başlatılamadı: %w", err)
	}
	if _, err := writer.Write(message); err != nil {
		writer.Close()
		return fmt.Errorf("mesaj verisi yazılamadı: %w", err)
	}
	return writer.Close()
}

// MailService'in IMailService arayüzünü uyguladığını derleme zamanında doğrula.
var _ IMailService = (*MailService)(nil)
