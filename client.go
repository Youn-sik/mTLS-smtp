package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// MailRequest는 서버의 /send_mail 엔드포인트에 보낼 요청 데이터 구조입니다.
type MailRequest struct {
	From    string `json:"from" binding:"required"`
	To      string `json:"to" binding:"required"`
	Subject string `json:"subject" binding:"required"`
	Msg     []byte `json:"msg" binding:"required"`
}

func main() {
	// 클라이언트 인증서와 키 로드
	clientCert, err := tls.LoadX509KeyPair("./certs/client.crt", "./certs/client.key")
	if err != nil {
		log.Fatalf("클라이언트 인증서 로드 실패: %v", err)
	}

	// CA 인증서를 읽어와서 신뢰할 CA 풀 생성
	caCert, err := ioutil.ReadFile("./certs/ca.crt")
	if err != nil {
		log.Fatalf("CA 인증서 파일 읽기 실패: %v", err)
	}
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		log.Fatalf("CA 인증서 추가 실패")
	}

	// TLS 클라이언트 설정: 클라이언트 인증서 제공 및 CA 인증서 검증
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool,
	}

	// HTTP 클라이언트에 커스텀 TLS 설정 적용
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	client := &http.Client{
		Transport: transport,
	}

	// 요청 데이터 생성: "to" 필드는 수신자 이메일, "msg"는 메시지의 바이트 배열
	reqBody := MailRequest{
		From:    "FROM <FROM@ai.kr>",
		To:      "yscho20@koolsign.net",
		Subject: "[PROJECT] 테스트 이메일 전송",
		Msg:     []byte("안녕하세요. [PROJECT] 서비스입니다."),
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatalf("JSON 인코딩 실패: %v", err)
	}

	// mTLS를 사용하여 서버에 POST 요청 전송
	resp, err := client.Post("https://localhost:25443/send_mail", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("POST 요청 실패: %v", err)
	}
	defer resp.Body.Close()

	// 서버 응답 읽기
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("응답 읽기 실패: %v", err)
	}
	fmt.Printf("Response: %s\n", body)
}
