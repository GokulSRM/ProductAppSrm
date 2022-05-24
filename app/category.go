package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// struct for storing data ,min=3,max=20,required
type category struct {
	CId         string `json:"cid" validate:"required,alphanum,min=4,max=10"`
	Cname       string `json:"cname" validate:"required,min=3,max=20"`
	Cdesc       string `json:"cdesc" validate:"required,min=5,max=100"`
	Ccreatedby  string `json:"ccreatedby" validate:"required,min=3,max=20"`
	Cmodifiedby string `json:"cmodifiedby" validate:"required,min=3,max=20"`
	Cstatus     bool   `json:"cstatus"`
}

var categoryCollection = db().Database("ProductApp").Collection("Category") // get collection "users" from db() which returns *mongo.Client

// Create Category

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	validate := validator.New()

	w.Header().Set("Content-Type", "application/json") // for adding Content-type

	var cat category

	err := json.NewDecoder(r.Body).Decode(&cat) // storing in category variable of type cat
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(400)
	}
	errv := validate.Struct(cat)
	if errv != nil {
		fmt.Println(errv)
		w.WriteHeader(401)
	} else {
		count, errc := categoryCollection.CountDocuments(context.TODO(), bson.D{{"cid", cat.CId}})
		fmt.Println("check count of cid:", count)
		if errc != nil {
			log.Fatal(err)
		} else {
			if count == 0 {
				insertResult, err := categoryCollection.InsertOne(context.TODO(), cat)
				if err != nil {
					log.Fatal(err)
					w.WriteHeader(500)
				} else {
					fmt.Println("Inserted a single document: ", insertResult)
					json.NewEncoder(w).Encode(insertResult.InsertedID) // return the mongodb ID of generated document
					w.WriteHeader(201)
				}
			} else {
				json.NewEncoder(w).Encode("Data redundancy!")
			}

		}

	}

}

// Get Category

func GetCategory(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)["id"] //get Parameter value as string

	// var body category
	// e := json.NewDecoder(r.Body).Decode(&body)
	// if e != nil {
	// 	fmt.Print(e)
	// }
	var result primitive.M //  an unordered representation of a BSON document which is a Map
	filter := bson.M{"cid": params}
	err := categoryCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {

		fmt.Println(err)
		w.WriteHeader(204)
		json.NewEncoder(w).Encode("No data found!")

	} else {
		json.NewEncoder(w).Encode(result) // returns a Map containing document
		w.WriteHeader(200)
	}

}

// Get All Category

func GetAllCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var results []primitive.M                                       //slice for multiple documents
	cur, err := categoryCollection.Find(context.TODO(), bson.D{{}}) //returns a *mongo.Cursor
	if err != nil {

		fmt.Println(err)
		w.WriteHeader(400)

	}
	for cur.Next(context.TODO()) { //Next() gets the next document for corresponding cursor

		var elem primitive.M
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(500)
		}

		results = append(results, elem) // appending document pointed by Next()
	}
	cur.Close(context.TODO()) // close the cursor once stream of documents has exhausted
	json.NewEncoder(w).Encode(results)
	w.WriteHeader(200)

}

//Update Category

func UpdateCategory(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	type updateBody struct {
		CId         string `json:"cid"`         //value that has to be matched
		Cname       string `json:"cname"`       // value that has to be modified
		Cdesc       string `json:"cdesc"`       // value that has to be modified
		Cmodifiedby string `json:"cmodifiedby"` // value that has to be modified
	}
	var body updateBody
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {

		fmt.Print(e)
		w.WriteHeader(400)
	}
	filter := bson.D{{"cid", body.CId}} // converting value to BSON type
	after := options.After              // for returning updated document
	returnOpt := options.FindOneAndUpdateOptions{

		ReturnDocument: &after,
	}
	update := bson.D{{"$set", bson.D{{"cname", body.Cname}, {"cdesc", body.Cdesc}, {"cmodifiedby", body.Cmodifiedby}}}}
	updateResult := categoryCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)

	var result primitive.M
	_ = updateResult.Decode(&result)

	json.NewEncoder(w).Encode(result)
	w.WriteHeader(200)

}

// Update Category Status

func UpdateCategoryStatus(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	type updateBody struct {
		CId     string `json:"cid"`     //value that has to be matched
		Cstatus bool   `json:"cstatus"` // value that has to be modified
	}
	var body updateBody
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {

		fmt.Print(e)
		w.WriteHeader(400)
	}
	filter := bson.D{{"cid", body.CId}} // converting value to BSON type
	after := options.After              // for returning updated document
	returnOpt := options.FindOneAndUpdateOptions{

		ReturnDocument: &after,
	}
	update := bson.D{{"$set", bson.D{{"cstatus", body.Cstatus}}}}
	updateResult := categoryCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)

	var result primitive.M
	_ = updateResult.Decode(&result)

	json.NewEncoder(w).Encode(result)
	w.WriteHeader(200)

}

//Delete Category

func DeleteCategory(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)["id"] //get Parameter value as string
	// _id, err := primitive.ObjectIDFromHex(params) // convert params to mongodb Hex ID
	// if err != nil {
	// 	fmt.Printf(err.Error())
	// }
	opts := options.Delete().SetCollation(&options.Collation{}) // to specify language-specific rules for string comparison, such as rules for lettercase
	res, err := categoryCollection.DeleteOne(context.TODO(), bson.D{{"cid", params}}, opts)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(400)
	} else {
		fmt.Printf("deleted %v documents\n", res.DeletedCount)
		json.NewEncoder(w).Encode(res.DeletedCount) // return number of documents deleted
		w.WriteHeader(200)
	}

}
