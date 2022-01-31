package hundles

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"math/rand"
	"net/http"
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

func GetUser(id int) (user User, err error) {
	user = User{}
	err = Db.QueryRow("select id, username, password, money, oak_fruits, thunder_fruits from users where id = $1", id).Scan(&user.Id, &user.Username, &user.Password, &user.Money, &user.OakFruits, &user.ThunderFruits)
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

// ---------------------- ハンドラ関数 ----------------------------

// ハンドラを生成する関数
func ReturnFileHandler() (fileHandler http.Handler) {
	dir := http.Dir(".")
	fileHandler = http.FileServer(dir)
	fileHandler = http.StripPrefix("/file/", fileHandler)
	return
}

func ReturnUsersInformation(w http.ResponseWriter, r *http.Request) {
	user, err := GetUser(1)

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

func Push(w http.ResponseWriter, r *http.Request) {
	user, err := GetUser(1)

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
	user, err := GetUser(1)

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
	user, err := GetUser(1)

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
	user, err := GetUser(1)

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
	user, err := GetUser(1)

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
	user, err := GetUser(1)

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
	user, err := GetUser(1)

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
