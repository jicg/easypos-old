package pos

// import (
// 	"fmt"
// 	"github.com/alexbrainman/printer"
// 	"io/ioutil"
// 	"log"
// 	"os"
// )

// func TestPrinter() {
// 	name, err := printer.Default()
// 	if err != nil {
// 		fmt.Printf("Default failed: %v", err)
// 	}

// 	p, err := printer.Open(name)
// 	if err != nil {
// 		fmt.Printf("Open failed: %v", err)
// 	}
// 	defer p.Close()

// 	err = p.StartDocument("pos", "RAW")
// 	if err != nil {
// 		fmt.Printf("StartDocument failed: %v", err)
// 	}
// 	defer p.EndDocument()
// 	err = p.StartPage()
// 	if err != nil {
// 		fmt.Printf("StartPage failed: %v", err)
// 	}

// 	var bytedata []byte
// 	bytedata, err = ioutil.ReadFile("log/xorm.log")

// 	fmt.Println(len(bytedata), " ")
// 	if err != nil {
// 		fmt.Printf("file failed: %v", err)
// 	}
// 	_, err = p.Write(bytedata)
// 	if err != nil {
// 		fmt.Printf("rite failed: %v", err)
// 	}
// 	err = p.EndPage()
// 	if err != nil {
// 		fmt.Printf("EndPage failed: %v", err)
// 	}
// }

// func TestReadNames() {
// 	names, err := printer.ReadNames()
// 	if err != nil {
// 		log.Fatalf("ReadNames failed: %v", err)
// 	}
// 	name, err := printer.Default()
// 	if err != nil {
// 		log.Fatalf("Default failed: %v", err)
// 	}
// 	// make sure default printer is listed
// 	for _, v := range names {
// 		if v == name {
// 			return
// 		}
// 	}
// 	log.Fatal("Default printed %q is not listed amongst printers returned by ReadNames %q", name, names)
// }
