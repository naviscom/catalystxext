package catalystxext

import (
	// "bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	// "os/exec"

	// "strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
)



type keyResults struct {
	rowExcelNew        []string
	locnRecUpdatedNew  int
	siteRecUpdatedNew  int
	cellRecUpdatedNew  int
	locnRecInsertedNew int
	siteRecInsertedNew int
	cellRecInsertedNew int
}



func excelSheet2(file *excelize.File, sheet string, c chan keyResults) {
	// var selectRows *sql.Rows
	// var sel, cmd, MASTER_REC_FLAG, TLM1_SYS_ID, RSM1_SYS_ID, RCM1_SYS_ID string                               //,
	var locnRecUpdated, siteRecUpdated, cellRecUpdated, locnRecInserted, siteRecInserted, cellRecInserted int //,  int
	// var rowsAffected int64
	var rowExcel []string
	//db := oracondb12cr.Odmart()
	allRowsExcel, err := file.Rows(sheet)
	if err != nil {
		fmt.Println(err)
		// log.Fatal(err)
		return
	}
	rowsProcessed := 0
	header := true
	// logfile, err := os.OpenFile("/home/ubuntu/catalyst/log/sheet2processed.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	log.Println(err)
	// }
	// defer logfile.Close()
	// logger := log.New(logfile, "prefix: ", log.LstdFlags)
	locnRecUpdated, locnRecInserted, siteRecUpdated, siteRecInserted, cellRecUpdated, cellRecInserted = 0, 0, 0, 0, 0, 0
	for allRowsExcel.Next() {
		if header {
			rowExcel = allRowsExcel.Columns()
			fmt.Print(rowExcel[0], "    ", rowExcel[54])
			if strings.ToLower(rowExcel[0]) == "cell ids_optima" && strings.ToLower(rowExcel[54]) == "monthly sites" {
				fmt.Print("File Format Validated for Sheet 2")
				// logger.Print("File Format Validated for Sheet 2: ", "\r\n")
				time.Sleep(10 * time.Second)
			} else {
				fmt.Print("!!!!!!!!!!!!File Format Not Valid for Sheet 2!!!!!!!!!!!!")
				// logger.Print("File Format Not Valid for Sheet 2: ", "\r\n")
				time.Sleep(10 * time.Second)
				break
			}
		}

		rowsProcessed++
		if !header {
			rowExcel = allRowsExcel.Columns()
			if len(rowExcel[46]) != 5 {
				// logger.Print("Erroneous Property_ID", rowExcel, "\r\n")
				continue
			}
			if len(rowExcel[31]) == 0 {
				// logger.Print("No Latitude, Replaced with 23.600000", rowExcel, "\r\n")
				// rowExcel[31] = "23.600000"
				continue

			}
			if len(rowExcel[32]) == 0 {
				// logger.Print("No Longitude, Replaced with 58.500000", rowExcel, "\r\n")
				// rowExcel[32] = "58.500000"
				continue
			}
			if rowExcel[46] == "M0011" {
				// logger.Print("Erroneous Latitude, Longitude, Replaced with 23.613703, 58.503529", rowExcel, "\r\n")
				rowExcel[31] = "23.613703"
				rowExcel[32] = "58.503529"
			}
			v_lat, _ := strconv.ParseFloat(rowExcel[31], 32)
			v_lon, _ := strconv.ParseFloat(rowExcel[32], 32)
			v_lat_truncated := float64(int(v_lat*1000000)) / 1000000
			v_lon_truncated := float64(int(v_lon*1000000)) / 1000000
			fmt.Println(reflect.TypeOf(v_lat_truncated), v_lat_truncated, reflect.TypeOf(v_lon_truncated), v_lon_truncated)
			//v_lat_truncated, v_lon_truncated
			for index := range rowExcel {
				rowExcel[index] = strings.ReplaceAll(rowExcel[index], ",", "-")
				if len(rowExcel[index]) == 0 {
					rowExcel[index] = "-"
				}
				if index == 5 {
					if len(rowExcel[index]) < 7 {
						// logger.Print("Erroneous NodeB_Name_New: ", rowExcel[index], "\r\n")
						fmt.Println("Erroneous NodeB_Name_New: ", rowExcel[index], "\r\n")
					}
					runes := []rune(rowExcel[index])
					v_NodeB_ID_New := string(runes[:1])
					if v_NodeB_ID_New == "B" {
						rowExcel[index] = "U" + string(runes[1:])
					}
				}
				if index == 30 {
					rowExcel[index] = models.ParseSiteType(rowExcel[index], sheet, rowsProcessed)
					//fmt.Print(rowsProcessed," ",rowExcel[index-29]," excelSheet2 SiteType--->", rowExcel[index])
				}
				if index == 34 {
					rowExcel[index] = models.ParseRegion(rowExcel[index], sheet, rowsProcessed)
					//fmt.Print(" excelSheet2 Region--->",rowExcel[index])
				}
				if index == 35 {
					rowExcel[index] = models.ParseWillayat(rowExcel[index], sheet, rowsProcessed)
					//fmt.Print(" excelSheet2 Willayat--->",rowExcel[index])
				}
				if index == 37 {
					rowExcel[index] = models.ParseArea(rowExcel[index], sheet, rowsProcessed)
					//fmt.Print(" excelSheet2 Area--->",rowExcel[index])
				}
				if index == 40 {
					rowExcel[index] = models.ParseClutter(rowExcel[index], sheet, rowsProcessed)
					//fmt.Print(" excelSheet2 Clutter--->",rowExcel[index])
				}
				if index == 41 {
					rowExcel[index] = models.ParseVendor(rowExcel[index], sheet, rowsProcessed)
					//fmt.Print(" excelSheet2 Vendor--->",rowExcel[index])
				}
				if index == 43 {
					rowExcel[index] = models.ParseBand(rowExcel[index], rowExcel[index+4], sheet, rowsProcessed)
					//fmt.Print(" excelSheet2 Band--->",rowExcel[index])
				}
				if index == 44 {
					rowExcel[index] = models.ParseCarrier(rowExcel[index], rowExcel[index-1], rowExcel[index+3], sheet, rowsProcessed)
					//fmt.Print(" excelSheet2 Carrier--->",rowExcel[index])
				}
				if index == 47 {
					rowExcel[index] = models.ParseTechnology(rowExcel[index], sheet, rowsProcessed)
					//fmt.Print(" excelSheet2 Technology--->",rowExcel[index], "\r\n")
				}
			}
			logger.Print(rowExcel, "\r\n")
			//check for property name in telco locn master1
			sel = "SELECT SYS_ID FROM odm_basic_tier_telco_locn_master1 WHERE PROPERTY_ID LIKE '" + strings.Trim(rowExcel[46], " ") + "'"
			//fmt.Println(sel)
			selectRows = models.Select(db, sel)
			MASTER_REC_FLAG = "Y"
			if selectRows.Next() {
				selectRows.Scan(&TLM1_SYS_ID)
				//oracle version
				//cmd = "UPDATE ODM_TELCO_LOCN_MASTER1 SET MASTER_REC = :1, LATITUDE = :2, LONGITUDE = :3, GOVERNORATE_ID = :4, WILLAYAT_ID = :5, TOWN = :6, AREA_ID = :7, CLUTTER_ID = :8 WHERE SYS_ID LIKE :9"
				//postgresql version
				cmd = "UPDATE odm_basic_tier_telco_locn_master1 SET MASTER_REC = $1, LATITUDE = $2, LONGITUDE = $3, GOVERNORATE_ID = $4, WILLAYAT_ID = $5, TOWN = $6, AREA_ID = $7, CLUTTER_ID = $8 WHERE SYS_ID = $9"
				rowsAffected = models.Odmtelcolocnmaster1_3g(db, cmd, MASTER_REC_FLAG, rowExcel, TLM1_SYS_ID, v_lat_truncated, v_lon_truncated)
				if rowsAffected < 1 {
					fmt.Println("Problem with rowsAffected")
				}
				locnRecUpdated++
				//fmt.Println("Ch2::Record updated successfully in odm_basic_tier_telco_locn_master1. rowsAffected -->", rowsAffected, "Accumulated updated records: ",locnRecUpdated)
				selectRows.Close()
			} else {
				//oracle version
				//cmd = "INSERT INTO ODM_TELCO_LOCN_MASTER1 (GOVERNORATE_ID, WILLAYAT_ID, TOWN, AREA_ID, CLUTTER_ID, LATITUDE, LONGITUDE, PROPERTY_ID, MASTER_REC) VALUES(:1, :2, :3, :4, :5, :6, :7, :8, :9)"
				//postgresql version
				cmd = "INSERT INTO odm_basic_tier_telco_locn_master1 (GOVERNORATE_ID, WILLAYAT_ID, TOWN, AREA_ID, CLUTTER_ID, LATITUDE, LONGITUDE, PROPERTY_ID, MASTER_REC) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)"
				//fmt.Println(rowExcel[34], rowExcel[35], rowExcel[36], rowExcel[37], rowExcel[40], strings.Trim(rowExcel[31], " "), strings.Trim(rowExcel[32], " "), rowExcel[46], MASTER_REC_FLAG)
				rowsAffected = models.ODMTELCOLOCNMASTER1_3g_Insert(db, cmd, MASTER_REC_FLAG, rowExcel, v_lat_truncated, v_lon_truncated)
				if rowsAffected < 1 {
					fmt.Println("Problem with rowsAffected")
				}
				locnRecInserted++
				//fmt.Println("Ch2::Record successfully inserted for PROPERTY_ID ", rowExcel[46], " in odm_basic_tier_telco_locn_master1", "   Records affected", rowsAffected, "   Accumulated inserted records: ",locnRecInserted)
				sel = "SELECT SYS_ID FROM odm_basic_tier_telco_locn_master1 WHERE PROPERTY_ID LIKE '" + strings.Trim(rowExcel[46], " ") + "'"
				selectRows = models.Select(db, sel)
				if selectRows.Next() {
					selectRows.Scan(&TLM1_SYS_ID)
				}
				selectRows.Close()
			}
			selectRows.Close()

			//check for site name in radio site master1
			sel = "SELECT SYS_ID FROM odm_basic_tier_radio_site_master1 WHERE NODEB_NAME_NEW LIKE '" + strings.Trim(rowExcel[5], " ") + "'" // AND TLM1_SYS_ID LIKE "+TLM1_SYS_ID
			selectRows = models.Select(db, sel)
			MASTER_REC_FLAG = "Y"
			if selectRows.Next() {
				selectRows.Scan(&RSM1_SYS_ID)
				//oracle version
				//cmd = "UPDATE ODM_RADIO_SITE_MASTER1 SET MASTER_REC = :1, SITE_NAME_NEW = :2, NODEB_NAME_NEW = :3, SITE_NAME = :4, NODEB_NAME = :5, SITE_ID_NEW = :6, SITE_ID = :7, RNC_OLD = :8, RNC_NEW = :9, LAC = :10, RAC = :11, SITE_TYPE_ID = :12, VENDOR_ID = :13, TECH_ID = :14, TLM1_SYS_ID = :15 WHERE SYS_ID LIKE :16"
				//postgresql version
				cmd = "UPDATE odm_basic_tier_radio_site_master1 SET MASTER_REC = $1, SITE_NAME_NEW = $2, NODEB_NAME_NEW = $3, SITE_NAME = $4, NODEB_NAME = $5, SITE_ID_NEW = $6, SITE_ID = $7, RNC_OLD = $8, RNC_NEW = $9, LAC = $10, RAC = $11, SITE_TYPE_ID = $12, VENDOR_ID = $13, TECH_ID = $14, TLM1_SYS_ID = $15 WHERE SYS_ID = $16"
				rowsAffected = models.Odmradiositemaster1_3g(db, cmd, MASTER_REC_FLAG, rowExcel, TLM1_SYS_ID, RSM1_SYS_ID)
				if rowsAffected < 1 {
					fmt.Println("Problem with rowsAffected")
				}
				siteRecUpdated++
				//fmt.Println("Ch2::Record updated successfully in odm_basic_tier_radio_site_master1. rowsAffected -->", rowsAffected, "Accumulated updated records: ",siteRecUpdated)
				selectRows.Close()
			} else {
				//oracle version
				//cmd = "INSERT INTO ODM_RADIO_SITE_MASTER1 (MASTER_REC, SITE_NAME_NEW, NODEB_NAME_NEW, SITE_NAME, NODEB_NAME, SITE_ID_NEW, SITE_ID, RNC_OLD, RNC_NEW, LAC, RAC, SITE_TYPE_ID, VENDOR_ID, TECH_ID, TLM1_SYS_ID) VALUES(:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12, :13, :14, :15)"
				//postgresql version
				cmd = "INSERT INTO odm_basic_tier_radio_site_master1 (MASTER_REC, SITE_NAME_NEW, NODEB_NAME_NEW, SITE_NAME, NODEB_NAME, SITE_ID_NEW, SITE_ID, RNC_OLD, RNC_NEW, LAC, RAC, SITE_TYPE_ID, VENDOR_ID, TECH_ID, TLM1_SYS_ID) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)"
				rowsAffected = models.Odmradiositemaster1_3g_Insert(db, cmd, MASTER_REC_FLAG, rowExcel, TLM1_SYS_ID)
				if rowsAffected < 1 {
					fmt.Println("Problem with rowsAffected")
				}
				siteRecInserted++
				//fmt.Println("Ch2::Record successfully inserted for NODEB_NAME_NEW ", rowExcel[5], " in odm_basic_tier_radio_site_master1", "   Records affected", rowsAffected, "   Accumulated inserted records: ",siteRecInserted)
				sel = "SELECT SYS_ID FROM odm_basic_tier_radio_site_master1 WHERE NODEB_NAME_NEW LIKE '" + strings.Trim(rowExcel[5], " ") + "'" // AND TLM1_SYS_ID LIKE "+TLM1_SYS_ID
				selectRows = models.Select(db, sel)
				if selectRows.Next() {
					selectRows.Scan(&RSM1_SYS_ID)
				}
				selectRows.Close()
			}
			selectRows.Close()

			//check for cell name in radio cell master1
			sel = "SELECT SYS_ID FROM odm_basic_tier_radio_cell_master1 WHERE CELL_NAME_NEW LIKE '" + strings.Trim(rowExcel[1], " ") + "'" // AND RSM1_SYS_ID LIKE "+RSM1_SYS_ID
			selectRows = models.Select(db, sel)
			MASTER_REC_FLAG = "Y"
			if selectRows.Next() {
				selectRows.Scan(&RCM1_SYS_ID)
				//oracle version
				// cmd = "UPDATE ODM_RADIO_CELL_MASTER1 SET MASTER_REC = :1, CELL_NEW = :2, CELL_NAME_NEW = :3, CELL = :4, CELL_NAME = :5, CELL_ID_NEW = :6, CELL_ID = :7, "+
				//  "SITE_NAME_NEW = :8, NODEB_NAME_NEW = :9, SITE_NAME = :10, NODEB_NAME = :11, SITE_ID_NEW = :12, SITE_ID = :13, RNC_OLD = :14, RNC_NEW = :15, SECTOR_ID = :16, "+
				//  "SECTOR_ID_NEW = :17, UPLINKUARFCN = :18, DOWNLINKUARFCN = :19, DLPRSCRAMBLECODE = :20, AZIMUTHS = :21, HEIGHTS = :22, ETILT = :23, MTILT = :24, "+
				//  "ANTENNATYPE = :25, BAND_ID = :26, CARRIER_ID = :27, TECH_ID = :28, RSM1_SYS_ID = :29 WHERE SYS_ID LIKE :30"
				//postgresql version
				cmd = "UPDATE odm_basic_tier_radio_cell_master1 SET MASTER_REC = $1, CELL_NEW = $2, CELL_NAME_NEW = $3, CELL = $4, CELL_NAME = $5, CELL_ID_NEW = $6, CELL_ID = $7, " +
					"SITE_NAME_NEW = $8, NODEB_NAME_NEW = $9, SITE_NAME = $10, NODEB_NAME = $11, SITE_ID_NEW = $12, SITE_ID = $13, RNC_OLD = $14, RNC_NEW = $15, SECTOR_ID = $16, " +
					"SECTOR_ID_NEW = $17, UPLINKUARFCN = $18, DOWNLINKUARFCN = $19, DLPRSCRAMBLECODE = $20, AZIMUTHS = $21, HEIGHTS = $22, ETILT = $23, MTILT = $24, " +
					"ANTENNATYPE = $25, BAND_ID = $26, CARRIER_ID = $27, TECH_ID = $28, RSM1_SYS_ID = $29 WHERE SYS_ID = $30"
				rowsAffected = models.Odmradiocellmaster1_3g(db, cmd, MASTER_REC_FLAG, rowExcel, RSM1_SYS_ID, RCM1_SYS_ID)
				if rowsAffected < 1 {
					fmt.Println("Problem with rowsAffected")
				}
				cellRecUpdated++
				//fmt.Println("Ch2::Record updated successfully in odm_basic_tier_radio_cell_master1. rowsAffected -->", rowsAffected, "Accumulated updated records: ",cellRecUpdated)
				selectRows.Close()
			} else {
				//oracle version
				//cmd = "INSERT INTO ODM_RADIO_CELL_MASTER1 (MASTER_REC, CELL_NEW, CELL_NAME_NEW, CELL, CELL_NAME, CELL_ID_NEW, CELL_ID, SITE_NAME_NEW, NODEB_NAME_NEW, SITE_NAME, "+
				//	 "NODEB_NAME, SITE_ID_NEW, SITE_ID, RNC_OLD, RNC_NEW, SECTOR_ID, SECTOR_ID_NEW, UPLINKUARFCN, DOWNLINKUARFCN, DLPRSCRAMBLECODE, AZIMUTHS, HEIGHTS, ETILT, "+
				//	 "MTILT, ANTENNATYPE, BAND_ID, CARRIER_ID, TECH_ID, RSM1_SYS_ID) "+
				//	 "VALUES(:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12, :13, :14, :15, :16, :17, :18, :19, :20, :21, :22, :23, :24, :25, :26, :27, :28, :29)"
				//postgresql version
				cmd = "INSERT INTO odm_basic_tier_radio_cell_master1 (MASTER_REC, CELL_NEW, CELL_NAME_NEW, CELL, CELL_NAME, CELL_ID_NEW, CELL_ID, SITE_NAME_NEW, NODEB_NAME_NEW, SITE_NAME, " +
					"NODEB_NAME, SITE_ID_NEW, SITE_ID, RNC_OLD, RNC_NEW, SECTOR_ID, SECTOR_ID_NEW, UPLINKUARFCN, DOWNLINKUARFCN, DLPRSCRAMBLECODE, AZIMUTHS, HEIGHTS, ETILT, " +
					"MTILT, ANTENNATYPE, BAND_ID, CARRIER_ID, TECH_ID, RSM1_SYS_ID) " +
					"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29)"
				rowsAffected = models.Odmradiocellmaster1_3g_Insert(db, cmd, MASTER_REC_FLAG, rowExcel, RSM1_SYS_ID)
				if rowsAffected < 1 {
					fmt.Println("Problem with rowsAffected")
				}
				cellRecInserted++
				//fmt.Println("Ch2::Record successfully inserted for CELL_NAME_NEW ", rowExcel[1], " in odm_basic_tier_radio_cell_master1", "   Records affected", rowsAffected, "   Accumulated inserted records: ",cellRecInserted)
				sel = "SELECT SYS_ID FROM odm_basic_tier_radio_cell_master1 WHERE CELL_NAME_NEW LIKE '" + strings.Trim(rowExcel[1], " ") + "'" // AND RSM1_SYS_ID LIKE "+RSM1_SYS_ID
				selectRows = models.Select(db, sel)
				if selectRows.Next() {
					selectRows.Scan(&RCM1_SYS_ID)
				}
				selectRows.Close()
			}
			selectRows.Close()
			res := new(keyResults)
			res.rowExcelNew = rowExcel
			res.locnRecUpdatedNew = locnRecUpdated
			res.siteRecUpdatedNew = siteRecUpdated
			res.cellRecUpdatedNew = cellRecUpdated
			res.locnRecInsertedNew = locnRecInserted
			res.siteRecInsertedNew = siteRecInserted
			res.cellRecInsertedNew = cellRecInserted
			logger.Print(res.locnRecUpdatedNew, res.siteRecUpdatedNew, res.cellRecUpdatedNew, res.locnRecInsertedNew, res.siteRecInsertedNew, res.cellRecInsertedNew, "\r\n")
			c <- *res // send res to c
		}
		header = false
	}
	close(c)
}



func (server *Server) UpdateNetworkDB(ctx *gin.Context) {
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
	filenamenew := fileHeader.Filename
	networkDB, err := os.Create("/app/downloads/"+filenamenew)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer networkDB.Close()		

	// //Delete temp file after importing
	// defer os.Remove("/app/downloads/"+filenamenew)

	//Write data from received file to temp file
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
	///////////////////////////
	f, err := excelize.OpenFile("/app/downloads/"+filenamenew)
	if err != nil {
		println(err.Error())
		time.Sleep(20 * time.Second)
	}
	sheetName2 := f.GetSheetName(2)
	sheetName3 := f.GetSheetName(3)
	sheetName4 := f.GetSheetName(4)

	ch2 := make(chan keyResults)
	ch3 := make(chan keyResults)
	ch4 := make(chan keyResults)
	if sheetName2 == "3G" {
		fmt.Println("Found 3G as sheet2: ", sheetName2)
		go excelSheet2(f, sheetName2, ch2)
	}
	if sheetName3 == "LTE" {
		fmt.Println("Found LTE as sheet3: ", sheetName3)
		go excelSheet3(f, sheetName3, ch3, db)
	}
	if sheetName4 == "5G" {
		fmt.Println("Found 5G as sheet4: ", sheetName4)
		go excelSheet4(f, sheetName4, ch4, db)
	}
	ch2Count := 0
	ch3Count := 0
	ch4Count := 0
	var res2New, res3New, res4New keyResults
	for {
		res2, ok2 := <-ch2
		res3, ok3 := <-ch3
		res4, ok4 := <-ch4
		if ok2 {
			ch2Count++
			res2New = res2
			//fmt.Println("Ch2-->", ok2, ch2Count, res2)
		}
		if ok3 {
			ch3Count++
			res3New = res3
			//fmt.Println("Ch3-->", ok3, ch3Count)
		}
		if ok4 {
			ch4Count++
			res4New = res4
			//fmt.Println("Ch4-->", ok4, ch4Count)
		}
		fmt.Println("ch2Count: ", " locnRecUpdatedNew: ", res2New.locnRecUpdatedNew, " siteRecUpdatedNew: ", res2New.siteRecUpdatedNew, " cellRecUpdatedNew: ", res2New.cellRecUpdatedNew)
		fmt.Println("ch2Count: ", " locnRecInsertedNew: ", res2New.locnRecInsertedNew, " siteRecInsertedNew: ", res2New.siteRecInsertedNew, " cellRecInsertedNew: ", res2New.cellRecInsertedNew)
		fmt.Println("ch3Count: ", " locnRecUpdatedNew: ", res3New.locnRecUpdatedNew, " siteRecUpdatedNew: ", res3New.siteRecUpdatedNew, " cellRecUpdatedNew: ", res3New.cellRecUpdatedNew)
		fmt.Println("ch3Count: ", " locnRecInsertedNew: ", res3New.locnRecInsertedNew, " siteRecInsertedNew: ", res3New.siteRecInsertedNew, " cellRecInsertedNew: ", res3New.cellRecInsertedNew)
		fmt.Println("ch4Count: ", " locnRecUpdatedNew: ", res4New.locnRecUpdatedNew, " siteRecUpdatedNew: ", res4New.siteRecUpdatedNew, " cellRecUpdatedNew: ", res4New.cellRecUpdatedNew)
		fmt.Println("ch4Count: ", " locnRecInsertedNew: ", res4New.locnRecInsertedNew, " siteRecInsertedNew: ", res4New.siteRecInsertedNew, " cellRecInsertedNew: ", res4New.cellRecInsertedNew)
		if !ok2 && !ok3 && !ok4 {
			break
		}
	}





	//////////////////////////
	ctx.JSON(http.StatusOK, "File "+filenamenew+" uploaded successfully")
	time.Sleep(1 * time.Second)
}
