package app

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	"github.com/sfu-teamproject/smartbuy/backend/apperrors"
	"github.com/sfu-teamproject/smartbuy/backend/models"
)

// SendTmpPassword creates and sends temporary passwords
// @Summary      Creates a temporary password for a user with a certain email, saves it to a database and sends to the user
// @Tags         auth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        smartphone_id path int true "Smartphone ID"
// @Param        email body models.TmpRequest true "Email of a user"
// @Success      204
// @Router       /users/restore [post]
func (app *App) SendTmpPassword(w http.ResponseWriter, r *http.Request) {
	t := models.TmpPassword{}
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("%w: error decoding email: %w", apperrors.ErrBadRequest, err))
		return
	}
	t.Email = strings.TrimSpace(t.Email)
	_, err = app.DB.GetUserByEmail(t.Email)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error getting user: %w", err))
		return
	}
	_, err = app.DB.DeleteTmpPassword(t.Email)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		app.ErrorJSON(w, r, fmt.Errorf("error deleting old tmp password: %w", err))
		return
	}
	newT, err := GenerateTmpPassword()
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error generating tmp password: %w", err))
		return
	}
	newT.Email = t.Email
	newT, err = app.DB.CreateTmpPassword(newT)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error saving tmp password into database: %w", err))
		return
	}
	err = SendTmpPassword(newT)
	if err != nil {
		app.ErrorJSON(w, r, fmt.Errorf("error sending tmp password: %w", err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func GenerateTmpPassword() (models.TmpPassword, error) {
	t := models.TmpPassword{}
	b := make([]byte, 5)
	if _, err := rand.Read(b); err != nil {
		return t, err
	}
	pass := hex.EncodeToString(b)
	t.Password = pass
	t.ExpiresAt = time.Now().Add(time.Hour * 24)
	return t, nil
}

func SendTmpPassword(tmpPassword models.TmpPassword) error {
	smtpHost := "smtp.mail.ru"
	from := "smartbuy.store@mail.ru"
	password := "qUuejB1p8REw83xShQDP"
	smtpPort := "465"
	subject := "Smartbuy temporary password"
	body := "Ваш одноразовый пароль для входа в систему: " + tmpPassword.Password +
		"\nДействителен до: " + tmpPassword.ExpiresAt.String() +
		"\nПосле входа в систему поменяйте пароль в личном кабинете"
	// 1. Формирование сообщения
	msg := []byte(
		"To: " + tmpPassword.Email + "\r\n" +
			"From: " + from + "\r\n" +
			"Subject: " + subject + "\r\n\r\n" +
			body + "\r\n")
	// 2. Аутентификация
	auth := smtp.PlainAuth("", from, password, smtpHost)
	// 3. Установка безопасного TLS соединения (Implicit TLS для порта 465)
	conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, &tls.Config{
		ServerName: smtpHost,
	})
	if err != nil {
		return fmt.Errorf("TLS Dial failed: %w", err)
	}
	defer conn.Close()
	// 4. Создание SMTP клиента
	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return fmt.Errorf("NewClient failed: %w", err)
	}
	defer client.Close()
	// 5. Аутентификация и отправка
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP Auth failed: %w", err)
	}
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("mail command failed: %w", err)
	}
	if err = client.Rcpt(tmpPassword.Email); err != nil {
		return fmt.Errorf("rcpt command failed: %w", err)
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("data command failed: %w", err)
	}
	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("write failed: %w", err)
	}
	err = w.Close()
	if err != nil {
		return fmt.Errorf("close failed: %w", err)
	}
	return client.Quit()
}

func (app *App) LoginWithTmpPassword(login models.LoginRequest) error {
	tmpPass, err := app.DB.GetTmpPassword(login.Email)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return err
	}
	if errors.Is(err, apperrors.ErrNotFound) {
		return fmt.Errorf("tmp password not present")
	}
	if tmpPass.Password != login.Password {
		return fmt.Errorf("tmp password not incorrect")
	}
	if time.Now().After(tmpPass.ExpiresAt) {
		return fmt.Errorf("tmp password expired")
	}
	_, err = app.DB.DeleteTmpPassword(tmpPass.Email)
	if err != nil {
		return fmt.Errorf("error deleting tmp password from database: %w", err)
	}
	return nil
}
