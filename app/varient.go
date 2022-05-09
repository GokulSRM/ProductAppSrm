package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/// struct for storing data
type varient struct {
	VId         string `json:"vid"`
	Vname       string `json:"vname"`
	Vdesc       string `json:"vdesc"`
	Vcreatedby  string `json:"vcreatedby"`
	Vmodifiedby string `json:"vmodifiedby"`
	Vstatus     bool   `json:"vstatus"`
}

var varientCollection = db().Database("ProductApp").Collection("Varient") // get collection "users" from db() which returns *mongo.Client

// Create Varient

func CreateVarient(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json") // for adding Content-type

	var varient varient
	err := json.NewDecoder(r.Body).Decode(&varient) // storing in person variable of type user
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(400)
	}
	insertResult, err := varientCollection.InsertOne(context.TODO(), varient)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(500)
	} else {
		fmt.Println("Inserted a single document: ", insertResult)
		json.NewEncoder(w).Encode(insertResult.InsertedID) // return the mongodb ID of generated document
		w.WriteHeader(200)
	}

}

// Get Varient

func GetVarient(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)["id"] //get Parameter value as string
	// var body varient
	// e := json.NewDecoder(r.Body).Decode(&body)
	// if e != nil {

	// 	fmt.Print(e)
	// }
	var result primitive.M //  an unordered representation of a BSON document which is a Map
	filter := bson.M{"vid": params}
	err := varientCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {

		fmt.Println(err)
		json.NewEncoder(w).Encode("No data found!")

	} else {
		json.NewEncoder(w).Encode(result) // returns a Map containing document
		w.WriteHeader(200)
	}

}

// Get All Varient

func GetAllVarient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var results []primitive.M                                      //slice for multiple documents
	cur, err := varientCollection.Find(context.TODO(), bson.D{{}}) //returns a *mongo.Cursor
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

//Update Varient of Varient Id

func UpdateVarient(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	type updateBody struct {
		VId         string `json:"vid"`         //value that has to be matched
		Vname       string `json:"vname"`       // value that has to be modified
		Vdesc       string `json:"vdesc"`       // value that has to be modified
		Vcreatedby  string `json:"vcreatedby"`  // value that has to be modified
		Vmodifiedby string `json:"vmodifiedby"` // value that has to be modified
	}
	var body updateBody
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {

		fmt.Print(e)
		w.WriteHeader(400)
	}
	filter := bson.D{{"vid", body.VId}} // converting value to BSON type
	after := options.After              // for returning updated document
	returnOpt := options.FindOneAndUpdateOptions{

		ReturnDocument: &after,
	}
	update := bson.D{{"$set", bson.D{{"vname", body.Vname}, {"vdesc", body.Vdesc}, {"vmodifiedby", body.Vmodifiedby}}}}
	updateResult := varientCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)

	var result primitive.M
	_ = updateResult.Decode(&result)

	json.NewEncoder(w).Encode(result)
	w.WriteHeader(200)

}

// Update Varient Status

func UpdateVarientStatus(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	type updateBody struct {
		VId     string `json:"vid"`     //value that has to be matched
		Vstatus bool   `json:"vstatus"` // value that has to be modified
	}
	var body updateBody
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {

		fmt.Print(e)
		w.WriteHeader(400)
	}
	filter := bson.D{{"vid", body.VId}} // converting value to BSON type
	after := options.After              // for returning updated document
	returnOpt := options.FindOneAndUpdateOptions{

		ReturnDocument: &after,
	}
	update := bson.D{{"$set", bson.D{{"vstatus", body.Vstatus}}}}
	updateResult := varientCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)

	var result primitive.M
	_ = updateResult.Decode(&result)

	json.NewEncoder(w).Encode(result)
	w.WriteHeader(200)

}

//Delete Varient

func DeleteVarient(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)["id"] //get Parameter value as string

	// _id, err := primitive.ObjectIDFromHex(params) // convert params to mongodb Hex ID
	// if err != nil {
	// 	fmt.Printf(err.Error())
	// }
	opts := options.Delete().SetCollation(&options.Collation{}) // to specify language-specific rules for string comparison, such as rules for lettercase
	res, err := varientCollection.DeleteOne(context.TODO(), bson.D{{"vid", params}}, opts)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(400)
	} else {
		fmt.Printf("deleted %v documents\n", res.DeletedCount)
		json.NewEncoder(w).Encode(res.DeletedCount) // return number of documents deleted
		w.WriteHeader(200)
	}

}
