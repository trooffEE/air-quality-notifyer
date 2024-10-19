package districts

const (
	center         = "center"
	kirovskii      = "kirovskii"
	circus         = "circus"
	boulevard      = "boulevard"
	yuzhinii       = "yuzhinii"
	metalploshadka = "metalploshadka"
	lesnayaPolyana = "lesnayaPolyana"
)

type DictionaryWithSensors struct {
	Name      string
	SensorIds []int
}

var DictionaryNames = map[string]string{
	center:         "Центральный",
	kirovskii:      "Кировский",
	circus:         `"Цирк"`,
	boulevard:      "Бульвар",
	yuzhinii:       "Южный",
	metalploshadka: "Металлплощадка",
	lesnayaPolyana: "Лесная Поляна",
}

var Dictionary []DictionaryWithSensors = []DictionaryWithSensors{
	DictionaryWithSensors{
		boulevard,
		[]int{7},
	},
	DictionaryWithSensors{
		lesnayaPolyana,
		[]int{11},
	},
	DictionaryWithSensors{
		metalploshadka,
		[]int{20, 53},
	},
	DictionaryWithSensors{
		center,
		[]int{73, 40, 39, 48},
	},
	DictionaryWithSensors{
		kirovskii,
		[]int{47},
	},
	DictionaryWithSensors{
		yuzhinii,
		[]int{59, 51, 56},
	},
	DictionaryWithSensors{
		circus,
		[]int{71},
	},
}
