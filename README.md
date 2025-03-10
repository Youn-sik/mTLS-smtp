mTLS를 구현하기 위해서는 서버와 클라이언트 모두 신뢰할 수 있는 CA(인증기관)로부터 발급받은 인증서를 사용해야 합니다. 일반적으로 다음 단계로 진행합니다.

---

### 1. CA(인증기관) 인증서 생성

먼저, 자체 CA를 생성하여 서버와 클라이언트 인증서를 서명할 수 있습니다.

```sh
# CA 개인키 생성
openssl genrsa -out ca.key 2048

# CA 인증서 생성 (유효기간 100년, 36500일)
openssl req -x509 -new -nodes -key ca.key -days 36500 -out ca.crt -subj "/CN=MyCustomCA"
```

---

### 2. 서버 인증서 생성

서버용 개인키와 CSR(인증서 서명 요청)을 생성한 후, CA로 서명합니다.

```sh
# 서버 개인키 생성
openssl genrsa -out server.key 2048

# 서버 CSR 생성
openssl req -new -key server.key -out server.csr -subj "/CN=server"

# 서버 인증서를 CA로 서명 (유효기간 100년, 36500일)
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 36500 -sha256
```

---

### 3. 클라이언트 인증서 생성

클라이언트용 개인키와 CSR을 생성한 후, CA로 서명합니다.

```sh
# 클라이언트 개인키 생성
openssl genrsa -out client.key 2048

# 클라이언트 CSR 생성
openssl req -new -key client.key -out client.csr -subj "/CN=unique_client"

# 클라이언트 인증서를 CA로 서명 (유효기간 100년, 36500일)
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 36500 -sha256
```

---

### 4. 요약

- **CA 인증서(ca.crt)와 개인키(ca.key)를 생성**하여 자체 서명한 CA를 구성합니다.
- **서버 인증서(server.crt)와 개인키(server.key)를 생성**하고, CA로 서명하여 서버의 신뢰성을 보장합니다.
- **클라이언트 인증서(client.crt)와 개인키(client.key)를 생성**하고, 동일한 CA로 서명하여 특정 클라이언트만 신뢰할 수 있게 합니다.

이제 서버에서는 CA 인증서를 이용해 클라이언트 인증서를 검증하고, 클라이언트는 서버 인증서를 검증할 수 있습니다. 이를 통해 양쪽 모두 상호 인증(mTLS)을 수행하게 됩니다.
