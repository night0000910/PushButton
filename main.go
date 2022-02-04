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
	http.HandleFunc("/", hundles.MoveToHomepageOrPushPage)
	http.HandleFunc("/homepage", hundles.DisplayHomepage)
	http.HandleFunc("/explanation", hundles.DisplayExplanation)
	http.HandleFunc("/signup_page", hundles.DisplaySignupPage)
	http.HandleFunc("/signup", hundles.Signup)
	http.HandleFunc("/succeed_in_signup", hundles.SucceedInSignup)
	http.HandleFunc("/login_page", hundles.DisplayLoginPage)
	http.HandleFunc("/login", hundles.Login)
	http.HandleFunc("/logout", hundles.Logout)
	http.HandleFunc("/push", hundles.Push)
	http.HandleFunc("/earn_money", hundles.EarnMoney)
	http.HandleFunc("/invest", hundles.Invest)
	http.HandleFunc("/reset", hundles.Reset)
	http.HandleFunc("/store", hundles.EnterStore)
	http.HandleFunc("/buy_oak_fruits", hundles.BuyOakFruits)
	http.HandleFunc("/buy_thunder_fruits", hundles.BuyThunderFruits)
	server.ListenAndServe()
}
