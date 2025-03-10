package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"mailHttpToSmtp/mail"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Gin Release 모드 사용
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.POST("/send_mail", func(c *gin.Context) {
		var req mail.MailRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		m := mail.NewMail()
		err := m.InitMailSenderAndRecipient(req.To)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = m.SetMailHeaderBodyToWc(req.From, req.To, req.Subject, string(req.Msg))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = m.SendMail()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "mail sent successfully"})
	})

	serverCert, err := tls.LoadX509KeyPair("./certs/server.crt", "./certs/server.key")
	if err != nil {
		log.Fatalf("서버 인증서/키 로드 실패: %v", err)
	}

	caCert, err := ioutil.ReadFile("./certs/ca.crt")
	if err != nil {
		log.Fatalf("CA 인증서 파일 읽기 실패: %v", err)
	}
	clientCAs := x509.NewCertPool()
	if ok := clientCAs.AppendCertsFromPEM(caCert); !ok {
		log.Fatalf("CA 인증서 추가 실패")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCAs,
	}

	server := &http.Server{
		Addr:      ":25443",
		Handler:   router,
		TLSConfig: tlsConfig,
	}

	log.Println("Starting TLS server: https://localhost:25443")
	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatalf("Failed to start TLS server: %v", err)
	}
}
