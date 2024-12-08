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

// TODO Remove it, use district table
var Dictionary = []DictionaryWithSensors{
	{
		boulevard,
		[]int{7},
	},
	{
		lesnayaPolyana,
		[]int{11},
	},
	{
		metalploshadka,
		[]int{20, 53},
	},
	{
		center,
		[]int{73, 40, 39, 48},
	},
	{
		kirovskii,
		[]int{47},
	},
	{
		yuzhinii,
		[]int{59, 51, 56},
	},
	{
		circus,
		[]int{71},
	},
}
