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
	Db, err = sql.Open("postgres", "user=postgres dbname=pushbutton password=jhn5dzi38 sslmode=disable")

	if err != nil {
		panic(err)
	}
}

type User struct {
	Id            int
	Username      string
	Password      string
	Money         int
	OakFruits     int
	ThunderFruits int
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

func (user *User) Update() (err error) {
	_, err = Db.Exec("update users set username=$1, password=$2, money=$3, oak_fruits=$4, thunder_fruits=$5 where id = $6", user.Username, user.Password, user.Money, user.OakFruits, user.ThunderFruits, user.Id)
	return
}

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

func CreateUser(username string, password string) (err error) {
	_, err = Db.Exec("insert into users (username, password) values ($1, $2)", username, password)
	return
}

func Authenticate(username string, password string) (user User, validity bool, err error) {
	user = User{}
	var numberOfRecords int
	err = Db.QueryRow("select count(*) from users where username = $1 and password = $2", username, password).Scan(&numberOfRecords)

	if err != nil {
		return
	}

	if numberOfRecords == 1 {
		validity = true
		err = Db.QueryRow("select id, username, password, money, oak_fruits, thunder_fruits from users where username = $1", username).Scan(&user.Id, &user.Username, &user.Password, &user.Money, &user.OakFruits, &user.ThunderFruits)

	} else if numberOfRecords == 0 {
		validity = false
	}

	return
}

func CheckUserIsAuthenticated(w http.ResponseWriter, r *http.Request) (user User, err error) {
	sessionIdCookie, err := r.Cookie("sessionId")

	if err != nil {
		MoveToLoginPage(w)
	}

	sessionId := ExtractSessionIdFromCookie(sessionIdCookie)
	var validity bool
	user, validity, err = GetUser(sessionId)

	if !validity {
		MoveToLoginPage(w)
	}

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

func CreateSession(userId int, sessionId string) (err error) {
	hashedSessionId := Hash(sessionId)
	_, err = Db.Exec("insert into session (user_id, session_id) values ($1, $2)", userId, hashedSessionId)
	return
}

func DeleteSession(sessionId string) (err error) {
	hashedSessionId := Hash(sessionId)
	_, err = Db.Exec("delete from session where session_id = $1", hashedSessionId)
	return
}

func returnTrueWithCertainProbability() (result bool) {
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

func MoveToLoginPage(w http.ResponseWriter) {
	w.Header().Set("Location", "http://127.0.0.1:8080/login_page")
	w.WriteHeader(302)
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
	user, err := CheckUserIsAuthenticated(w, r)

	if err != nil {
		panic(err)
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
	var jsonData []byte
	jsonData, err = json.Marshal(information)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, string(jsonData))
}

func DisplaySignupPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/signup.html")

	if err != nil {
		panic(err)
	}

	t.Execute(w, nil)
}

func Signup(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm["username"][0]
	password := r.PostForm["password"][0]
	hashedPassword := Hash(password)

	err := CreateUser(username, hashedPassword)

	var result Information

	if err == nil {
		result = Information{
			Message: "success",
		}

	} else {
		fmt.Println(err)
		result = Information{
			Message: "failed",
		}
	}

	var jsonData []byte
	jsonData, err = json.Marshal(result)

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, string(jsonData))
}

func SucceedInSignup(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/succeedInSignup.html")

	if err != nil {
		panic(err)
	}

	t.Execute(w, nil)
}

func DisplayLoginPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/login.html")

	if err != nil {
		panic(err)
	}

	t.Execute(w, nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm["username"][0]
	password := r.PostForm["password"][0]
	hashedPassword := Hash(password)

	user, validity, err := Authenticate(username, hashedPassword)

	if err != nil {
		panic(err)
	}

	if validity {
		sessionId := MakeRandomString()
		err = CreateSession(user.Id, sessionId)

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
		var jsonData []byte
		jsonData, err = json.Marshal(result)

		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(jsonData))

	} else {
		result := Information{
			Message: "failed",
		}
		var jsonData []byte
		jsonData, err = json.Marshal(result)

		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(jsonData))
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	sessionIdCookie, err := r.Cookie("sessionId")

	if err != nil {
		panic(err)
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
	MoveToLoginPage(w)

}

func Push(w http.ResponseWriter, r *http.Request) {
	user, err := CheckUserIsAuthenticated(w, r)

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
	user, err := CheckUserIsAuthenticated(w, r)

	if err != nil {
		panic(err)
	}

	user.Money += int(math.Pow(2, float64(user.OakFruits)/10.0))
	err = user.Update()

	if err != nil {
		panic(err)
	}

	result := Information{
		Money: user.Money,
	}
	var jsonData []byte
	jsonData, err = json.Marshal(result)

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, string(jsonData))
}

func Reset(w http.ResponseWriter, r *http.Request) {
	user, err := CheckUserIsAuthenticated(w, r)

	if err != nil {
		panic(err)
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
	var jsonData []byte
	jsonData, err = json.Marshal(result)

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, string(jsonData))
}

func Invest(w http.ResponseWriter, r *http.Request) {
	user, err := CheckUserIsAuthenticated(w, r)

	if err != nil {
		panic(err)
	}

	if user.ThunderFruits >= 1 {
		user.ThunderFruits -= 1
		err := user.Update()

		if err != nil {
			panic(err)
		}

		if returnTrueWithCertainProbability() {
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
			var jsonData []byte
			jsonData, err = json.Marshal(result)

			if err != nil {
				panic(err)
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, string(jsonData))

		} else {
			result := Information{
				Message:       "invested-failed",
				Money:         user.Money,
				ThunderFruits: user.ThunderFruits,
			}
			var jsonData []byte
			jsonData, err = json.Marshal(result)

			if err != nil {
				panic(err)
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, string(jsonData))
		}

	} else {
		result := Information{
			Message:       "not invested",
			Money:         user.Money,
			ThunderFruits: user.ThunderFruits,
		}
		var jsonData []byte
		jsonData, err = json.Marshal(result)

		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(jsonData))
	}
}

func EnterStore(w http.ResponseWriter, r *http.Request) {
	user, err := CheckUserIsAuthenticated(w, r)

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
	user, err := CheckUserIsAuthenticated(w, r)

	if err != nil {
		panic(err)
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
		var jsonData []byte
		jsonData, err = json.Marshal(result)

		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(jsonData))

	} else {
		result := Information{Message: "failed"}
		var jsonData []byte
		jsonData, err = json.Marshal(result)

		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(jsonData))
	}
}

func BuyThunderFruits(w http.ResponseWriter, r *http.Request) {
	user, err := CheckUserIsAuthenticated(w, r)

	if err != nil {
		panic(err)
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
		var jsonData []byte
		jsonData, err = json.Marshal(result)

		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(jsonData))

	} else {
		result := Information{
			Message: "failed",
		}
		var jsonData []byte
		jsonData, err = json.Marshal(result)

		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(jsonData))
	}
}
