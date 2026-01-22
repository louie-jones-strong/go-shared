package chartjs

import (
	"fmt"

	"github.com/louie-jones-strong/go-shared/collections/dataframe"
)

type ChartJSData struct {
	Labels   []any            `json:"labels"`
	Datasets []ChartJSDataset `json:"datasets"`
}

type ChartJSDataset struct {
	Label string `json:"label"`
	Data  []any  `json:"data"`
}

func DfToChartJS(
	df *dataframe.DataFrame,
	xAxisName string,
	yAxesNames []string,
) (*ChartJSData, error) {
	switch {
	case df == nil:
		return nil, fmt.Errorf("DfToChartJS called with nil df")
	case len(yAxesNames) == 0:
		return nil, fmt.Errorf("DfToChartJS called with empty x axes")
	}

	labels, err := df.GetColumnByName(xAxisName)
	if err != nil {
		return nil, err
	}

	axes := make([]ChartJSDataset, len(yAxesNames))
	for i, yAxisName := range yAxesNames {
		data, err := df.GetColumnByName(yAxisName)
		if err != nil {
			return nil, err
		}

		axes[i] = ChartJSDataset{
			Label: yAxisName,
			Data:  data.Values(),
		}
	}

	chartData := &ChartJSData{
		Labels:   labels.Values(),
		Datasets: axes,
	}

	return chartData, nil
}
