package forms

type CatForm struct {
	Name string `json:"name" bson:"name,omitempty"`
	Age  int    `json:"age" bson:"age,omitempty"`
}
