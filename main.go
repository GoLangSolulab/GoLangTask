package main

import (
	
	"fmt"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"encoding/json"
	"strings"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	
)

//var db *sql.DB
var sku_code string
var keys []string

// structures for coffee pod
type pods struct{
	SKey string       `json:"s_key"`
	PodType string    `json:"pod_type"`
	Size int          `json:"size"`
	Flavour string    `json:"flavour"`
}

// structure for coffee machine
type machine struct{
	MKey string        `json:"m_key"`
	MachineType string `json:"machine_type"`
	Model string       `json:"model"`
	WaterLine int      `json:"water_line"`
}

// function to return cross sell products
func get_products(parsed_skucode string , t string, db *sql.DB) []pods {

	var cross_sell_pods[] pods
	// condition to return coffee pods as cross sell products for coffee machine

	if parsed_skucode == "CM" { //search by coffee machine
		var key string
		var m_type string
		var model string
		var comp int
		//defer stmt_m.Close()
		stmt_m, err1:= db.Prepare("select * from info_cm where m_key=?")
		
		if err1!= nil{
			fmt.Println("error")
			panic(err1.Error())	
		}
		result_m, err1 := stmt_m.Query(sku_code)
		if err1!= nil{
			fmt.Println("error")
			panic(err1.Error())	
		}
		for result_m.Next(){
			err1 = result_m.Scan(&key, &m_type, &model, &comp)	
			if err1!= nil{
				panic(err1.Error())	
			}
		}
		// conditions to fetch the pods from db	
		stmt, err1 := db.Prepare("select * from info_cp where pod_type=? and size=?")
		if err1!= nil{
			fmt.Println("error")
			panic(err1.Error())	
		}
		defer stmt.Close()
		result_pods, err1 := stmt.Query(m_type,1)
		if err1!= nil{
			fmt.Println("error")
			panic(err1.Error())	
		}
		defer result_pods.Close()
		var pod_key string
		var pod_typ string
		var size int
		var flav string

		for result_pods.Next(){
			err1 = result_pods.Scan(&pod_key, &pod_typ, &size, &flav)
			if err1!= nil{
				panic(err1.Error())	
			}
			temp_pod := pods{
				SKey : pod_key,
				PodType : pod_typ,
				Size : size,
				Flavour : flav,
			}
			cross_sell_pods = append(cross_sell_pods, temp_pod)	
		}
		
		fmt.Println("values")
	
	} else if sku_code == "espresso vanilla"{ //search by espresso vanilla flavour
		
		result_p, err1 := db.Query("select * from info_cp where pod_type = ? and flavour = ?", keys[0], keys[1])
		if err1!= nil{
			panic(err1.Error())	
		}
		
		var pod_key string
		var pod_typ string
		var pod_size int
		var pod_flav string
		
		for result_p.Next(){
			result_p.Scan(&pod_key,&pod_typ,&pod_size,&pod_flav)
			if err1!= nil{
				panic(err1.Error())	
			}
			temp_pod := pods{
				SKey : pod_key,
				PodType : pod_typ,
				Size : pod_size,
				Flavour : pod_flav,
			}
			cross_sell_pods = append(cross_sell_pods, temp_pod)
		}
			
	} else if sku_code == "espresso machine"{ //search by espresso machine type
		
		var pod_key string
		var pod_typ string
		var pod_size int
		var pod_flav string
		
		result_m, err1 := db.Query("select * from info_cp where pod_type=? and size=?", keys[0], 3)
		
		if err1!= nil{
			panic(err1.Error())	
		}
		
		for result_m.Next(){
			result_m.Scan(&pod_key,&pod_typ,&pod_size,&pod_flav)	
			temp_pod := pods{
				SKey : pod_key,
				PodType : pod_typ,
				Size : pod_size,
				Flavour : pod_flav,
			}
			cross_sell_pods = append(cross_sell_pods, temp_pod)
		}
		
	} else if sku_code == "vanilla"{ // search by general vanilla
		
		var pod_key string
		var pod_typ string
		var pod_size int
		var pod_flav string

		result_m, err1 := db.Query("select * from info_cp where pod_type != ? and size=? and flavour=?","espresso",1,sku_code)
		if err1!= nil{
			panic(err1.Error())	
		}
		
		for result_m.Next(){
			result_m.Scan(&pod_key,&pod_typ,&pod_size,&pod_flav)	
			temp_pod := pods{
				SKey : pod_key,
				PodType : pod_typ,
				Size : pod_size,
				Flavour : pod_flav,
			}
			cross_sell_pods = append(cross_sell_pods, temp_pod)
		}
		
		result_m, err1 = db.Query("select * from info_cp where pod_type = ? and size=? and flavour=?","espresso",3,sku_code)
		
		if err1!= nil{
			panic(err1.Error())	
		}
		
		for result_m.Next(){
			result_m.Scan(&pod_key,&pod_typ,&pod_size,&pod_flav)	
			temp_pod := pods{
				SKey : pod_key,
				PodType : pod_typ,
				Size : pod_size,
				Flavour : pod_flav,
			}
			cross_sell_pods = append(cross_sell_pods, temp_pod)
		}
		
	} else if parsed_skucode == "CP"{   // search by pod_key
		
		var pod_key string
		var pod_typ string
		var pod_size int
		var pod_flav string
		err1:= db.QueryRow("select * from info_cp where s_key=?", sku_code).Scan(&pod_key,&pod_typ,&pod_size,&pod_flav)
		if err1!= nil{
			panic(err1.Error())	
		}
		
		result_m, err1:= db.Query("select * from info_cp where s_key!=? and pod_type=? and size=?",pod_key,pod_typ, 1)
		if err1!= nil{
			panic(err1.Error())	
		}
		
		for result_m.Next(){
			result_m.Scan(&pod_key,&pod_typ,&pod_size,&pod_flav)	
			temp_pod := pods{
				SKey : pod_key,
				PodType : pod_typ,
				Size : pod_size,
				Flavour : pod_flav,
			}
			cross_sell_pods = append(cross_sell_pods, temp_pod)
		}
		
	} else if parsed_skucode == "EP"{
		
		var pod_key string
		var pod_typ string
		var pod_size int
		var pod_flav string
		err1:= db.QueryRow("select * from info_cp where s_key=?",sku_code).Scan(&pod_key,&pod_typ,&pod_size,&pod_flav)
		if err1!= nil{
			panic(err1.Error())	
		}
		
		result_m, err1:= db.Query("select * from info_cp where s_key!=? and pod_type=? and size=?",pod_key,pod_typ, 3)
		if err1!= nil{
			panic(err1.Error())	
		}
		
		for result_m.Next(){
			result_m.Scan(&pod_key,&pod_typ,&pod_size,&pod_flav)	
			temp_pod := pods{
				SKey : pod_key,
				PodType : pod_typ,
				Size : pod_size,
				Flavour : pod_flav,
			}
			cross_sell_pods = append(cross_sell_pods, temp_pod)
		}
	}
	return cross_sell_pods
}

// common function to return error
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

//common function to return the response in json format
func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)  
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func get_pods(w http.ResponseWriter, r *http.Request){
	
	params := mux.Vars(r)
	sku_code = params["id"]
	sku_code = strings.TrimSuffix(sku_code,"\n")
	fmt.Println("sku_code",sku_code)
	parsed_skucode := sku_code[0:2]
	keys = strings.Fields(sku_code)
	db, err := sql.Open("mysql", "root:@/demo")
	
	if err!= nil{
		panic(err.Error())	
	}
	cross_sell_pods := get_products(parsed_skucode, sku_code,db);

	respondWithJson(w, http.StatusCreated, cross_sell_pods)
}

// our main function
func main() {

	//code for making route
	r := mux.NewRouter()
	r.HandleFunc("/products/{id}", get_pods).Methods("GET")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}


}
