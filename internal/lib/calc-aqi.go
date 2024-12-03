package lib

func CalcAQI(particlePM, particlePMReferenceHigh, particlePMReferenceLow, pmReferenceIndexHigh, pmReferenceIndexLow float64) float64 {
	return ((pmReferenceIndexHigh-pmReferenceIndexLow)/(particlePMReferenceHigh-particlePMReferenceLow))*(particlePM-particlePMReferenceLow) + pmReferenceIndexLow
}
