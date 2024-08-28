package service

import (
	"encoding/json"
	"fmt"
	"hapemu/model"
	"log"
	"net/http"
	"net/smtp"
)

type EmailService struct {
	smtpPort int
	smtpHost string
	username string
	password string
}

func NewEmailService(port int, host, username, password string) *EmailService {
	return &EmailService{
		smtpPort: port,
		smtpHost: host,
		username: username,
		password: password,
	}
}

func (s *EmailService) sendEmail(to, subject, body string) (bool, error) {
	auth := smtp.PlainAuth("", s.username, s.password, s.smtpHost)
	from := s.username

	headers := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n", from, to, subject)
	htmlBody := fmt.Sprintf(`<html><body style="font-family: 'Arial', sans-serif; font-size: 16px;">%s</body></html>`, body)
	msg := []byte(headers + htmlBody)

	err := smtp.SendMail(fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort), auth, from, []string{to}, msg)
	if err != nil {
		return false, err
	}
	return true, nil
}

// for testing
// func EmailRecommendation() {
// 	emailService := NewEmailService(587, "smtp.gmail.com", "hapemu.id@gmail.com", "fdvtmvobhemhxvvi")

// 	var messages = `<p>Here are top 5 recommendations for you based on the quiz. We hope you find them useful:</p>
//     <ol>
//         <li>%s</li>
//         <li>smartphone 2</li>
//         <li>smartphone 3</li>
//         <li>smartphone 4</li>
//         <li>smartphone 5</li>
//     </ol>`
// 	formattedMessage := fmt.Sprintf(messages, "Samsung S24")
// 	sent, err := emailService.SendEmail("ard00243@gmail.com", "Email recommendations from hapemu", formattedMessage)

// 	if err != nil {
// 		log.Fatalf("Error when sending email: %s", err)
// 	}

// 	if sent {
// 		log.Println("Email sent successfully")
// 	} else {
// 		log.Println("fail to send email")
// 	}
// }

func (es *EmailService) EmailRecommendations(w http.ResponseWriter, r *http.Request) {
	var emailRequest model.EmailRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&emailRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var message = `<p>Here are top 5 recommendations for you based on the quiz. We hope you find them useful:</p>
    <ol>
        <li>%s</li>
        <li>%s</li>
        <li>%s</li>
        <li>%s</li>
        <li>%s</li>
    </ol>`

	var recommendations = emailRequest.Recommendations
	formattedMessage := fmt.Sprintf(message, recommendations[0], recommendations[1], recommendations[2], recommendations[3], recommendations[4])
	sent, err := es.sendEmail(emailRequest.UserEmail, "Email recommendations from hapemu", formattedMessage)

	if err != nil {
		log.Fatalf("Error when sending email: %s", err)
	}

	w.Header().Set("Content-Type", "application/json")

	var emailResponse model.EmailResponse
	if sent {
		emailResponse.Success = true
		emailResponse.Message = "Email sent successfully"
	} else {
		emailResponse.Success = false
		emailResponse.Message = "Fail to send email!!"
	}
	response, err := json.Marshal(emailResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}
