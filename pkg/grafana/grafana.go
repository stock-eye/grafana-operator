package grafana

import (
	"fmt"

	"github.com/linclaus/stock/pkg/cache"
	"github.com/linclaus/stock/pkg/model"

	"github.com/linclaus/stock/pkg/gografana"
	"github.com/sirupsen/logrus"
)

var (
	INCREASE_EXPR_FORMAT           = "stock_increase_gauge{code=\"%s\"}"
	CURRENT_EXPR_FORMAT            = "stock_current_gauge{code=\"%s\"}"
	INCREASE7_MAX_EXPR_FORMAT      = "((stock_current_gauge{code=\"%s\"}-max_over_time(stock_current_gauge{code=\"%s\"}[7d]))/max_over_time(stock_current_gauge{code=\"%s\"}[7d]))*100"
	INCREASE7_MIN_EXPR_FORMAT      = "((stock_current_gauge{code=\"%s\"}-min_over_time(stock_current_gauge{code=\"%s\"}[7d]))/min_over_time(stock_current_gauge{code=\"%s\"}[7d]))*100"
	TRADEVOLUME_EXPR_FORMAT        = "irate(stock_trade_volume_total{code=\"%s\"}[10m])"
	WEIBI_EXPR_FORMAT              = "stock_weibi_gauge{code=\"%s\"}"
	LASTBUY3_EXPR_FORMAT           = "stock_last_buy_3_gauge{code=\"%s\"}"
	LASTSELL3_EXPR_FORMAT          = "stock_last_sell_3_gauge{code=\"%s\"}"
	DASHBOARD_NAME_FORMAT          = "%s:%s"
	INCREASE_DASHBOARD_TITLE       = "今日涨跌幅"
	CURRENT_DASHBOARD_TITLE        = "股票价格"
	INCREASE7_MAX_DASHBOAARD_TITLE = "七日内相对于最大值跌幅"
	INCREASE7_MIN_DASHBOAARD_TITLE = "七日内相对于最小值涨幅"
	TRADEVOLUME_DASHBOARD_TITLE    = "交易量变化趋势"
	WEIBI_DASHBOARD_TITLE          = "委比"
	LASTBUY3_DASHBOARD_TITLE       = "买一至买三之和"
	LASTSELL3_DASHBOARD_TITLE      = "卖一至卖三之和"
)

func CreateDashboardByStockMap(sm *model.StockMap) {
	for _, stock := range sm.List() {
		CreateDashboard(stock.Code, stock.Name)
	}
}

func CreateCustomDashboard(folderName, dashboardName, expr, pannelType string) error {
	fs, err := client.GetAllFolders()
	if err != nil {
		return err
	}
	//find all folders.
	var folderId int
	for i := 0; i < len(fs); i++ {
		if fs[i].Title == folderName {
			folderId = fs[i].ID
		}
	}
	//create new folder
	if folderId == 0 {
		fId, _, err := client.EnsureFolderExists(-1, "", folderName)
		if err != nil {
			return err
		}
		folderId = fId
	}
	rows := []*gografana.Row{}
	row := &gografana.Row{Panels: []gografana.Panel_5_0{}}
	row.Panels = append(row.Panels, pannelFactory(1, dashboardName, expr, pannelType))
	rows = append(rows, row)
	board := &gografana.Board{
		Title:    dashboardName,
		Editable: true,
		Rows:     rows,
	}
	return internalCreateDashboard(dashboardName, folderId, board)
}

func CreateBasicDashboard(folderName, dashboardName, name string) error {
	fs, err := client.GetAllFolders()
	if err != nil {
		return err
	}
	//find all folders.
	var folderId int
	for i := 0; i < len(fs); i++ {
		if fs[i].Title == folderName {
			folderId = fs[i].ID
		}
	}
	//create new folder
	if folderId == 0 {
		fId, _, err := client.EnsureFolderExists(-1, "", folderName)
		if err != nil {
			return err
		}
		folderId = fId
	}

	rows := []*gografana.Row{
		{Panels: []gografana.Panel_5_0{
			// CPU Usage
			{
				ID:              1,
				Datasource:      dataSourceName,
				DashLength:      10,
				Pointradius:     5,
				Linewidth:       1,
				SeriesOverrides: []interface{}{},
				Type:            "graph",
				Title:           "CPU使用率",
				Legend: struct {
					Avg          bool `json:"avg"`
					Current      bool `json:"current"`
					Max          bool `json:"max"`
					Min          bool `json:"min"`
					Show         bool `json:"show"`
					Total        bool `json:"total"`
					Values       bool `json:"values"`
					AlignAsTable bool `json:"alignAsTable"`
				}{Show: true, Avg: true, Current: true, Values: true, AlignAsTable: true},
				Lines: true,
				Targets: []struct {
					Expr           string `json:"expr"`
					Format         string `json:"format"`
					Instant        bool   `json:"instant"`
					IntervalFactor int    `json:"intervalFactor"`
					LegendFormat   string `json:"legendFormat"`
					RefID          string `json:"refId"`
				}{
					{
						Expr:           fmt.Sprintf("avg(sum(irate(container_cpu_usage_seconds_total{name=\"%s\"}[1h])) by (name)*100) by (name)", name),
						Format:         "time_series",
						LegendFormat:   "{{name}}",
						Instant:        false,
						RefID:          "A",
						IntervalFactor: 1,
					},
				},
				Xaxis: struct {
					Buckets interface{}   `json:"buckets"`
					Mode    string        `json:"mode"`
					Name    interface{}   `json:"name"`
					Show    bool          `json:"show"`
					Values  []interface{} `json:"values"`
				}{Mode: "time", Show: true, Values: []interface{}{}},
				Yaxes: []struct {
					Format   string      `json:"format"`
					Label    interface{} `json:"label"`
					LogBase  int         `json:"logBase"`
					Max      interface{} `json:"max"`
					Min      interface{} `json:"min"`
					Show     bool        `json:"show"`
					Decimals int         `json:"decimals"`
				}{
					{Format: "percent", Show: true, LogBase: 1},
					{Format: "short", Show: true, LogBase: 1},
				}},
			// Memory Usage
			{
				ID:              2,
				Datasource:      dataSourceName,
				DashLength:      10,
				Pointradius:     5,
				Linewidth:       1,
				SeriesOverrides: []interface{}{},
				Type:            "graph",
				Title:           "Memory使用量",
				Legend: struct {
					Avg          bool `json:"avg"`
					Current      bool `json:"current"`
					Max          bool `json:"max"`
					Min          bool `json:"min"`
					Show         bool `json:"show"`
					Total        bool `json:"total"`
					Values       bool `json:"values"`
					AlignAsTable bool `json:"alignAsTable"`
				}{Show: true, Avg: true, Current: true, Values: true, AlignAsTable: true},
				Lines: true,
				Targets: []struct {
					Expr           string `json:"expr"`
					Format         string `json:"format"`
					Instant        bool   `json:"instant"`
					IntervalFactor int    `json:"intervalFactor"`
					LegendFormat   string `json:"legendFormat"`
					RefID          string `json:"refId"`
				}{
					{
						Expr:           fmt.Sprintf("sum(container_memory_usage_bytes{name=\"%s\"})  by (name)", name),
						Format:         "time_series",
						LegendFormat:   "{{name}}",
						Instant:        false,
						RefID:          "A",
						IntervalFactor: 1,
					},
				},
				Xaxis: struct {
					Buckets interface{}   `json:"buckets"`
					Mode    string        `json:"mode"`
					Name    interface{}   `json:"name"`
					Show    bool          `json:"show"`
					Values  []interface{} `json:"values"`
				}{Mode: "time", Show: true, Values: []interface{}{}},
				Yaxes: []struct {
					Format   string      `json:"format"`
					Label    interface{} `json:"label"`
					LogBase  int         `json:"logBase"`
					Max      interface{} `json:"max"`
					Min      interface{} `json:"min"`
					Show     bool        `json:"show"`
					Decimals int         `json:"decimals"`
				}{
					{Format: "decbytes", Show: true, LogBase: 1},
					{Format: "short", Show: true, LogBase: 1},
				}},
			// Network Usage(Receive)
			{
				ID:              3,
				Datasource:      dataSourceName,
				DashLength:      10,
				Pointradius:     5,
				Linewidth:       1,
				SeriesOverrides: []interface{}{},
				Type:            "graph",
				Title:           "网络使用率(数据接收)",
				Legend: struct {
					Avg          bool `json:"avg"`
					Current      bool `json:"current"`
					Max          bool `json:"max"`
					Min          bool `json:"min"`
					Show         bool `json:"show"`
					Total        bool `json:"total"`
					Values       bool `json:"values"`
					AlignAsTable bool `json:"alignAsTable"`
				}{Show: true, Avg: true, Current: true, Values: true, AlignAsTable: true},
				Lines:       true,
				Renderer:    "flot",
				SpaceLength: 10,
				Targets: []struct {
					Expr           string `json:"expr"`
					Format         string `json:"format"`
					Instant        bool   `json:"instant"`
					IntervalFactor int    `json:"intervalFactor"`
					LegendFormat   string `json:"legendFormat"`
					RefID          string `json:"refId"`
				}{
					{
						Expr:           fmt.Sprintf("sum(irate(container_network_receive_bytes_total{name=\"%s\"}[1h])) by (name)", name),
						Format:         "time_series",
						LegendFormat:   "{{name}}",
						Instant:        false,
						RefID:          "A",
						IntervalFactor: 1,
					},
				},
				Xaxis: struct {
					Buckets interface{}   `json:"buckets"`
					Mode    string        `json:"mode"`
					Name    interface{}   `json:"name"`
					Show    bool          `json:"show"`
					Values  []interface{} `json:"values"`
				}{Mode: "time", Show: true, Values: []interface{}{}},
				Yaxes: []struct {
					Format   string      `json:"format"`
					Label    interface{} `json:"label"`
					LogBase  int         `json:"logBase"`
					Max      interface{} `json:"max"`
					Min      interface{} `json:"min"`
					Show     bool        `json:"show"`
					Decimals int         `json:"decimals"`
				}{
					{Format: "Bps", Show: true, LogBase: 1},
					{Format: "short", Show: true, LogBase: 1},
				}},
			// Network Usage(Transmit)
			{
				ID:              4,
				Datasource:      dataSourceName,
				DashLength:      10,
				Pointradius:     5,
				Linewidth:       1,
				SeriesOverrides: []interface{}{},
				Type:            "graph",
				Title:           "网络使用率(数据发送)",
				Legend: struct {
					Avg          bool `json:"avg"`
					Current      bool `json:"current"`
					Max          bool `json:"max"`
					Min          bool `json:"min"`
					Show         bool `json:"show"`
					Total        bool `json:"total"`
					Values       bool `json:"values"`
					AlignAsTable bool `json:"alignAsTable"`
				}{Show: true, Avg: true, Current: true, Values: true, AlignAsTable: true},
				Lines:       true,
				Renderer:    "flot",
				SpaceLength: 10,
				Targets: []struct {
					Expr           string `json:"expr"`
					Format         string `json:"format"`
					Instant        bool   `json:"instant"`
					IntervalFactor int    `json:"intervalFactor"`
					LegendFormat   string `json:"legendFormat"`
					RefID          string `json:"refId"`
				}{
					{
						Expr:           fmt.Sprintf("sum(irate(container_network_transmit_bytes_total{name=\"%s\"}[1h])) by (name)", name),
						Format:         "time_series",
						LegendFormat:   "{{name}}",
						Instant:        false,
						RefID:          "A",
						IntervalFactor: 1,
					},
				},
				Xaxis: struct {
					Buckets interface{}   `json:"buckets"`
					Mode    string        `json:"mode"`
					Name    interface{}   `json:"name"`
					Show    bool          `json:"show"`
					Values  []interface{} `json:"values"`
				}{Mode: "time", Show: true, Values: []interface{}{}},
				Yaxes: []struct {
					Format   string      `json:"format"`
					Label    interface{} `json:"label"`
					LogBase  int         `json:"logBase"`
					Max      interface{} `json:"max"`
					Min      interface{} `json:"min"`
					Show     bool        `json:"show"`
					Decimals int         `json:"decimals"`
				}{
					{Format: "Bps", Show: true, LogBase: 1},
					{Format: "short", Show: true, LogBase: 1},
				}},
			// I/O R
			{
				ID:              5,
				Datasource:      dataSourceName,
				DashLength:      10,
				Pointradius:     5,
				Linewidth:       1,
				SeriesOverrides: []interface{}{},
				Type:            "graph",
				Title:           "I/O使用率(数据读取)",
				Legend: struct {
					Avg          bool `json:"avg"`
					Current      bool `json:"current"`
					Max          bool `json:"max"`
					Min          bool `json:"min"`
					Show         bool `json:"show"`
					Total        bool `json:"total"`
					Values       bool `json:"values"`
					AlignAsTable bool `json:"alignAsTable"`
				}{Show: true, Avg: true, Current: true, Values: true, AlignAsTable: true},
				Lines:       true,
				Renderer:    "flot",
				SpaceLength: 10,
				Targets: []struct {
					Expr           string `json:"expr"`
					Format         string `json:"format"`
					Instant        bool   `json:"instant"`
					IntervalFactor int    `json:"intervalFactor"`
					LegendFormat   string `json:"legendFormat"`
					RefID          string `json:"refId"`
				}{
					{
						Expr:           fmt.Sprintf("sum(irate(container_fs_reads_bytes_total{name=\"%s\"}[1h])) by (name)", name),
						Format:         "time_series",
						LegendFormat:   "{{name}}",
						Instant:        false,
						RefID:          "A",
						IntervalFactor: 1,
					},
				},
				Xaxis: struct {
					Buckets interface{}   `json:"buckets"`
					Mode    string        `json:"mode"`
					Name    interface{}   `json:"name"`
					Show    bool          `json:"show"`
					Values  []interface{} `json:"values"`
				}{Mode: "time", Show: true, Values: []interface{}{}},
				Yaxes: []struct {
					Format   string      `json:"format"`
					Label    interface{} `json:"label"`
					LogBase  int         `json:"logBase"`
					Max      interface{} `json:"max"`
					Min      interface{} `json:"min"`
					Show     bool        `json:"show"`
					Decimals int         `json:"decimals"`
				}{
					{Format: "Bps", Show: true, LogBase: 1},
					{Format: "short", Show: true, LogBase: 1},
				}},
			// I/O W
			{
				ID:              6,
				Datasource:      dataSourceName,
				DashLength:      10,
				Pointradius:     5,
				Linewidth:       1,
				SeriesOverrides: []interface{}{},
				Type:            "graph",
				Title:           "I/O使用率(数据写入)",
				Legend: struct {
					Avg          bool `json:"avg"`
					Current      bool `json:"current"`
					Max          bool `json:"max"`
					Min          bool `json:"min"`
					Show         bool `json:"show"`
					Total        bool `json:"total"`
					Values       bool `json:"values"`
					AlignAsTable bool `json:"alignAsTable"`
				}{Show: true, Avg: true, Current: true, Values: true, AlignAsTable: true},
				Lines:       true,
				Renderer:    "flot",
				SpaceLength: 10,
				Targets: []struct {
					Expr           string `json:"expr"`
					Format         string `json:"format"`
					Instant        bool   `json:"instant"`
					IntervalFactor int    `json:"intervalFactor"`
					LegendFormat   string `json:"legendFormat"`
					RefID          string `json:"refId"`
				}{
					{
						Expr:           fmt.Sprintf("sum(irate(container_fs_writes_bytes_total{name=\"%s\"}[1h])) by (name)", name),
						Format:         "time_series",
						LegendFormat:   "{{name}}",
						Instant:        false,
						RefID:          "A",
						IntervalFactor: 1,
					},
				},
				Xaxis: struct {
					Buckets interface{}   `json:"buckets"`
					Mode    string        `json:"mode"`
					Name    interface{}   `json:"name"`
					Show    bool          `json:"show"`
					Values  []interface{} `json:"values"`
				}{Mode: "time", Show: true, Values: []interface{}{}},
				Yaxes: []struct {
					Format   string      `json:"format"`
					Label    interface{} `json:"label"`
					LogBase  int         `json:"logBase"`
					Max      interface{} `json:"max"`
					Min      interface{} `json:"min"`
					Show     bool        `json:"show"`
					Decimals int         `json:"decimals"`
				}{
					{Format: "Bps", Show: true, LogBase: 1},
					{Format: "short", Show: true, LogBase: 1},
				}},
		}},
	}
	board := &gografana.Board{
		Title:    dashboardName,
		Editable: true,
		Rows:     rows,
	}
	return internalCreateDashboard(dashboardName, folderId, board)
}

func RemoveDashboard(dashboardName string) error {
	fmt.Printf("Remove Dashboard name:%s\n", dashboardName)
	ok, board, err := client.IsBoardExists(dashboardName)
	if err != nil {
		return err
	}
	if !ok {
		logrus.Debugf("Is Board: %s Exists: %t", dashboardName, ok)
		return nil
	}
	_, err = client.DeleteDashboard(board.UID)
	logrus.Debugf("DeleteDashboard UID: %s Name: %s", board.UID, dashboardName)
	return err
}

func ChangeDashboard(code, dashboardName string) error {
	fmt.Printf("Change Dashboard code:%s, name:%s\n", code, dashboardName)
	err := RemoveDashboard(dashboardName)
	if err != nil {
		return err
	}
	return CreateDashboard(code, dashboardName)
}

func CreateDashboard(code, name string) error {
	fmt.Printf("Create Dashboard code:%s, name:%s\n", code, name)
	rows := []*gografana.Row{}
	row1 := &gografana.Row{Panels: []gografana.Panel_5_0{}}
	row1.Panels = append(row1.Panels, pannelFactory(1, CURRENT_DASHBOARD_TITLE, fmt.Sprintf(CURRENT_EXPR_FORMAT, code), "graph"))
	row1.Panels = append(row1.Panels, pannelFactory(2, CURRENT_DASHBOARD_TITLE, fmt.Sprintf(CURRENT_EXPR_FORMAT, code), "gauge"))
	rows = append(rows, row1)
	row2 := &gografana.Row{Panels: []gografana.Panel_5_0{}}
	row2.Panels = append(row2.Panels, pannelFactory(3, INCREASE_DASHBOARD_TITLE, fmt.Sprintf(INCREASE_EXPR_FORMAT, code), "graph"))
	row2.Panels = append(row2.Panels, pannelFactory(4, INCREASE_DASHBOARD_TITLE, fmt.Sprintf(INCREASE_EXPR_FORMAT, code), "gauge"))
	rows = append(rows, row2)

	row3 := &gografana.Row{Panels: []gografana.Panel_5_0{}}
	row3.Panels = append(row3.Panels, pannelFactory(5, TRADEVOLUME_DASHBOARD_TITLE, fmt.Sprintf(TRADEVOLUME_EXPR_FORMAT, code), "graph"))
	row3.Panels = append(row3.Panels, pannelFactory(6, WEIBI_DASHBOARD_TITLE, fmt.Sprintf(WEIBI_EXPR_FORMAT, code), "graph"))
	rows = append(rows, row3)

	row4 := &gografana.Row{Panels: []gografana.Panel_5_0{}}
	row4.Panels = append(row4.Panels, pannelFactory(7, LASTBUY3_DASHBOARD_TITLE, fmt.Sprintf(LASTBUY3_EXPR_FORMAT, code), "graph"))
	row4.Panels = append(row4.Panels, pannelFactory(8, LASTSELL3_DASHBOARD_TITLE, fmt.Sprintf(LASTSELL3_EXPR_FORMAT, code), "graph"))
	rows = append(rows, row4)

	row5 := &gografana.Row{Panels: []gografana.Panel_5_0{}}
	row5.Panels = append(row5.Panels, pannelFactory(9, INCREASE7_MAX_DASHBOAARD_TITLE, fmt.Sprintf(INCREASE7_MAX_EXPR_FORMAT, code, code, code), "graph"))
	row5.Panels = append(row5.Panels, pannelFactory(10, INCREASE7_MIN_DASHBOAARD_TITLE, fmt.Sprintf(INCREASE7_MIN_EXPR_FORMAT, code, code, code), "graph"))
	rows = append(rows, row5)

	dashboardName := generateDashboardName(code, name)
	board := &gografana.Board{
		Title:    generateDashboardName(code, name),
		Editable: true,
		Rows:     rows,
	}
	return internalCreateDashboard(dashboardName, folderId, board)
}

func generateDashboardName(code, name string) string {
	return fmt.Sprintf(DASHBOARD_NAME_FORMAT, code, name)
}

func pannelFactory(panelId int, title, exp, pannelType string) gografana.Panel_5_0 {
	var panel gografana.Panel_5_0
	if pannelType == "graph" {
		panel = gografana.Panel_5_0{
			Datasource:      dataSourceName,
			ID:              panelId,
			DashLength:      10,
			Pointradius:     5,
			Linewidth:       1,
			SeriesOverrides: []interface{}{},
			Type:            "graph",
			Title:           title,
			// GridPos: struct {
			// 	H int `json:"h"`
			// 	W int `json:"w"`
			// 	X int `json:"x"`
			// 	Y int `json:"y"`
			// }{W: 13, H: 10},
			Legend: struct {
				Avg          bool `json:"avg"`
				Current      bool `json:"current"`
				Max          bool `json:"max"`
				Min          bool `json:"min"`
				Show         bool `json:"show"`
				Total        bool `json:"total"`
				Values       bool `json:"values"`
				AlignAsTable bool `json:"alignAsTable"`
			}{Show: true, Max: true, Min: true, Current: true, Values: true, AlignAsTable: true},
			Lines: true,
			Targets: []struct {
				Expr           string `json:"expr"`
				Format         string `json:"format"`
				Instant        bool   `json:"instant"`
				IntervalFactor int    `json:"intervalFactor"`
				LegendFormat   string `json:"legendFormat"`
				RefID          string `json:"refId"`
			}{
				{
					Expr:           exp,
					Format:         "time_series",
					LegendFormat:   "{{code}}: {{name}}",
					Instant:        false,
					RefID:          "A",
					IntervalFactor: 1,
				},
			},
			Xaxis: struct {
				Buckets interface{}   `json:"buckets"`
				Mode    string        `json:"mode"`
				Name    interface{}   `json:"name"`
				Show    bool          `json:"show"`
				Values  []interface{} `json:"values"`
			}{Mode: "time", Show: true, Values: []interface{}{}},
			Yaxes: []struct {
				Format   string      `json:"format"`
				Label    interface{} `json:"label"`
				LogBase  int         `json:"logBase"`
				Max      interface{} `json:"max"`
				Min      interface{} `json:"min"`
				Show     bool        `json:"show"`
				Decimals int         `json:"decimals"`
			}{
				{Format: "none", Show: true, LogBase: 1, Decimals: 3},
				{Format: "short", Show: true, LogBase: 1, Decimals: 3},
			},
		}
	} else if pannelType == "gauge" {
		panel = gografana.Panel_5_0{
			Datasource:      dataSourceName,
			ID:              panelId,
			DashLength:      10,
			Pointradius:     5,
			Linewidth:       1,
			SeriesOverrides: []interface{}{},
			Type:            "gauge",
			Title:           title,
			// GridPos: struct {
			// 	H int `json:"h"`
			// 	W int `json:"w"`
			// 	X int `json:"x"`
			// 	Y int `json:"y"`
			// }{W: 6, H: 10},
			Targets: []struct {
				Expr           string `json:"expr"`
				Format         string `json:"format"`
				Instant        bool   `json:"instant"`
				IntervalFactor int    `json:"intervalFactor"`
				LegendFormat   string `json:"legendFormat"`
				RefID          string `json:"refId"`
			}{
				{
					Expr:         exp,
					Format:       "time_series",
					LegendFormat: "{{code}}: {{name}}",
					RefID:        "A",
				},
			},
			FieldConfig: struct {
				Defaults struct {
					Unit     string `json:"unit"`
					Decimals int    `json:"decimals"`
				} `json:"defaults"`
			}{
				Defaults: struct {
					Unit     string `json:"unit"`
					Decimals int    `json:"decimals"`
				}{
					Unit:     "none",
					Decimals: 3,
				},
			},
			Options: struct {
				Orientation   string `json:"orientation"`
				ReduceOptions struct {
					Calcs  []string `json:"calcs"`
					Fields string   `json:"fields"`
					Values bool     `json:"values"`
				} `json:"reduceOptions"`
				ShowThresholdLabels  bool `json:"showThresholdLabels"`
				ShowThresholdMarkers bool `json:"showThresholdMarkers"`
			}{
				Orientation: "auto",
				ReduceOptions: struct {
					Calcs  []string `json:"calcs"`
					Fields string   `json:"fields"`
					Values bool     `json:"values"`
				}{
					Calcs:  []string{"lastNotNull"},
					Fields: "",
					Values: false,
				},
				ShowThresholdLabels:  false,
				ShowThresholdMarkers: true,
			},
			Xaxis: struct {
				Buckets interface{}   `json:"buckets"`
				Mode    string        `json:"mode"`
				Name    interface{}   `json:"name"`
				Show    bool          `json:"show"`
				Values  []interface{} `json:"values"`
			}{Mode: "time", Show: true, Values: []interface{}{}},
		}
	} else if pannelType == "bar" {
		panel = gografana.Panel_5_0{
			Datasource: dataSourceName,
			ID:         panelId,
			Type:       "bargauge",
			Title:      title,
			Targets: []struct {
				Expr           string `json:"expr"`
				Format         string `json:"format"`
				Instant        bool   `json:"instant"`
				IntervalFactor int    `json:"intervalFactor"`
				LegendFormat   string `json:"legendFormat"`
				RefID          string `json:"refId"`
			}{
				{
					Expr:         exp,
					Format:       "heatmap",
					LegendFormat: "{{le}}",
					RefID:        "A",
				},
			},
			FieldConfig: struct {
				Defaults struct {
					Unit     string `json:"unit"`
					Decimals int    `json:"decimals"`
				} `json:"defaults"`
			}{
				Defaults: struct {
					Unit     string `json:"unit"`
					Decimals int    `json:"decimals"`
				}{
					Unit:     "none",
					Decimals: 0,
				},
			},
			Options: struct {
				Orientation   string `json:"orientation"`
				ReduceOptions struct {
					Calcs  []string `json:"calcs"`
					Fields string   `json:"fields"`
					Values bool     `json:"values"`
				} `json:"reduceOptions"`
				ShowThresholdLabels  bool `json:"showThresholdLabels"`
				ShowThresholdMarkers bool `json:"showThresholdMarkers"`
			}{
				Orientation: "auto",
				ReduceOptions: struct {
					Calcs  []string `json:"calcs"`
					Fields string   `json:"fields"`
					Values bool     `json:"values"`
				}{
					Calcs:  []string{"lastNotNull"},
					Fields: "",
					Values: false,
				},
				ShowThresholdLabels:  false,
				ShowThresholdMarkers: true,
			},
		}
	} else if pannelType == "heatmap" {
		panel = gografana.Panel_5_0{
			Datasource: dataSourceName,
			ID:         panelId,
			Type:       "heatmap",
			Title:      title,
			Targets: []struct {
				Expr           string `json:"expr"`
				Format         string `json:"format"`
				Instant        bool   `json:"instant"`
				IntervalFactor int    `json:"intervalFactor"`
				LegendFormat   string `json:"legendFormat"`
				RefID          string `json:"refId"`
			}{
				{
					Expr:         exp,
					Format:       "heatmap",
					LegendFormat: "{{le}}",
					RefID:        "A",
				},
			},
			FieldConfig: struct {
				Defaults struct {
					Unit     string `json:"unit"`
					Decimals int    `json:"decimals"`
				} `json:"defaults"`
			}{
				Defaults: struct {
					Unit     string `json:"unit"`
					Decimals int    `json:"decimals"`
				}{
					Unit:     "none",
					Decimals: 0,
				},
			},
			Options: struct {
				Orientation   string `json:"orientation"`
				ReduceOptions struct {
					Calcs  []string `json:"calcs"`
					Fields string   `json:"fields"`
					Values bool     `json:"values"`
				} `json:"reduceOptions"`
				ShowThresholdLabels  bool `json:"showThresholdLabels"`
				ShowThresholdMarkers bool `json:"showThresholdMarkers"`
			}{
				Orientation: "auto",
				ReduceOptions: struct {
					Calcs  []string `json:"calcs"`
					Fields string   `json:"fields"`
					Values bool     `json:"values"`
				}{
					Calcs:  []string{"lastNotNull"},
					Fields: "",
					Values: false,
				},
				ShowThresholdLabels:  false,
				ShowThresholdMarkers: true,
			},
			Xaxis: struct {
				Buckets interface{}   `json:"buckets"`
				Mode    string        `json:"mode"`
				Name    interface{}   `json:"name"`
				Show    bool          `json:"show"`
				Values  []interface{} `json:"values"`
			}{Mode: "time", Show: true, Values: []interface{}{}},
			Yaxes: []struct {
				Format   string      `json:"format"`
				Label    interface{} `json:"label"`
				LogBase  int         `json:"logBase"`
				Max      interface{} `json:"max"`
				Min      interface{} `json:"min"`
				Show     bool        `json:"show"`
				Decimals int         `json:"decimals"`
			}{
				{Format: "short", Show: true, LogBase: 1},
			},
			Color: struct {
				CardColor   string  `json:"cardColor"`
				ColorScale  string  `json:"colorScale"`
				ColorScheme string  `json:"colorScheme"`
				Exponent    float64 `json:"exponent"`
				Mode        string  `json:"mode"`
			}{CardColor: "#C4162A", ColorScale: "sqrt", ColorScheme: "interpolateOranges", Exponent: 0.3, Mode: "opacity"},
			DataFormat: "tsbuckets",
		}
	}
	return panel
}

func internalCreateDashboard(dashboardName string, folderId int, board *gografana.Board) error {
	_, err := client.NewDashboard(board, uint(folderId), true)

	if err == nil {
		fmt.Printf("Create Dashboard name:%s successful\n", dashboardName)
		return nil
	}

	internalErr, ok := err.(*gografana.NewDashboardError)
	if ok && internalErr != nil && internalErr.Status == "name-exists" {
		return nil
	}
	fmt.Printf("Create Dashboard name:%s failed for reason:%s\n", dashboardName, err)
	return err
}

func findOrCreateFolder() error {
	fs, err := client.GetAllFolders()
	if err != nil {
		return err
	}
	//find all folders.
	for i := 0; i < len(fs); i++ {
		if fs[i].Title == folderName {
			folderId = fs[i].ID
			return nil
		}
	}
	//create new folder
	logrus.Debugf("Preparing create new Grafana foldere: %s...", folderName)
	fId, result, err := client.EnsureFolderExists(-1, "", folderName)
	if err != nil {
		return err
	}
	if !result {
		return fmt.Errorf("Failed to create Grafana folder(%s)!", folderName)
	}
	folderId = fId
	return nil
}

func RefreshDashboardsByFolderId() {
	page := 1
	for {
		boards, err := client.GetDashboardsByFolderId(folderId, page)
		if err != nil {
			fmt.Printf("get dashboards failed for reason: %s\n", err)
		}
		if len(boards) == 0 {
			break
		}
		for _, board := range boards {
			cache.BoardMap[board.Title] = board
		}
		page = page + 1
	}
}
