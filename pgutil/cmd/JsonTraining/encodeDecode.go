package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
)

type UseAll struct {
	Name    string `json:"username"`
	Surname string `json:"surname"`
	Year    int    `json:"created"`
}

type NoEmpty struct {
	Name    string `json:"username"`
	Surname string `json:"surname"`
	Year    int    `json:"creationyear,omitepmty"`
}

type WithPassword struct {
	Name     string `json:"username"`
	Surname  string `json:"surname,omitepmty"`
	Year     int    `json:"creationyear,omitepmty"`
	Password string `json:"-"`
}

type Data struct {
	Key string `json:"key"`
	Val int    `json:"value"`
}

var DataRecords []Data

func Deserialize(e *json.Decoder, slice interface{}) error {
	return e.Decode(slice)
}

// Serializes slice/fragment of data to json
func Serialize(e *json.Encoder, slice interface{}) error {
	return e.Encode(slice)
}

func JSONstream(data interface{}) (string, error) {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	//encoder.SetIndent("", "\t")

	err := encoder.Encode(data)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func DecodeStream(data interface{}) ([]Data, error) {
	result := []Data{}
	decoder := json.NewDecoder(strings.NewReader(data.(string)))

	//	err := decoder.Decode(data)
	// Reading in loop, assigning if decoding is successfull
	for {
		var d []Data
		if err := decoder.Decode(&d); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		result = d
	}
	return result, nil
}

func test() {
	// Create random records
	var i int
	var t Data
	for i = 0; i < 2; i++ {
		t = Data{
			Key: "JUJ" + string(i),
			Val: i,
		}
		DataRecords = append(DataRecords, t)
	}

	//	fmt.Println("Last record:", t)
	//	_ = PrettyPrint(t)

	val, _ := JSONstream(DataRecords)
	fmt.Printf("%T, %v", val, val)

	check, err := DecodeStream(val)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Type: %T, Value: %v", check, check)
}

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "\t")
	if err == nil {
		fmt.Println(string(b))
	}
	return err
}

func main() {
	test()
	//Marshal
	/*
		useall := UseAll{"John", "Doe", 1908}
		useall2 := UseAll{"Doe", "Hoe", 1958}
		slc := []UseAll{useall, useall2}

		b := new(bytes.Buffer)
		jencoder := json.NewEncoder(b)
		err := Serialize(jencoder, slc)

		fmt.Println(b.String())

		jdecoder := json.NewDecoder(b)
		err = Deserialize(jdecoder, slc)

		fmt.Println(b.String())

		t, err := json.Marshal(&useall)
		if err != nil {
			fmt.Println(err)
			return
		} else {
			fmt.Printf("Value: %s\n", t)
		}
		// Unmarshal
		str := `{"username": "M.", "surname": "TS", "created": 2020}`
		jsonRecord := []byte(str)
		temp := UseAll{}
		err = json.Unmarshal(jsonRecord, &temp)
		if err != nil {
			fmt.Println(err)
			return
		} else {
			fmt.Printf("Data type: %v, Data value: %v\n", temp, temp)
		}
	*/
}
