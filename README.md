
# 360° VISUALMOMILITY

This project collects data from two XLSX files about traffic in Rouen and returns a JSON database.

## Usage/Examples
First, clone this repository to your laptop using the following command:
```bash
git clone <repository URL>
```
After that, add your trafic.xlsx file, which contains all the traffic information, and the info.xlsx file, which contains information about all the stations with their GPS coordinates, as shown in the example.

Then, run the following command:
```go
go run main.go
```
This command will generate the file `test.json` that contain all the data needed.

Next, open your browser and navigate to the following URL:
```arduino
http://localhost:8080/
```
Finally, you will have access to all the traffic information."





## Tech Stack

**Client:**  HTML, CSS, js

**Server:** GO

