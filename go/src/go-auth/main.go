package main

import (
	"bytes"
	"encoding/base64"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	//"github.com/humamfauzi/go-registration/utils"
	//"github.com/humamfauzi/go-registration/exconn"
	"github.com/gorilla/mux"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	// pb "github.com/humamfauzi/go-registration/protobuf"
)

var arrayBody Bodies

const (
	letters               = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alphaNumeric          = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	JwtSignatureAlgorithm = jwa.HS256
	EncryptionSalt        = "falantagantengjagojava"
)

func main() {
	// Initialize DB connection
	//db = exconn.ConnectToMySQL()
	//errorMap = utils.InitError("./error.json")
	arrayBody = Bodies{}
	r := mux.NewRouter()
	r.HandleFunc("/auth/register", RegisterHandler).Methods(http.MethodPost)
	r.HandleFunc("/auth/getData", GetDataHandler).Methods(http.MethodPost)
	r.HandleFunc("/auth/verify", VerifyHandler).Methods(http.MethodPost)
	//r.HandleFunc("/auth/loginjwt", LoginJWTHandler).Methods(http.MethodGet)
	content, err := ioutil.ReadFile("./server_url.txt")
	if err != nil {
		log.Fatal(err)
	}
	addr := string(content)
	srv := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Print("STARTING SERVER at ", addr)
	log.Fatal(srv.ListenAndServe())

}

type Body struct {
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
	Password  string `json:"password,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

type ErrorReply struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Bodies []Body

type Token struct {
	Jwt string `json:"jwt"`
}

func Registry(payload []Body) {

}

func GenerateRandomString(length int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	stringBytes := make([]byte, length)
	var randomInteger int
	for i := 0; i < length; i++ {
		randomInteger = rand.Intn(len(alphaNumeric))
		stringBytes[i] = alphaNumeric[randomInteger]
	}
	return string(stringBytes)
}

func LookupArray(b Bodies, s string) bool {
	length := len(b)

	for i := 0; i < length; i++ {
		if b[i].Name == s {
			return true
		}
	}

	return false
}

func LookupData(b Bodies, phone string, password string) Body {
	length := len(b)

	for i := 0; i < length; i++ {
		fmt.Println("Not Found", b[i])
		if b[i].Phone == phone && b[i].Password == password {
			fmt.Println("Found", b[i])
			return b[i]
		}
	}

	return Body{}
}

func GenerateWebToken(name, phone, role string) ([]byte, error) {
	// log := loggerFactory.CreateLog().SetFunctionName("GenerateWebToken").SetStartTime()
	// defer log.SetFinishTime().WriteAndDeleteLog()

	tokenJwt := jwt.New()
	tokenJwt.Set(`name`, name)
	tokenJwt.Set(`phone`, phone)
	tokenJwt.Set(`role`, role)
	tokenJwt.Set(`timestamp`, time.Now().UTC())

	payload, err := jwt.Sign(tokenJwt, JwtSignatureAlgorithm, []byte(EncryptionSalt))
	if err != nil {
		return []byte{}, err
	}
	return payload, nil
}

func GetDataHandler(w http.ResponseWriter, r *http.Request) {

	var regBody Body

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errReply := ErrorReply{
			Code:    "ERR_UNREADBLE_PAYLOAD",
			Message: "Cannot parse incoming payload",
		}

		result, _ := json.Marshal(errReply)
		w.Write(result)
		return
	}
	fmt.Println(body)
	err = json.Unmarshal(body, &regBody)
	fmt.Println(regBody)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errReply := ErrorReply{
			Code:    "ERR_UNREADBLE_PAYLOAD",
			Message: "Cannot parse incoming payload",
		}
		result, _ := json.Marshal(errReply)
		w.Write(result)
		return
	}

	emptyBody := Body{}
	returnedData := LookupData(arrayBody, regBody.Phone, regBody.Password)
	fmt.Println("ReturnedData ", returnedData)
	var payload []byte
	if returnedData != emptyBody {
		payload, _ = GenerateWebToken(returnedData.Name, returnedData.Phone, returnedData.Role)
	}
	fmt.Println("Payload ", payload)
	sEnc := b64.StdEncoding.EncodeToString(payload)
	fmt.Println("sEnc ", sEnc)
	structToken := Token{Jwt: sEnc}
	fmt.Println("structToken ", structToken)
	output, _ := json.Marshal(structToken)
	fmt.Println("output ", output)
	w.Write(output)

}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	var regBody Body

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errReply := ErrorReply{
			Code:    "ERR_UNREADBLE_PAYLOAD",
			Message: "Cannot parse incoming payload",
		}

		result, _ := json.Marshal(errReply)
		w.Write(result)
		return
	}
	fmt.Println(body)
	err = json.Unmarshal(body, &regBody)
	fmt.Println(regBody)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errReply := ErrorReply{
			Code:    "ERR_UNREADBLE_PAYLOAD",
			Message: "Cannot parse incoming payload",
		}
		result, _ := json.Marshal(errReply)
		w.Write(result)
		return
	}
	regBody.Password = GenerateRandomString(4)

	if !LookupArray(arrayBody, regBody.Name) {
		arrayBody = append(arrayBody, regBody)
	}

	fmt.Println(arrayBody)
	output, _ := json.Marshal(regBody)
	w.Write(output)

}

////////////////////ENDPOINT KE 3/////////////////////////////////////
//////////////////////////////////////////////////////////////////////

func VerifyHandler(w http.ResponseWriter, r *http.Request) {

	body, _ := GetWebToken(r)
	output, _ := json.Marshal(body)
	w.Write(output)
}

func GetWebToken(r *http.Request) (Body, error) {
	var err error
	auth, ok := r.Header["Authorization"]
	if !ok {
		err = errors.New("ERR_CANNOT_PARSE_HEADER")
		return Body{}, err
	}

	splitAuth := strings.Split(auth[0], " ")
	if splitAuth[0] != "Bearer" {
		err = errors.New("ERR_WRONG_AUTHORIZATION")
		return Body{}, err
	}
	return VerifyToken(splitAuth[1])
}

func VerifyData(b Bodies, name, phone, role string) Body {
	length := len(b)

	for i := 0; i < length; i++ {
		fmt.Println("Not Found", b[i])
		if b[i].Phone == phone && b[i].Role == role && b[i].Name == name {
			fmt.Println("Found", b[i])
			return b[i]
		}
	}

	return Body{}
}

func VerifyToken(incomingToken string) (Body, error) {
	convertedToken, err := base64.StdEncoding.DecodeString(incomingToken)
	if err != nil {
		err = errors.New("ERR_WRONG_AUTHORIZATION")
		return Body{}, err
	}

	options := jwt.WithVerify(JwtSignatureAlgorithm, []byte(EncryptionSalt))
	token, err := jwt.Parse(bytes.NewReader(convertedToken), options)
	if err != nil {
		err = errors.New("ERR_WRONG_AUTHORIZATION")
		return Body{}, err
	}
	name, _ := token.Get("name")
	phone, _ := token.Get("phone")
	role, _ := token.Get("role")
	tstamp, _ := token.Get("timestamp")
	var outputBody Body
	body := VerifyData(arrayBody, name.(string), phone.(string), role.(string))
	emptyBody := Body{}
	if body != emptyBody {
		outputBody = Body{Name: name.(string), Phone: phone.(string), Role: role.(string), Timestamp: tstamp.(string)}
	}

	return outputBody, nil
}
