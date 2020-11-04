package main

// We'll use the standard HTTP library as well as the gorrila router for this app
import (
	"encoding/json"
	"errors"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	// Middleware that allows our Go application secure the application using JWT token
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// Verify 'aud' claim
			aud 		:= "https://thuocdongy.com/" //"https://dev--njhv5y3.au.auth0.com/api/v2/" //"YOUR_API_IDENTIFIER"
			checkAud 	:= token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAud {
				return token, errors.New("Invalid audience.")
			}
			// Verify 'iss'	claim
			iss 		:= "https://dev--njhv5y3.au.auth0.com/" // "https://YOUR_DOMAIN/"
			checkIss 	:= token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("Invalid issuer.")
			}

			cert, err := getPemCert(token)
			if err != nil {
				panic(err.Error())
			}

			result, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})

	// Instantiating the gorilla/mux router
	r := mux.NewRouter()

	// On the default page we will simply serve our static index page.
	r.Handle("/", http.FileServer(http.Dir("./views/")))
	// We will setup our server so we can serve static assets like images, css
	// from the /static/{file} route
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", // strip prefix
			http.FileServer(http.Dir("./static/")))) // serve from file directory

	// Our application will run on port 8080. Here we declare the port and pass in our router

	// See our API explanation in Readme.md
	// Add handler for products
	r.Handle("/status", StatusHandler).Methods("GET")
	//r.Handle("/products", ProductsHandler).Methods("GET")
	r.Handle("/products", jwtMiddleware.Handler(ProductsHandler)).Methods("GET")
	r.Handle("/products/{slug}/feedback", jwtMiddleware.Handler(AddFeedbackHandler)).Methods("POST")

	// For dev only - Set up CORS so React client can consume our API
	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"}	,
	})

	http.ListenAndServe(":8080", corsWrapper.Handler(r))
}

var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
})

func getPemCert(token *jwt.Token) (string, error) {
	cert		:= ""
	// "https://YOUR_DOMAIN/.well-known/jwks.json"
	resp, err 	:= http.Get("https://dev--njhv5y3.au.auth0.com/.well-known/jwks.json")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)
	//err = json.NewDecoder(resp.Body).Decode(&jwks)
	if err != nil {
		return cert, err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}

	return cert, nil
}