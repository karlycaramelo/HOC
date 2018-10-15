package main

import (
    "database/sql"
    "fmt"
    "strconv"
    "strings"
    "sort"
    "math"
    "math/rand"


    "log"
    "os"

    _ "github.com/mattn/go-sqlite3"
)

func FloatToString(input_num float64) string {
    // to convert a float number to a string
    return strconv.FormatFloat(input_num, 'f', 6, 64)
}


//Estructura para las ciudades
type city struct {  
    Id int 
    Name string
    Country string
    Population int 
    Latitude float64
    Longitude float64
}

//Estructura param para calular temp inicial
type temp_parametes struct {
    teInicial float64
    pmayus  float64
    epsilomTe float64
    epsilomPe float64
}
//Estructura param para calular tsp
type tsp_parameters struct {
    maximaDistancia float64
    normalizador float64
    citiesDistance map[int]map[int]float64
    ene int64
    te float64
    ele float64
    phi float64
    epsilon float64
}

//func float64 ciudadesDistancias(cities []float64){
    
//}


//Funcion que da la lista de ciudad obtenida de la base de datos
func listaCiudades(entradaSize int, ciudadesIds string) []city{
    //Conexion a la base de datos 
    database, _ := sql.Open("sqlite3", "./tsp.db")

    //Creamos el arreglo de las TSP
    var ciudadesIdsArray = strings.Split(ciudadesIds, ",")

    //Inicializamos un slice/arreglo de ciudades
    cities := []city{}

    // Iteramos en el arreglo de los id de las ciudades
    for _, ciudad := range ciudadesIdsArray{
        //Creamos la query para sacar de la base de datos la ciudad con el id correspondiente
        var sqlQuery = "SELECT id, name, country, population, latitude, longitude FROM cities WHERE id="+ciudad
        //Consulta a la base de datos
        rows, _ := database.Query(sqlQuery)
        //Inicializamos valores para obtener los datos de la consulta
        var id int 
        var name string
        var country string
        var population int 
        var latitude float64
        var longitude float64
        //Para cada uno de los resultados de la query
        for rows.Next() {
            //Obtenemos datos de la consulta
            rows.Scan(&id, &name, &country, &population, &latitude, &longitude)
            //Creamos la estructura para la ciudad
            cid := city{id, name, country, population, latitude, longitude}
            //Agregamos la ciudad al slice/arreglo de ciudades
            cities = append(cities,cid)
        }
    }
    defer database.Close()
    return cities
}

//Dada una lista de ciudades regresa una lista de ciudades
func listaDistancias(cities []city) []float64{
    //Conexion a la base de datos 
    database, _ := sql.Open("sqlite3", "./tsp.db")
    //Creamos un slice/arreglo para guardar todas las distancias
    var distancias []float64

    //2 for anidados para crear las query a la base de datos para preguntar por todos los pares de ciudades
    for _, cid1 := range cities{
        for _, cid2 := range cities{
            //Query para ver si existe la conexuon de la cid1 con cid2
            var sqlQueryDistance = "SELECT distance FROM connections WHERE id_city_1="+ strconv.Itoa(cid1.Id) +" AND id_city_2="+ strconv.Itoa(cid2.Id)
            rows, _ := database.Query(sqlQueryDistance)
            //Variable para guardar la distancia consultada
            var distance float64
            //Este for solo se correra una vez si es que la consulta regresa existencias 
            for rows.Next() {
                //Obtiene la dinstancia
                rows.Scan(&distance)
                //Imprimimos la query y la distancia obtenida
                //fmt.Println(sqlQueryDistance)
                //fmt.Println(distance)
                //Agregamos la distancia obtenida al slice/arreglo de distancias
                distancias = append(distancias, distance)
                //si la distanciaMaxima es menor a la dinstancia obtenida actualizaos el valor
                //if maximaDistancia < distance {
                //    maximaDistancia = distance
                //}
            }
        }
    }
    //Ordenamos el slice/arreglo de distancias
    sort.Float64s(distancias)
    defer database.Close()
    return distancias
}

//dada la lista de distancias obtenemos la lista para el normalizador
func listaNormalizador(distancias []float64, entradaSize int) []float64{
    //Guardamos el ultimo indice del slice/arreglo de distancias
    var sliceLastIndex = len(distancias)-1

    //Inicializamos un slice/arreglo para guardar los valores para el normalizador
    var listaNormaliza []float64
    //Inicializamos un index que contara de 0 a entradaSize-1 para obtner la lista paa el normalizador
    //de tama;o S-1 donde S es entradaSize
    var indexMinus = 0
    for indexMinus < entradaSize-1{
        //Al ultimo indice sliceLastIndex le vamos restando indexMinus y lo guardamos en listaNormalizador
        listaNormaliza = append(listaNormaliza, distancias[sliceLastIndex-indexMinus])
        //fmt.Println(distancias[sliceLastIndex-indexMinus])
        //fmt.Println(indexMinus)
        indexMinus += 1 
    }
    return listaNormaliza
}

//dada la lista para el normalizador obtenemos el normalizador
func getNormalizador(listaNormaliza []float64) float64{
    //Inicializamos una variable para sumar los valores de la lista para calcular 
    //el normalizador
    var normalizador = 0.0
    //Sumamos los valores de todos los elementos del arreglo/slice para tener el
    //Normalizador
    for _, value := range listaNormaliza{
        normalizador += value
    }
    return normalizador
}


//Funcion para obtener la distancia natural acorde al PDF
func distanciaNatural(ciudad1 city, ciudad2 city) float64 {
    //cid.Latitude) + "|" + FloatToString(cid.Longitude
    var latv float64
    var latu float64
    var lonv float64
    var lonu float64
    var r float64
    var c float64
    var dNat float64
    latv = (math.Pi * ciudad1.Latitude)/180.0
    latu = (math.Pi * ciudad2.Latitude)/180.0
    lonv = (math.Pi * ciudad1.Longitude)/180.0
    lonu = (math.Pi * ciudad2.Longitude)/180.0
    r = 6373000.0
    var aVal float64
    aVal = (math.Pow((math.Sin((latv-latu)/2.0)),2.0))+(math.Cos(latu)*math.Cos(latv)*(math.Pow((math.Sin((lonv-lonu)/2.0)),2.0)))
    //No estoy seguro si es la funcion Atan2 o cual funcion es creo que es por esto que 
    //no funciono correctamente
    c = 2 * math.Atan2(math.Sqrt(aVal),math.Sqrt(1-aVal))
    dNat = r*c
    return dNat

}

//Funcion para obetner el pesoAumentado si es que existe la conexion en la base de datos regresa 
//la distancia si no regresa la distancia natural multiplicado por el normalizador
func pesoAumentado(ciudad1 city, ciudad2 city, maximaDistancia float64) float64 {
    //Conexion a la base de datos 
    database, _ := sql.Open("sqlite3", "./tsp.db")

    var sqlQueryDistance1 = "SELECT distance FROM connections WHERE id_city_1="+ strconv.Itoa(ciudad1.Id) +" AND id_city_2="+ strconv.Itoa(ciudad2.Id)
    //fmt.Println(sqlQueryDistance1) 
    rows, _ := database.Query(sqlQueryDistance1)

    var distance float64
    var numRows = 0
    //Este for solo se correra una vez si es que la consulta regresa existencias

    for rows.Next() {
        rows.Scan(&distance)
        numRows = numRows +1
    }
    defer database.Close()
    //fmt.Print("Rows: ")
    //fmt.Print(numRows)
    //fmt.Print("\n")

    //Si hay rows entonces hay un raw y la conexion existe por lo que regresamos la distancia
    if(numRows>0){
        return distance
    //Si no hay rows entonces regresa la distancia natural multiplicada por el normalizador
    }else{
        distNat := distanciaNatural(ciudad1, ciudad2)
        return distNat*maximaDistancia
    }

}


//Dada una lista de ciudades calcula la funcion costo como lo indica el pdf
//Sumamos los pesos aumentos de los pares de nodos (vi-1,vi) y los dividimos entre
//el normalizador 
func funcionCosto(cities []city, tspParams tsp_parameters) float64 {
    var normalizador = tspParams.normalizador
    var citiesDistance = tspParams.citiesDistance
    var index = 1
    var eval = 0.0
    for (index < len(cities)){
        //var pesAuOr float64
        var pesAu float64
        //pesAu = pesoAumentado(cities[index-1], cities[index], maximaDistancia)
        pesAu = citiesDistance[cities[index-1].Id][cities[index].Id]
        //######Prints incluidos pare verificar paso a paso como se hace el calculo de la funcion costo
        /*
        fmt.Print("Distancia de ")
        fmt.Print(cities[index-1])
        fmt.Print(" a: ")
        fmt.Print(cities[index])
        fmt.Print(": ")
        fmt.Print(FloatToString(pesAu))
        fmt.Print("::")
        fmt.Print(FloatToString(pesAuOr))
        fmt.Print("\n")
        */
        //######Prints incluidos pare verificar paso a paso como se hace el calculo de la funcion costo
        //fmt.Println(FloatToString(pesAu))
        eval +=  pesAu
        index = index +1
    }
    eval = eval / normalizador
    return eval
}



//Funcion que calcula un vecion d eforma aletorioa 
//Simplemente hacea el swap de 2 ciudade de manera aleatoria
func vecino(random *rand.Rand, cities []city) []city{
    var numChanges = random.Intn(9) + 1 
    numChanges = 1
    var iter = 1
    //Probe hacer sawp de un numero random entre 1-9 pero 
    //DEspues de probar un par de veces funciona mejor hacer el sawp de una sola ciudad
    //Por eso numchange se queda en 1 para que solo haga un swap
    for (iter <= numChanges){
        var index1 = random.Intn(len(cities))
        var index2 = random.Intn(len(cities))
        var swapCitie city
        swapCitie = cities[index1]
        cities[index1] = cities[index2]
        cities[index2] = swapCitie
        //fmt.Println(cities)
        iter = iter +1
    }
    return cities
}

//Hace tspParams.ene intentos para ver si cuantos acepta con la te actual 
func porcentajeAceptados(random *rand.Rand, cities []city, te float64, tspParams tsp_parameters) float64{
    var c = 0
    var i = 1 
    var ene = tspParams.ene
    var ese = cities
    var efeese = funcionCosto(cities, tspParams)
    for (i < int(ene)){
        var eseprima = vecino(random, ese)
        var efeeseprima = funcionCosto(eseprima, tspParams)
        if (efeeseprima < efeese + float64(te)){
            c = c +1
            ese = eseprima
            efeese = efeeseprima
        }       
        i = i+1
    }
    fmt.Print("Porcentaje Aceptados: ")
    fmt.Print(FloatToString(float64(c)/float64(ene)))   
    fmt.Print("\n")
    return float64(c)/float64(ene)
}


//Busqueda binaria de una te 
func busquedaBinaria(random *rand.Rand, cities []city, te1 float64, te2 float64, tempParams temp_parametes, tspParams tsp_parameters) float64 {
    var pmayus = tempParams.pmayus
    var teeme = float64(te1+te2)/2.0
    var epsilomTe = tempParams.epsilomTe
    var epsilomPe = tempParams.epsilomPe
    if (float64(te2-te1) < epsilomTe){
        return float64(teeme)
    }
    var pminus=porcentajeAceptados(random, cities, float64(teeme), tspParams)
    if (math.Abs(pmayus-pminus) < epsilomPe){
        return float64(teeme)
    }
    if (pminus > pmayus){
        return busquedaBinaria(random, cities, te1, float64(teeme), tempParams, tspParams)
    }else{
        return busquedaBinaria(random, cities, float64(teeme), te2, tempParams, tspParams)
    }
}


//Funcion para buscar la temparatura inicial correcta
func temperaturaInicial(random *rand.Rand, cities []city, tempParams temp_parametes, tspParams tsp_parameters) float64 {
    var te = tempParams.teInicial
    var pmayus = tempParams.pmayus
    var epsilomPe = tempParams.epsilomPe
    var pminus=porcentajeAceptados(random, cities, float64(te), tspParams)
    if (math.Abs(pmayus-pminus) < epsilomPe){
        return float64(te)
    }
    var te1 float64
    var te2 float64
    if (pminus < pmayus){
        for(pminus < pmayus){
            te = te*2
             pminus =porcentajeAceptados(random, cities, float64(te), tspParams)
        }
        te1 = te/2
        te2 = te 
    }else{
         for(pminus > pmayus){
            te = te/2
             pminus =porcentajeAceptados(random, cities, float64(te), tspParams)
        }
        te1 = te
        te2 = te*2        
    }
    return busquedaBinaria(random, cities, te1, te2, tempParams, tspParams)
}


//Funcion calculaLote definida como en el PDF pero se pasan todos lo parametros necesario y bestCities para i guardando la mejor solucion
func calculaLote(random *rand.Rand, te float64, cities []city,  bestCities []city, tspParams tsp_parameters)(float64, []city, bool, []city){
    var maximaDistancia = tspParams.maximaDistancia
    var c = 0.0
    var r = 0.0
    var ele = tspParams.ele
    var ese = cities
    var efese = funcionCosto(ese, tspParams)
    var efebestcities = funcionCosto(bestCities, tspParams)
    //Variable que nos dice cuadno el agoritmos va a continuar pordefualt es verdadero
    var continua = true
    //El limite para saber cuabndo se detienes es ele^2
    var stopLimit = (math.Pow(float64(ele),2.0))
    //var stopLimit = 1000.0
    //Se inicializa el conteo para deternse
    var stopCount = 0.0
    var returnbestCities []city
    for c < ele{
        stopCount = stopCount + 1.0
        var eseprima = vecino(random, ese)
        var efeeseprima = funcionCosto(eseprima, tspParams)
        if (efeeseprima < efese + float64(te)){
            ese = eseprima
            efese = efeeseprima
            c = c +1
            r = r + efeeseprima
            //fmt.Print("c= ")
            //fmt.Print(c)   
            //fmt.Print("\n")
            //fmt.Print("f(s)= ")
            //fmt.Print(FloatToString(efeeseprima))   
            //fmt.Print("\n")
            //Si el resultado que se encuenta es mejor al bestCities actual
            //actualizamos el valor
            if (efese < efebestcities){
                efebestcities = efese
                bestCities = make([]city, len(ese))
                copy(bestCities, ese)
                fmt.Print("MejorSolucionFactibilidad: ")
                fmt.Print(FloatToString(efebestcities/maximaDistancia))   
                fmt.Print("\n")
            }
            stopCount = 0.0
        }
        //Si el conteo es mayor al limite se va a deetener y la variable continua se pone en falso
        if (stopCount >= stopLimit){
            continua = false
            break
        }
        returnbestCities = bestCities
    }
    var salidabestcities = funcionCosto(returnbestCities, tspParams)
    fmt.Print("MejorSolucionFactibilidadSalida: ")
    fmt.Print(FloatToString(salidabestcities/maximaDistancia))   
    fmt.Print("\n")

    //Regresamos el promedio de la soluciones encontradas, la solucion ese que se encontra como ultima
    //continua que nos dice si algoritmo va a continuar esta es false solo cuando se pasa el limite de intentos stopLomit
    //Y la mejor solucion encontrada hasta el momento (bestCities)
    return (float64(r)/float64(ele)), ese, continua, returnbestCities
} 


//Funcion de aceptacion por umbrales acorde al PDF 
func aceptacionPorUmbrales(random *rand.Rand, te float64, cities []city, tspParams tsp_parameters){
    var maximaDistancia = tspParams.maximaDistancia
    var phi = tspParams.phi
    var epsilon float64
    epsilon = tspParams.epsilon
    var p = 0.0
    //Variable que va guardando la mejor solucion 
    //Esta variable se pasa como parametro a cada vez que se llamanda a llmar funcionCosto y se actualiza con los resultados
    //De la funcion
    var bestCities []city
    //Esta varaibles nos dice cuando debemos detener el algoritmo (que es cuando se hacen N  (muy grande en relacion a L) intentos
    //De mejorar lel resultado y no se avanza
    var continuar bool
    //bestCities se empieza con la lista de entrada
    bestCities = make([]city, len(cities))
    copy(bestCities, cities)
    fmt.Print(bestCities)
    fmt.Print(cities)
    for(te > epsilon){
        fmt.Print("T value: ")
        fmt.Print(te)
        fmt.Print("\n")
        var q = math.MaxFloat64
        for(p <= q){
            q = p
            //Mandamos a llamar calculaLote con todo los parametros necesarios, usnado bestCities para guardar el mejor resultado
            //la variabl continuar nos dice si debemos de continuar con el algoritmo (si ya no se pudo avanzar)
            p, cities, continuar, bestCities = calculaLote(random, te, cities, bestCities, tspParams)
            /*
            var efebestcities = funcionCosto(bestCities, normalizador, maximaDistancia, citiesDistance)
            fmt.Print("MejorSolucionCiudades: ")
            fmt.Print(bestCities)   
            fmt.Print("\n")
            fmt.Print("MejorSolucionFactibilidad: ")
            fmt.Print(FloatToString(efebestcities/maximaDistancia))   
            fmt.Print("\n")
            */
            if (!continuar){
                break
            }
        }
        fmt.Print("Promedio aceptados: ")
        fmt.Print(FloatToString(p))   
        fmt.Print("\n")
        //fmt.Print("Lista ciudades: ")    
        //fmt.Print(ciudadesRes)
        te = float64(float64(te)*phi) 

        //Cuando ya no se va a continur imprimos el mejero resultado y su funcion de cost en la consola y lo 
        //Guardamos en un archivo
        if (!continuar){
            var efebestcities = funcionCosto(bestCities, tspParams)
            fmt.Print("MejorSolucionCiudades: ")
            fmt.Print(bestCities)   
            fmt.Print("\n")
            fmt.Print("MejorSolucionFactibilidad: ")
            fmt.Print(FloatToString(efebestcities/maximaDistancia))   
            fmt.Print("\n")

            var strMejorSolCiudades = citiesACadenaIndices(bestCities)+"\n"
            appendFile("fileResultsMEM.txt", strMejorSolCiudades)

            var strMejorSolFacti = "MejorSolucionFactibilidad: "+FloatToString(efebestcities/maximaDistancia)+"\n"
            appendFile("fileResultsMEM.txt", strMejorSolFacti)



            break

        }
    }
}

func appendFile(file_name string, string_to_write string) {  
    os.OpenFile(file_name, os.O_RDONLY|os.O_CREATE, 0666)   
    file, err := os.OpenFile(file_name, os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        log.Fatalf("failed opening file: %s", err)
    }
    defer file.Close()
 
    len, err := file.WriteString(string_to_write)
    if err != nil {
        log.Fatalf("failed writing to file: %s", err)
    }
    fmt.Printf("\nLength: %d bytes", len)
    fmt.Printf("\nFile Name: %s", file.Name())
}

func citiesACadenaIndices(cities []city) string {
    var sizeCities = len(cities)
    var indexCities = 0
    var cadenaCities = ""
    for (indexCities < sizeCities){
        if (indexCities != 0){
            cadenaCities += ","
        } 
        cadenaCities += strconv.Itoa(cities[indexCities].Id)
        indexCities = indexCities +1
    }
    return cadenaCities
}

func  main() {  

    //Instancia del TSP
    //Tama;o de la entrada
    var entradaSize = 40
    var ciudadesIds = "1,2,3,28,74,163,164,165,166,167,169,326,327,328,329,330,489,490,491,492,493,494,495,653,654,655,658,666,814,815,816,817,818,819,978,979,980,981,1037,1073" 
    var cities = listaCiudades(entradaSize, ciudadesIds)

    //Inicializamos una variable para la guardar la distancia maxima
    var maximaDistancia float64
    var distancias []float64
    distancias = listaDistancias(cities)
    maximaDistancia = distancias[len(distancias)-1]
    //Imprimos la distancia maxima calculada
    fmt.Print("Dinstancia Maxima: ")
    fmt.Print(FloatToString(maximaDistancia))
    fmt.Print("\n")

    //Inicializamos un slice/arreglo para guardar los valores para el normalizador
    var listaNormaliza []float64
    listaNormaliza = listaNormalizador(distancias, entradaSize)
    //Inicializamos una variable para el normalizador
    var normalizador float64
    normalizador = getNormalizador(listaNormaliza)
    //Imprimimos el normalizador calculado
    fmt.Print("Normalizador: ")
    fmt.Print(FloatToString(normalizador))
    fmt.Print("\n")

    //Calculamos el pesoAumentadoParaTodasLasCiudades
    var citiesDistance map[int]map[int]float64
    citiesDistance = make(map[int]map[int]float64)
    var indexCI = 0 
    for (indexCI < len(cities)){
        var indexCJ = 0 
        citiesDistance[cities[indexCI].Id] = make(map[int]float64)
        for (indexCJ < len(cities)){
            citiesDistance[cities[indexCI].Id][cities[indexCJ].Id] = pesoAumentado(cities[indexCI], cities[indexCJ], maximaDistancia)
            indexCJ = indexCJ + 1
        }
        indexCI = indexCI +1
    }



    //###########CODIGO PARA PROBAR LA LISTA QUE CANEK DIJO QUE TGENIAMOS MAL EL VALOR y se modifico funcionCosto para imprimr cada paso y tratar de ver que esta mal
    /*
    //Tama;o de la entrada para la funcion costo 10
    var entradaSizeFunCosto = 40
    var pruebaCanek = "330,3,817,653,816,493,163,1,815,490,329,165,978,2,492,654,164,979,981,167,326,814,818,28,1037,655,666,658,169,328,495,166,327,489,494,980,74,819,1073,491"
    var citiesPruebaCanek = listaCiudades(entradaSizeFunCosto, pruebaCanek)
    var funCostoCanek float64
    funCostoCanek = funcionCosto(citiesPruebaCanek, normalizador, maximaDistancia, citiesDistance)
    //Imprimos el resultado de la funcion costo
    fmt.Print("Funcion Costo Canek: ")
    fmt.Print(FloatToString(funCostoCanek))
    fmt.Print("\n")
    */
    //###########CODIGO PARA PROBAR LA LISTA QUE CANEK DIJO QUE TGENIAMOS MAL EL VALOR y se modifico funcionCosto para imprimr cada paso y tratar de ver que esta mal
    
/*
//Estructura param para calular temp inicial
type temp_parametes struct {
    teInicial float64
    pmayus  float64
    epsilomTe float64
    epsilomPe float64
}
//Estructura param para calular tsp
type tsp_parameters struct {
    maximaDistancia float64
    normalizador float64
    citiesDistance map[int]map[int]float64
    ene int64
    te float64
    ele float64
    phi float64
    epsilon float64
}
*/

    tempParams := temp_parametes{8, 0.9, 0.02, 0.04}    
    tspParams := tsp_parameters{maximaDistancia, normalizador, citiesDistance, 1000, 0.0, 400.0, 0.9, 1000}
    
    //Definimos las variables para hacer un loob para correr varias veces el algoritmo
    //tomando diferentes semillas para el random
    var intInicio float64
    intInicio = 1
    var intFinal float64
    intFinal = 9999
    for (intInicio <= intFinal){
       //Incializamos la funcion random
        random := rand.New(rand.NewSource(int64(intInicio)))
        var semillaRandom = "Semilla para random: "
        semillaRandom += strconv.Itoa(int(intInicio))
        semillaRandom += "\n"
        appendFile("fileResultsMEM.txt", semillaRandom)

        //Vamos a calcular temp inicial:
        var te = temperaturaInicial(random, cities, tempParams, tspParams)
        //Actualizamos la temperatura inicial para le valor random
        tspParams.te = te
        var valorInicial = "T value inicial: "
        valorInicial += strconv.Itoa(int(tspParams.te))
        valorInicial += "\n"
        appendFile("fileResultsMEM.txt", valorInicial)

 
        //Corremos la funcion de aceptaron por umbrales
        aceptacionPorUmbrales(random, te, cities, tspParams)
        intInicio = intInicio +1
    }
    

}
