package main

import (
	"database/sql"
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"

	"html/template"
	"net/http"

	_ "github.com/lib/pq"
)

var tpl *template.Template
var id int

func init() {
	tpl = template.Must(template.ParseGlob(".gohtml"))
}

func main() {

	log.Info("Connecting to SQL DB ...")
	http.HandleFunc("/", index)
	http.HandleFunc("/recetas", recetasf)
	http.HandleFunc("/ingredientes", ingredientesf)
	http.HandleFunc("/calcular", calcular)
	http.HandleFunc("/calculardos", calculardos)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func recetasf(w http.ResponseWriter, r *http.Request) {

	tpl.ExecuteTemplate(w, "recetas.gohtml", nil)
}

func processorRf(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		fmt.Printf(r.Method)
		return
	}

	connStr := "user=victoraguila dbname=victoraguila host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	nReceta := r.FormValue("recetaN")
	nIngrediente := r.FormValue("ingredienteN")
	cIngrediente := r.FormValue("cantidadI")

	qr, err := db.Query(`INSERT INTO public.receta(nombre) VALUES ($1)`, nReceta)

	

	if err != nil {
		log.Warn(err)
	}

	rows, err := db.Query(`SELECT id FROM receta WHERE nombre like $1`, nReceta)

	for rows.Next() {
		rows.Scan(&id)
		fmt.Printf("dice: %v \n", id)
	}

	db.Query(`INSERT INTO public.ingredientes(name, id_receta, cantidad) VALUES ($1, $2, $3)`, nIngrediente, id, cIngrediente)

	tpl.ExecuteTemplate(w, "ingredientes.gohtml", nil)

}

func menuf(w http.ResponseWriter, r *http.Request) {

	connStr := "user=victoraguila dbname=victoraguila host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	var nom string
	var aid int

	rows, err := db.Query(`SELECT * FROM receta`)

	var b [15]string
	var conter int

	for rows.Next() {
		rows.Scan(&aid, &nom)
		b[conter] = nom
		conter++

	}

	html := `<!DOCTYPE html>
	<html>
	<body>
	<div id="marcom">
	<form method = "GET" action = "/processorM"> 
	Ingrese un nombre para este menu <input type="text" name="menuN"><br>
	<select name="LUNES"> 
	<option selected disabled>LUNES</option>
	{{range $value := .}}
    <option value="{{ $value }}">{{ $value }}</option>
   	{{end}}
   </select>
   <select name="MARTES"> 
   <option selected disabled>MARTES</option>
	{{range $value := .}}
    <option value="{{ $value }}">{{ $value }}</option>
   	{{end}}
   </select>
   <select name="miercoles"> 
   <option selected disabled>Miercoles</option>
	{{range $value := .}}
    <option value="{{ $value }}">{{ $value }}</option>
   	{{end}}
   </select>
   <select name="jueves"> 
   <option selected disabled>Jueves</option>
	{{range $value := .}}
    <option value="{{ $value }}">{{ $value }}</option>
   	{{end}}
   </select>
   <select name="viernes"> 
   <option selected disabled>Viernes</option>
	{{range $value := .}}
    <option value="{{ $value }}">{{ $value }}</option>
   	{{end}}
   </select>
   <input id="botonx" type="submit" value="Submit">
   </form>
   </div>
	</body>
	<style>
		#botonx{
			background-color: rgba(255,255,255,0.6);
			padding: 20px 20px;
			margin-top: 80px;
			border-radius: 20px;
			font-size: 25px;
			margin-left: 310px;
		}
		body{
			background-color:#053752 ;
		}
		#slogan{
			color:rgba(255,255,255,0.6);
			font-size: 40px;
		}
		#marcom{
			padding-top: 30px;
			background-color: rgba(255,255,255,0.6);
			width: 800px;
			height: 400px;
			position: relative;
			left: 550px;
			top: 100px;  
			border-radius: 20px;
		}
		input[type=text] {
		width: 100%;
		padding: 12px 20px;
		margin: 8px 0;
		box-sizing: border-box;
		padding-left: 10px;
		background-color: #E2F0FF;
			
		}
		select{
			background-color:  rgba(255,255,255,0.6);
			border-radius: 10px;
			padding: 20px 30px;
			font-size: 18px;
			margin-left: 8px;
			margin-top: 30px;
		}
</style>
	</html>`

	dropdownTemplate, err := template.New("dropdownexample").Parse(string(html))

	dropdownTemplate.Execute(w, b)
	
}

func ingredientesf(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		fmt.Printf(r.Method)
		return
	}

	connStr := "user=victoraguila dbname=victoraguila host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	nIngrediente := r.FormValue("ingredienteN")
	cIngrediente := r.FormValue("cantidadI")
	db.Query(`INSERT INTO public.ingredientes(name, id_receta, cantidad) VALUES ($1, $2, $3)`, nIngrediente, id, cIngrediente)
	tpl.ExecuteTemplate(w, "ingredientes.gohtml", nil)
}

func processorM(w http.ResponseWriter, r *http.Request) {

	connStr := "user=victoraguila dbname=victoraguila host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	menu := r.FormValue("menuN")

	LUNES := r.FormValue("LUNES")
	MARTES := r.FormValue("MARTES")
	miercoles := r.FormValue("miercoles")
	jueves := r.FormValue("jueves")
	viernes := r.FormValue("viernes")

	db.Query(`INSERT INTO public.menu(nombre) VALUES ($1)`, menu)
	db.Query(`INSERT INTO public.menu_dias(menu, dia, receta) VALUES ($1, 'LUNES', $2)`, menu, LUNES)
	db.Query(`INSERT INTO public.menu_dias(menu, dia, receta ) VALUES ($1, 'MARTES', $2)`, menu, MARTES)
	db.Query(`INSERT INTO public.menu_dias(menu, dia, receta) VALUES ($1, 'miercoles', $2)`, menu, miercoles)
	db.Query(`INSERT INTO public.menu_dias(menu, dia, receta) VALUES ($1, 'jueves', $2)`, menu, jueves)
	db.Query(`INSERT INTO public.menu_dias(menu, dia, receta) VALUES ($1, 'viernes', $2)`, menu, viernes)

	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func calcular(w http.ResponseWriter, r *http.Request) {
	connStr := "user=victoraguila dbname=victoraguila host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	var nom string
	var aid int

	rows, err := db.Query(`SELECT * FROM menu`)

	var b [5]string
	var conter int

	for rows.Next() {
		rows.Scan(&aid, &nom)
		b[conter] = nom
		conter++

	}

	html := `<!DOCTYPE html>
	<html>
	<body>
	<div id="marcoin">
	<form method = "GET" action = "/calculardos"> 
	Ingrese el numero de personas que comenran en la semana <input type="text" name="personas"><br>
	<select name="menug"> 
	<option selected disabled>Elija una opcion</option>
	{{range $value := .}}
    <option value="{{ $value }}">{{ $value }}</option>
   	{{end}}
   <input id="botonx" type="submit" value="Submit">
   </form>
   </div>
	</body>
	<style>
		#botonx{
			background-color: rgba(255,255,255,0.6);
			padding: 20px 20px;
			margin-top: 20px;
			border-radius: 20px;
			font-size: 25px;
			margin-left: 310px;
		}
		body{
			background-color:#053752 ;
		}
		#slogan{
			color:rgba(255,255,255,0.6);
			font-size: 40px;
		}
		#marcoin{
			padding-top: 30px;
			background-color: rgba(255,255,255,0.6);
			width: 800px;
			height: 200px;
			position: relative;
			left: 550px;
			top: 100px;  
			border-radius: 20px;
		}
		input[type=text] {
		width: 100%;
		padding: 12px 20px;
		margin: 8px 0;
		box-sizing: border-box;
		padding-left: 10px;
		background-color: #E2F0FF;
			
		}
		select{
			background-color:  rgba(255,255,255,0.6);
			border-radius: 10px;
			padding: 20px 30px;
			font-size: 18px;
			margin-left: 8px;
			margin-top: 30px;
		}
	</style>
	</html>`
	dropdownTemplate, err := template.New("dropdownexample").Parse(string(html))

	dropdownTemplate.Execute(w, b)
	//tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func calculardos(w http.ResponseWriter, r *http.Request) {

	connStr := "user=victoraguila dbname=victoraguila host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	nPersonas := r.FormValue("personas")
	menug := r.FormValue("menug")

	i1, err := strconv.Atoi(nPersonas)
	if err == nil {
		fmt.Println(i1)
	}

	rows, err := db.Query(`SELECT * FROM menu_dias where menu like $1`, menug)

	var m [15]string
	var contero int
	var aiditemp int
	var aidi int
	var oide int
	var menup, diap, recetap string
	var namou string

	var uno int
	var dos string
	var tres int
	var cuatro int

	var bar int
	//var multiplier int

	//multiplier = 4

	for rows.Next() {
		rows.Scan(&oide, &menup, &diap, &recetap)
		m[contero] = recetap
		contero++
	}

	fmt.Printf("------- LISTA DE COMPRAS ------- \n")

	for i := 0; i < 5; i++ {
		fmt.Printf("Receta: %v \n", m[i])
		rowset, err := db.Query(`SELECT * FROM receta where nombre like $1`, m[i])
		for rowset.Next() {
			rowset.Scan(&aidi, &namou)
			aiditemp = aidi

		}
		rowso, err := db.Query(`SELECT * FROM ingredientes where id_receta = $1`, aiditemp)
		for rowso.Next() {
			rowso.Scan(&uno, &dos, &tres, &cuatro)

			bar = cuatro * i1

			fmt.Printf("%v: %v \n", dos, bar)
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("------------------------ \n")

	}

	tpl.ExecuteTemplate(w, "index.gohtml", nil)

	//rows, err := db.Query(`SELECT * FROM menu`)

}
