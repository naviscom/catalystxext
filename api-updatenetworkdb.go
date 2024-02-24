package catalystxext

import (
	// "bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	// "reflect"
	// "strconv"
	// "strings"
	// "os/exec"
	// "strings"
	"time"

	// "github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
)

// type keyResults struct {
// //	rowExcelNew        []string
// 	locnRecUpdatedNew  int
// 	siteRecUpdatedNew  int
// 	cellRecUpdatedNew  int
// 	locnRecInsertedNew int
// 	siteRecInsertedNew int
// 	cellRecInsertedNew int
// }

// func excelSheet2(file *excelize.File, sheet string, c chan keyResults) {
// 	var rowExcel []string
// 	allRowsExcel, err := file.Rows(sheet)
// 	if err != nil {
// 		fmt.Println(err)
// 		// log.Fatal(err)
// 		return
// 	}
// 	rowsProcessed := 0
// 	for allRowsExcel.Next() {
// 		rowExcel = allRowsExcel.Columns()
// 		fmt.Println(rowExcel)
// 		rowsProcessed++
// 	}
// 	close(c)
// }

func  UpdateNetworkDB(ctx *gin.Context) {
	os.Mkdir("/app/downloads", 0700)
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		ctx.Error(err)
		return
	}
	//Open received file
	fileToImport, err := fileHeader.Open()
	if err != nil {
		ctx.Error(err)
		return
	}
	defer fileToImport.Close()

	//Reading the name of received file and creating a new file with the same name
	filenamenew := fileHeader.Filename
	networkDB, err := os.Create("/app/downloads/"+filenamenew)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer networkDB.Close()		
	
	// //Delete temp file after importing
	// defer os.Remove("/app/downloads/"+filenamenew)

	//Write data from received file to the newly created file
	fileBytes, err := io.ReadAll(fileToImport)
	if err != nil {
		ctx.Error(err)
		return
	}
	_, err = networkDB.Write(fileBytes)
	if err != nil {
		ctx.Error(err)
		return
	}
	networkDB.Close()
	// ///////////////////////////
	// f, err := excelize.OpenFile("/app/downloads/"+filenamenew)
	// if err != nil {
	// 	println(err.Error())
	// 	time.Sleep(20 * time.Second)
	// }
	// sheetName2 := f.GetSheetName(2)

	// ch2 := make(chan keyResults)
	// if sheetName2 == "3G" {
	// 	fmt.Println("Found 3G as sheet2: ", sheetName2)
	// 	go excelSheet2(f, sheetName2, ch2)
	// }
	// ch2Count := 0
	// var res2New keyResults
	// for {
	// 	res2, ok2 := <-ch2
	// 	if ok2 {
	// 		ch2Count++
	// 		res2New = res2
	// 		//fmt.Println("Ch2-->", ok2, ch2Count, res2)
	// 	}
	// 	fmt.Println("ch2Count: ", " locnRecUpdatedNew: ", res2New.locnRecUpdatedNew, " siteRecUpdatedNew: ", res2New.siteRecUpdatedNew, " cellRecUpdatedNew: ", res2New.cellRecUpdatedNew)
	// 	fmt.Println("ch2Count: ", " locnRecInsertedNew: ", res2New.locnRecInsertedNew, " siteRecInsertedNew: ", res2New.siteRecInsertedNew, " cellRecInsertedNew: ", res2New.cellRecInsertedNew)
	// 	if !ok2  {
	// 		break
	// 	}
	// }
	//////////////////////////
	ctx.JSON(http.StatusOK, "File "+filenamenew+" uploaded successfully")
	// ctx.JSON(http.StatusOK, "File uploaded successfully")
	time.Sleep(1 * time.Second)
}
