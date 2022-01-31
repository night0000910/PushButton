package main

import (
	"net/http"

	"pushbutton/hundles"
)

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	http.Handle("/file/", hundles.ReturnFileHandler())
	http.HandleFunc("/users_information", hundles.ReturnUsersInformation)
	http.HandleFunc("/signup", hundles.Signup)
	http.HandleFunc("/push", hundles.Push)
	http.HandleFunc("/earn_money", hundles.EarnMoney)
	http.HandleFunc("/invest", hundles.Invest)
	http.HandleFunc("/reset", hundles.Reset)
	http.HandleFunc("/store", hundles.EnterStore)
	http.HandleFunc("/buy_oak_fruits", hundles.BuyOakFruits)
	http.HandleFunc("/buy_thunder_fruits", hundles.BuyThunderFruits)
	server.ListenAndServe()
}
