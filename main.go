package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"

	"os"
	"strings"

	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

//create a struct of data

type Userpayment struct{
Firstname string `json:"first_name"`
Lastname string `json:"last_name"`
Phonenumber int `json:"phone_number"`
Amount int `json:"amount"`
Email string `json:"email"`
}

 
//strcut of response
type Responseinstasend struct{
Id string `json:"id"`
CheckoutUrl string `json:"url"`
Email string `json:"email"`

}


func main(){
Serversetup()


}
var Checkout Responseinstasend



func Generatecheckoutlink(res http.ResponseWriter,req *http.Request){
	var Payloadinfo Userpayment
dotenv := godotenv.Load()
if dotenv != nil {
	log.Fatal(dotenv)
	return
}

apiKeyinstasend := os.Getenv("publicKey")








client := http.Client{}
//decode json
reqjson := json.NewDecoder(req.Body).Decode(&Payloadinfo)
if reqjson != nil {
	log.Fatal(reqjson)
}



requestBody,err :=json.Marshal(Payloadinfo)
actualRequestBody := strings.NewReader(string(requestBody)) 


configReq,err := http.NewRequest("POST","https://sandbox.intasend.com/api/v1/checkout/",actualRequestBody)


if err != nil {
log.Fatal(err)
return	
}


//add auhtorization header

configReq.Header.Add("X-IntaSend-Public-API-Key",apiKeyinstasend)
configReq.Header.Add("content-type","application/json")
response,errr := client.Do(configReq)
if errr != nil {
	panic(errr.Error())

}


errorr := json.NewDecoder(response.Body).Decode(&Checkout)
if errorr != nil {
	panic(err)
}
Sendcheckoutlink(Checkout.CheckoutUrl,Checkout.Email)

json.NewEncoder(res).Encode(map[string]string{"Check out link":Checkout.CheckoutUrl})
fmt.Println(Checkout.Email)

defer response.Body.Close()

}




func Serversetup(){
dotenv := godotenv.Load()
if dotenv != nil {
log.Fatal(dotenv)	
return
}

PORT := os.Getenv("port")



Router := mux.NewRouter()
Router.HandleFunc("/post/details",Generatecheckoutlink).Methods("POST")


fmt.Printf("Server listening for requests at %s",PORT)
http.ListenAndServe(":"+PORT,Router)

}





func Sendcheckoutlink(checkoutlink string, checkoutemail string){
dotenv := godotenv.Load()
if dotenv != nil {
	log.Fatal(dotenv.Error())
	return
}
gmailPassword := os.Getenv("instasend")


mail := gomail.NewMessage()
mail.SetHeader("From","jamesmukumu03@gmail.com")
mail.SetHeader("To",checkoutemail)
mail.SetHeader("Subject","Check Out Link")
mail.SetBody("text/plain",checkoutlink)


//dialer
dialer := gomail.NewDialer("smtp.gmail.com",587,"jamesmukumu03@gmail.com",gmailPassword)

dialer.TLSConfig = &tls.Config{
	InsecureSkipVerify: false,
	ServerName: "smtp.gmail.com",
}
 
 
errorSendingmail := dialer.DialAndSend(mail)
if errorSendingmail !=nil {
	log.Fatal(errorSendingmail.Error())
}


 
}