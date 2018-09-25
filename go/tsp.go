package main

import (
    "database/sql"
    "fmt"
    "strconv"
    "strings"
    "sort"
    "math"
    "math/rand"

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
func funcionCosto(cities []city, normalizador float64, maximaDistancia float64) float64 {
    var index = 1
    var eval = 0.0
    for (index < len(cities)){
        var pesAu float64
        pesAu = pesoAumentado(cities[index-1], cities[index], maximaDistancia)
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

//Hace 1000 intentos para ver si cuantos acepta con la te actual 
func porcentajeAceptados(random *rand.Rand, cities []city, te int64, normalizador float64, maximaDistancia float64) float64{
    var c = 0
    var i = 1 
    var ene = 1000
    var ese = cities
    var efeese = funcionCosto(cities, normalizador, maximaDistancia)
    for (i < ene){
        var eseprima = vecino(random, ese)
        var efeeseprima = funcionCosto(eseprima, normalizador, maximaDistancia)
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
func busquedaBinaria(random *rand.Rand, cities []city, te1 int64, te2 int64, pmayus float64, normalizador float64, maximaDistancia float64) int64 {
    var teeme = float64(te1+te2)/2.0
    var epsilomTe = 0.02
    var epsilomPe = 0.04
    if (float64(te2-te1) < epsilomTe){
        return int64(teeme)
    }
    var pminus=porcentajeAceptados(random, cities, int64(teeme), normalizador, maximaDistancia)
    if (math.Abs(pmayus-pminus) < epsilomPe){
        return int64(teeme)
    }
    if (pminus > pmayus){
        return busquedaBinaria(random, cities, te1, int64(teeme), pmayus, normalizador, maximaDistancia)
    }else{
        return busquedaBinaria(random, cities, int64(teeme), te2, pmayus, normalizador, maximaDistancia)
    }
}


//Funcion para buscar la temparatura inicial correcta
func temperaturaInicial(random *rand.Rand, cities []city, te int64, pmayus float64, normalizador float64, maximaDistancia float64) int64 {
    var epsilomPe = 0.04
    var pminus=porcentajeAceptados(random, cities, int64(te), normalizador, maximaDistancia)
    if (math.Abs(pmayus-pminus) < epsilomPe){
        return int64(te)
    }
    var te1 int64
    var te2 int64
    if (pminus < pmayus){
        for(pminus < pmayus){
            te = te*2
             pminus =porcentajeAceptados(random, cities, int64(te), normalizador, maximaDistancia)
        }
        te1 = te/2
        te2 = te 
    }else{
         for(pminus > pmayus){
            te = te/2
             pminus =porcentajeAceptados(random, cities, int64(te), normalizador, maximaDistancia)
        }
        te1 = te
        te2 = te*2        
    }
    return busquedaBinaria(random, cities, te1, te2, pmayus, normalizador, maximaDistancia)
}

func calculaLote(random *rand.Rand, te int64, cities []city, normalizador float64, maximaDistancia float64)(float64, []city){
    var c = 0
    var r = 0.0
    var ele = 400
    var ese = cities
    var efese = funcionCosto(ese, normalizador, maximaDistancia)
    for c < ele{
        var eseprima = vecino(random, ese)
        var efeeseprima = funcionCosto(eseprima, normalizador, maximaDistancia)
        if (efeeseprima < efese + float64(te)){
            ese = eseprima
            efese = efeeseprima
            c = c +1
            r = r + efeeseprima
            fmt.Print("c= ")
            fmt.Print(c)   
            fmt.Print("\n")
            fmt.Print("f(s)= ")
            fmt.Print(FloatToString(efeeseprima))   
            fmt.Print("\n")
        }
    }
    return (float64(r)/float64(ele)), ese
} 

func aceptacionPorUmbrales(random *rand.Rand, te int64, cities []city, normalizador float64, maximaDistancia float64){
    var phi = 0.9
    var epsilon int64
    epsilon = 1000
    var p = 0.0
    for(te > epsilon){
        fmt.Print("T value: ")
        fmt.Print(te)
        fmt.Print("\n")
        var q = math.MaxFloat64
        for(p <= q){
            q = p
            p, cities = calculaLote(random, te, cities, normalizador, maximaDistancia)
        }
        fmt.Print("Promedio aceptados: ")
        fmt.Print(FloatToString(p))   
        fmt.Print("\n")
        //fmt.Print("Lista ciudades: ")    
        //fmt.Print(ciudadesRes)
        te = int64(float64(te)*phi)
        
    
  }
}

func  main() {  
    //Tama;o de la entrada
    var entradaSize = 40
    //Instancia del TSP
    var ciudadesIds = "1,2,3,28,74,163,164,165,166,167,169,326,327,328,329,330,489,490,491,492,493,494,495,653,654,655,658,666,814,815,816,817,818,819,978,979,980,981,1037,1073" 
    var cities = listaCiudades(entradaSize, ciudadesIds)

    //Para cada una de las ciudades imprimos sus valores en la terminal
    //for _, cid := range cities{
    //    fmt.Println(strconv.Itoa(cid.Id) + "|" + cid.Name + 
    //        "|" + cid.Country + "|" + strconv.Itoa(cid.Population) + 
    //        "|" + FloatToString(cid.Latitude) + "|" + FloatToString(cid.Longitude))
    //}

    var distancias []float64
    distancias = listaDistancias(cities)   

    //Inicializamos una variable para la guardar la distancia maxima
    var maximaDistancia float64
    maximaDistancia = distancias[len(distancias)-1]

    //Inicializamos un slice/arreglo para guardar los valores para el normalizador
    var listaNormaliza []float64
    listaNormaliza = listaNormalizador(distancias, entradaSize)

    //Imprimos la lista para el normalizador
    //fmt.Println(listaNormaliza)

    //Inicializamos una variable para el normalizador
    var normalizador float64
    normalizador = getNormalizador(listaNormaliza)
    
    //Imprimos la distancia maxima calculada
    fmt.Print("Dinstancia Maxima: ")
    fmt.Print(FloatToString(maximaDistancia))
    fmt.Print("\n")
    //Imprimimos el normalizador calculado
    fmt.Print("Normalizador: ")
    fmt.Print(FloatToString(normalizador))
    fmt.Print("\n")

    //Tama;o de la entrada para la funcion costo 10
    var entradaSizeFunCosto = 40
    //Instancia del TSP
    var ciudadesFunCostoIds = "1,2,3,28,74,163,164,165,166,167,169,326,327,328,329,330,489,490,491,492,493,494,495,653,654,655,658,666,814,815,816,817,818,819,978,979,980,981,1037,1073" 
    var citiesSolIni = listaCiudades(entradaSizeFunCosto, ciudadesFunCostoIds)

    //Inicializamos una variable para la guardar la funcion costo
    var funCosto float64
    funCosto = funcionCosto(citiesSolIni, normalizador, maximaDistancia)
    //Imprimos el resultado de la funcion costo
    fmt.Print("Funcion Costo: ")
    fmt.Print(FloatToString(funCosto))
    fmt.Print("\n")


    //Incializamos la funcion random
    random := rand.New(rand.NewSource(1))

    //Inicializamos la t
    var teinicial int64
    //La te inicial la tomamos como 8 por lo recomendando en el PDF
    teinicial = 8 
    var pmayus = 0.9
    fmt.Print("Vamos a calcular temp inicial:\n ")
    var te = temperaturaInicial(random, citiesSolIni, teinicial, pmayus, normalizador, maximaDistancia)
    fmt.Print("T value inicial: ")
    fmt.Print(te)
    fmt.Print("\n")

    //Corremos la funcion de aceptaron por umbrales
    aceptacionPorUmbrales(random, te, citiesSolIni, normalizador, maximaDistancia)

}
