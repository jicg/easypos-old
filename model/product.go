package model

func AddProduct(pro *Product) error {
	_, err := Engine.InsertOne(pro)
	return err
}


func AddProductType(pro *Producttype) error {
	_, err := Engine.InsertOne(pro)
	return err
}