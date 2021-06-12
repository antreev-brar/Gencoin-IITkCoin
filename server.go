package main 

import (
	"database/sql"
	"time"
	"net/http"
	"io"
	"log"
	"fmt"
	//"io/ioutil"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	_ "golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
)

// Create the JWT key used to create the signature
var jwtKey = []byte("Dassvi_tu_tappi_ni_try_maare_jatt_te")

////////////////////////////////////////////////////////////////////////////////////////////////
func main(){
	database, _ := sql.Open("sqlite3" , "./database.db")
	defer database.Close()
	CreateTable(database)
	Get(database)

	http.HandleFunc("/" , servepage)
	http.HandleFunc("/signup" , signup)
	http.HandleFunc("/login" , login)
	http.Handle("/secretpage" , isAuthorized(secretpage))
    log.Println("Server listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000" , nil))
    //log.Println("Server listening on port 3000")
}
func servepage(w http.ResponseWriter , r *http.Request){
	io.WriteString(w , "hello world")
}
func secretpage(w http.ResponseWriter , r *http.Request){
	w.Write([]byte("Welcome to the secret-page Amigo!\n"))
	
}

func signup(w http.ResponseWriter, req *http.Request) {
	    if req.URL.Path != "/signup" {
             http.NotFound(w, req)
             return
		}

		db, _ := sql.Open("sqlite3" , "./database.db")
		defer db.Close()
		
	    fmt.Println(req.URL.Path)
	    switch req.Method {
		  case "GET":
			      Get(db)
				  w.Write([]byte("Received a Get request\n"))
				  
	      case "POST":

				  var newUser SignupJSON
				  decoder := json.NewDecoder(req.Body)
				  err := decoder.Decode(&newUser)
				  CheckError(err)
				  log.Println(newUser)
				  Add(db, newUser)
				  w.Header().Set("Content-Type", "application/json")
				  json.NewEncoder(w).Encode(newUser)	
	      default:
	              w.WriteHeader(http.StatusNotImplemented)
	              w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
	     }
	
	  }


func login(w http.ResponseWriter, req *http.Request) {
	    if req.URL.Path != "/login" {
             http.NotFound(w, req)
             return
		}

		db, _ := sql.Open("sqlite3" , "./database.db")
		defer db.Close()
		
	    fmt.Println(req.URL.Path)
	    switch req.Method {
		  case "GET":
			      //Get(db)
				  w.Write([]byte("Received a Get request on login \n"))
				  
	      case "POST":

				  var newUser LoginJSON
				  decoder := json.NewDecoder(req.Body)
				  err := decoder.Decode(&newUser)
				  CheckError(err)
				  log.Println(newUser)
				  if Verify(db ,newUser){
					

					// Declare the expiration time of the token
					// here, we have kept it as 5 minutes
					expirationTime := time.Now().Add(10* time.Minute)
					// Create the JWT claims, which includes the username and expiry time
					claims := &CustomClaims{
						Rollno: newUser.Rollno,
						StandardClaims: jwt.StandardClaims{
							// In JWT, the expiry time is expressed as unix milliseconds
							ExpiresAt: expirationTime.Unix(),
							},
					}

					// Declare the token with the algorithm used for signing, and the claims
					token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
					// Create the JWT string
					tokenString, err := token.SignedString(jwtKey)
					CheckError(err)
					log.Println("TOKEN:",tokenString)
					http.SetCookie(w, &http.Cookie{
						Name:    "token",
						Value:   tokenString,
						Expires: expirationTime,
					})
					w.Write([]byte("U son of a bitch. I am in \n")) 
				  } else {
					w.Write([]byte("Who are you mf ? identify urself nigga \n")) 
				  }
				  
	      default:
	              w.WriteHeader(http.StatusNotImplemented)
	              w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
	     }
	
	  }


func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// We can obtain the session token from the requests cookies, which come with every request
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				w.WriteHeader(http.StatusUnauthorized)
				return 
			}
			// For any other type of error, return a bad request status
			w.WriteHeader(http.StatusBadRequest)
			return 
		}
	
		// Get the JWT string from the cookie
		tknStr := cookie.Value
	
		// Initialize a new instance of `Claims`
		claims := &CustomClaims{}
	
		// Parse the JWT string and store the result in `claims`.
		// Note that we are passing the key in this method as well. This method will return an error
		// if the token is invalid (if it has expired according to the expiry time we set on sign in),
		// or if the signature does not match
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} 
	
		// Finally, return the welcome message to the user, along with their
		// username given in the token
		//w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Rollno)))
		endpoint(w, r)
	})
}

