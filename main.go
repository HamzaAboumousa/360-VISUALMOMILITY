package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/xuri/excelize/v2"
)

type Reseau struct {
	Alltrajet []Trajet
	fileName  string
}
type Trajet struct {
	Allbus      []Bus
	Name_trajet string
}
type Bus struct {
	Id         int
	Course_id  string
	Depart     string
	Heur_depar string
	Fr         []Frequence
}
type Frequence struct {
	La      string
	Lo      string
	Poid    string
	Station string
}

const Filecord = "stop_times_-_Copie.xlsx"
const Filefreq = "TAUX CHARGE _ 20230404.xlsx"

func (R *Reseau) Fillfr() {
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
	for k := 0; k < len(R.Alltrajet); k++ {
		for i := 0; i < len(R.Alltrajet[k].Allbus); i++ {
			for j := 0; j < len(R.Alltrajet[k].Allbus[i].Fr); j++ {
				for _, row := range rows {
					a, _ := strconv.Atoi(row[4])
					if row[0] == R.Alltrajet[k].Name_trajet && a == j+1 {

						R.Alltrajet[k].Allbus[i].Fr[j].La = row[8]
						R.Alltrajet[k].Allbus[i].Fr[j].Lo = row[9]
						R.Alltrajet[k].Allbus[i].Fr[j].Station = row[7]

					}
				}
			}
		}
	}

	fmt.Printf("Read Successfully:%s\n", Filecord)
}

//ReadXlsxFile read given xlsx file using golang
func (R *Reseau) ReadXlsxFile() {
	xlsx, err := excelize.OpenFile(Filefreq)
	if err != nil {
		log.Fatalf("File Not exists: %s", err.Error())
		return
	}
	defer xlsx.Close()
	// Get value from cell by given worksheet name and axis.
	i := 0
	for _, name := range xlsx.GetSheetMap() {
		R.Alltrajet = append(R.Alltrajet, Trajet{Name_trajet: name})
		rows, err := xlsx.GetRows(name)
		if err != nil {
			log.Fatalf("RowsRead rows:%s", err.Error())
			return
		}
		for j, row := range rows {
			for k, colCell := range row {
				if j != 0 {
					if k == 0 {
						a, _ := strconv.Atoi(colCell)
						R.Alltrajet[i].Allbus = append(R.Alltrajet[i].Allbus, Bus{Id: a})
						continue
					}
					if k == 1 {
						R.Alltrajet[i].Allbus[j-1].Course_id = colCell
						continue
					}
					if k == 2 {
						R.Alltrajet[i].Allbus[j-1].Depart = colCell
						continue
					}
					if k == 3 {
						R.Alltrajet[i].Allbus[j-1].Heur_depar = colCell
						continue
					}
					if k > 3 {
						R.Alltrajet[i].Allbus[j-1].Fr = append(R.Alltrajet[i].Allbus[j-1].Fr, Frequence{Poid: colCell})
					}
				}
			}
		}
		i++
	}
	fmt.Printf("Read Successfully:%s\n", Filefreq)
	R.Fillfr()

}
func main() {
	var r *Reseau
	r.ReadXlsxFile()
}
