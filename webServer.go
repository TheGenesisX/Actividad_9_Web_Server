package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

var materias = map[string]informacion{}
var alumnos = map[string]informacion{}

type informacion struct {
	nombre       string // Nombre de materia o alumno, dependiendo del map.
	calificacion float64
}

func obtenerPromedioIndividual(nombreAlumno string) float64 {
	promedio, contadorMaterias := 0.0, 0.0

	for _, value := range alumnos {
		if value.nombre == nombreAlumno {
			promedio += value.calificacion
			contadorMaterias++
		}
	}
	promedio /= contadorMaterias
	return promedio

}

func loadHTML(htmlDocument string) string {
	html, _ := ioutil.ReadFile(htmlDocument)

	return string(html)
}

func index(response http.ResponseWriter, request *http.Request) {
	response.Header().Set(
		"Content-Type",
		"text.html",
	)

	fmt.Fprintf(
		response,
		loadHTML("index.html"),
	)
}

func postReceiver(response http.ResponseWriter, request *http.Request) {
	var message string

	switch request.Method {
	case "POST":
		if err := request.ParseForm(); err != nil {
			fmt.Fprintf(response, "ParseForm() error %v", err)
			return
		}

		aux := request.FormValue("calificacion")
		alumno := request.FormValue("nombreAlumno")
		materia := request.FormValue("materia")

		// Convertimos el string de la calificacion del form en float.
		if calificacion, err := strconv.ParseFloat(aux, 32); err == nil {
			alumnos[alumno] = informacion{nombre: materia, calificacion: calificacion}
			materias[materia] = informacion{nombre: alumno, calificacion: calificacion}
			message = "Registro realizado con exito"

		} else {
			fmt.Println(err)
			return
		}

		response.Header().Set(
			"Content-Type",
			"text.html",
		)

		fmt.Fprintf(
			response,
			loadHTML("postReceiver.html"),
			message,
		)
	}
}

func promedioIndividual(response http.ResponseWriter, request *http.Request) {
	var message string

	switch request.Method {
	case "GET":
		if err := request.ParseForm(); err != nil {
			fmt.Fprintf(response, "ParseForm() error %v", err)
			return
		}

		alumno := request.FormValue("promedioIndividual")
		promedio := obtenerPromedioIndividual(alumno)

		if promedio == -1 {
			message = "Alumno no existente"
		} else {
			message = fmt.Sprintf("%f", promedio)
		}

		response.Header().Set(
			"Content-Type",
			"text.html",
		)

		fmt.Fprintf(
			response,
			loadHTML("promedioIndividual.html"),
			message,
		)
	}
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/postReceiver", postReceiver)
	http.HandleFunc("/promedioIndividual", promedioIndividual)
	fmt.Println("Servidor en ejecucion...")
	http.ListenAndServe(":9000", nil)
}
