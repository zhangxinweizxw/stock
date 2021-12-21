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

var ZtStockDb []*stocks_db.ZtStockDB
var ZjlxStockDb []*stocks_db.ZjlxStockDb
var QgqpStockDb []*stocks_db.QgqpStockDb
var DxStockDb []*stocks_db.DxStockDb
var XqFxStockDb []*stocks_db.XQ_Stock_FX
var ZtStockDb01 []*stocks_db.ZtStockDB01
