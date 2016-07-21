package main

import(
  "github.com/gorilla/mux"
  "log"
  "net/http"
  "encoding/json"
  "sync"
)

type Response struct{
  Message string `json:"message"`
  Status  int   `json:"status"`
  IsValid bool   `json:"isvalid"`
}

var Users = struct{
  m map[string] User
  sync.RWMutex

}{m: make(map[string] User)}

type User struct{
  User_name string
}

func HolaMundo(w http.ResponseWriter, r *http.Request)  {
  w.Write([]byte ("Hola Mundo desde Go"))
}

func HolaJson(w http.ResponseWriter, r *http.Request)  {
  response := CreateResponse("hola json",200, true)
  json.NewEncoder(w).Encode(response)
}

func CreateResponse(message string, status int, valid bool) Response{
  return Response{message,status,valid}
}

func LoadStatic(w http.ResponseWriter, r *http.Request)  {
  http.ServeFile(w,r,"./Front/index.html")
}

func UserExist(user_name string)bool{
  Users.RLock()
  defer Users.RUnlock()

  if _, ok := Users.m[user_name]; ok{
    return true
  }
  return false
}

func Validate(w http.ResponseWriter, r *http.Request)  {
  r.ParseForm()
  user_name := r.FormValue("user_name")

  response := Response{}

  if UserExist(user_name){
    //Vamos a no permitir el ingreso
    response.IsValid = false
  }else{
    //Vamos a permitir el ingreso
    response.IsValid = true
  }
  json.NewEncoder(w).Encode(response)
}

func main() {

  cssHandle := http.FileServer(http.Dir("./Front/CSS/"))
  Js_Handle := http.FileServer(http.Dir("./Front/JS/"))

  mux := mux.NewRouter()
  mux.HandleFunc("/Hola", HolaMundo).Methods("GET")
  mux.HandleFunc("/HolaJson",HolaJson).Methods("GET")
  mux.HandleFunc("/", LoadStatic).Methods("GET")
  mux.HandleFunc("/validate", Validate).Methods("POST")

  http.Handle("/", mux)
  http.Handle("/CSS/", http.StripPrefix("/CSS/", cssHandle))
  http.Handle("/JS/", http.StripPrefix("/JS/", Js_Handle))
  log.Println("El servidor se encuentra en el puerto 8000")
  log.Fatal(http.ListenAndServe(":8000", nil))

}
