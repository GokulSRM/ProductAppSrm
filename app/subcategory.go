package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// struct for storing data
type subcategory struct {
	ScId         string `json:"scid"`
	CId          string `json:"cid"`
	Scname       string `json:"scname"`
	Scdesc       string `json:"scdesc"`
	Sccreatedby  string `json:"sccreatedby"`
	Scmodifiedby string `json:"scmodifiedby"`
	Scstatus     bool   `json:"scstatus"`
}

var subCategoryCollection = db().Database("ProductApp").Collection("Subcategory") // get collection "users" from db() which returns *mongo.Client

// Create SubCategory

func CreateSubCategory(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json") // for adding Content-type

	var subc subcategory
	err := json.NewDecoder(r.Body).Decode(&subc) // storing in person variable of type user
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(400)
	}

	ccount, errc := categoryCollection.CountDocuments(context.TODO(), bson.D{{"cid", subc.CId}})
	scount, errs := subCategoryCollection.CountDocuments(context.TODO(), bson.D{{"scid", subc.ScId}})

	if errc != nil || errs != nil {
		log.Fatal(errc)
		log.Fatal(errs)
	} else {
		if ccount == 0 {
			// w.WriteHeader(204)
			json.NewEncoder(w).Encode("Invalid Category Id")

		} else if scount > 0 {
			fmt.Print(scount)
			json.NewEncoder(w).Encode("Duplicate Subcategory!")
		} else if scount == 0 {
			insertResult, err := subCategoryCollection.InsertOne(context.TODO(), subc)
			if err != nil {
				log.Fatal(err)
				w.WriteHeader(500)
			} else {
				fmt.Println("Inserted a single document: ", insertResult)
				json.NewEncoder(w).Encode(insertResult.InsertedID) // return the mongodb ID of generated document
				w.WriteHeader(200)
			}
		} else {
			json.NewEncoder(w).Encode("Invalid Subcategory!")
		}
	}

}

// Get SubCategory

func GetSubCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)["id"] //get Parameter value as string
	is_alphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(params)

	// var body subcategory
	// e := json.NewDecoder(r.Body).Decode(&body)
	// if e != nil {

	// 	fmt.Print(e)
	// }
	if is_alphanumeric {
		var result primitive.M //  an unordered representation of a BSON document which is a Map
		filter := bson.M{"scid": params}
		err := subCategoryCollection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {

			fmt.Println(err)
			json.NewEncoder(w).Encode("No data found!")

		} else {
			json.NewEncoder(w).Encode(result) // returns a Map containing document
			w.WriteHeader(200)
		}
	} else {
		json.NewEncoder(w).Encode("Please enter the correct SubCategory ID!")
		w.WriteHeader(204)
	}
}

// Get All SubCategory
func GetAllSubCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var results []primitive.M                                          //slice for multiple documents
	cur, err := subCategoryCollection.Find(context.TODO(), bson.D{{}}) //returns a *mongo.Cursor
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

//Update SubCategory

func UpdateSubCategory(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	validate := validator.New()

	type updateBody struct {
		ScId         string `json:"scid" validate:"required,alphanum,min=4,max=10"` // value that has to be matched
		CId          string `json:"cid" validate:"required,alphanum,min=4,max=10"`  // value that has to be modified
		Scname       string `json:"scname" validate:"required,min=3,max=20"`        // value that has to be modified
		Scdesc       string `json:"scdesc" validate:"required,min=3,max=20"`        // value that has to be modified
		Scmodifiedby string `json:"scmodifiedby" validate:"required,min=3,max=20"`  // value that has to be modified
	}
	var body updateBody
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {

		fmt.Print(e)
		w.WriteHeader(400)
	}
	errv := validate.Struct(body) // update struct validation
	if errv != nil {
		fmt.Println(errv)
		w.WriteHeader(401)
		json.NewEncoder(w).Encode("Validation error")
	} else {
		filter := bson.D{{"scid", body.ScId}} // converting value to BSON type
		after := options.After                // for returning updated document
		returnOpt := options.FindOneAndUpdateOptions{

			ReturnDocument: &after,
		}
		update := bson.D{{"$set", bson.D{{"cid", body.CId}, {"scname", body.Scname}, {"scdesc", body.Scdesc}, {"scmodifiedby", body.Scmodifiedby}}}}
		updateResult := subCategoryCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)

		var result primitive.M
		_ = updateResult.Decode(&result)

		json.NewEncoder(w).Encode(result)
		w.WriteHeader(200)
	}

}

// Update SubCategory Status

func UpdateSubCategoryStatus(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	validate := validator.New()

	type updateBodystatus struct {
		ScId     string `json:"scid" validate:"required,alphanum,min=4,max=10"` //value that has to be matched
		Scstatus bool   `json:"scstatus"`                                       // value that has to be modified
	}
	var bodys updateBodystatus
	e := json.NewDecoder(r.Body).Decode(&bodys)
	if e != nil {

		fmt.Print(e)
		w.WriteHeader(400)
	}
	errv := validate.Struct(bodys) // update status struct validation
	if errv != nil {
		fmt.Println(errv)
		w.WriteHeader(401)
		json.NewEncoder(w).Encode("Validation error")
	} else {
		filter1 := bson.D{{"scid", bodys.ScId}} // converting value to BSON type
		after1 := options.After                 // for returning updated document
		returnOpt1 := options.FindOneAndUpdateOptions{

			ReturnDocument: &after1,
		}
		update1 := bson.D{{"$set", bson.D{{"scstatus", bodys.Scstatus}}}}
		updateResult1 := subCategoryCollection.FindOneAndUpdate(context.TODO(), filter1, update1, &returnOpt1)

		var result1 primitive.M
		_ = updateResult1.Decode(&result1)

		json.NewEncoder(w).Encode(result1)
		w.WriteHeader(200)
	}
}

//Delete SubCategory

func DeleteSubCategory(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)["id"] //get Parameter value as string
	is_alphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(params)

	// _id, err := primitive.ObjectIDFromHex(params) // convert params to mongodb Hex ID
	// if err != nil {
	// 	fmt.Printf(err.Error())
	// }
	if is_alphanumeric {
		opts := options.Delete().SetCollation(&options.Collation{}) // to specify language-specific rules for string comparison, such as rules for lettercase
		res, err := subCategoryCollection.DeleteOne(context.TODO(), bson.D{{"scid", params}}, opts)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(400)
		} else {
			fmt.Printf("deleted %v documents\n", res.DeletedCount)
			json.NewEncoder(w).Encode(res.DeletedCount) // return number of documents deleted
			w.WriteHeader(200)
		}
	} else {
		json.NewEncoder(w).Encode("Please enter the correct SubCategory ID!")
		w.WriteHeader(204)
	}

}
