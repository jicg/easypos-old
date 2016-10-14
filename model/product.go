package model

func AddProduct(pro *Product) error {
	_, err := Engine.InsertOne(pro)
	return err
}
