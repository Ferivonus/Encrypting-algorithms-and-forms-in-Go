package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"unicode"
	//"database/sql"
)

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

		/*
				var count int
			err := db.QueryRow("SELECT count(*) FROM users WHERE username = ?", username).Scan(&count)
			if err != nil {
				http.Error(w, "Internal server error.", http.StatusInternalServerError)
				return
			}
			if count > 0 {
				http.Error(w, "Kullanıcı adı zaten alınmış.", http.StatusBadRequest)
				return
			}

			// add user to the database
			_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, password)
			if err != nil {
				http.Error(w, "Internal server error.", http.StatusInternalServerError)
				return
			}
		*/

		// Cookie'leri oluşturun ve kaydedin
		cookieUsername := http.Cookie{Name: "username", Value: username}
		cookiePassword := http.Cookie{Name: "password", Value: password}
		http.SetCookie(w, &cookieUsername)
		http.SetCookie(w, &cookiePassword)

		// Yönlendirme yapın
		http.Redirect(w, r, "/welcome", http.StatusSeeOther)
	}
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
		/*
			db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/Go database")
			if err != nil {
				panic(err.Error())
			}
			defer db.Close()

			// INSERT INTO EncryptedFiles (PersonID, plaintext, key_of_plaintext, ciphertext, decryptedPlaintext) VALUES (1, 'hello world', 'my key', 'encrypted text', 'decrypted text')
			insert, err := db.Prepare("INSERT INTO EncryptedFiles(PersonID, plaintext, key_of_plaintext, ciphertext, decryptedPlaintext) VALUES(?,?,?,?,?)")
			if err != nil {
				panic(err.Error())
			}
			defer insert.Close()

			_, err = insert.Exec(1, "hello world", "my key", "encrypted text", "decrypted text")
			if err != nil {
				panic(err.Error())
			}

			fmt.Println("Kayıt eklendi.")
		*/

	}

}
