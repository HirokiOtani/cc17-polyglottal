// Sample run-helloworld is a minimal Cloud Run service.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	//Google Cloud Natural Language API
	language "cloud.google.com/go/language/apiv1"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"

	//Connecting to PostgreSQL
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Patient struct {
	gorm.Model
	Name       string `gorm:"not null"`
	Complaints []Complaint
}

type Complaint struct {
	gorm.Model
	Text       string `gorm:"not null"`
	PostedTime time.Time
	PatientID  uint
	Score      float32
	Magnitude  float32
	Word1      string
	Salience1  float32
	Score1     float32
	Magnitude1 float32
	Word2      string
	Salience2  float32
	Score2     float32
	Magnitude2 float32
	Word3      string
	Salience3  float32
	Score3     float32
	Magnitude3 float32
}

//loc, _ := time.LoadLocation("Asia/Tokyo")

var (
	loc, _   = time.LoadLocation("Asia/Tokyo")
	patients = []Patient{
		{Name: "Tarou"},
		{Name: "John"},
	}
	complaints = []Complaint{
		{PostedTime: time.Date(2021, 2, 24, 0, 0, 0, 0, loc), PatientID: 1, Text: "2階から変な音が聞こえます。絶対に自分を監視している人がいるんです。これは幻覚なんかではありません。"},
		{PostedTime: time.Date(2021, 2, 25, 0, 0, 0, 0, loc), PatientID: 1, Text: "私は悪くないです。不当な言いがかりをつけられています。誰もわかってくれません。"},
		{PostedTime: time.Date(2021, 2, 25, 0, 0, 0, 0, loc), PatientID: 1, Text: "すごく調子がいいです。世界が自分のもののような気がします。みんなに感謝しています。ありがとう。先生にもありがとう。"},
		{PostedTime: time.Date(2021, 2, 25, 0, 0, 0, 0, loc), PatientID: 1, Text: "特にないです。"},
		{PostedTime: time.Date(2021, 2, 25, 0, 0, 0, 0, loc), PatientID: 1, Text: "特にないです。"},
		{PostedTime: time.Date(2021, 2, 25, 0, 0, 0, 0, loc), PatientID: 1, Text: "普通です。"},
		{PostedTime: time.Date(2021, 2, 25, 0, 0, 0, 0, loc), PatientID: 1, Text: "できる限り頑張ろうと思います。"},
		{PostedTime: time.Date(2021, 2, 25, 0, 0, 0, 0, loc), PatientID: 1, Text: "この世界の真実を知りました。それは秘密にしないといけないことなので先生にも言えません。言えば刺客が私を殺しにやってくるからです。"},
		{PostedTime: time.Date(2021, 2, 25, 0, 0, 0, 0, loc), PatientID: 1, Text: "夢の中で両親に会いました。何を言っているかわかりませんでした。"},
		{PostedTime: time.Date(2021, 2, 25, 0, 0, 0, 0, loc), PatientID: 1, Text: "こんな状態も悪くないと思います。でも病院に連れて行かれたりするのをやめてほしいです。"},
		{PostedTime: time.Date(2021, 2, 25, 0, 0, 0, 0, loc), PatientID: 1, Text: "私は悪くないです。"},
		{PostedTime: time.Date(2021, 2, 25, 0, 0, 0, 0, loc), PatientID: 1, Text: "普通です。"},
		{PostedTime: time.Date(2021, 2, 25, 0, 0, 0, 0, loc), PatientID: 1, Text: "普通です。"},
		{PostedTime: time.Date(2021, 2, 25, 0, 0, 0, 0, loc), PatientID: 1, Text: "気分がいいです。"},
		{PostedTime: time.Date(2021, 2, 25, 0, 0, 0, 0, loc), PatientID: 1, Text: "幸せってなんでしょうか。"},
		{PostedTime: time.Date(2021, 2, 25, 0, 0, 0, 0, loc), PatientID: 1, Text: "わからないことが増えてきています。"},
		{PostedTime: time.Date(2021, 2, 25, 0, 0, 0, 0, loc), PatientID: 1, Text: "いい加減言いたいことがあります。どうして私を異常だととらえるのですか。私は正常です。むしろ先生のほうがおかしいです。その証拠だってあります。先生は私の言うことじゃなくて母親の言うことを信用しましたよね。本人の言うことじゃなくて家族の言うことを信じるなんて異常だとわかりますよね。許せないので訴訟を起こそうかと思います。"},
		{PostedTime: time.Date(2021, 2, 25, 0, 0, 0, 0, loc), PatientID: 2, Text: "I'm nervous. I hate everything. I want to kill myself."},
		{PostedTime: time.Date(2021, 2, 25, 0, 0, 0, 0, loc), PatientID: 2, Text: "I'm tired, but not so bad."},
	}
)

var db *gorm.DB
var err error

func migration() {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=fukuzawakoki dbname=complaints sslmode=disable")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	//Reflash database & Migration
	db.DropTableIfExists(&Patient{})
	db.DropTableIfExists(&Complaint{})

	db.AutoMigrate(&Patient{})
	db.AutoMigrate(&Complaint{})
	for index := range patients {
		db.Create(&patients[index])
	}
	for index := range complaints {
		entity := createEntitySentiment(complaints[index].Text)
		senti := createSentiment(complaints[index].Text)
		num := len(entity)
		complaints[index].Magnitude = senti.Magnitude
		complaints[index].Score = senti.Score
		if num > 0 {
			complaints[index].Word1 = entity[0].Name
			complaints[index].Salience1 = entity[0].Salience
			complaints[index].Score1 = entity[0].Sentiment.Score
			complaints[index].Magnitude1 = entity[0].Sentiment.Magnitude
		}
		if num > 1 {
			complaints[index].Word2 = entity[1].Name
			complaints[index].Salience2 = entity[1].Salience
			complaints[index].Score2 = entity[1].Sentiment.Score
			complaints[index].Magnitude2 = entity[1].Sentiment.Magnitude
		}
		if num > 2 {
			complaints[index].Word3 = entity[2].Name
			complaints[index].Salience3 = entity[2].Salience
			complaints[index].Score3 = entity[2].Sentiment.Score
			complaints[index].Magnitude3 = entity[2].Sentiment.Magnitude
		}
		db.Create(&complaints[index])
	}
}

func main() {
	//router := mux.NewRouter()
	log.Print("starting migration and seeding...")
	migration()

	log.Print("starting server...")
	//fileServer := http.FileServer(http.Dir("./static"))
	//http.Handle("/", fileServer)
	http.HandleFunc("/", handler)
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/getcomplaints", getComplaintsHandler)
	http.HandleFunc("/getpatients", getPatientsHandler)
	http.HandleFunc("/addcomplaint", postComplaintHandler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	//fmt.Println(createSentiment("2階から変な音が聞こえます。絶対に自分を監視している人がいるんです。これは幻覚なんかではありません。そもそも私は幻覚なんて見たことはありません。何が悪いかと言うと先生が悪いと思います。どうしてかというと私の言うことを信じてくれなくて母の言う事ばっかり信じているからです。私は正常です。早く退院させてください。"))

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	name := os.Getenv("NAME")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello %s!\n", name)
}

type Profile struct {
	Name    string
	Hobbies []string
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	profile := Profile{"Alex", []string{"snowboaring", "programing"}}
	js, err := json.Marshal(profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(js)
}

func getPatientsHandler(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=fukuzawakoki dbname=complaints sslmode=disable")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	var patients []Patient
	db.Find(&patients)
	json.NewEncoder(w).Encode(&patients)
}

func getComplaintsHandler(w http.ResponseWriter, r *http.Request) {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=fukuzawakoki dbname=complaints sslmode=disable")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	var complaints []Complaint
	db.Find(&complaints)
	json.NewEncoder(w).Encode(&complaints)
}

func postComplaintHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	var newComplaint Complaint
	//newComplaint.PatientID = r.FormValue("id")
	newComplaint.PatientID = 2 // authentication needed
	newComplaint.Text = r.FormValue("complaint")
	newComplaint.PostedTime = time.Now()
	entity := createEntitySentiment(newComplaint.Text)
	senti := createSentiment(newComplaint.Text)
	num := len(entity)
	newComplaint.Magnitude = senti.Magnitude
	newComplaint.Score = senti.Score
	newComplaint.Word1 = entity[0].Name
	newComplaint.Salience1 = entity[0].Salience
	newComplaint.Score1 = entity[0].Sentiment.Score
	newComplaint.Magnitude1 = entity[0].Sentiment.Magnitude

	if num > 1 {
		newComplaint.Word2 = entity[1].Name
		newComplaint.Salience2 = entity[1].Salience
		newComplaint.Score2 = entity[1].Sentiment.Score
		newComplaint.Magnitude2 = entity[1].Sentiment.Magnitude
	}
	if num > 2 {
		newComplaint.Word3 = entity[2].Name
		newComplaint.Salience3 = entity[2].Salience
		newComplaint.Score3 = entity[2].Sentiment.Score
		newComplaint.Magnitude3 = entity[2].Sentiment.Magnitude
	}

	db, err = gorm.Open("postgres", "host=localhost port=5432 user=fukuzawakoki dbname=complaints sslmode=disable")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.Create(&newComplaint)
	//create entity and insert them
	//fmt.Fprintf(w, "Name = %s", name)
	//fmt.Fprintf(w, "Complaint = %s", complaint)
}

func createEntitySentiment(text string) []*languagepb.Entity {
	ctx := context.Background()
	c, err := language.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	req := &languagepb.AnalyzeEntitySentimentRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	}
	resp, err := c.AnalyzeEntitySentiment(ctx, req)
	if err != nil {
		log.Fatalf("Failed to analyze text: %v", err)
	}
	return resp.Entities
}

func createSentiment(text string) *languagepb.Sentiment {
	ctx := context.Background()
	c, err := language.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	req := &languagepb.AnalyzeSentimentRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	}
	resp, err := c.AnalyzeSentiment(ctx, req)
	if err != nil {
		log.Fatalf("Failed to analyze text: %v", err)
	}
	return resp.DocumentSentiment
}
