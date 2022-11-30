package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"personal-web/connection"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()

	connection.DatabaseConnect()

	// Route untuk menginisialisasi folder public
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/project", project).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/project-detail/{id}", projectDetail).Methods("GET")
	route.HandleFunc("/add-project", addProjects).Methods("POST")
	route.HandleFunc("/delete-project/{index}", deleteProjects).Methods("GET")
	route.HandleFunc("/edit-project/{in}", editProject).Methods("GET")
	route.HandleFunc("/edit-project/{in}", formEditProject).Methods("POST")

	fmt.Println("Server sedang berjalan pada port 5000")
	http.ListenAndServe("localhost:5000", route)
}

type Project struct {
	Id              int
	Title           string
	Description     string
	NodeJs          string
	ReactJs         string
	JavaScript      string
	TypeScript      string
	CheckNodeJs     string
	CheckReactJs    string
	CheckJavaScript string
	CheckTypeScript string
}

var projects = []Project{
	{
		Title:       "Bootcamp Dumbways",
		Description: "Membuat Website Bootcamp",
	},
}

func addProjects(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	title := r.PostForm.Get("title")
	description := r.PostForm.Get("description")
	nodejs := r.PostForm.Get("nodejs")
	reactjs := r.PostForm.Get("reactjs")
	javascript := r.PostForm.Get("javascript")
	typescript := r.PostForm.Get("typescript")

	// Menyimpan String Img
	iconnodejs := ""
	iconreactjs := ""
	iconjavascipt := ""
	icontypescript := ""

	checknodejs := ""
	checkreactjs := ""
	checkjavascript := ""
	checktypescript := ""

	// Mengecek kondisi apakah sudah di centang?
	if nodejs == "true" {
		iconnodejs = "/public/img/nodejs.png"
		checknodejs = "checked"
	}
	if reactjs == "true" {
		iconreactjs = "/public/img/reactJs.png"
		checkreactjs = "checked"

	}
	if javascript == "true" {
		iconjavascipt = "/public/img/javaScript.png"
		checkjavascript = "checked"
	}
	if typescript == "true" {
		icontypescript = "/public/img/typeScript.png"
		checktypescript = "checked"
	}

	var newProject = Project{
		Title:           title,
		Description:     description,
		NodeJs:          iconnodejs,
		ReactJs:         iconreactjs,
		JavaScript:      iconjavascipt,
		TypeScript:      icontypescript,
		CheckNodeJs:     checknodejs,
		CheckReactJs:    checkreactjs,
		CheckJavaScript: checkjavascript,
		CheckTypeScript: checktypescript,
	}

	// Untuk Push ke Array projects
	projects = append(projects, newProject)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/index.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	dataProject, errQuery := connection.Conn.Query(context.Background(), "SELECT title, description FROM public.tb_projects")

	if errQuery != nil {
		fmt.Println("Message : " + errQuery.Error())
		return
	}

	var result []Project

	for dataProject.Next() {
		var each = Project{}

		err := dataProject.Scan(&each.Title, &each.Description)
		if err != nil {
			fmt.Println("Message : " + err.Error())
			return
		}
		result = append(result, each)
	}

	resData := map[string]interface{}{
		"Projects": result,
	}

	// dataProject := map[string]interface{}{
	// 	"Projects": projects,
	// }

	tmpt.Execute(w, resData)
}

func projectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/projectDetail.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// w.Write([]byte("Message : " + err.Error()))

	var ProjectDetail = Project{}

	for index, data := range projects {
		if index == id {
			ProjectDetail = Project{
				Title:       data.Title,
				Description: data.Description,
			}
		}
	}
	fmt.Println(ProjectDetail)

	dataDetail := map[string]interface{}{
		"Project": ProjectDetail,
	}

	tmpt.Execute(w, dataDetail)
}

// ---------------------

func project(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/addProject.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}
	tmpt.Execute(w, nil)
}

func editProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/editProject.html")
	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	in, _ := strconv.Atoi(mux.Vars(r)["in"])

	var EditProject = Project{}

	for index, edit := range projects {
		if index == in {
			EditProject = Project{
				Id:              in,
				Title:           edit.Title,
				Description:     edit.Description,
				NodeJs:          edit.NodeJs,
				ReactJs:         edit.ReactJs,
				JavaScript:      edit.JavaScript,
				TypeScript:      edit.TypeScript,
				CheckNodeJs:     edit.CheckNodeJs,
				CheckReactJs:    edit.CheckReactJs,
				CheckJavaScript: edit.CheckJavaScript,
				CheckTypeScript: edit.CheckTypeScript,
			}
		}
	}

	dataEdit := map[string]interface{}{
		"Project": EditProject,
	}

	tmpt.Execute(w, dataEdit)
}
func formEditProject(w http.ResponseWriter, r *http.Request) {
	in, _ := strconv.Atoi(mux.Vars(r)["in"])
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	title := r.PostForm.Get("title")
	description := r.PostForm.Get("description")
	nodejs := r.PostForm.Get("nodejs")
	reactjs := r.PostForm.Get("reactjs")
	javascript := r.PostForm.Get("javascript")
	typescript := r.PostForm.Get("typescript")

	// Menyimpan String Img
	iconnodejs := ""
	iconreactjs := ""
	iconjavascipt := ""
	icontypescript := ""

	checknodejs := ""
	checkreactjs := ""
	checkjavascript := ""
	checktypescript := ""

	// Mengecek kondisi apakah sudah di centang?
	if nodejs == "true" {
		iconnodejs = "/public/img/nodejs.png"
		checknodejs = "checked"
	}
	if reactjs == "true" {
		iconreactjs = "/public/img/reactJs.png"
		checkreactjs = "checked"
	}
	if javascript == "true" {
		iconjavascipt = "/public/img/javaScript.png"
		checkjavascript = "checked"
	}
	if typescript == "true" {
		icontypescript = "/public/img/typeScript.png"
		checktypescript = "checked"
	}

	var newEdit = Project{
		Title:           title,
		Description:     description,
		NodeJs:          iconnodejs,
		ReactJs:         iconreactjs,
		JavaScript:      iconjavascipt,
		TypeScript:      icontypescript,
		CheckNodeJs:     checknodejs,
		CheckReactJs:    checkreactjs,
		CheckJavaScript: checkjavascript,
		CheckTypeScript: checktypescript,
	}

	projects[in] = newEdit
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func deleteProjects(w http.ResponseWriter, r *http.Request) {

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	projects = append(projects[:index], projects[index+1:]...)

	http.Redirect(w, r, "/", http.StatusFound)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/contact.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}
	tmpt.Execute(w, nil)
}
