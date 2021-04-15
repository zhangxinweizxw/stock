package models

import "stock/models/stocks_db"

// App Setting
//---------------------------------------------------------------------------------
const (
	APP_NAME    = "stock"
	APP_VERSION = "0.1.10.1"
	APP_PID     = "stock"
)

var XQStock []*stocks_db.XQ_Stock
var AvsHStockl []*stocks_db.AvsHStock
var ZjlxStockDb []*stocks_db.ZjlxStockDb
var QgqpStockDb []*stocks_db.QgqpStockDb
var DxStockDb []*stocks_db.DxStockDb
