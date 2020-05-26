package helpers

import (
	"go.mongodb.org/mongo-driver/bson"
)

// MapStringToBSonD ...
func MapStringToBSonD(args map[string]interface{}) bson.D {
	f := bson.D{}

	for k, val := range args {
		f = append(f, bson.E{k, val})
	}
	return f
}

// MapStringToBSonM ...
func MapStringToBSonM(args map[string]interface{}) bson.M {
	f := bson.M{}

	for k, val := range args {
		f[k] = val
	}
	return f
}

func InterfaceTostruct(i bson.M, m interface{}) error {
	// convert m to s
	bsonBytes, err := bson.Marshal(i)

	if err != nil {
		return err
	}

	bson.Unmarshal(bsonBytes, m)

	return nil
}

func MapToBson(i bson.M, m map[string]interface{}) error {
	// convert m to s
	bsonBytes, err := bson.Marshal(m)

	if err != nil {
		return err
	}

	bson.Unmarshal(bsonBytes, i)

	return nil
}

func MapToBsonM(m map[string]interface{}) (bson.M, error) {

	var i bson.M

	bsonBytes, err := bson.Marshal(m)

	if err != nil {
		return i, err
	}

	bson.Unmarshal(bsonBytes, i)

	return i, nil
}
