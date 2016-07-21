package main

import(
  "github.com/gorilla/websocket"
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
  WebSocket *websocket.Conn
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

func CreateUser(user_name string, ws *websocket.Conn) User{
  return User{user_name, ws}

}

func AddUser(user User){
  Users.Lock()
  defer Users.Unlock()

  Users.m[user.User_name] = user

}

func RemoveUser(user_name string){
  Users.Lock()
  defer Users.Unlock()
  delete(Users.m, user_name)
}

func ToArrayBite(value string) []byte{
  return []byte(value)
}

func ConactMessage(user_name string, arreglo []byte) string {
  return user_name + " : " + string(arreglo[:])
}

func SendMessage(type_message int, message []byte)  {
  Users.RLock()
  defer Users.RUnlock()

  for _, user := range Users.m{
    if err := user.WebSocket.WriteMessage(type_message, message); err != nil{
      return
    }
  }
}

func WebSocket(w http.ResponseWriter, r *http.Request)  {
  vars := mux.Vars(r)
  user_name := vars["user_name"]

  ws, err := websocket.Upgrade(w,r,nil,1024,1024)
  if err != nil{
    log.Println (err)
    return
  }

  current_user := CreateUser(user_name,ws)
  AddUser(current_user)
  log.Println("Nuevo Usuario Agregado")

  for{

    type_message, message, err := ws.ReadMessage()
    if err != nil{
      RemoveUser(user_name)
      return
    }
    final_message := ConactMessage(user_name,message)
    SendMessage(type_message, ToArrayBite(final_message))
  }
}

func main() {

  cssHandle := http.FileServer(http.Dir("./Front/CSS/"))
  Js_Handle := http.FileServer(http.Dir("./Front/JS/"))

  mux := mux.NewRouter()
  mux.HandleFunc("/Hola", HolaMundo).Methods("GET")
  mux.HandleFunc("/HolaJson",HolaJson).Methods("GET")
  mux.HandleFunc("/Chat/{user_name}", WebSocket).Methods("GET")
  mux.HandleFunc("/", LoadStatic).Methods("GET")
  mux.HandleFunc("/validate", Validate).Methods("POST")

  http.Handle("/", mux)
  http.Handle("/CSS/", http.StripPrefix("/CSS/", cssHandle))
  http.Handle("/JS/", http.StripPrefix("/JS/", Js_Handle))
  log.Println("El servidor se encuentra en el puerto 8000")
  log.Fatal(http.ListenAndServe(":8000", nil))

}
