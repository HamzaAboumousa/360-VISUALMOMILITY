package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

type Reseau struct {
	Trajets []Trajet `json:"trajets"`
}

type Trajet struct {
	Name  string `json:"name"`
	Buses []Bus  `json:"buses"`
}

type Bus struct {
	Id         int         `json:"id"`
	CourseID   string      `json:"CourseID"`
	Depart     string      `json:"depart"`
	HeureDepar string      `json:"heure_depar"`
	Frequence  []Frequence `json:"frequence"`
}

type Frequence struct {
	Poid      string `json:"poid"`
	Station   string `json:"station"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type Allbus struct {
	Allbuss []buss
}
type buss struct {
	Nom  string `json:"nom"`
	Data []Info `json:"data"`
}
type Info struct {
	Heure       string       `json:"heure"`
	Coordinates []Coordinate `json:"coordinates"`
}
type Coordinate struct {
	W float64   `json:"w"`
	C []float64 `json:"c"`
}

const Filecord = "T1_T4.xlsx"
const Filefreq = "TAUX CHARGE _ 20230404.xlsx"

func (R Reseau) Fillfr() {
	aa, err := excelize.OpenFile(Filecord)
	if err != nil {
		log.Fatalf("File Not exists: %s", err.Error())
		return
	}
	defer aa.Close()
	// Get value from cell by given worksheet name and axis.

	rows, err := aa.GetRows("Worksheet")
	if err != nil {
		log.Fatalf("RowsRead rows:%s", err.Error())
		return
	}
	for k := 0; k < len(R.Trajets); k++ {
		for i := 0; i < len(R.Trajets[k].Buses); i++ {
			for j := 0; j < len(R.Trajets[k].Buses[i].Frequence); j++ {
				for _, row := range rows {
					a, _ := strconv.Atoi(row[4])
					if row[0] == R.Trajets[k].Name && a == j+1 {

						R.Trajets[k].Buses[i].Frequence[j].Latitude = row[8]
						R.Trajets[k].Buses[i].Frequence[j].Longitude = row[9]
						R.Trajets[k].Buses[i].Frequence[j].Station = row[7]

					}
				}
			}
		}
	}

	fmt.Printf("Read Successfully:%s\n", Filecord)
}

//ReadXlsxFile read given xlsx file using golang
func (R Reseau) ReadXlsxFile() Reseau {
	xlsx, err := excelize.OpenFile(Filefreq)
	if err != nil {
		log.Fatalf("File Not exists: %s", err.Error())
		return Reseau{}
	}
	defer xlsx.Close()
	// Get value from cell by given worksheet name and axis.
	i := 0
	for _, name := range xlsx.GetSheetMap() {
		R.Trajets = append(R.Trajets, Trajet{Name: name})
		rows, err := xlsx.GetRows(name)
		if err != nil {
			log.Fatalf("RowsRead rows:%s", err.Error())
			return Reseau{}
		}
		for j, row := range rows {
			for k, colCell := range row {
				if j != 0 {
					if k == 0 {
						a, _ := strconv.Atoi(colCell)
						R.Trajets[i].Buses = append(R.Trajets[i].Buses, Bus{Id: a})
						continue
					}
					if k == 1 {
						R.Trajets[i].Buses[j-1].CourseID = colCell
						continue
					}
					if k == 2 {
						R.Trajets[i].Buses[j-1].Depart = colCell
						continue
					}
					if k == 3 {
						R.Trajets[i].Buses[j-1].HeureDepar = colCell
						continue
					}
					if k > 3 {
						R.Trajets[i].Buses[j-1].Frequence = append(R.Trajets[i].Buses[j-1].Frequence, Frequence{Poid: colCell})
					}
				}
			}
		}
		i++
	}
	fmt.Printf("Read Successfully:%s\n", Filefreq)
	R.Fillfr()
	return R

}
func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(404)
		w.Write([]byte("Error 404 not found"))
		return
	}

	var R Reseau
	R = R.ReadXlsxFile()
	var Data Allbus
	for i, value := range R.Trajets {
		Data.Allbuss = append(Data.Allbuss, buss{Nom: value.Name})
		for j, v := range value.Buses {
			Data.Allbuss[i].Data = append(Data.Allbuss[i].Data, Info{Heure: v.HeureDepar})
			for k, a := range v.Frequence {
				a.Poid = strings.Replace(a.Poid, ",", ".", 1)
				z, _ := strconv.ParseFloat(a.Poid, 64)
				fmt.Println(z)
				z = z/5.0 + 1.0
				Data.Allbuss[i].Data[j].Coordinates = append(Data.Allbuss[i].Data[j].Coordinates, Coordinate{
					W: z,
				})
				f, _ := strconv.ParseFloat(a.Latitude, 64)
				d, _ := strconv.ParseFloat(a.Longitude, 64)
				Data.Allbuss[i].Data[j].Coordinates[k].C = append(Data.Allbuss[i].Data[j].Coordinates[k].C, f)
				Data.Allbuss[i].Data[j].Coordinates[k].C = append(Data.Allbuss[i].Data[j].Coordinates[k].C, d)
			}
		}
	}
	file, err := json.Marshal(Data)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	_ = ioutil.WriteFile("test.json", file, 0644)
	w.Header().Set("Content-Type", "application/json")
	w.Write(file)

}
func main() {

	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
}
