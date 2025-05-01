package signSystem

import (
	"SmartSafe/database"
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/argon2"
	gomail "gopkg.in/mail.v2"
)

var (
	//Argon2 errors
	ErrInvalidHash         = errors.New("argon2: invalid hash")
	ErrIncompatibleVersion = errors.New("argon2: incompatible version of argon2")

	//JWT Secret
	jwtSecret = os.Getenv("JWT_SECRET")
	clientURL = os.Getenv("CLIENT_URL")
)

var rdb = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS_URL"),
	DB:       0, // use default DB
	Password: os.Getenv("REDIS_PASSWORD"),
})

type accounts struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type verifytokens struct {
	Token string `json:"token"`
}

type Payload struct {
	RememberMe    bool   `json:"rememberMe"`
	ResetToken    string `json:"resetToken"`
	ResetPWDtoken string `json:"resetPWDtoken"`
	RefreshToken  string `json:"refreshToken"`
	LoginToken    string `json:"loginToken"`
}

type OTPPayload struct {
	Email string `json:"email"`
	OTP   string `json:"oneTimeCode"`
}

type params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

type Token struct {
	Raw       string
	Method    jwt.SigningMethod
	Header    map[string]interface{}
	Claims    jwt.Claims
	Signature []byte
	Valid     bool
}

func hashData(Data string, p *params) (string, error) {
	salt, err := generateRandomBytes(p.SaltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(Data), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.Memory, p.Iterations, p.Parallelism,
		base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(hash))

	return encodedHash, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func compareDataAndHash(password, encodedHash string) (match bool, err error) {
	p, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func decodeHash(encodedHash string) (p *params, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	p = &params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Iterations, &p.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.SaltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.KeyLength = uint32(len(hash))

	return p, salt, hash, nil
}

func sendJsonResponse(w http.ResponseWriter, status int, message string, ResetToken ...string) {
	responeToken := ""
	if len(ResetToken) > 0 {
		responeToken = ResetToken[0]
	}

	if responeToken == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]string{"message": message})
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]string{"message": message, "refreshToken": responeToken})
	}
}

func sendJsonError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"message": message,
	})
}

func Accounts(w http.ResponseWriter, r *http.Request) {
	conn := db.Connect()
	defer conn.Close()

	var accounts []accounts
	err := conn.Model(&accounts).Select()
	if err != nil {
		log.Println("Database Query Error:", err)
		sendJsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}

func Register(w http.ResponseWriter, r *http.Request) {
	p := &params{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}

	var exists bool

	conn := db.Connect()
	defer conn.Close()

	var req accounts
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("Json Decode Error:", err)
		sendJsonError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Username == "" || req.Password == "" {
		log.Println("username, email or password is empty")
		sendJsonResponse(w, http.StatusBadRequest, "username, email or password is empty")
		return
	}

	exists, err = conn.Model(new(accounts)).Where("email = ?", req.Email).Exists()
	if err != nil {
		log.Println("Database Query Error:", err)
		sendJsonError(w, "Database error", http.StatusInternalServerError)
		return
	}
	if exists {
		log.Println("Email already exists")
		sendJsonResponse(w, http.StatusConflict, "Email already exists")
		return
	}

	exists, err = conn.Model(new(accounts)).Where("username = ?", req.Username).Exists()
	if err != nil {
		log.Println("Database Query Error:", err)
		sendJsonError(w, "Database error", http.StatusInternalServerError)
		return
	}
	if exists {
		log.Println("Username already exists")
		sendJsonResponse(w, http.StatusConflict, "Username already exists")
		return
	}

	hash, err := hashData(req.Password, p)
	if err != nil {
		log.Println("Password Hash Error:", err)
		sendJsonError(w, "Password Hash Error", http.StatusInternalServerError)
		return
	}

	req.Password = hash

	_, err = conn.Model(&req).Insert()
	if err != nil {
		log.Println("Database Insert Error:", err)
		sendJsonError(w, "Database Insert Error", http.StatusInternalServerError)
		return
	}

	log.Println("User registered successfully")
	sendJsonResponse(w, http.StatusOK, "User registered successfully")
}

func SendEmailOTP(w http.ResponseWriter, r *http.Request) {
	var exists bool

	conn := db.Connect()
	defer conn.Close()

	var req OTPPayload
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("JSON Decode Error:", err)
		sendJsonError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	email := req.Email
	if email == "" {
		log.Println("Email is empty")
		sendJsonError(w, "Email is required", http.StatusBadRequest)
		return
	}

	exists, err = conn.Model(new(accounts)).Where("email = ?", email).Exists()
	if err != nil {
		log.Println("Database Query Error:", err)
		sendJsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	if exists {
		log.Println("Email already exists")
		sendJsonError(w, "Email already exists", http.StatusConflict)
		return
	}

	log.Println("Received OTP request for:", email)

	otp, err := GenerateRandomOTP()
	if err != nil {
		log.Println("OTP Generation Error:", err)
		sendJsonError(w, "OTP generation error", http.StatusInternalServerError)
		return
	}

	err = StoreOTP(email, otp)
	if err != nil {
		log.Println("OTP Storage Error:", err)
		sendJsonError(w, "OTP storage error", http.StatusInternalServerError)
		return
	}
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("WEBMAIL"))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Your OTP Code")
	m.SetBody("text/html", `
		<html>
			<body>
				<div style="text-align: center; width: 50%; margin: 0 auto;">
					<h1>OTP</h1>
					<p>We received a request to send you an OTP code</p>
					<p>Use the OTP code below to verify your email</p>
					<p>OTP Code: <strong>`+otp+`</strong></p>
					</div>
				<div style="text-align: start; margin-top: 20px;">
					<td style="font-size: 14px; line-height: 170%; font-weight: 400; color: #000000; letter-spacing: 0.01em;">
                        Best regards, <br><strong>SmartSafe Team</strong>
                    </td>
				</div>
			</body>
		</html>
		`)

	d := gomail.NewDialer("live.smtp.mailtrap.io", 587, "api", os.Getenv("MAILTRAP_TOKEN"))
	if err := d.DialAndSend(m); err != nil {
		log.Println("Email Sending Error:", err)
		sendJsonError(w, "Email sending error", http.StatusInternalServerError)
		return
	}

	log.Println("OTP sent successfully")
	sendJsonResponse(w, http.StatusOK, "OTP sent successfully")
}

func VerifyEmailOTP(w http.ResponseWriter, r *http.Request) {
	var req OTPPayload
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("Json Decode Error:", err)
		sendJsonError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	email := req.Email
	otp := req.OTP

	if email == "" || otp == "" {
		log.Println("email or otp is empty")
		sendJsonError(w, "email or otp is empty", http.StatusBadRequest)
		return
	}

	storedOTP, err := RetrieveOTP(email)
	if err != nil {
		log.Println("OTP Retrieval Error:", err)
		sendJsonError(w, "OTP retrieval error", http.StatusInternalServerError)
		return
	}

	if storedOTP != otp {
		log.Println("Invalid OTP")
		sendJsonError(w, "Invalid OTP", http.StatusUnauthorized)
		return
	}

	err = DeleteOTP(email)
	if err != nil {
		log.Println("OTP Deletion Error:", err)
		sendJsonError(w, "OTP deletion error", http.StatusInternalServerError)
		return
	}

	log.Println("OTP verified successfully")
	sendJsonResponse(w, http.StatusOK, "OTP verified successfully")

}

func GenerateRandomOTP() (string, error) {
	const OTPLength = 6

	otpBytes := make([]byte, OTPLength)
	_, err := rand.Read(otpBytes)
	if err != nil {
		return "", err
	}

	otp := ""

	for _, b := range otpBytes {
		otp += fmt.Sprintf("%d", b%10)
	}

	return otp, nil
}

func StoreOTP(email, otp string) error {
	ctx := context.Background()
	expiresAt := time.Now().Add(5 * time.Minute)
	return rdb.Set(ctx, email, otp, expiresAt.Sub(time.Now())).Err()
}

func RetrieveOTP(email string) (string, error) {
	ctx := context.Background()
	return rdb.Get(ctx, email).Result()
}

func DeleteOTP(email string) error {
	ctx := context.Background()
	return rdb.Del(ctx, email).Err()
}

func Login(w http.ResponseWriter, r *http.Request) {
	p := &params{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}

	var exists bool

	conn := db.Connect()
	defer conn.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Read Body Error:", err)
		sendJsonError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var req accounts
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Println("JSON Decode Error for Accounts:", err)
		sendJsonError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var payload Payload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Println("JSON Decode Error for Payload:", err)
		sendJsonError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var verifyToken verifytokens
	err = json.Unmarshal(body, &verifyToken)
	if err != nil {
		log.Println("JSON Decode Error for Token:", err)
		sendJsonError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		log.Println("email or password is empty")
		sendJsonError(w, "email or password is empty", http.StatusBadRequest)
		return
	}

	exists, err = conn.Model(new(accounts)).Where("email = ?", req.Email).Exists()
	if err != nil {
		log.Println("Database Query Error:", err)
		sendJsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	if !exists {
		log.Println("Email does not exist")
		sendJsonError(w, "Email does not exist", http.StatusUnauthorized)
		return
	}

	var account accounts
	err = conn.Model(&account).Where("email = ?", req.Email).Select()
	if err != nil {
		log.Println("Database Query Error:", err)
		sendJsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	verify, err := compareDataAndHash(req.Password, account.Password)
	if err != nil {
		log.Println("Password Hash Error:", err)
		sendJsonError(w, "Password Hash Error", http.StatusInternalServerError)
		return
	}

	if verify {
		if payload.RememberMe {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"email": account.Email,
				"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
			})
			tokenString, err := token.SignedString([]byte(jwtSecret))
			if err != nil {
				log.Println("Token Signing Error:", err)
				sendJsonError(w, "Token signing error", http.StatusInternalServerError)
				return
			}

			tokenHash, err := hashData(tokenString, p)
			if err != nil {
				log.Println("Token Hash Error:", err)
				sendJsonError(w, "Token Hash Error", http.StatusInternalServerError)
				return
			}

			verifyToken.Token = tokenHash

			_, err = conn.Model(&verifyToken).Insert()
			if err != nil {
				log.Println("Database Insert Error:", err)
				sendJsonError(w, "Database Insert Error", http.StatusInternalServerError)
				return
			}

			log.Println("Login successful (30 Day)")
			sendJsonResponse(w, http.StatusOK, "Login successful", tokenString)
			return
		} else {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"email": account.Email,
				"exp":   time.Now().Add(time.Hour * 24).Unix(),
			})
			tokenString, err := token.SignedString([]byte(jwtSecret))
			if err != nil {
				log.Println("Token Signing Error:", err)
				sendJsonError(w, "Token signing error", http.StatusInternalServerError)
				return
			}

			tokenHash, err := hashData(tokenString, p)
			if err != nil {
				log.Println("Token Hash Error:", err)
				sendJsonError(w, "Token Hash Error", http.StatusInternalServerError)
				return
			}

			verifyToken.Token = tokenHash

			_, err = conn.Model(&verifyToken).Insert()
			if err != nil {
				log.Println("Database Insert Error:", err)
				sendJsonError(w, "Database Insert Error", http.StatusInternalServerError)
				return
			}

			log.Println("Login successful (1 Day)")
			sendJsonResponse(w, http.StatusOK, "Login successful", tokenString)
			return
		}

	} else {
		log.Println("Password is invalid")
		sendJsonError(w, "Password is invalid", http.StatusUnauthorized)
		return
	}
}

func ForgetPassword(w http.ResponseWriter, r *http.Request) {
	var exists bool

	conn := db.Connect()
	defer conn.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Read Body Error:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var req accounts
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Println("JSON Decode Error for Accounts:", err)
		sendJsonError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		log.Println("email is empty")
		sendJsonError(w, "email is empty", http.StatusBadRequest)
		return
	}

	exists, err = conn.Model(new(accounts)).Where("email = ?", req.Email).Exists()
	if err != nil {
		log.Println("Database Query Error:", err)
		sendJsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	if !exists {
		log.Println("Email does not exist")
		sendJsonError(w, "Email does not exist", http.StatusUnauthorized)
		return
	}

	var account accounts
	err = conn.Model(&account).Where("email = ?", req.Email).Select()
	if err != nil {
		log.Println("Database Query Error:", err)
		sendJsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": account.Email,
		"exp":   time.Now().Add(time.Minute * 10).Unix(),
	})
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Println("Token Signing Error:", err)
		sendJsonError(w, "Token signing error", http.StatusInternalServerError)
		return
	}

	clientURL := os.Getenv("CLIENT_URL")

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("WEBMAIL"))
	m.SetHeader("To", account.Email)
	m.SetHeader("Subject", "Reset Password")
	m.SetBody("text/html", `
		<html>
			<body>
				<div style="text-align: center; width: 50%; margin: 0 auto;">
					<h1>Reset Password</h1>
					<p>We received a request to reset your password from `+req.Email+`</p>
					<p>Click the link below to reset your password</p>
		<a class="button" title="Reset Password" href="`+clientURL+`/resetpassword?token=`+tokenString+`" style="width: 100%; background: #22D172; text-decoration: none; display: inline-block; padding: 10px 0; color: #fff; font-size: 14px; line-height: 21px; text-align: center; font-weight: bold; border-radius: 7px;">Reset Password</a>
				</div>
				<div style="text-align: start; margin-top: 20px;">
					<td style="font-size: 14px; line-height: 170%; font-weight: 400; color: #000000; letter-spacing: 0.01em;">
                        Best regards, <br><strong>SmartSafe Team</strong>
                    </td>
				</div>
			</body>
		</html>
		`)

	d := gomail.NewDialer("live.smtp.mailtrap.io", 587, "api", os.Getenv("MAILTRAP_TOKEN"))
	if err := d.DialAndSend(m); err != nil {
		log.Println("Email Sending Error:", err)
		sendJsonError(w, "Email sending error", http.StatusInternalServerError)
		return
	}

	sendJsonResponse(w, http.StatusOK, "Email sent successfully")
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	p := &params{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}

	var exists bool

	conn := db.Connect()
	defer conn.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Read Body Error:", err)
		sendJsonError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var req accounts
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Println("JSON Decode Error for Accounts:", err)
		sendJsonError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var payload Payload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Println("JSON Decode Error for Payload:", err)
		sendJsonError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		log.Println("password is empty")
		http.Error(w, "password is empty", http.StatusBadRequest)
		return
	}

	if payload.ResetPWDtoken == "" {
		log.Println("Error: Token is empty")
		sendJsonError(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	cleanToken := strings.TrimSpace(payload.ResetPWDtoken)

	token, err := jwt.ParseWithClaims(cleanToken, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok || !token.Valid {
		log.Println("Token is invalid or claims are malformed")
		sendJsonError(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	email, ok := (*claims)["email"].(string)
	if !ok || email == "" {
		log.Println("Email claim missing or invalid")
		sendJsonError(w, "Invalid token data", http.StatusUnauthorized)
		return
	}

	exists, err = conn.Model(new(accounts)).Where("email = ?", email).Exists()
	if err != nil {
		log.Println("Database Query Error:", err)
		sendJsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	if !exists {
		log.Println("Email does not exist")
		sendJsonError(w, "Email does not exist", http.StatusUnauthorized)
		return
	}

	hash, err := hashData(req.Password, p)
	if err != nil {
		log.Println("Password Hash Error:", err)
		sendJsonError(w, "Password Hash Error", http.StatusInternalServerError)
		return
	}

	_, err = conn.Model(new(accounts)).Set("password = ?", hash).Where("email = ?", email).Update()
	if err != nil {
		log.Println("Database Update Error:", err)
		sendJsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	log.Println("Password reset successfully")
	sendJsonResponse(w, http.StatusOK, "Password reset successfully")
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	var exists bool

	conn := db.Connect()
	defer conn.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Read Body Error:", err)
		sendJsonError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var payload Payload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Println("JSON Decode Error for Payload:", err)
		sendJsonError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if payload.RefreshToken == "" {
		log.Println("Error: Token is empty")
		sendJsonError(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	cleanToken := strings.TrimSpace(payload.RefreshToken)

	var verifyTokens []verifytokens
	err = conn.Model(new(verifytokens)).Select(&verifyTokens)
	if err != nil {
		log.Println("Database query error:", err)
		sendJsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	//TO-DO: change the logic for checking the token in the database to be more efficient and secure
	var isValid bool
	for _, token := range verifyTokens {
		isValid, err = compareDataAndHash(cleanToken, token.Token)
		if err != nil {
			log.Println("Failed to compare token hash:", err)
			sendJsonError(w, "Token comparison error", http.StatusInternalServerError)
			return
		}

		if isValid {
			break
		}
	}

	if !isValid {
		log.Println("Invalid token: no match found in the database")
		sendJsonError(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	token, err := jwt.ParseWithClaims(cleanToken, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		log.Println("Token parsing error:", err)
		sendJsonError(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok || !token.Valid {
		log.Println("Token is invalid or claims are malformed")
		sendJsonError(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	email, ok := (*claims)["email"].(string)
	if !ok || email == "" {
		log.Println("Email claim missing or invalid")
		sendJsonError(w, "Invalid token data", http.StatusUnauthorized)
		return
	}

	exists, err = conn.Model(new(accounts)).Where("email = ?", email).Exists()
	if err != nil {
		log.Println("Database Query Error Accounts:", err)
		sendJsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	if !exists {
		log.Println("Email does not exist")
		sendJsonError(w, "Email does not exist", http.StatusUnauthorized)
		return
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := newToken.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Println("Token Signing Error:", err)
		sendJsonError(w, "Token signing error", http.StatusInternalServerError)
		return
	}

	log.Println("Token refreshed successfully")
	sendJsonResponse(w, http.StatusOK, "Token refreshed successfully", tokenString)
}

func VerifyToken(w http.ResponseWriter, r *http.Request) {
	var exists bool

	conn := db.Connect()
	defer conn.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Read Body Error:", err)
		sendJsonError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var payload struct {
		LoginToken   string `json:"loginToken"`
		RefreshToken string `json:"refreshToken"`
	}

	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Println("JSON Decode Error for Payload:", err)
		sendJsonError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if payload.LoginToken == "" {
		log.Println("Error: Login token is empty")
		sendJsonError(w, "Invalid login token", http.StatusUnauthorized)
		return
	}

	cleanLoginToken := strings.TrimSpace(payload.LoginToken)

	token, err := jwt.ParseWithClaims(cleanLoginToken, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err == nil && token.Valid {
		claims, ok := token.Claims.(*jwt.MapClaims)
		if !ok {
			log.Println("Claims are malformed")
			sendJsonError(w, "Invalid token data", http.StatusUnauthorized)
			return
		}

		email, ok := (*claims)["email"].(string)
		if !ok || email == "" {
			log.Println("Email claim missing or invalid")
			sendJsonError(w, "Invalid token data", http.StatusUnauthorized)
			return
		}

		exists, err = conn.Model(new(accounts)).Where("email = ?", email).Exists()
		if err != nil {
			log.Println("Database Query Error Accounts:", err)
			sendJsonError(w, "Database error", http.StatusInternalServerError)
			return
		}

		if !exists {
			log.Println("Email does not exist")
			sendJsonError(w, "Email does not exist", http.StatusUnauthorized)
			return
		}

		log.Println("Login token is valid")
		sendJsonResponse(w, http.StatusOK, "Token is valid")
		return
	}

	if payload.RefreshToken == "" {
		log.Println("Login token is invalid and no refresh token provided")
		sendJsonError(w, "Invalid token and no refresh token", http.StatusUnauthorized)
		return
	}

	log.Println("Login token is invalid. Attempting to refresh with refresh token.")

	cleanRefreshToken := strings.TrimSpace(payload.RefreshToken)

	var verifyTokens []verifytokens
	err = conn.Model(new(verifytokens)).Select(&verifyTokens)
	if err != nil {
		log.Println("Database query error:", err)
		sendJsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	//TO-DO: change the logic for checking the token in the database to be more efficient and secure
	var isValid bool
	for _, token := range verifyTokens {
		isValid, err = compareDataAndHash(cleanRefreshToken, token.Token)
		if err != nil {
			log.Println("Failed to compare token hash:", err)
			sendJsonError(w, "Token comparison error", http.StatusInternalServerError)
			return
		}

		if isValid {
			break
		}
	}

	if !isValid {
		log.Println("Invalid refresh token: no match found in the database")
		sendJsonError(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	refreshToken, err := jwt.ParseWithClaims(cleanRefreshToken, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		log.Println("Refresh token parsing error:", err)
		sendJsonError(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	refreshClaims, ok := refreshToken.Claims.(*jwt.MapClaims)
	if !ok || !refreshToken.Valid {
		log.Println("Refresh token is invalid or claims are malformed")
		sendJsonError(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	email, ok := (*refreshClaims)["email"].(string)
	if !ok || email == "" {
		log.Println("Email claim missing or invalid in refresh token")
		sendJsonError(w, "Invalid refresh token data", http.StatusUnauthorized)
		return
	}

	exists, err = conn.Model(new(accounts)).Where("email = ?", email).Exists()
	if err != nil {
		log.Println("Database Query Error Accounts:", err)
		sendJsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	if !exists {
		log.Println("Email does not exist")
		sendJsonError(w, "Email does not exist", http.StatusUnauthorized)
		return
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := newToken.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Println("Token Signing Error:", err)
		sendJsonError(w, "Token signing error", http.StatusInternalServerError)
		return
	}

	log.Println("Token refreshed successfully")
	sendJsonResponse(w, http.StatusOK, "Token refreshed successfully", tokenString)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	var exists bool

	conn := db.Connect()
	defer conn.Close()

	var payload Payload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Println("JSON Decode Error:", err)
		sendJsonError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if payload.RefreshToken == "" {
		log.Println("Error: Token is empty")
		sendJsonError(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	cleanToken := strings.TrimSpace(payload.RefreshToken)
	var verifyTokens []verifytokens
	err = conn.Model(new(verifytokens)).Select(&verifyTokens)
	if err != nil {
		log.Println("Database query error:", err)
		sendJsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	for _, token := range verifyTokens {
		exists, err = compareDataAndHash(cleanToken, token.Token)
		if err != nil {
			log.Println("Failed to compare token hash:", err)
			sendJsonError(w, "Token comparison error", http.StatusInternalServerError)
			return
		}

		if exists {
			break
		}
	}

	if !exists {
		log.Println("Invalid token: no match found in the database")
		sendJsonError(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	for _, token := range verifyTokens {
		if exists, err = conn.Model(new(verifytokens)).Where("token = ?", token.Token).Exists(); err != nil {
			log.Println("Database Query Error:", err)
			sendJsonError(w, "Database error", http.StatusInternalServerError)
			return
		}
		if exists {
			_, err = conn.Model(new(verifytokens)).Where("token = ?", token.Token).Delete()
			if err != nil {
				log.Println("Database Delete Error:", err)
				sendJsonError(w, "Database error", http.StatusInternalServerError)
				return
			}
		}
	}
	if err != nil {
		log.Println("Database Delete Error:", err)
		sendJsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	log.Println("Logout successful")
	sendJsonResponse(w, http.StatusOK, "Logout successful")
}

func NewWithClaims(claims jwt.Claims, method jwt.SigningMethod) *Token {
	return &Token{
		Method: method,
		Claims: claims,
	}
}
