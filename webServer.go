package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

var materias = map[string]map[string]float64{}
var alumnos = map[string]map[string]float64{}

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
			if _, alumnoExists := alumnos[alumno]; alumnoExists {
				if _, materiaExists := alumnos[alumno][materia]; materiaExists {
					message = "Calificacion existente. Imposible modificar."
					// Si existe el alumno, y la materia *para ese alumno*, quiere decir que ya existe calificacion y no es modificable.
				} else {
					// Veamos si la materia ya existe en el mapa; que haya sido generada por el registro de otro alumno.
					if _, materiaExists := materias[materia]; materiaExists {
						alumnos[alumno][materia] = calificacion
						materias[materia][alumno] = calificacion
					} else {
						alumnos[alumno][materia] = calificacion
						materias[materia] = make(map[string]float64)
						materias[materia][alumno] = calificacion
						// Creamos la nueva materia.
					}
					message = "Registro realizado con exito"
				}
			} else {
				// Si no existe el alumno, necesitamos revisar si existe o no la materia.
				if _, materiaExists := materias[materia]; materiaExists {
					alumnos[alumno] = make(map[string]float64)
					alumnos[alumno][materia] = calificacion

					materias[materia][alumno] = calificacion
				} else {
					// Si no existen ninguno de los dos, los creamos y les damos sus valores.
					alumnos[alumno] = make(map[string]float64)
					alumnos[alumno][materia] = calificacion

					materias[materia] = make(map[string]float64)
					materias[materia][alumno] = calificacion
				}
				message = "Registro realizado con exito"
			}
		} else {
			fmt.Println(err)
			return
		}

		fmt.Println("Alumnos:", alumnos)
		fmt.Println("Materias:", materias)

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

func obtenerPromedioIndividual(nombreAlumno string) float64 {
	promedio, contadorMaterias := 0.0, 0.0

	if _, alumnoExists := alumnos[nombreAlumno]; alumnoExists {
		for _, value := range alumnos[nombreAlumno] {
			promedio += value
			contadorMaterias++
		}
		promedio /= contadorMaterias
		return promedio
	}
	return -1
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
			aux := fmt.Sprintf("%f", promedio)
			message = "<tr>" +
				"<td>" + alumno + "</td>" +
				"<td>" + aux + "</td>" +
				"</tr>"
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

func promedioGeneral(response http.ResponseWriter, request *http.Request) {
	var message string
	promedio := 0.0

	if len(alumnos) == 0 {
		message = "No hay alumnos registrados"
	} else {
		for alumno := range alumnos {
			promedio += obtenerPromedioIndividual(alumno)
		}
		promedio /= float64(len(alumnos))
		aux := fmt.Sprintf("%f", promedio)

		message = "<tr>" +
			"<td>" + "General" + "</td>" +
			"<td>" + aux + "</td>" +
			"</tr>"
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

func obtenerPromedioMateria(nombreMateria string) float64 {
	promedio, contadorAlumnos := 0.0, 0.0

	if _, materiaExists := materias[nombreMateria]; materiaExists {
		for _, value := range materias[nombreMateria] {
			promedio += value
			contadorAlumnos++
		}
		promedio /= contadorAlumnos
		return promedio
	}
	return -1
}

func promedioMateria(response http.ResponseWriter, request *http.Request) {
	var message string

	objetivo := request.FormValue("promedioMateria")
	promedio := obtenerPromedioMateria(objetivo)

	if promedio == -1 {
		message = "Materia no existente"
	} else {
		aux := fmt.Sprintf("%f", promedio)

		message = "<tr>" +
			"<td>" + objetivo + "</td>" +
			"<td>" + aux + "</td>" +
			"</tr>"
	}

	response.Header().Set(
		"Content-Type",
		"text.html",
	)

	fmt.Fprintf(
		response,
		loadHTML("promedioMateria.html"),
		message,
	)
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/postReceiver", postReceiver)
	http.HandleFunc("/promedioIndividual", promedioIndividual)
	http.HandleFunc("/promedioGeneral", promedioGeneral)
	http.HandleFunc("/promedioMateria", promedioMateria)
	fmt.Println("Servidor en ejecucion...")
	http.ListenAndServe(":9000", nil)
}
