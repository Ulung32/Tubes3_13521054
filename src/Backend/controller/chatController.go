package controller

import (
	"Backend/models"
	"Backend/regex"
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type message struct{
	Query string
	HistoryId string
	Method string
}


func parseQuery(query string) (string, string) {
    var x, y string
    i := strings.Index(query, "pertanyaan") + len("pertanyaan")
    j := strings.Index(query, "dengan jawaban")
    if i >= 0 && j >= 0 {
        x = strings.TrimSpace(query[i:j])
        y = strings.TrimSpace(query[j+len("dengan jawaban"):])
    }
    return x, y
}


func getTop3Indexes(numbers []float64) []int {
    var indexes []int

    // Sorting slice in descending order
    sort.Slice(numbers, func(i, j int) bool {
        return numbers[i] > numbers[j]
    })

    // Getting the indexes of the top 3 elements
    for i, n := range numbers {
        if i == 3 {
            break
        }
        for j, m := range numbers {
            if n == m {
                indexes = append(indexes, j)
                break
            }
        }
    }

    return indexes
}

func GetAnswers(c echo.Context) error{
	GetQnA()
	
	Message := new(message)

	if err := c.Bind(Message); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	fmt.Println(Message)
	client := models.MongoConnect()
	defer client.Disconnect(context.TODO())

	coll := models.MongoCollection("Chat", client)
	_, err := coll.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hisId, err := primitive.ObjectIDFromHex(Message.HistoryId)

	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}	

	_, errInsert := coll.InsertOne(ctx, models.Chat{
		ID : primitive.NewObjectID(),
		HistoryId: hisId,
		Chat : Message.Query,
		IsBot: false,
	})
	if errInsert != nil {
		return c.JSON(http.StatusBadRequest, "Error inserting chat")
	}

	questions := regex.SplitQuestion(Message.Query)
	for i:=0; i<len(questions); i++ {
		fmt.Println(questions[i])
	}
	
	answers := ""
	for i:=0;i<len(questions);i++ {
		classificationQuery := regex.QueryClassification(questions[i])
		if(classificationQuery == "1"){
			index, similarity := Processor.QuerySearch(strings.ToLower(Message.Method), strings.ToLower(questions[i]))
			if(similarity > 90){
				answers = answers  + Processor.QnAList[index].Answer
			}else{
				answers = answers + "tidak ada pertanyaan dalam database, mungkin maksud anda:\n"
				similarity := Processor.GetSimilarityList(questions[i])
				if(len(Processor.QnAList) > 3){
					top3 := getTop3Indexes(similarity)
					for j:=0; j<3; j++{
						answers = answers + Processor.QnAList[top3[j]].Question	+ "\n"
					}
				}
				answers = answers + "\n"
			}
			
		}else if(classificationQuery == "4"){
			question, answer := parseQuery(questions[i])
			fmt.Println(answer)

			qnacoll := models.MongoCollection("QnA", client)

			index, similarity := Processor.QuerySearch(Message.Method, question)

			if(similarity > 90){
				answers = answers + "Pertanyaan serupa sudah ada di database! sukses mengupdate jawaban menjadi " + answer
				_, cancel := context.WithCancel(context.Background())
				defer cancel()

				fmt.Println(Processor.QnAList[index].ID)
				filter := bson.M{"_id": Processor.QnAList[index].ID}
				update := bson.M{"$set": bson.M{"answer": answer}}
				_, err := qnacoll.UpdateOne(context.Background(), filter, update)
				if err != nil {
					log.Fatal(err)
				}
			}else{
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				_, errInsert := qnacoll.InsertOne(ctx, models.QnA{
					ID: primitive.NewObjectID(),
					Question: question,
					Answer:  answer,
				})

				if errInsert != nil {
					fmt.Println("Error Create QnA")
				}
				answers = answers + "Sukses menambahkan pertanyaan"
			}
		}else if(classificationQuery == "5"){
			deletedStr := strings.Replace(questions[i], "hapus pertanyaan ", "", -1)

			qnacoll := models.MongoCollection("QnA", client)
			index, similarity := Processor.QuerySearch(Message.Method, deletedStr)
			if(similarity < 90){
				answers = answers + "Tidak ada pertanyaan dalam database"
			}else{
				filter := bson.M{"_id": Processor.QnAList[index].ID}
				result, err := qnacoll.DeleteOne(context.Background(), filter)
				if err != nil {
					return c.String(http.StatusInternalServerError, "Error deleting document")
				}

				if result.DeletedCount == 0 {
					return c.String(http.StatusNotFound, "Document not found")
				}
				
				answers = answers + "Sukses menghapus pertanyaan"
			}
		}else{
			answers = answers + classificationQuery
		}
		if(i < len(questions)-1){
			answers = answers + "\n"
		}
		
		
	}
	_, errInsert = coll.InsertOne(ctx, models.Chat{
		ID : primitive.NewObjectID(),
		HistoryId: hisId,
		Chat : answers,
		IsBot: true,
	})
	if errInsert != nil {
		return c.JSON(http.StatusBadRequest, "Error inserting chat")
	}
	return c.JSON(http.StatusOK, answers)
}
func GetChatHistory(c echo.Context) error{
	hisID := c.QueryParam("HisID")

	client := models.MongoConnect()
	defer client.Disconnect(context.TODO())

	coll := models.MongoCollection("Chat", client)
	
	objID, err := primitive.ObjectIDFromHex(hisID)

	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	cursor, err := coll.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	var ChatHis = make([]models.Chat, 0)
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var tempChat models.Chat
		if err := cursor.Decode(&tempChat); err != nil {
			log.Fatal(err)
		}
		if(tempChat.HistoryId == objID){
			ChatHis = append(ChatHis, tempChat)
		}	
	}
	return c.JSON(http.StatusOK, ChatHis)
}