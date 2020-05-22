/**
 * @author Sven Sørensen
 * @version 02/04/2020 (1)
 *
 * Dette biblioteket har som hensikt å implementere funksjonalitet for å ta et datasett og generere
 * en eksponensiell funksjon, en linjær funksjon eller en polynomfunksjon som passet best for dette datasettet.
 * Dette datasettet representeres som et nøkkel-verdi sett der nøkkelen er X og verdien er Y,
 * som i X og Y punkter på en graf
 *
 *  Eksempel på et datasett:
 *  let datasett = {
 *  0: 1,
 *  1: 4,
 *  2: 6,
 *  3: 15,
 *  4: 19,
 *  5: 25,
 *  6: 33,
 *  7: 56,
 * }
 */

/**
 * Denne funksjonen tar et datasett og returnerer summen av alle X verdiene i datasettet
 * @param dataSet et nøkkel-verdi sett der nøkkelen er X og verdien er Y, som
 * i X og Y punkter på en graf.
 * @returns {number} summen av alle X verdiene i datasettet
 */
function sumX(dataSet) {
    let sumX = 0.0;
    for(let i in dataSet) {
        sumX += parseFloat(i);
    }
    return parseFloat(sumX.toFixed(17));
}

/**
 * Denne funksjonen tar et datasett og returnerer summen av alle Y verdiene i datasettet
 * @param dataSet et nøkkel-verdi sett der nøkkelen er X og verdien er Y, som
 * i X og Y punkter på en graf.
 * @returns {number} summen av alle Y verdiene i datasettet
 */
function sumY(dataSet) {
    let sumY = 0.0;
    for(let i in dataSet) {
        sumY += parseFloat(dataSet[i]);
    }
    return parseFloat(sumY.toFixed(17));
}

/**
 * Denne funksjonen tar et datasett og returnerer summen av alle X verdiene ganget med Y verdiene
 * @param dataSet et nøkkel-verdi sett der nøkkelen er X og verdien er Y, som
 * i X og Y punkter på en graf.
 * @returns {number} summen av alle X verdiene ganget med Y verdiene
 */
function sumXY(dataSet) {
    let sumXY = 0.0;
    for (let i in dataSet) {
        sumXY += parseFloat(i) * parseFloat(dataSet[i]);
    }
    return parseFloat(sumXY.toFixed(17));
}

/**
 * Denne funksjonen tar et datasett og returnerer alle verdienene opphøyd i n
 * @param dataSet et nøkkel-verdi sett der nøkkelen er X og verdien er Y, som
 * i X og Y punkter på en graf.
 * @param n eksponenten
 * @returns {number} summen av alle x verdiene opphøyd i n
 */
function sumXNthPower(dataSet, n) {
    let sumXNthPower = 0.0;
    for(let i in dataSet) {
        sumXNthPower += Math.pow(parseFloat(i), n);
    }
    return parseFloat(sumXNthPower.toFixed(17));
}

/**
 * Denne funksjonen tar inn et datasett og returnerer summen av alle x verdiene opphøyd i n
 * ganger y verdien
 * @param dataSet et nøkkel-verdi sett der nøkkelen er X og verdien er Y, som
 * i X og Y punkter på en graf
 * @param n eksponenten du ønsker å opphøye x med
 * @returns {number}
 */
function sumXNthPowerTimesY(dataset, n) {
    let sumXNthPowerTimesY = 0.0;
    for(let i in dataset) {
        sumXNthPowerTimesY += parseFloat((Math.pow(i, n))*dataset[i]);
    }
    return parseFloat(sumXNthPowerTimesY.toFixed(17));
}

/**
 * Denne funksjonens tar inn relevante parameteret og gir det stigningstallet til 'linja'
 * @param n hvor mange elementer i datasettet
 * @param sumxy summen av alle x verdiene ganget med y verdiene
 * @param sumx summen av x verdiene
 * @param sumy summen av y verdiene
 * @param sumx2 summen av x verdiene i kvadrat
 * @returns {number} stigningstallet
 */
function slope(n, sumxy, sumx, sumy, sumx2) {
    // For å kalkulere sitgningstallet til datasettet vårt har vi brukt følgende formel
    // m = (NΣ(xy))−ΣxΣy) / (NΣ(x2))−(Σx)^2
    let m = ((n*sumxy)-sumx*sumy)/((n*sumx2)-Math.pow(sumx, 2));
    return  parseFloat(m.toFixed(17));
}

/**
 * Denne funksjonen tar et datasett og lager en linjær funksjon som passer best for dataene
 * @param dataSet et verdi-nøkkel sett der nøkkelen er X og verdien er Y, som
 * i X og Y punkter på en graf.
 * @returns {{b: number, m: number}} kryssningspunktet og stigningstallet til datasettet
 */
function getLinearRegressionFunction(dataSet) {
    // n = antall elementer i datasettet
    let n = parseFloat(Object.keys(dataSet).length);
    let sumxy = sumXY(dataSet);
    let sumx = sumX(dataSet);
    let sumy = sumY(dataSet);
    let sumx2 = sumXNthPower(dataSet, 2);

    let m = slope(n, sumxy, sumx, sumy, sumx2);

    // For å kalkulere kryssningspunktet til datasettet vårt har vi brukt følgende formel:
    // b = (Σy-mΣx)/n;
    let b = (sumy-m*sumx)/n;

    return {
        b: parseFloat(b.toFixed(17)),
        m: parseFloat(m),
    }
}

/**
 * Denne funksjonen tar en linjær funksjon og inndata som parametere
 * og utfører funksjonen
 * @param func funksjonen
 * @param x inndataen
 * @returns {number} utdataen til funksjonen
 */
function doLinearFunction(func, x) {
    return func.m*x+func.b;
}

/**
 * Denne funksjonen tar et datasett og finner de naturlige logaritmene til Y verdiene i datasettet,
 * altså hvilken eksponent må e konstanten ha for å få Y verdien
 * @param dataSet et verdi-nøkkel sett der nøkkelen er X og verdien er Y, som
 * i X og Y punkter på en graf
 * @returns {{X: number, Y: number}} et nøkkel-verdi sett der X representerer den opprinnelige X verdien
 * og Y representerer den nye Y verdien
 */
function lognDataSet(dataSet) {
    let retData = {};

    for(let item in dataSet) {
        retData[item] = Math.log(dataSet[item]);
    }
    return retData;
}


/**
 * Denne funksjonen tar et datasett og lager en eksponensiell funksjon som passer best til dataene
 * @param dataSet et nøkkel-verdi sett der nøkkelen er X og verdien er Y, som
 * i X og Y punkter på en graf
 * @returns {{b: number, m: number}} kryssningspunktet og stigningstallet til datasettet
 */
function getExponentialRegressionFunction(dataSet) {
    /* Generere et nytt datasett basert på de naturlige logaritmene til
    de opprinnelige Y verdiene*/
    let lognData = lognDataSet(dataSet);
    let lognDataLinear = getLinearRegressionFunction(lognData); // linjær regressjon basert på lognData

    return {
        b: Math.exp(lognDataLinear.b), // Datasettets kryssningspunkt
        m: lognDataLinear.m, // Datasettets stigningstall
    }
}

/**
 * Denne funksjonen tar en exsponentiell funksjon og inndata som parametere
 * og utfører funksjonen
 * @param func funksjonen
 * @param x inndataen
 * @returns {number} utdataen til funksjonen
 */
function doExponentialFunction(func, x) {
    return func.b * Math.exp(func.m*x);
}

/**
 * Denne funksjonen tar et datasett og lager en polynom funksjon som passer best til dataene
 * @param dataSet et nøkkel-verdi sett der nøkkelen er X og verdien er Y, som
 * i X og Y punkter på en graf
 * @returns {{a1: number, a2: number, a0: number}} polynom koeffisientene for a1, a2 og a0
 */
function get2DegreePolynomialRegressionFunction(dataSet) {
    let n = Object.keys(dataSet).length;
    let sumx = sumX(dataSet);
    const k = 2;
    let sumxPowerK = sumXNthPower(dataSet, k);
    let sumxSquared = sumXNthPower(dataSet, 2);
    let sumxPowerK1 = sumXNthPower(dataSet, k+1);
    let sumxPowerKX2 = sumXNthPower(dataSet, k*2);

    let a0 = sumY(dataSet);
    let a1 = sumXY(dataSet);
    let a2 = sumXNthPowerTimesY(dataSet, k);

    let baseM = [];
    for(let i = 0; i < 3; i++) {
        baseM[i] = new Array(3);
    }

    baseM[0][0] = n;
    baseM[1][0] = sumx;
    baseM[2][0]  = sumxPowerK;

    baseM[0][1] = sumx;
    baseM[1][1] = sumxSquared;
    baseM[2][1] = sumxPowerK1;

    baseM[0][2] = sumxPowerK;
    baseM[1][2] = sumxPowerK1;
    baseM[2][2] = sumxPowerKX2;

    let columnVector = [a0, a1, a2];

    let M0 = replaceColumnVector(baseM, 1, columnVector);
    let M1 = replaceColumnVector(baseM, 2, columnVector);
    let M2 = replaceColumnVector(baseM, 3, columnVector);

    return {
        a0: parseFloat(determinant(M0)/determinant(baseM)),
        a1: parseFloat(determinant(M1)/determinant(baseM)),
        a2: parseFloat(determinant(M2)/determinant(baseM)),
    }
}

/**
 * Denne funksjonen tar inn en 3x3 matrise og gir deg matrisens
 * determinant basert på Cramers regel
 * @param M matrisen
 * @returns {number} matrisens determinant
 */
function determinant(m) {
    let a = m[0][0];
    let b = m[1][0];
    let c = m[2][0];

    let d = m[0][1];
    let e = m[1][1];
    let f = m[2][1];

    let g = m[0][2];
    let h = m[1][2];
    let i = m[2][2];

    return (a*((e*i)-(f*h)))-(b*((d*i)-(f*g)))+(c*((d*h)-(e*g)));
}

/**
 * Denne funksjonen tar inn en 3x3 matrise, et kolonnenummer og en 3x1 vektor
 * og returnerer en ny matrise der den c. kolonnen er erstattet med d data
 * @param m matrisen
 * @param c kolonnenummeret, starter på index 1
 * @param d den nye dataen
 */
function replaceColumnVector(m, c, d) {
    let retMatrix = [];
    for(var i=0; i<m.length; i++) {
        retMatrix[i] = new Array(m.length);
    }

    for(let x = 0; x < m.length; x++) {
        for(let y = 0; y < m.length; y++) {
            y == c-1 ? retMatrix[x][y] = d[x] :
                retMatrix[x][y] = m[x][y];
        }
    }
    return retMatrix;
}

/**
 * Denne funksjonen tar en polynomfunksjon og inndata som parametere
 * og utfører funksjonen
 * @param func funksjonen
 * @param x inndataen
 * @returns {number} utdataen til funksjonen
 */
function doPolynomialFunction(func, x) {
    return (func.a2*Math.pow(x, 2))+func.a1*x+func.a0;
}
