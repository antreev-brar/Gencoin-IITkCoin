package main 

import (
	"database/sql"
	"time"
	"net/http"
	"io"
	"log"
	"fmt"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	_ "golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
)

//////////////////////////////////////////////////////////////////////////////////////////////////////////
func Servepage(w http.ResponseWriter , r *http.Request){
	io.WriteString(w , "hello world")
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////
func Secretpage(w http.ResponseWriter , r *http.Request){
	w.Write([]byte("Welcome to the secret-page Amigo!\n"))
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////
func Signup(w http.ResponseWriter, req *http.Request) {
	    //check if path is correct
	    if req.URL.Path != "/signup" {
             http.NotFound(w, req)
             return
		}
        // Reference database to a variable
		db, _ := sql.Open("sqlite3" , "./database.db")
		defer db.Close()
		
	    fmt.Println(req.URL.Path)
	    switch req.Method {
		  case "GET":
				 //returns the entire database to log (not required really)
			      Get(db)
				  w.Write([]byte("Received a Get request\n"))
				  
	      case "POST":
				 // parse the input JSON into JSON object struct and add it to database (if not already exists)
				  var newUser SignupJSON
				  decoder := json.NewDecoder(req.Body)
				  err := decoder.Decode(&newUser)
				  CheckError(err)
				  log.Println(newUser)
				  Add(db, newUser)

				  //return the struct instance passed to add() function
				  w.Header().Set("Content-Type", "application/json")
				  json.NewEncoder(w).Encode(newUser)	
		  default:
				  //endpoint can't be accessed via other Request methods
	              w.WriteHeader(http.StatusNotImplemented)
	              w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
	     }
	
	  }

//////////////////////////////////////////////////////////////////////////////////////////////////////////
func Login(w http.ResponseWriter, req *http.Request) {
		//check if path is correct
	    if req.URL.Path != "/login" {
             http.NotFound(w, req)
             return
		}
		//Open database
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

				  // condition to check if the username and password is valid 
				  // If valid return JWT token
				  if Verify(db ,newUser){
					

					// Declare the expiration time of the token
					expirationTime := time.Now().Add(5* time.Minute)
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

					//set params of cookie as we are passing our JWT thorugh it 
					http.SetCookie(w, &http.Cookie{
						Name:    "token",
						Value:   tokenString,
						Expires: expirationTime,
						HttpOnly: true,
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

//////////////////////////////////////////////////////////////////////////////////////////////////////////
//wrapper to be used to secure any endpoint
func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
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
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		//log.Println("TOKEN:",tkn)
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
	
		//If everything is fine take user to desired endpoint
		endpoint(w, r)
	})
}


//////////////////////////////////////////////////////////////////////////////////////////////////////////
//Refresh ur token expiration time
func Refresh(w http.ResponseWriter, r *http.Request) {
    //check if path is indeed correct
	if r.URL.Path != "/refresh" {
		 http.NotFound(w, r)
		 return
	}
    log.Println("Refreshing Token")
    cookie, err := r.Cookie("token")
	tknStr := cookie.Value
	claims := &CustomClaims{}
    // parse JWT string and store it in claims , checking not required since its already done
	tkn , err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	log.Println("TOKEN:",tkn)
	//refresh is allowed only if token would expire in 30s
	if time.Unix( claims.ExpiresAt , 0 ).Sub(time.Now()) > 30*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//set params of cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   tokenString,
		Expires: expirationTime,
		HttpOnly: true,
	})
			  
	w.Write([]byte("Expiration time extended\n"))
  }