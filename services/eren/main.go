package main

import (
	"log"
	"os"

	"deviloza.com.mx/auth"
	"deviloza.com.mx/misc"
	"deviloza.com.mx/web"
)

func main() {
	log.Println("Starting web server")
	baseDir := os.Getenv("MICRO_SERVICES_MODS")

	hostsFile := baseDir + "/availableHosts.json"
	hosts, err := web.GetHostsFromFile(hostsFile)
	if err != nil {
		panic(err)
	}

	secrets, err := misc.NewSecrets(baseDir + "/secrets.json")
	if err != nil {
		panic(err)
	}

	srv := web.NewServer(secrets, "1239")
	authHandler := auth.NewAuthHandler(
		secrets.GenerateDSN(),
		"deviloza.com.mx:1239/",
		secrets.SendGridAPIKey,
		&auth.AuthOptions{},
	)
	webHandler := web.NewWebHandler(
		secrets.GenerateDSN(),
		hosts,
	)

	rs := web.Routes{
		web.Route{
			Method:       "POST",
			Url:          "/signup",
			Handler:      authHandler.SignUpUser,
			IsRestricted: false,
		},
		web.Route{
			Method:       "POST",
			Url:          "/login",
			Handler:      authHandler.LoginUser,
			IsRestricted: false,
		},
		web.Route{
			Method:       "GET",
			Url:          "/verify",
			Handler:      authHandler.ValidateEmail,
			IsRestricted: false,
		},
		web.Route{
			Method:       "GET",
			Url:          "/profile",
			Handler:      webHandler.Profile,
			IsRestricted: true,
		},
		web.Route{
			Method:       "POST",
			Url:          "/upload",
			Handler:      webHandler.UploadFile,
			IsRestricted: true,
		},
		web.Route{
			Method:       "GET",
			Url:          "/document/:docId",
			Handler:      webHandler.GetDocument,
			IsRestricted: true,
		},
		web.Route{
			Method: "GET",
			Url: "/documents",
			Handler: webHandler.GetUserDocuments,
			IsRestricted: true,
		},
	}

	err = srv.Start(rs)
	log.Println(err)
}
