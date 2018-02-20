package kolmogorovSmirnovTest2

import (
	"fmt"
	"github.com/dgryski/go-onlinestats"
	"github.com/go-graphite/carbonapi/expr/helper"
	"github.com/go-graphite/carbonapi/expr/interfaces"
	"github.com/go-graphite/carbonapi/expr/metadata"
	"github.com/go-graphite/carbonapi/expr/types"
	"github.com/go-graphite/carbonapi/pkg/parser"
	"math"
)

func init() {
	f := &kolmogorovSmirnovTest2{}
	metadata.RegisterFunction("kolmogorovSmirnovTest2", f)
	metadata.RegisterFunction("ksTest2", f)
}

type kolmogorovSmirnovTest2 struct {
	interfaces.FunctionBase
}

// ksTest2(series, series, points|"interval")
// https://en.wikipedia.org/wiki/Kolmogorov%E2%80%93Smirnov_test
func (f *kolmogorovSmirnovTest2) Do(e parser.Expr, from, until int32, values map[parser.MetricRequest][]*types.MetricData) ([]*types.MetricData, error) {
	arg1, err := helper.GetSeriesArg(e.Args()[0], from, until, values)
	if err != nil {
		return nil, err
	}

	arg2, err := helper.GetSeriesArg(e.Args()[1], from, until, values)
	if err != nil {
		return nil, err
	}

	if len(arg1) != 1 || len(arg2) != 1 {
		return nil, types.ErrWildcardNotAllowed
	}

	a1 := arg1[0]
	a2 := arg2[0]

	windowSize, err := e.GetIntArg(2)
	if err != nil {
		return nil, err
	}

	w1 := &types.Windowed{Data: make([]float64, windowSize)}
	w2 := &types.Windowed{Data: make([]float64, windowSize)}

	r := *a1
	r.Name = fmt.Sprintf("kolmogorovSmirnovTest2(%s,%s,%d)", a1.Name, a2.Name, windowSize)
	r.Values = make([]float64, len(a1.Values))
	r.IsAbsent = make([]bool, len(a1.Values))
	r.StartTime = from
	r.StopTime = until

	d1 := make([]float64, windowSize)
	d2 := make([]float64, windowSize)

	for i, v1 := range a1.Values {
		v2 := a2.Values[i]
		if a1.IsAbsent[i] || a2.IsAbsent[i] {
			// make sure missing values are ignored
			v1 = math.NaN()
			v2 = math.NaN()
		}
		w1.Push(v1)
		w2.Push(v2)

		if i >= windowSize {
			// need a copy here because KS is destructive
			copy(d1, w1.Data)
			copy(d2, w2.Data)
			r.Values[i] = onlinestats.KS(d1, d2)
		} else {
			r.Values[i] = 0
			r.IsAbsent[i] = true
		}
	}
	return []*types.MetricData{&r}, nil
}

// TODO: Implement normal description
// Description is auto-generated description, based on output of https://github.com/graphite-project/graphite-web
func (f *kolmogorovSmirnovTest2) Description() map[string]*types.FunctionDescription {
	return map[string]*types.FunctionDescription{
		"kolmogorovSmirnovTest2": {
			Description: "Nonparametric test of the equality of continuous, one-dimensional probability distributions that can be used to compare a sample with a reference probability distribution (one-sample K–S test), or to compare two samples (two-sample K–S test). https://en.wikipedia.org/wiki/Kolmogorov%E2%80%93Smirnov_test",
			Function:    "kolmogorovSmirnovTest2(seriesList, seriesList, windowSize)",
			Group:       "Transform",
			Module:      "graphite.render.functions",
			Name:        "kolmogorovSmirnovTest2",
			Params: []types.FunctionParam{
				{
					Name:     "seriesList",
					Required: true,
					Type:     types.SeriesList,
				},
				{
					Name:     "seriesList",
					Required: true,
					Type:     types.SeriesList,
				},
				{
					Name:     "window",
					Required: true,
					Type:     types.Integer,
				},
			},
		},
		"ksTest2": {
			Description: "Nonparametric test of the equality of continuous, one-dimensional probability distributions that can be used to compare a sample with a reference probability distribution (one-sample K–S test), or to compare two samples (two-sample K–S test). https://en.wikipedia.org/wiki/Kolmogorov%E2%80%93Smirnov_test",
			Function:    "ksTest2(seriesList, seriesList, windowSize)",
			Group:       "Transform",
			Module:      "graphite.render.functions.custom",
			Name:        "ksTest2",
			Params: []types.FunctionParam{
				{
					Name:     "seriesList",
					Required: true,
					Type:     types.SeriesList,
				},
				{
					Name:     "seriesList",
					Required: true,
					Type:     types.SeriesList,
				},
				{
					Name:     "window",
					Required: true,
					Type:     types.Integer,
				},
			},
		},
	}
}