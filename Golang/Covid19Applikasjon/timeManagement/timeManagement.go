/*
* @Author Sindre Fredriksen & Sven Sørensen
* @Version 28.04.2020
 */

/*
 * timeManagement pakken har ansvar for å kontrollere
 * tids messige aspekter som å utsette tid til nærmeste time.
 */

package timeManagement

import (
	"time"
)

/*
* TruncateToNearestHour returnerer tid av type int64
* ved bruk av pakke time tar den det nåværende tidspunktet
* og runder frem til nærmeste time og returnerer det.
 */

func TruncateToNearestHour() time.Time {

	nextTime := time.Now().Add(time.Hour)

	nextTime = nextTime.Truncate(time.Hour)

	return nextTime
}
