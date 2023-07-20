//API keyfor weather api-->37440b76c2e61030c97bdafcb82dd091

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	//"sort"

	"github.com/gorilla/mux"
)

type Temp struct {

	//T is the temperature in kelvin
	Temperature float64 `json:"temp"`
}

type weather struct {

	//main in structure jiske andar temperature hai
	Main Temp `json:"main"`
}

const apiKey = "37440b76c2e61030c97bdafcb82dd091"

var wg sync.WaitGroup

func ApiHandler(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// cityInfo := vars["city"]
	param := r.URL.Query()

	// city := []string{"mohali","chandigarh","mumbai"} // will be provided by sir
	city := []string{"Mumbai", "Delhi", "Bangalore", "Hyderabad", "Ahmedabad", "Chennai", "Kolkata", "Surat", "Pune", "Jaipur", "Lucknow", "Kanpur", "Nagpur", "Visakhapatnam", "Indore", "Thane", "Bhopal", "Pimpri-Chinchwad", "Patna", "Vadodara", "Ghaziabad", "Ludhiana", "Agra", "Nashik", "Faridabad", "Meerut", "Rajkot", "Varanasi", "Srinagar", "Aurangabad", "Dhanbad", "Amritsar", "Navi Mumbai", "Allahabad", "Howrah", "Ranchi", "Gwalior", "Jabalpur", "Coimbatore", "Vijayawada", "Jodhpur", "Madurai", "Raipur", "Kota", "Chandigarh", "Guwahati", "Solapur", "Hubli-Dharwad", "Mysore"}

	for _, v := range param {

		city = append(city, v[0])
	}
	var url string
	//cityTEMPS map for storing temperature of different cities

	cityTEMPS := make(map[string]float64)

	
	for i,_ := range city {

		wg.Add(1)
		fmt.Println("IN FOR LOOP ")
		i:=i
		go func() {

			fmt.Println("IN GO FUNC ")
			defer wg.Done()

			//fmt.Println("information of ", val)
			url = fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city[i], apiKey)
			//fmt.Println("url : ", url)
			// get request
			res, err := http.Get(url)
			if err != nil {
				http.Error(w, "error querrying the url "+url, http.StatusInternalServerError)
			}
			//defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				http.Error(w, "error reading response", http.StatusInternalServerError)
			}
			// fmt.Println("data : ", data)
			var temp weather

			json.Unmarshal(data, &temp)
			cityTEMPS[city[i]] = temp.Main.Temperature

		}()

		
	}
	wg.Wait()

	fmt.Println("cityTEMPS:", cityTEMPS)

	type particularCity struct {
		CITY string  `json:"CITY"`
		TEMP float64 `json:"TEMP"`
		Diff string  `json:"diff"`
	}
	var bigMap = make(map[string]particularCity)

	for i, v := range cityTEMPS {

		var diffString string
		for j, val := range cityTEMPS {

			if j == i {
				continue
			}

			diff := v - val
			if diff < 0 {

				diff = -diff
				var formattteddiffString = fmt.Sprintf(" %s is cooler than %s by %f degrees ", i, j, diff)

				diffString = diffString + "," + formattteddiffString

			} else {
				var formattteddiffString = fmt.Sprintf(" %s is hotter than  %s by  %f degrees ", i, j, diff)

				diffString = diffString + formattteddiffString

			}

			//fmt.Println("byteddiffstring",formattteddiffString)

		}

		var citiStruct = particularCity{i, v, diffString}

		//create BIG map

		// w.Write([]byte(citiStruct))
		// jsoncitiStrcut,_:=json.Marshal(citiStruct)
		// // fmt.Fprintf(w,"%s",citiStruct.CITY)
		//  w.Write(jsoncitiStrcut)

		bigMap[citiStruct.CITY] = citiStruct

	}
	bytemap, _ := json.Marshal(bigMap)
	w.Write(bytemap)

	w.Header().Set("Content-Type", "application/json")

}

func main() {
	fmt.Println("Working with weather api request ....")

	router := mux.NewRouter()

	router.HandleFunc("/weather/api", ApiHandler).Methods("GET")
	// router.HandleFunc("/weather/{city}", ApiHandlerChan).Methods("GET")

	http.ListenAndServe(":8888", router)

}
