package main

import (
    "database/sql"
    "fmt"
    "strconv"
    "strings"
    "sort"


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

//Main donde por el momento calculamos el maximo y el normalizador hasta reestructurar el codigo
func main() {
    //Conexion a la base de datos 
    database, _ := sql.Open("sqlite3", "./tsp.db")
   
    //Tama;o de la entrada
    var entradaSize = 40
    //Instancia del TSP
    var ciudadesIds = "1,2,3,28,74,163,164,165,166,167,169,326,327,328,329,330,489,490,491,492,493,494,495,653,654,655,658,666,814,815,816,817,818,819,978,979,980,981,1037,1073" 

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

    //Para cada una de las ciudades imprimos sus valores en la terminal
    for _, cid := range cities{
        fmt.Println(strconv.Itoa(cid.Id) + "|" + cid.Name + 
            "|" + cid.Country + "|" + strconv.Itoa(cid.Population) + 
            "|" + FloatToString(cid.Latitude) + "|" + FloatToString(cid.Longitude))
    }

    //Inicializamos una variable para la guardar la distancia maxima
    var maximaDistancia float64
    maximaDistancia = 0
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
                fmt.Println(sqlQueryDistance)
                fmt.Println(distance)
                //Agregamos la distancia obtenida al slice/arreglo de distancias
                distancias = append(distancias, distance)
                //si la distanciaMaxima es menor a la dinstancia obtenida actualizaos el valor
                if maximaDistancia < distance {
                    maximaDistancia = distance
                }
            }
        }
    }
    //Ordenamos el slice/arreglo de distancias
    sort.Float64s(distancias)
    //fmt.Println(distancias)
    //fmt.Println(len(distancias))

    //Guardamos el ultimo indice del slice/arreglo de distancias
    var sliceLastIndex = len(distancias)-1

    //Inicializamos un slice/arreglo para guardar los valores para el normalizador
    var listaNormalizador []float64
    //Inicializamos un index que contara de 0 a entradaSize-1 para obtner la lista paa el normalizador
    //de tama;o S-1 donde S es entradaSize
    var indexMinus = 0
    for indexMinus < entradaSize-1{
        //Al ultimo indice sliceLastIndex le vamos restando indexMinus y lo guardamos en listaNormalizador
        listaNormalizador = append(listaNormalizador, distancias[sliceLastIndex-indexMinus])
        //fmt.Println(distancias[sliceLastIndex-indexMinus])
        //fmt.Println(indexMinus)
        indexMinus += 1 
    }
    
    //Imprimos la lista para el normalizador
    fmt.Println(listaNormalizador)

    //Inicializamos una variable para sumar los valores de la lista para calcular 
    //el normalizador
    var normalizador = 0.0
    //Sumamos los valores de todos los elementos del arreglo/slice para tener el
    //Normalizador
    for _, value := range listaNormalizador{
        normalizador += value
    }
    
    //Imprimos la distancia maxima calculada
    fmt.Println(FloatToString(maximaDistancia))
    //Imprimos el normalizador calculado
    fmt.Println(FloatToString(normalizador))
}
