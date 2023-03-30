package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode"

	//"database/sql"

	"database/sql"

	"os"

	//"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type users struct {
	name     string
	password string
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	// attention: If you do not call ParseForm method, the following data can not be obtained form
	fmt.Println(r.Form) // print information on server side.
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	//fmt.Fprintf(w, "Hello astaxie!") // write data to response
}

func login(w http.ResponseWriter, r *http.Request) {

	cookieUsername, Usernameerr := r.Cookie("username")

	cookiePassword, ePassworderr := r.Cookie("password")

	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		//if there is 2 cookies
		if Usernameerr == nil && ePassworderr == nil {

			println(w, "hoşgeldin ", cookieUsername)
			println(w, "Duyduğuma göre şifren de şuymuş: ", cookiePassword)

			http.Redirect(w, r, "/welcome", http.StatusSeeOther)
			return
		}

		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
		return
	} else {
		r.ParseForm()

		//if there is not any particular account of user
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])

		//Geting information from form

		var username = r.FormValue("username")
		var password = r.FormValue("password")

		println(w, "Girdiğiniz bilgiler şunlar: ")
		println(w, "Hoşeldiniz !", username)
		println(w, "Şifreniz de !", password)

		//check on the system if it's true:

		//if true

		cookieUsername := http.Cookie{Name: "username", Value: username}
		cookiePassword := http.Cookie{Name: "password", Value: password}
		http.SetCookie(w, &cookieUsername)
		http.SetCookie(w, &cookiePassword)

		http.Redirect(w, r, "/welcome", http.StatusSeeOther)
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method

	cookieUsername, Usernameerr := r.Cookie("username")

	cookiePassword, ePassworderr := r.Cookie("password")
	if r.Method == "GET" {
		//if there is 2 cookies
		if Usernameerr == nil && ePassworderr == nil {

			println(w, "hoşgeldin ", cookieUsername)
			println(w, "Duyduğuma göre şifren de şuymuş: ", cookiePassword)

			http.Redirect(w, r, "/welcome", http.StatusSeeOther)
			return
		}

		t, _ := template.ParseFiles("register.gtpl")
		t.Execute(w, nil)
		return
	} else {
		// Form verilerini alın
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		err := godotenv.Load()
		if err != nil {
			fmt.Println("Error loading .env file")
		}

		// MySQL bağlantısı için gerekli değişkenleri .env dosyasından okuyun
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")

		fmt.Println(dbUser, dbPassword, dbName, dbHost, dbPort)

		dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
		db, err := sql.Open("mysql", dataSourceName)

		if err != nil {
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}

		var count int

		p := users{
			name:     username,
			password: password,
		}

		errrr := db.QueryRow("SELECT count(*) FROM users WHERE username = ?", username).Scan(&count)

		fmt.Println(errrr)
		if errrr != nil {
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			//return
		}

		if count > 0 {
			http.Error(w, "Kullanıcı adı zaten alınmış.", http.StatusBadRequest)
			return
		} else {
			fmt.Println("This account could be taken from u")
		}

		// add user to the database
		err = insert(db, p)

		if err != nil {
			log.Printf("Insert user failed with error %s", err)
		}

		db.Close()

		fmt.Println(err)

		// Cookie'leri oluşturun ve kaydedin
		cookieUsername := http.Cookie{Name: "username", Value: username}
		cookiePassword := http.Cookie{Name: "password", Value: password}
		http.SetCookie(w, &cookieUsername)
		http.SetCookie(w, &cookiePassword)

		// Yönlendirme yapın
		http.Redirect(w, r, "/welcome", http.StatusSeeOther)
	}
}

func insert(db *sql.DB, p users) error {
	query := "INSERT INTO users(username, userpassword) VALUES (?, ?)"
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, p.name, p.password)
	if err != nil {
		log.Printf("Error %s when inserting row into products table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d products created ", rows)
	return nil
}

func index(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.gtpl")
	t.Execute(w, nil)
}

func welcome(w http.ResponseWriter, r *http.Request) {

	// cookie'den kullanıcı adını ve şifreyi al
	cookieUsername, err := r.Cookie("username")
	cookiePassword, err := r.Cookie("password")
	if err != nil {
		http.Error(w, "Kullanıcı adı veya şifre cookie'si bulunamadı.", http.StatusNotFound)
		return
	}

	// kullanıcı adını ve şifreyi ekrana yazdır
	t, _ := template.ParseFiles("welcome.gtpl")
	println(w, "Hoş geldin, &s! Şifren: &s", cookieUsername.Value, cookiePassword.Value)
	t.Execute(w, nil)
}

func main() {

	http.HandleFunc("/", index) // setting router rule
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.HandleFunc("/welcome", welcome)

	http.HandleFunc("/Encription", Encription)

	var err = http.ListenAndServe(":8080", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func encrypt(plaintext string, key string) string {
	var ciphertext strings.Builder
	keyLen := len(key)

	for i, currentChar := range plaintext {
		// Get the current character and convert it to uppercase
		currentChar = unicode.ToUpper(currentChar)

		// Calculate the index of the current character in the alphabet
		alphaIndex := byte(currentChar - 'A')

		// Calculate the index of the current character in the key
		keyIndex := i % keyLen
		keyChar := unicode.ToUpper(rune(key[keyIndex]))
		keyOffset := byte(keyChar - 'A')

		// Calculate the new index of the current character in the alphabet
		newAlphaIndex := (alphaIndex + keyOffset) % 26

		// Calculate the new character and add it to the ciphertext
		newChar := rune(newAlphaIndex + 'A')
		ciphertext.WriteRune(newChar)
	}

	return ciphertext.String()
}

func decrypt(ciphertext string, key string) string {
	var plaintext strings.Builder
	keyLen := len(key)

	for i, currentChar := range ciphertext {
		// Get the current character and convert it to uppercase
		currentChar = unicode.ToUpper(currentChar)

		// Calculate the index of the current character in the alphabet
		alphaIndex := byte(currentChar - 'A')

		// Calculate the index of the current character in the key
		keyIndex := i % keyLen
		keyChar := unicode.ToUpper(rune(key[keyIndex]))
		keyOffset := byte(keyChar - 'A')

		// Calculate the new index of the current character in the alphabet
		newAlphaIndex := (alphaIndex + 26 - keyOffset) % 26

		// Calculate the new character and add it to the plaintext
		newChar := rune(newAlphaIndex + 'A')
		plaintext.WriteRune(newChar)
	}

	return plaintext.String()
}

func Encription(w http.ResponseWriter, r *http.Request) {

	cookieUsername, err := r.Cookie("username")
	if err != nil {
		http.Error(w, "Kullanıcı adı veya şifre cookie'si bulunamadı.", http.StatusNotFound)
		return
	}

	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("Encription.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()

		// logic part of log in
		plaintext := r.Form.Get("plaintext")
		key := r.Form.Get("key")
		fmt.Println("plaintext:", plaintext)
		fmt.Println("key:", key)

		ciphertext := encrypt(plaintext, key)
		decryptedPlaintext := decrypt(ciphertext, key)

		fmt.Fprintf(w, "Merhaba bay %s", cookieUsername.Value)
		fmt.Fprintf(w, `
        <p>Plaintext: %s</p>
        <p>Key: %s</p>
        <p>Ciphertext: %s</p>
        <p>Decrypted plaintext: %s</p>
    `, plaintext, key, ciphertext, decryptedPlaintext)

		// MySQL bağlantısı için gerekli değişkenleri .env dosyasından okuyun
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")

		fmt.Println(dbUser, dbPassword, dbName, dbHost, dbPort)

		dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
		db, err := sql.Open("mysql", dataSourceName)

		if err != nil {
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		// INSERT INTO EncryptedFiles (PersonID, plaintext, key_of_plaintext, ciphertext, decryptedPlaintext) VALUES (1, 'hello world', 'my key', 'encrypted text', 'decrypted text')
		insert, err := db.Prepare("INSERT INTO encryptedfiles(PersonID, plaintext, key_of_plaintext, ciphertext, decryptedPlaintext) VALUES(?,?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		_, err = insert.Exec(1, "hello world", "my key", "encrypted text", "decrypted text")
		if err != nil {
			panic(err.Error())
		}

		fmt.Println("Kayıt eklendi.")

	}

}
