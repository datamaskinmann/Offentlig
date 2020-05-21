/*
* @Author Sindre Fredriksen & Sven Sørensen
* @Version 28.04.2020
 */

/*
* threadManagement pakken brukes for å synkronisere
* funksjoner i programmet.
 */
package threadManagement

import "sync"

/*
* MainSync synkroniserer main thread med
* dataManagement threaden.
 */
var MainSync sync.WaitGroup
