package main

import(
  "github.com/gorilla/mux"
  "log"
  "net/http"
  "encoding/json"
)

type Response struct{
  Message string `json:"message"`
  Status  bool `json:"status"`
}

func HolaMundo(w http.ResponseWriter, r *http.Request)  {
  w.Write([]byte ("Hola Mundo desde Go"))
}

func HolaJson(w http.ResponseWriter, r *http.Request)  {
  response := CreateResponse()
  json.NewEncoder(w).Encode(response)
}

func CreateResponse() Response{
  return Response{"Esto esta en formato Json", true}
}

func LoadStatic(w http.ResponseWriter, r *http.Request)  {
  http.ServeFile(w,r,"./Front/index.html")
}

func main() {

  cssHandle := http.FileServer(http.Dir("./Front/CSS/"))

  mux := mux.NewRouter()
  mux.HandleFunc("/Hola", HolaMundo).Methods("GET")
  mux.HandleFunc("/HolaJson",HolaJson).Methods("GET")
  mux.HandleFunc("/Static", LoadStatic).Methods("GET")

  http.Handle("/", mux)
  http.Handle("/CSS/", http.StripPrefix("/CSS/", cssHandle))
  log.Println("El servidor se encuentra en el puerto 8000")
  log.Fatal(http.ListenAndServe(":8000", nil))

}
