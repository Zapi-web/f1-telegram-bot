package racing

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Race struct {
	CountryName      string `json:"country_name"`
	DateStart        string `json:"date_start"`
	CountryCode      string `json:"country_code"`
	CircuitShortName string `json:"circuit_short_name"`
}

func GetNewRace() (Race, int, error) {
	var nextRace Race
	url := "https://api.openf1.org/v1/meetings?year=2026"
	dateTodayExample := time.Now()

	ans, err := http.Get(url)

	if err != nil {
		return nextRace, 0, err
	}

	defer ans.Body.Close()

	body, err := io.ReadAll(ans.Body)

	if err != nil {
		return nextRace, 0, err
	}

	var races []Race

	err = json.Unmarshal(body, &races)
	if err != nil {
		return nextRace, 0, err
	}

	if len(races) > 0 {
		var daysBeforeRace int

		for _, race := range races {
			parsTime, err := time.Parse(time.RFC3339, race.DateStart)

			if err != nil {
				return nextRace, 0, err
			}

			result := parsTime.Compare(dateTodayExample)

			if result == -1 {
				continue
			}

			if result == 1 || result == 0 {
				duration := parsTime.Sub(dateTodayExample)
				daysBeforeRace = int(duration.Hours() / 24)

				nextRace = race
				nextRace.DateStart = parsTime.Format("2006-01-02")

				return nextRace, daysBeforeRace, nil
			}
		}

		return nextRace, 0, fmt.Errorf("Похоже в этом году больше гонок нет... ждём вас в новом году гонщики!")
	}

	return nextRace, 0, fmt.Errorf("Гонка не найдена")
}
