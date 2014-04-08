package captcha

import (
	"image"
	"captcha/cv"
	"sort"
)

type MetaPredictor struct {
	predictors []*Decoder
}

func NewMetaPredictor() *MetaPredictor {
	ret := MetaPredictor{}

	ret.predictors = []*Decoder{
		&Decoder{
			ImageProcessors: []cv.ImageProcessor {
				&cv.MeanShift{K:1},
			},
			BiColorProcessor: &cv.PeakAverageBasedBiColor{},
			BinaryImageProcessors: []cv.BinaryImageProcessor{
				&cv.RemoveBinaryImageBorder{},
				&cv.RemoveIsolatePoints{},
				&cv.BoundBinaryImage{},
				&cv.ScaleBinaryImage{Height : SCALE_HEIGHT},
			},
			BinaryImagePredictor: &BinaryImageConnectedComponentPredictor{Dx : 1, Dy : 6},
		},
		&Decoder{
			BiColorProcessor: &cv.PeakAverageBasedBiColor{},
			BinaryImageProcessors: []cv.BinaryImageProcessor{
				&cv.RemoveBinaryImageBorder{},
				&cv.RemoveIsolatePoints{},
				&cv.BoundBinaryImage{XMinOpen: 1, YMinOpen : 3},
				&cv.RemoveXAxis{K : 2},
				&cv.ScaleBinaryImage{Height : SCALE_HEIGHT},
			},
			BinaryImagePredictor: &BinaryImageConnectedComponentPredictor{Dx : 1, Dy : 6},
		},
		&Decoder{
			ImageProcessors: []cv.ImageProcessor {
				&cv.MeanShift{K:1},
			},
			ImagePredictor: &ConnectedComponentPredictor{Dx : 1, Dy : 6},
		},
		&Decoder{
			BiColorProcessor: &cv.PeakAverageBasedBiColor{},
			BinaryImageProcessors: []cv.BinaryImageProcessor{
				&cv.RemoveBinaryImageBorder{},
				&cv.RemoveIsolatePoints{},
				&cv.BoundBinaryImage{XMinOpen: 2, YMinOpen : 5},
				&cv.RemoveXAxis{K : 2},
				&cv.ScaleBinaryImage{Height : SCALE_HEIGHT},
			},
			BinaryImagePredictor: &BinaryImageConnectedComponentPredictor{Dx : 1, Dy : 1},
		},
		&Decoder{
			ImageProcessors: []cv.ImageProcessor {
				&cv.MeanShift{K:1},
			},
			ImagePredictor: &ConnectedComponentPredictor{Dx : 2, Dy : 2},
		},
		&Decoder{
			BiColorProcessor: &cv.PeakAverageBasedBiColor{},
			BinaryImageProcessors: []cv.BinaryImageProcessor{
				&cv.Erosion{
					Mask:[]cv.Point{
						cv.Point{X:1,Y:0},
						cv.Point{X:0,Y:1},
					},
				},
				&cv.RemoveBinaryImageBorder{},
				&cv.RemoveIsolatePoints{},
				&cv.BoundBinaryImage{},
				&cv.ScaleBinaryImage{Height : SCALE_HEIGHT},
			},
			BinaryImagePredictor: &BinaryImageConnectedComponentPredictor{Dx : 1, Dy : 6},
		},
		&Decoder{
			BiColorProcessor: &cv.PeakAverageBasedBiColor{},
			BinaryImageProcessors: []cv.BinaryImageProcessor{
				&cv.RemoveBinaryImageBorder{},
				&cv.RemoveIsolatePoints{},
				&cv.BoundBinaryImage{XMinOpen: 1, YMinOpen : 3},
				&cv.ScaleBinaryImage{Height : SCALE_HEIGHT},
			},
			BinaryImagePredictor: &BinaryImageConnectedComponentPredictor{Dx : 2, Dy : 6},
		},
		/*
		&Decoder{
			BiColorProcessor: &cv.PeakAverageBasedBiColor{},
			BinaryImageProcessors: []cv.BinaryImageProcessor{
				&cv.RemoveBinaryImageBorder{},
				&cv.RemoveIsolatePoints{},
				&cv.BoundBinaryImage{XMinOpen: 1, YMinOpen : 3},
				&cv.RemoveXAxis{K : 2},
				&cv.ScaleBinaryImage{Height : SCALE_HEIGHT},
			},
			BinaryImagePredictor: &FastCutBasedPredictor{},
		},
		*/
	}
	return &ret
}

func (self *MetaPredictor) Predict(img image.Image, mki *MaskIndex, chType int) []*Result {
	results := []*Result{}
	rank := make(map[string]*Result)
	for _, predictor := range self.predictors{
		rets := predictor.Predict(img, mki, chType)
		for _, ret := range rets {
			if len(ret.Label) < 4 {
				continue
			}
			existRet, ok := rank[ret.Label]
			if ok {
				if existRet.Weight < ret.Weight {
					rank[ret.Label] = ret
				}
			} else {
				rank[ret.Label] = ret
			}
		}
	}
	for _, ret := range rank {
		results = append(results, ret)
	}
	sort.Sort(ResultSorter(results))
	return results
}
