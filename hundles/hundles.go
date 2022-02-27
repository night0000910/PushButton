package hundles

import (
	crand "crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("postgres", "user=postgres dbname=pushbutton password=postgres sslmode=disable")

	if err != nil {
		panic(err)
	}
}

type User struct {
	Id            int
	Username      string
	Password      string
	CsrfToken     string
	Money         int
	OakFruits     int
	ThunderFruits int
}

func (user *User) Update() (err error) {
	hashedCsrfToken := Hash(user.CsrfToken)
	_, err = Db.Exec("update users set username=$1, password=$2, csrf_token=$3, money=$4, oak_fruits=$5, thunder_fruits=$6 where id = $7", user.Username, user.Password, hashedCsrfToken, user.Money, user.OakFruits, user.ThunderFruits, user.Id)
	return
}

func (user *User) CreateSession(sessionId string) (err error) {
	hashedSessionId := Hash(sessionId)
	_, err = Db.Exec("insert into session (user_id, session_id) values ($1, $2)", user.Id, hashedSessionId)
	return
}

func (user *User) CreateCsrfToken() (err error) {
	user.CsrfToken = MakeRandomString()
	err = user.Update()
	return
}

func (user *User) VerifyCsrfToken() (validity bool, err error) {
	hashedCsrfToken := Hash(user.CsrfToken)
	var numberOfRecords int
	err = Db.QueryRow("select count(*) from users where id = $1 and csrf_token = $2", user.Id, hashedCsrfToken).Scan(&numberOfRecords)
	if err != nil {
		return
	}

	if numberOfRecords == 1 {
		validity = true
	} else if numberOfRecords == 0 {
		validity = false
	}
	return
}

type Information struct {
	Message              string
	Money                int
	OakFruits            int
	ThunderFruits        int
	PriceOfOakFruits     int
	PriceOfThunderFruits int
	Profit               int
}

// ----------------------------- 関数 ------------------------------

func Hash(str string) (s string) {
	hashedStr := sha256.Sum256([]byte(str))
	s = base64.URLEncoding.EncodeToString(hashedStr[:])
	return
}

func MakeRandomString() (randomString string) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, 64)
	crand.Read(b)

	for _, byteValue := range b {
		value := int(byteValue)
		index := value % len(letters)
		randomString += string(letters[index])
	}

	return
}

func ExtractSessionIdFromCookie(cookie *http.Cookie) (sessionId string) {
	cookieString := cookie.String()
	sessionId = strings.Split(cookieString, "=")[1]
	return
}

func CreateUser(username string, password string) (user User, err error) {
	user = User{}
	hashedPassword := Hash(password)
	err = Db.QueryRow("insert into users (username, password) values ($1, $2) returning id", username, hashedPassword).Scan(&user.Id)

	if err != nil {
		return
	}

	err = Db.QueryRow("select id, username, password, csrf_token, money, oak_fruits, thunder_fruits from users where id = $1", user.Id).Scan(&user.Id, &user.Username, &user.Password, &user.CsrfToken, &user.Money, &user.OakFruits, &user.ThunderFruits)
	return
}

func Authenticate(username string, password string) (user User, validity bool, err error) {
	var numberOfRecords int
	hashedPassword := Hash(password)
	err = Db.QueryRow("select count(*) from users where username = $1 and password = $2", username, hashedPassword).Scan(&numberOfRecords)

	if err != nil {
		return
	}

	if numberOfRecords == 1 {
		validity = true
		err = Db.QueryRow("select id, username, password, csrf_token, money, oak_fruits, thunder_fruits from users where username = $1", username).Scan(&user.Id, &user.Username, &user.Password, &user.CsrfToken, &user.Money, &user.OakFruits, &user.ThunderFruits)

	} else if numberOfRecords == 0 {
		validity = false
	}

	return
}

func CheckUserIsAuthenticated(r *http.Request) (user User, validity bool, err error) {
	sessionIdCookie, err := r.Cookie("sessionId")

	if err != nil {
		return
	}

	sessionId := ExtractSessionIdFromCookie(sessionIdCookie)
	user, validity, err = GetUser(sessionId)

	return
}

func GetUser(sessionId string) (user User, validity bool, err error) {
	user = User{}
	var numberOfRecords int
	hashedSessionId := Hash(sessionId)
	err = Db.QueryRow("select count(*) from session where session_id = $1", hashedSessionId).Scan(&numberOfRecords)

	if err != nil {
		return
	}

	if numberOfRecords >= 1 {
		validity = true
		var userId int
		err = Db.QueryRow("select user_id from session where session_id = $1", hashedSessionId).Scan(&userId)

		if err != nil {
			return
		}

		err = Db.QueryRow("select id, username, password, money, oak_fruits, thunder_fruits from users where id = $1", userId).Scan(&user.Id, &user.Username, &user.Password, &user.Money, &user.OakFruits, &user.ThunderFruits)

	} else if numberOfRecords == 0 {
		validity = false
	}

	return
}

func DeleteSession(sessionId string) (err error) {
	hashedSessionId := Hash(sessionId)
	_, err = Db.Exec("delete from session where session_id = $1", hashedSessionId)
	return
}

func ReturnTrueWithCertainProbability() (result bool) {
	seed := time.Now().UnixNano()
	src := rand.NewSource(seed)
	rnd := rand.New(src)
	number := rnd.Intn(150)
	fmt.Println(number)

	if number == 0 {
		result = true
	} else {
		result = false
	}

	return
}

func MoveTo(page string, w http.ResponseWriter) {
	w.Header().Set("Location", "http://13.113.131.106/"+page)
	w.WriteHeader(302)
}

func WriteInformationAsJson(w http.ResponseWriter, information Information) (err error) {
	var jsonData string
	jsonData, err = ChangeInformationToJson(information)

	WriteJson(w, jsonData)
	return
}

func ChangeInformationToJson(information Information) (jsonString string, err error) {
	var jsonByte []byte
	jsonByte, err = json.Marshal(information)

	if err != nil {
		return
	}

	jsonString = string(jsonByte)
	return
}

func WriteJson(w http.ResponseWriter, jsonData string) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, jsonData)
}

// ---------------------- ハンドラ関数 ----------------------------

// ハンドラを生成する関数
func ReturnFileHandler() (fileHandler http.Handler) {
	dir := http.Dir(".")
	fileHandler = http.FileServer(dir)
	fileHandler = http.StripPrefix("/file/", fileHandler)
	return
}

func ReturnUsersInformation(w http.ResponseWriter, r *http.Request) {
	user, validity, err := CheckUserIsAuthenticated(r)

	if (err != nil) || (!validity) {
		return
	}

	priceOfOakFruits := int(100.0 * (math.Pow(2, (float64(user.OakFruits)/10.0)+1.0)))
	priceOfThunderFruits := user.Money

	information := Information{
		Money:                user.Money,
		OakFruits:            user.OakFruits,
		ThunderFruits:        user.ThunderFruits,
		PriceOfOakFruits:     priceOfOakFruits,
		PriceOfThunderFruits: priceOfThunderFruits,
	}

	err = WriteInformationAsJson(w, information)

	if err != nil {
		panic(err)
	}
}

func MoveToHomepageOrPushPage(w http.ResponseWriter, r *http.Request) {
	_, validity, err := CheckUserIsAuthenticated(r)

	if (err == nil) && validity {
		MoveTo("push", w)

	} else {
		MoveTo("homepage", w)
	}
}

func DisplayHomepage(w http.ResponseWriter, r *http.Request) {
	_, validity, err := CheckUserIsAuthenticated(r)
	if (err == nil) && validity {
		MoveTo("push", w)
		return
	}

	var t *template.Template
	t, err = template.ParseFiles("templates/homepage.html")
	if err != nil {
		panic(err)
	}

	t.Execute(w, nil)
}

func DisplayExplanation(w http.ResponseWriter, r *http.Request) {
	_, validity, err := CheckUserIsAuthenticated(r)

	var authentication bool

	if (err == nil) && validity {
		authentication = true
	} else {
		authentication = false
	}

	var t *template.Template
	t, err = template.ParseFiles("templates/explanation.html")
	if err != nil {
		panic(err)
	}

	t.Execute(w, authentication)
}

func DisplaySignupPage(w http.ResponseWriter, r *http.Request) {
	_, validity, err := CheckUserIsAuthenticated(r)
	if (err == nil) && validity {
		MoveTo("push", w)
		return
	}

	var t *template.Template
	t, err = template.ParseFiles("templates/signup.html")
	if err != nil {
		panic(err)
	}

	t.Execute(w, nil)
}

func Signup(w http.ResponseWriter, r *http.Request) {
	_, validity, err := CheckUserIsAuthenticated(r)

	if (err == nil) && validity {
		result := Information{
			Message: "authenticated",
		}

		err = WriteInformationAsJson(w, result)

		if err != nil {
			panic(err)
		}
		return
	}

	r.ParseForm()
	username := r.PostForm["username"][0]
	password := r.PostForm["password"][0]
	var user User
	user, err = CreateUser(username, password)

	if err != nil {
		result := Information{
			Message: "failed",
		}

		err = WriteInformationAsJson(w, result)
		if err != nil {
			panic(err)
		}
		return
	}

	err = user.CreateCsrfToken()
	if err != nil {
		panic(err)
	}

	result := Information{
		Message: "success",
	}

	err = WriteInformationAsJson(w, result)
	if err != nil {
		panic(err)
	}
}

func SucceedInSignup(w http.ResponseWriter, r *http.Request) {
	_, validity, err := CheckUserIsAuthenticated(r)

	if (err == nil) && validity {
		MoveTo("push", w)
		return
	}

	var t *template.Template
	t, err = template.ParseFiles("templates/succeedInSignup.html")

	if err != nil {
		panic(err)
	}

	t.Execute(w, nil)
}

func DisplayLoginPage(w http.ResponseWriter, r *http.Request) {
	_, validity, err := CheckUserIsAuthenticated(r)

	if (err == nil) && validity {
		MoveTo("push", w)
		return
	}

	var t *template.Template
	t, err = template.ParseFiles("templates/login.html")

	if err != nil {
		panic(err)
	}

	t.Execute(w, nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	_, validity, err := CheckUserIsAuthenticated(r)

	if (err == nil) && validity {
		result := Information{
			Message: "authenticated",
		}

		err = WriteInformationAsJson(w, result)

		if err != nil {
			panic(err)
		}
		return
	}

	r.ParseForm()
	username := r.PostForm["username"][0]
	password := r.PostForm["password"][0]

	user, validity, err := Authenticate(username, password)

	if err != nil {
		panic(err)
	}

	if validity {
		sessionId := MakeRandomString()
		err = user.CreateSession(sessionId)

		if err != nil {
			panic(err)
		}

		cookie := http.Cookie{
			Name:   "sessionId",
			Value:  sessionId,
			MaxAge: 60 * 60 * 24 * 7, // 期間は1週間
		}
		http.SetCookie(w, &cookie)

		result := Information{
			Message: "success",
		}

		err = WriteInformationAsJson(w, result)

		if err != nil {
			panic(err)
		}

	} else {
		result := Information{
			Message: "failed",
		}

		err = WriteInformationAsJson(w, result)

		if err != nil {
			panic(err)
		}
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	sessionIdCookie, err := r.Cookie("sessionId")

	if err != nil {
		MoveTo("homepage", w)
		return
	}

	sessionId := ExtractSessionIdFromCookie(sessionIdCookie)
	err = DeleteSession(sessionId)

	if err != nil {
		panic(err)
	}

	cookie := http.Cookie{
		Name:   "sessionId",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)
	MoveTo("homepage", w)

}

func Push(w http.ResponseWriter, r *http.Request) {
	user, validity, err := CheckUserIsAuthenticated(r)

	if (err != nil) || (!validity) {
		MoveTo("homepage", w)
		return
	}

	err = user.CreateCsrfToken()
	if err != nil {
		panic(err)
	}

	t, err := template.ParseFiles("templates/push.html")

	if err != nil {
		panic(err)
	}

	t.Execute(w, user)
}

func EarnMoney(w http.ResponseWriter, r *http.Request) {
	user, validity, err := CheckUserIsAuthenticated(r)

	if (err != nil) || (!validity) {
		result := Information{
			Message: "not authenticated",
		}

		err = WriteInformationAsJson(w, result)

		if err != nil {
			panic(err)
		}
		return
	}

	r.ParseForm()
	user.CsrfToken = r.PostForm["csrfToken"][0]
	validity, err = user.VerifyCsrfToken()

	if err != nil {
		panic(err)
	}

	if !validity {
		return
	}

	user.Money += int(math.Pow(2, float64(user.OakFruits)/10.0))
	err = user.Update()

	if err != nil {
		panic(err)
	}

	result := Information{
		Money: user.Money,
	}

	err = WriteInformationAsJson(w, result)

	if err != nil {
		panic(err)
	}
}

func Reset(w http.ResponseWriter, r *http.Request) {
	user, validity, err := CheckUserIsAuthenticated(r)

	if (err != nil) || (!validity) {
		result := Information{
			Message: "not authenticated",
		}

		err = WriteInformationAsJson(w, result)

		if err != nil {
			panic(err)
		}
		return
	}

	r.ParseForm()
	user.CsrfToken = r.PostForm["csrfToken"][0]
	validity, err = user.VerifyCsrfToken()

	if err != nil {
		panic(err)
	}

	if !validity {
		return
	}

	user.Money = 0
	user.OakFruits = 0
	user.ThunderFruits = 0
	err = user.Update()

	if err != nil {
		panic(err)
	}

	result := Information{
		Money:         user.Money,
		OakFruits:     user.OakFruits,
		ThunderFruits: user.ThunderFruits,
	}

	err = WriteInformationAsJson(w, result)

	if err != nil {
		panic(err)
	}
}

func Invest(w http.ResponseWriter, r *http.Request) {
	user, validity, err := CheckUserIsAuthenticated(r)

	if (err != nil) || (!validity) {
		result := Information{
			Message: "not authenticated",
		}

		err = WriteInformationAsJson(w, result)

		if err != nil {
			panic(err)
		}
		return
	}

	r.ParseForm()
	user.CsrfToken = r.PostForm["csrfToken"][0]
	validity, err = user.VerifyCsrfToken()

	if err != nil {
		panic(err)
	}

	if !validity {
		return
	}

	if user.ThunderFruits >= 1 {
		user.ThunderFruits -= 1
		err := user.Update()

		if err != nil {
			panic(err)
		}

		if ReturnTrueWithCertainProbability() {
			profit := user.Money * 100
			user.Money += profit
			err := user.Update()

			if err != nil {
				panic(err)
			}

			result := Information{
				Message:       "invested-success",
				Money:         user.Money,
				ThunderFruits: user.ThunderFruits,
				Profit:        profit,
			}

			err = WriteInformationAsJson(w, result)

			if err != nil {
				panic(err)
			}

		} else {
			result := Information{
				Message:       "invested-failed",
				Money:         user.Money,
				ThunderFruits: user.ThunderFruits,
			}

			err = WriteInformationAsJson(w, result)

			if err != nil {
				panic(err)
			}
		}

	} else {
		result := Information{
			Message:       "not invested",
			Money:         user.Money,
			ThunderFruits: user.ThunderFruits,
		}

		err = WriteInformationAsJson(w, result)

		if err != nil {
			panic(err)
		}
	}
}

func EnterStore(w http.ResponseWriter, r *http.Request) {
	user, validity, err := CheckUserIsAuthenticated(r)

	if (err != nil) || (!validity) {
		MoveTo("homepage", w)
		return
	}

	err = user.CreateCsrfToken()
	if err != nil {
		panic(err)
	}

	t, err := template.ParseFiles("templates/store.html")

	if err != nil {
		panic(err)
	}

	t.Execute(w, user)
}

func BuyOakFruits(w http.ResponseWriter, r *http.Request) {
	user, validity, err := CheckUserIsAuthenticated(r)

	if (err != nil) || (!validity) {
		result := Information{
			Message: "not authenticated",
		}

		err = WriteInformationAsJson(w, result)

		if err != nil {
			panic(err)
		}
		return
	}

	r.ParseForm()
	user.CsrfToken = r.PostForm["csrfToken"][0]
	validity, err = user.VerifyCsrfToken()

	if err != nil {
		panic(err)
	}

	if !validity {
		return
	}

	priceOfOakFruits := int(100.0 * (math.Pow(2, (float64(user.OakFruits)/10.0)+1.0)))

	if user.Money >= priceOfOakFruits {
		user.Money -= priceOfOakFruits
		user.OakFruits += 10
		err = user.Update()

		if err != nil {
			panic(err)
		}

		priceOfOakFruits = int(100.0 * (math.Pow(2, (float64(user.OakFruits)/10.0)+1.0)))
		priceOfThunderFruits := user.Money

		result := Information{
			Message:              "success",
			Money:                user.Money,
			OakFruits:            user.OakFruits,
			ThunderFruits:        user.ThunderFruits,
			PriceOfOakFruits:     priceOfOakFruits,
			PriceOfThunderFruits: priceOfThunderFruits,
		}

		err = WriteInformationAsJson(w, result)

		if err != nil {
			panic(err)
		}

	} else {
		result := Information{Message: "failed"}
		err = WriteInformationAsJson(w, result)

		if err != nil {
			panic(err)
		}
	}
}

func BuyThunderFruits(w http.ResponseWriter, r *http.Request) {
	user, validity, err := CheckUserIsAuthenticated(r)

	if (err != nil) || (!validity) {
		result := Information{
			Message: "not authenticated",
		}

		err = WriteInformationAsJson(w, result)

		if err != nil {
			panic(err)
		}
		return
	}

	r.ParseForm()
	user.CsrfToken = r.PostForm["csrfToken"][0]
	validity, err = user.VerifyCsrfToken()

	if err != nil {
		panic(err)
	}

	if !validity {
		return
	}

	requiredOakFruits := 100

	if user.OakFruits >= requiredOakFruits {
		user.OakFruits -= requiredOakFruits
		user.Money = 0
		user.ThunderFruits += 30
		err = user.Update()

		if err != nil {
			panic(err)
		}

		priceOfOakFruits := int(100.0 * (math.Pow(2, (float64(user.OakFruits)/10.0)+1.0)))
		priceOfThunderFruits := user.Money

		result := Information{
			Message:              "success",
			Money:                user.Money,
			OakFruits:            user.OakFruits,
			ThunderFruits:        user.ThunderFruits,
			PriceOfOakFruits:     priceOfOakFruits,
			PriceOfThunderFruits: priceOfThunderFruits,
		}

		err = WriteInformationAsJson(w, result)

		if err != nil {
			panic(err)
		}

	} else {
		result := Information{
			Message: "failed",
		}

		err = WriteInformationAsJson(w, result)

		if err != nil {
			panic(err)
		}
	}
}
