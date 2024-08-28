package vector

import (
	"fmt"
	"hapemu/model"
	"math"
	"strconv"
	"strings"
)

type VectorGeneratorService struct{}

func NewVectorGeneratorService() *VectorGeneratorService {
	return &VectorGeneratorService{}
}

// region convert smartphone from database to vector
func getValueForPrice(price string) float64 {
	if strings.Compare(price, "essential") == 0 {
		return 1
	} else if strings.Compare(price, "mid") == 0 {
		return 2
	} else if strings.Compare(price, "high") == 0 {
		return 3
	}
	return 4
}

func getValueForProcessor(processor string) float64 {
	sTierProcessors := []string{
		"Snapdragon 8 Gen 3",
		"Snapdragon 8s Gen 3",
		"Snapdragon 8 Gen 2",
		"Dimensity 9300",
		"Dimensity 9200",
		"Dimensity 9000",
		"Exynos 2200",
		"Exynos 2100",
		"Kirin 9000",
		"Kirin 9000e",
		"Apple A17 Bionic",
		"Apple A17 Pro",
		"Apple A16 Bionic",
		"Apple A15 Bionic",
		"Dimensity 9200+",
	}

	aTierProcessors := []string{
		"Snapdragon 8+ Gen 1",
		"Snapdragon 8 Gen 1",
		"Snapdragon 7+ Gen 2",
		"Snapdragon 888+",
		"Snapdragon 888",
		"Snapdragon 870",
		"Snapdragon 865+",
		"Dimensity 8300",
		"Dimensity 8200",
		"Dimensity 8200 Ultra",
		"Dimensity 8100",
		"Dimensity 8100 Ultra",
		"Dimensity 1080",
		"Exynos 1080",
		"Exynos 990",
		"Exynos 9825",
		"Exynos 9820",
		"Kirin 990",
		"Kirin 985",
		"Kirin 980",
		"Apple A14 Bionic",
		"Apple A13 Bionic",
	}

	bTierProcessors := []string{
		"Snapdragon 7 Gen 3",
		"Snapdragon 7s Gen 2",
		"Snapdragon 7 Gen 1",
		"Snapdragon 782",
		"Snapdragon 780",
		"Snapdragon 778G",
		"Snapdragon 778G 5G",
		"Snapdragon 768",
		"Snapdragon 765G",
		"Snapdragon 765G 5G",
		"Snapdragon 750",
		"Snapdragon 750G 5G",
		"Snapdragon 732",
		"Snapdragon 732G",
		"Dimensity 8050",
		"Dimensity 8020",
		"Dimensity 8000",
		"Dimensity 7200",
		"Dimensity 7000",
		"Dimensity 700",
		"Exynos 980",
		"Exynos 9810",
		"Exynos 9611",
		"Exynos 9610",
		"Exynos 1380",
		"Helio G99",
		"Helio G96",
		"Helio G95",
		"Helio G90T",
		"Kirin 970",
		"Kirin 810",
		"Kirin 820 5G",
		"Apple A12 Bionic",
		"Apple A11 Bionic",
	}

	cTierProcessors := []string{
		"Snapdragon 6 Gen 1",
		"Snapdragon 4 Gen 2",
		"Snapdragon 4 Gen 1",
		"Snapdragon 730G",
		"Snapdragon 730",
		"Snapdragon 720",
		"Snapdragon 712",
		"Snapdragon 710",
		"Snapdragon 695",
		"Snapdragon 695 5G",
		"Snapdragon 690",
		"Snapdragon 685",
		"Snapdragon 680",
		"Snapdragon 680 4G",
		"Snapdragon 678",
		"Snapdragon 675",
		"Snapdragon 670",
		"Snapdragon 665",
		"Snapdragon 662",
		"Snapdragon 660",
		"Dimensity 7020",
		"Dimensity 7030",
		"Dimensity 6100",
		"Dimensity 6100+",
		"Dimensity 6080",
		"Exynos 850",
		"Exynos 8895",
		"Helio G88",
		"Helio G85",
		"Helio G80",
		"Helio G70",
		"Helio P90",
		"Helio P70",
		"Helio P65",
		"Helio P60",
		"Kirin 710",
		"Kirin 659",
		"Apple A10 Fusion",
		"Snapdragon 865 5G+",
		"Snapdragon 865 5G",
		"Snapdragon 855",
		"Snapdragon 855+",
		"Snapdragon 845",
		"Snapdragon 835",
		"Exynos 1280",
		"Exynos 1480",
		"Google Tensor G3",
		"Google Tensor G2",
		"Google Tensor",
		"MT6785V",
		"MT6769V",
		"Apple A9",
	}

	for _, cur := range sTierProcessors {
		if strings.Contains(processor, cur) {
			return 4
		}
	}

	for _, cur := range aTierProcessors {
		if strings.Contains(processor, cur) {
			return 3
		}
	}
	for _, cur := range bTierProcessors {
		if strings.Contains(processor, cur) {
			return 2
		}
	}
	for _, cur := range cTierProcessors {
		if strings.Contains(processor, cur) {
			return 1
		}
	}
	return 4
}

func getValueForCamera(camera int) float64 {
	if camera <= 75 {
		return 1
	} else if camera <= 104 {
		return 2
	} else if camera <= 134 {
		return 3
	}
	return 4
}

func getValueForBattery(battery string) float64 {
	var batteryValue, err = strconv.Atoi(battery)
	if err != nil {
		fmt.Println("error converting battery string to integer")
	}
	if batteryValue < 4000 {
		return 1
	} else if batteryValue < 4500 {
		return 2
	} else if batteryValue < 5000 {
		return 3
	}
	return 4
}

func getVecValueFromRam(ram int) float64 {
	if ram <= 4 {
		return 1
	} else if ram <= 8 {
		return 2
	} else if ram <= 12 {
		return 3
	} else {
		return 4
	}
}

func getValueForRam(ram string, ramVec float64) float64 {
	var minVec, maxVec float64
	minVec = 100
	maxVec = -1
	var curr int
	for i := 0; i < len(ram); i++ {
		if ram[i] > '0' && ram[i] <= '9' {
			curr = curr*10 + int(ram[i]-'0')
		} else {
			if curr != 0 {
				minVec = math.Min(minVec, getVecValueFromRam(curr))
				maxVec = math.Max(maxVec, getVecValueFromRam(curr))
				curr = 0
			}
		}
	}
	if curr != 0 {
		minVec = math.Min(minVec, getVecValueFromRam(curr))
		maxVec = math.Max(maxVec, getVecValueFromRam(curr))
		curr = 0
	}
	if ramVec < minVec {
		return minVec
	} else if ramVec >= minVec && ramVec <= maxVec {
		return ramVec
	} else {
		return maxVec
	}
}

func getVecValueFromStorage(storage string) float64 {
	if storage == "32GB" || storage == "64GB" || storage == "128GB" {
		return 1
	} else if storage == "256GB" {
		return 2
	} else if storage == "512GB" {
		return 3
	} else {
		return 4
	}
}

func getValueForStorage(storage string, storageVec float64) float64 {
	var storageList = []string{"32GB", "64GB", "128GB", "256GB", "512GB", "1TB"}
	var minVec, maxVec float64
	for _, cur := range storageList {
		if strings.Contains(storage, cur) {
			minVec = getVecValueFromStorage(cur)
			break
		}
	}
	for _, cur := range storageList {
		if strings.Contains(storage, cur) {
			maxVec = getVecValueFromStorage(cur)
		}
	}
	if storageVec < minVec {
		return minVec
	} else if storageVec >= minVec && storageVec <= maxVec {
		return storageVec
	} else {
		return maxVec
	}
}

func (vgs *VectorGeneratorService) ConvertSmartphoneToVec(smartphone model.Smartphone, userPreferenceVector []float64) []float64 {
	var smartphonesVecs []float64
	smartphonesVecs = append(smartphonesVecs, getValueForPrice(smartphone.SegmentPrice))                       // price
	smartphonesVecs = append(smartphonesVecs, getValueForProcessor(smartphone.Processor))                      // processor
	smartphonesVecs = append(smartphonesVecs, getValueForCamera(smartphone.DxomarkScore))                      // camera
	smartphonesVecs = append(smartphonesVecs, getValueForBattery(smartphone.Battery))                          // battery
	smartphonesVecs = append(smartphonesVecs, getValueForRam(smartphone.Ram, userPreferenceVector[4]))         // ram
	smartphonesVecs = append(smartphonesVecs, getValueForStorage(smartphone.Storage, userPreferenceVector[5])) // storage
	return smartphonesVecs
}

// region convert user quiz to vector
func getPriceValue(price string) float64 {
	if strings.Contains(price, "essensial") {
		return 1
	} else if strings.Contains(price, "midrange") {
		return 2
	} else if strings.Contains(price, "premium") {
		return 3
	} else if strings.Contains(price, "ultra") {
		return 4
	}
	return 2.5 // default value
}

func getValue(str string) float64 {
	if strings.Contains(str, "tidak") {
		return 1
	} else if strings.Contains(str, "cukup") {
		return 2
	} else if strings.Contains(str, "penting") {
		return 3
	} else if strings.Contains(str, "sangat") {
		return 4
	}
	return 2.5 // default value
}

func (vgs *VectorGeneratorService) ConvertUserPreferenceToTargetVec(request model.RecommendationsRequest) []float64 {
	var vec []float64

	vec = append(vec, getPriceValue(request.Price))
	vec = append(vec, getValue(request.Processor))
	vec = append(vec, getValue(request.Camera))
	vec = append(vec, getValue(request.Battery))
	vec = append(vec, getValue(request.Ram))
	vec = append(vec, getValue(request.Storage))

	return vec
}
