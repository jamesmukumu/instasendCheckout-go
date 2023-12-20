package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

//create a struct of data

type Userpayment struct{
Firstname string `json:"first_name"`
Lastname string `json:"last_name"`
Phonenumber int `json:"phone_number"`
Amount int `json:"amount"`

}

 
//strcut of response
type Responseinstasend struct{
Id string `json:"id"`
CheckoutUrl string `json:"url"`

}


func main(){
Serversetup()


}




func Generatecheckoutlink(res http.ResponseWriter,req *http.Request){
dotenv := godotenv.Load()
if dotenv != nil {
	log.Fatal(dotenv)
	return
}

apiKeyinstasend := os.Getenv("publicKey")





var Checkout Responseinstasend
var Payloadinfo Userpayment

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

json.NewEncoder(res).Encode(map[string]string{
	"checkOut":Checkout.CheckoutUrl,
})

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