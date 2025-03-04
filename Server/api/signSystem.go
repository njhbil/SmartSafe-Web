package signSystem

import (
	"SmartSafe/database"
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
	"golang.org/x/crypto/argon2"
)

var (
	//Argon2 errors
	ErrInvalidHash         = errors.New("argon2: invalid hash")
	ErrIncompatibleVersion = errors.New("argon2: incompatible version of argon2")

	//JWT Secret
	jwtSecret = os.Getenv("JWT_SECRET")
)

type accounts struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Payload struct {
	RememberMe bool   `json:"rememberMe"`
	ResetToken string `json:"resetToken"`
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

func hashPassword(password string, p *params) (string, error) {
	salt, err := generateRandomBytes(p.SaltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)
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

func comparePasswordAndHash(password, encodedHash string) (match bool, err error) {
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
		json.NewEncoder(w).Encode(map[string]string{"message": message, "RefreshToken": responeToken})
	}
}

func Accounts(w http.ResponseWriter, r *http.Request) {
	conn := db.Connect()
	defer conn.Close()

	var accounts []accounts
	err := conn.Model(&accounts).Select()
	if err != nil {
		log.Println("Database Query Error:", err)
		http.Error(w, "Database Query Error", http.StatusInternalServerError)
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
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
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
		http.Error(w, "Database Query Error", http.StatusInternalServerError)
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
		http.Error(w, "Database Query Error", http.StatusInternalServerError)
		return
	}
	if exists {
		log.Println("Username already exists")
		sendJsonResponse(w, http.StatusConflict, "Username already exists")
		return
	}

	hash, err := hashPassword(req.Password, p)
	if err != nil {
		log.Println("Password Hash Error:", err)
		http.Error(w, "Password Hash Error", http.StatusInternalServerError)
		return
	}

	req.Password = hash

	_, err = conn.Model(&req).Insert()
	if err != nil {
		log.Println("Database Insert Error:", err)
		http.Error(w, "Database Insert Error", http.StatusInternalServerError)
		sendJsonResponse(w, http.StatusInternalServerError, "Database Insert Error")
		return
	}

	log.Println("User registered successfully")
	sendJsonResponse(w, http.StatusOK, "User registered successfully")
}

func Login(w http.ResponseWriter, r *http.Request) {

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
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var payload Payload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Println("JSON Decode Error for Payload:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		log.Println("email or password is empty")
		http.Error(w, "email or password is empty", http.StatusBadRequest)
		return
	}

	exists, err = conn.Model(new(accounts)).Where("email = ?", req.Email).Exists()
	if err != nil {
		log.Println("Database Query Error:", err)
		http.Error(w, "Database Query Error", http.StatusInternalServerError)
		return
	}

	if !exists {
		log.Println("Email does not exist")
		http.Error(w, "Email does not exist", http.StatusUnauthorized)
		return
	}

	var account accounts
	err = conn.Model(&account).Where("email = ?", req.Email).Select()
	if err != nil {
		log.Println("Database Query Error:", err)
		http.Error(w, "Database Query Error", http.StatusInternalServerError)
		return
	}

	verify, err := comparePasswordAndHash(req.Password, account.Password)
	if err != nil {
		log.Println("Password Hash Error:", err)
		http.Error(w, "Password Hash Error", http.StatusInternalServerError)
		return
	}

	if verify {
		if payload.RememberMe {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"email": account.Email,
				"exp":   time.Now().Add(time.Minute * 3).Unix(),
			})
			tokenString, err := token.SignedString([]byte(jwtSecret))
			if err != nil {
				log.Println("Token Signing Error:", err)
				http.Error(w, "Token Signing Error", http.StatusInternalServerError)
				return
			}
			log.Println("Login successful (30 Day)")
			sendJsonResponse(w, http.StatusOK, "Login successful", tokenString)
			return
		} else {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"email": account.Email,
				"exp":   time.Now().Add(time.Minute * 3).Unix(),
			})
			tokenString, err := token.SignedString([]byte(jwtSecret))
			if err != nil {
				log.Println("Token Signing Error:", err)
				http.Error(w, "Token Signing Error", http.StatusInternalServerError)
				return
			}
			log.Println("Login successful (1 Day)")
			sendJsonResponse(w, http.StatusOK, "Login successful", tokenString)
			return
		}

	} else {
		log.Println("Password is invalid")
		http.Error(w, "Password is invalid", http.StatusUnauthorized)
		return
	}
}

func NewWithClaims(claims jwt.Claims, method jwt.SigningMethod) *Token {
	return &Token{
		Method: method,
		Claims: claims,
	}
}
