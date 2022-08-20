package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Menampilkan data
type Product struct {
	Id    int    `json:"id"`
	Kode  string `json:"kode"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
}

// Menampilkan kode response
type Hasil struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// Koneksi Ke Database
var db *gorm.DB
var err error

func main() {
	db, err = gorm.Open("mysql", "root:@/restapi?charset=utf8&parseTime=True")
	if err != nil {
		log.Printf("Koneksi gagal")
	}

	// Migrasi table struct ke database
	db.AutoMigrate(&Product{})

	// Menghandle request dari URL
	handleRequest()
}

// Handle request dari URL
func handleRequest() {
	log.Println("Menjalankan di port htttp://127.0.0.1:8000")

	// Routing
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homepage)
	myRouter.HandleFunc("/api/products", createProducts).Methods("POST")
	myRouter.HandleFunc("/api/products", getProducts).Methods("GET")
	myRouter.HandleFunc("/api/products/{id}", getProduct).Methods("GET")
	myRouter.HandleFunc("/api/products/{id}", updateProducts).Methods("PUT")
	myRouter.HandleFunc("/api/products/{id}", deleteProduct).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", myRouter))
}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Berhasil")
}

func createProducts(w http.ResponseWriter, r *http.Request) {
	// Membuat payload untuk menangkap yang diinput
	payloads, _ := ioutil.ReadAll(r.Body)

	var product Product
	// Data yang di tangkap di casting ke struct Product
	json.Unmarshal(payloads, &product)

	// Setelah ditangkap, lalu dimasukkan ke table product
	db.Create(&product)

	// Susunan response
	res := Hasil{Code: 200, Data: product, Message: "Sukses membuat product"}
	result, err := json.Marshal(res)

	// Deklarasi variable err
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Response json
	// "w" untuk Response dan "R untuk Request"
	w.Header().Set("Content-type", "application/json")
	// Status untuk HTTP
	w.WriteHeader(http.StatusOK)
	// Lempar ke response body
	w.Write(result)
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	// Mengambil products dengan array []Product{}
	products := []Product{}

	// Perbedaan mengambil semua data
	db.Find(&products)

	res := Hasil{Code: 200, Data: products, Message: "Sukses membaca product"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	// Cara mengambil dengan ID
	vars := mux.Vars(r)
	productID := vars["id"]

	var product Product

	// Perbedaan untuk mengambil dengan ID
	db.First(&product, productID)

	res := Hasil{Code: 200, Data: product, Message: "Sukses membaca product"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func updateProducts(w http.ResponseWriter, r *http.Request) {
	// Untuk mengupdate product, regulasi peraturan hampir sama dengan create product by ID
	vars := mux.Vars(r)
	productID := vars["id"]

	var productUpdate Product

	// Menangkap sama dengan seperti membuat
	payloads, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(payloads, &productUpdate)

	// membuat variable untuk existing
	var product Product

	// Get data product dari table product bedasarkan ID
	db.First(&product, productID)

	// Melakukan update
	db.Model(&product).Updates(productUpdate)

	res := Hasil{Code: 200, Data: product, Message: "Data berhasil di update"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	var product Product
	db.First(&product, productID)

	// Melakukan delete
	db.Delete(&product)

	res := Hasil{Code: 200, Message: "Data berhasil di delete"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
