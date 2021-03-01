package grafana

import (
	"strings"

	grafanav1 "github.com/linclaus/grafana-operator/api/v1"
	"github.com/linclaus/grafana-operator/pkg/gografana"
	"github.com/sirupsen/logrus"
	log "k8s.io/klog"
)

var folderMap map[string]int

type Grafana struct {
	client gografana.GrafanaClienter
}

func NewGrafana(version, host, adminUser, adminPass, token string) (*Grafana, error) {
	var auth gografana.Authenticator
	if token != "" {
		auth = gografana.NewAPIKeyAuthenticator(token)
	} else {
		auth = gografana.NewBasicAuthenticator(adminUser, adminPass)
	}
	client, err := gografana.GetClientByVersion(version, host, auth)
	if err != nil {
		return nil, err
	}
	fs, err := client.GetAllFolders()
	folderMap = make(map[string]int, 0)
	for _, f := range fs {
		folderMap[f.Title] = f.ID
	}
	return &Grafana{
		client: client,
	}, nil
}

func (grafana *Grafana) DeleteDashboard(dashboardName string) error {
	log.Infof("Remove Dashboard name:%s\n", dashboardName)
	ok, board, err := grafana.client.IsBoardExists(dashboardName)
	if err != nil {
		return err
	}
	if !ok {
		logrus.Debugf("Is Board: %s Exists: %t", dashboardName, ok)
		return nil
	}
	_, err = grafana.client.DeleteDashboard(board.UID)
	log.Infof("DeleteDashboard UID: %s Name: %s", board.UID, dashboardName)
	return err
}

func (grafana *Grafana) UpsertDashboard(gd *grafanav1.GrafanaDashboard) error {
	//find folder.
	folderId := folderMap[gd.Spec.Folder]

	//create new folder
	if folderId == 0 {
		logrus.Println(folderMap)
		logrus.Printf("Create Folder: %s", gd.Spec.Folder)
		fId, _, err := grafana.client.EnsureFolderExists(-1, "", gd.Spec.Folder)
		if err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				return err
			} else {
				fs, _ := grafana.client.GetAllFolders()
				for _, f := range fs {
					folderMap[f.Title] = f.ID
				}
			}
		}
		folderMap[gd.Spec.Folder] = fId
		folderId = fId
	}
	rows := []*gografana.Row{}
	for _, r := range gd.Spec.Rows {
		row := &gografana.Row{Panels: []gografana.Panel_5_0{}}
		row.Title = r.Name
		for i, p := range r.Panels {
			row.Panels = append(row.Panels, pannelFactory(i+1, p.Title, p.Targets[0].Query, p.Targets[0].Legend, p.Targets[0].Ref, p.Type, p.Datasource))
		}
		rows = append(rows, row)
	}
	board := &gografana.Board{
		Title:    gd.Spec.Title,
		Editable: gd.Spec.Editable,
		Rows:     rows,
	}
	return grafana.internalCreateDashboard(folderId, board)
}

func pannelFactory(panelId int, title, exp, legendFormat, refId, pannelType, dataSourceName string) gografana.Panel_5_0 {
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
					LegendFormat:   legendFormat,
					Instant:        false,
					RefID:          refId,
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
					LegendFormat: legendFormat,
					RefID:        refId,
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
					LegendFormat: legendFormat,
					RefID:        refId,
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
					LegendFormat: legendFormat,
					RefID:        refId,
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

func (grafana *Grafana) internalCreateDashboard(folderId int, board *gografana.Board) error {
	_, err := grafana.client.NewDashboard(board, uint(folderId), true)

	if err == nil {
		log.Infof("Create Dashboard name:%s successful\n", board.Title)
		return nil
	}

	internalErr, ok := err.(*gografana.NewDashboardError)
	if ok && internalErr != nil && internalErr.Status == "name-exists" {
		return nil
	}
	log.Infof("Create Dashboard name:%s failed for reason:%s\n", board.Title, err)
	return err
}
