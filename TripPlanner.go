package main
import (
    "fmt"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "strconv"
    "sort"
    "bytes"
    "log"
    "github.com/julienschmidt/httprouter"
    "net/http"
    "encoding/json"
    "strings"
    "io/ioutil"
)


const (
    APIUrl string = "https://sandbox-api.uber.com/v1/%s%s"
    ServerTokenValue string =""
    AccessTokenValue string =""
)

type PriceEstimates struct {
    StartLatitude  float64
    StartLongitude float64
    EndLatitude    float64
    EndLongitude   float64
    EstimatedPrice         []PriceEstimate `json:"prices"`
}
type PriceEstimate struct {
    ProductId       string  `json:"product_id"`
    Currency    string  `json:"currency_code"`
    DName     string  `json:"display_name"`
    Est        string  `json:"estimate"`
    LEst     int     `json:"low_estimate"`
    HEst    int     `json:"high_estimate"`
    SurgeMultiplier float64 `json:"surge_multiplier"`
    Duration        int     `json:"duration"`
    Distance        float64 `json:"distance"`
}
type Getter interface {
    get(c *Client) error
}
type Products struct {
    Latitude  float64
    Longitude float64
    Products  []Product `json:"products"`
}
type Product struct {
    ProductId   string `json:"product_id"`
    Description string `json:"description"`
    DisplayName string `json:"display_name"`
    Capacity    int    `json:"capacity"`
    Image       string `json:"image"`
}
type RequestOptions struct {
    ServerToken    string
    AuthorizeUrl   string
    AccessTokenUrl string
    AccessToken string
    BaseUrl        string
}
type modReqObj struct{
    Address string `json:"address"`
    City string `json:"city"`
    State string `json:"state"`
    Zip string `json:"zip"`
}
type reqObj struct{
Id int
Name string `json:"Name"`
Address string `json:"Address"`
City string `json:"City"`
State string `json:"State"`
Zip string `json:"Zip"`
Coordinates struct{
    Lat float64
    Lng float64
}
}
type Responz struct {
    Results []struct {
        AddressComponents []struct {
            LongName  string   `json:"long_name"`
            ShortName string   `json:"short_name"`
            Types     []string `json:"types"`
        } `json:"address_components"`
        FormattedAddress string `json:"formatted_address"`
        Geometry         struct {
            Location struct {
                Lat float64 `json:"lat"`
                Lng float64 `json:"lng"`
            } `json:"location"`
            LocationType string `json:"location_type"`
            Viewport     struct {
                Northeast struct {
                    Lat float64 `json:"lat"`
                    Lng float64 `json:"lng"`
                } `json:"northeast"`
                Southwest struct {
                    Lat float64 `json:"lat"`
                    Lng float64 `json:"lng"`
                } `json:"southwest"`
            } `json:"viewport"`
        } `json:"geometry"`
        PartialMatch bool     `json:"partial_match"`
        PlaceID      string   `json:"place_id"`
        Types        []string `json:"types"`
    } `json:"results"`
    Status string `json:"status"`
}
type TripResponse struct {
    ID                     string   `json:"id"`
    Status                 string   `json:"status"`
    StartingFromLocationID string   `json:"starting_from_location_id"`
    BestRouteLocationIds   []string `json:"best_route_location_ids"`
    TotalUberCosts         int      `json:"total_uber_costs"`
    TotalUberDuration      int      `json:"total_uber_duration"`
    TotalDistance          float64  `json:"total_distance"`
}
type userUber struct {
    StartingFromLocationID string   `json:"starting_from_location_id"`
    LocationIds            []string `json:"location_ids"`
}
type RideRequest struct {
    ProductID      string `json:"product_id"`
    StartLatitude  string `json:"start_latitude"`
    StartLongitude string `json:"start_longitude"`
    EndLatitude    string `json:"end_latitude"`
    EndLongitude   string `json:"end_longitude"`
}
type OnGoingTrip struct {
    ID                        string   `json:"id"`
    TotalDistance             float64  `json:"total_distance"`
    TotalUberCosts            int      `json:"total_uber_costs"`
    TotalUberDuration         int      `json:"total_uber_duration"`
    NextDestinationLocationID string   `json:"next_destination_location_id"`
    StartingFromLocationID    string   `json:"starting_from_location_id"`
    Status                    string   `json:"status"`
    BestRouteLocationIds      []string `json:"best_route_location_ids"`
    UberWaitTimeEta           int      `json:"uber_wait_time_eta"`
}
type ReqResponse struct {
    Driver          interface{} `json:"driver"`
    Eta             int         `json:"eta"`
    Location        interface{} `json:"location"`
    RequestID       string      `json:"request_id"`
    Status          string      `json:"status"`
    SurgeMultiplier int         `json:"surge_multiplier"`
    Vehicle         interface{} `json:"vehicle"`
}
type resObj struct{
Greeting string
}
type Client struct {
    Options *RequestOptions
}
func (pe *PriceEstimates) get(c *Client) error {
    priceEstimateParams := map[string]string{
        "start_latitude":  strconv.FormatFloat(pe.StartLatitude, 'f', 2, 32),
        "start_longitude": strconv.FormatFloat(pe.StartLongitude, 'f', 2, 32),
        "end_latitude":    strconv.FormatFloat(pe.EndLatitude, 'f', 2, 32),
        "end_longitude":   strconv.FormatFloat(pe.EndLongitude, 'f', 2, 32),
    }
    data := c.getRequest("estimates/price", priceEstimateParams)
    if e := json.Unmarshal(data, &pe); e != nil {
        return e
    }
    return nil
}
func Create(options *RequestOptions) *Client {
    return &Client{options}
}
func (c *Client) Get(getter Getter) error {
    if e := getter.get(c); e != nil {
        return e
    }
    return nil
}
func (c *Client) getRequest(endpoint string, params map[string]string) []byte {
    urlParams := "?"
    params["server_token"] = c.Options.ServerToken
    for k, v := range params {
        if len(urlParams) > 1 {
            urlParams += "&"
        }
        urlParams += fmt.Sprintf("%s=%s", k, v)
    }
    url := fmt.Sprintf(APIUrl, endpoint, urlParams)
    res, err := http.Get(url)
    if err != nil {
    }
    data, err := ioutil.ReadAll(res.Body)
    res.Body.Close()
    return data
}
func (pl *Products) get(c *Client) error {
    productParams := map[string]string{
        "latitude":  strconv.FormatFloat(pl.Latitude, 'f', 2, 32),
        "longitude": strconv.FormatFloat(pl.Longitude, 'f', 2, 32),
    }
    data := c.getRequest("products", productParams)
    if e := json.Unmarshal(data, &pl); e != nil {
        return e
    }
    return nil
}
var id int;
var tripId int;

func postlocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    id=id+1;
    decoder := json.NewDecoder(req.Body)
    var t reqObj 
    t.Id = id; 
    err := decoder.Decode(&t)
    if err != nil {
        fmt.Println("Error")
    }
    st:=strings.Join(strings.Split(t.Address," "),"+");
    fmt.Println(st);
    constr := []string {strings.Join(strings.Split(t.Address," "),"+"),strings.Join(strings.Split(t.City," "),"+"),t.State}
    lstringplus := strings.Join(constr,"+")
    locstr := []string{"http://maps.google.com/maps/api/geocode/json?address=",lstringplus}
    resp, err := http.Get(strings.Join(locstr,""))
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
       fmt.Println("Error: Wrong address");
     }
     var data Responz
    err = json.Unmarshal(body, &data)
    fmt.Println(data.Status)
    t.Coordinates.Lat=data.Results[0].Geometry.Location.Lat;
    t.Coordinates.Lng=data.Results[0].Geometry.Location.Lng;
 	conn, err := mgo.Dial("mongodb://cmpe273:1234@ds045454.mongolab.com:45454/cmpe273")
    if err != nil {
        panic(err)
    }
    defer conn.Close();
conn.SetMode(mgo.Monotonic,true);
c:=conn.DB("cmpe273").C("addressBook");
err = c.Insert(t);
    js,err := json.Marshal(t)
    if err != nil{
	   fmt.Println("Error")
	   return
	}
    rw.Header().Set("Content-Type","application/json")
    rw.Write(js)
}
func getlocationdetails(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
fmt.Println(p.ByName("locid"));
id ,err1:= strconv.Atoi(p.ByName("locid"))
if err1 != nil {
        panic(err1)
    }
 conn, err := mgo.Dial("mongodb://cmpe273:1234@ds045454.mongolab.com:45454/cmpe273")

    if err != nil {
        panic(err)
    }
    defer conn.Close();
conn.SetMode(mgo.Monotonic,true);
c:=conn.DB("cmpe273").C("addressBook");
result:=reqObj{}
err = c.Find(bson.M{"id":id}).One(&result)
if err != nil {
                fmt.Println(err)
        }
        js,err := json.Marshal(result)
    if err != nil{
       fmt.Println("Error")
       return
    }
    rw.Header().Set("Content-Type","application/json")
    rw.Write(js)
}

func updatelocationdetails(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
 	id ,err1:= strconv.Atoi(p.ByName("locid"))
 	if err1 != nil {
        panic(err1)
    }
  	conn, err := mgo.Dial("mongodb://cmpe273:1234@ds045454.mongolab.com:45454/cmpe273")
    if err != nil {
        panic(err)
    }
    defer conn.Close();
	conn.SetMode(mgo.Monotonic,true);
 	c:=conn.DB("cmpe273").C("addressBook");
    decoder := json.NewDecoder(req.Body)
    var t modReqObj  
    err = decoder.Decode(&t)
    if err != nil {
        fmt.Println("Error")
    }
    colQuerier := bson.M{"id": id}
    change := bson.M{"$set": bson.M{"address": t.Address, "city":t.City,"state":t.State,"zip":t.Zip}}
    err = c.Update(colQuerier, change)
    if err != nil {
        panic(err)
    }
}
func deletelocationdetails(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
    id ,err1:= strconv.Atoi(p.ByName("locid"))
 	if err1 != nil {
        panic(err1)
    }
  	conn, err := mgo.Dial("mongodb://cmpe273:1234@ds045454.mongolab.com:45454/cmpe273")
  	conn.SetMode(mgo.Monotonic,true);
	c:=conn.DB("cmpe273").C("addressBook");
    if err != nil {
    	panic(err)
    }
    defer conn.Close();
    err=c.Remove(bson.M{"id":id})
    if err != nil { fmt.Printf("Could not find kitten %s to delete", id)}
    rw.WriteHeader(http.StatusNoContent)
}
func plan(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
    decoder := json.NewDecoder(req.Body)
    var uUD userUber 
    err := decoder.Decode(&uUD)
    if err != nil {
        log.Println("Error")
    }
    log.Println(uUD.StartingFromLocationID);
    var options RequestOptions;
    options.ServerToken= ServerTokenValue;
    options.BaseUrl= "https://sandbox-api.uber.com/v1/";
    client :=Create(&options); 
    sid ,err1:= strconv.Atoi(uUD.StartingFromLocationID)
	fmt.Println(uUD.StartingFromLocationID);
 	fmt.Println(sid);
 	if err1 != nil {
        panic(err1)
    }
    conn, err := mgo.Dial("mongodb://cmpe273:1234@ds045454.mongolab.com:45454/cmpe273");
    if err != nil {
        panic(err)
    }
    defer conn.Close();
    conn.SetMode(mgo.Monotonic,true);
    c:=conn.DB("cmpe273").C("addressBook");
    result:=reqObj{}
    err = c.Find(bson.M{"id":sid}).One(&result)
    if err != nil {
        fmt.Println(err)
    }
    index:=0;
    totalPrice := 0;
    totalDistance :=0.0;
    totalDuration :=0;
    bestroute:=make([]float64,len(uUD.LocationIds));
    m := make(map[float64]string)
    for _,ids := range uUD.LocationIds{
        lid,err1:= strconv.Atoi(ids)
        if err1 != nil {
            panic(err1)
        }
        resultLID:=reqObj{}
        err = c.Find(bson.M{"id":lid}).One(&resultLID)
        if err != nil {
            fmt.Println(err)
        }
        pe := &PriceEstimates{}
        pe.StartLatitude = result.Coordinates.Lat;
        pe.StartLongitude = result.Coordinates.Lng;
        pe.EndLatitude = resultLID.Coordinates.Lat;
        pe.EndLongitude = resultLID.Coordinates.Lng;
        if e := client.Get(pe); e != nil {
            fmt.Println(e);
        }
        totalDistance=totalDistance+pe.EstimatedPrice[0].Distance;
        totalDuration=totalDuration+pe.EstimatedPrice[0].Duration;
        totalPrice=totalPrice+pe.EstimatedPrice[0].LEst;
        bestroute[index]=pe.EstimatedPrice[0].Distance;
        m[pe.EstimatedPrice[0].Distance]=ids;
        index=index+1;
    }
    sort.Float64s(bestroute);
    var tripres TripResponse;
    tripId=tripId+1;
    tripres.ID=strconv.Itoa(tripId);
    tripres.TotalDistance=totalDistance;
    tripres.TotalUberCosts=totalPrice;
    tripres.TotalUberDuration=totalDuration;
    tripres.Status="Planning";
    tripres.StartingFromLocationID=strconv.Itoa(sid);
    tripres.BestRouteLocationIds=make([]string,len(uUD.LocationIds));
    index=0;
    for _, ind := range bestroute{
       tripres.BestRouteLocationIds[index]=m[ind];
       index=index+1;
    }
    fmt.Println(tripres.BestRouteLocationIds[1]);
    c1:=conn.DB("cmpe273").C("trips");
    err = c1.Insert(tripres);
        js,err := json.Marshal(tripres)
    if err != nil{
       fmt.Println("Error")
       return
    }
    rw.Header().Set("Content-Type","application/json")
    rw.Write(js)
    }
	func gettripdetails(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
    conn, err := mgo.Dial("mongodb://cmpe273:1234@ds045454.mongolab.com:45454/cmpe273")
    if err != nil {
        panic(err)
    }
    defer conn.Close();
    conn.SetMode(mgo.Monotonic,true);
    c:=conn.DB("cmpe273").C("trips");
    result:=TripResponse{}
    err = c.Find(bson.M{"id":p.ByName("tripid")}).One(&result)
    if err != nil {
        fmt.Println(err)
    }
    js,err := json.Marshal(result)
    if err != nil{
       fmt.Println("Error")
       return
    }
    rw.Header().Set("Content-Type","application/json")
    rw.Write(js)
}
	var currentPos int;
	var ogtID int;
func reqtrip(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
    kid ,err1:= strconv.Atoi(p.ByName("tripid"))
    var siD int;
    if err1 != nil {
        panic(err1)
    }
    var ogt OnGoingTrip;
    result1:=reqObj{}
    result2:=reqObj{}
    conn, err := mgo.Dial("mongodb://cmpe273:1234@ds045454.mongolab.com:45454/cmpe273")
    if err != nil {
        panic(err)
    }
    defer conn.Close();
    conn.SetMode(mgo.Monotonic,true);
    c:=conn.DB("cmpe273").C("trips");
    result:=TripResponse{}
    err = c.Find(bson.M{"id":strconv.Itoa(kid)}).One(&result)
    if err != nil {
        fmt.Println(err)
    }else{
    var iD int;
    c1:=conn.DB("cmpe273").C("addressBook");
    if currentPos==0{
        iD, err = strconv.Atoi(result.StartingFromLocationID)
        siD=iD;
        if err != nil {
            fmt.Println(err)
        }
    }else
    {
        iD, err = strconv.Atoi(result.BestRouteLocationIds[currentPos-1])
        siD, err = strconv.Atoi(result.StartingFromLocationID)
        
        if err != nil {
            fmt.Println(err)
        }
    }
    err = c1.Find(bson.M{"id":iD}).One(&result1)
    if err != nil {
        fmt.Println(err)
    }
    iD, err = strconv.Atoi(result.BestRouteLocationIds[currentPos])
    if err != nil {
        fmt.Println(err)
    }
    err = c1.Find(bson.M{"id":iD}).One(&result2)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(result2.Coordinates.Lat);
    }
    ogt.ID=strconv.Itoa(ogtID);
    ogt.BestRouteLocationIds=result.BestRouteLocationIds;
    ogt.StartingFromLocationID=strconv.Itoa(siD);
    ogt.NextDestinationLocationID=result.BestRouteLocationIds[currentPos];
    ogt.TotalDistance=result.TotalDistance;
    ogt.TotalUberCosts=result.TotalUberCosts;
    ogt.TotalUberDuration=result.TotalUberDuration;
    ogt.Status="requesting";
    var options RequestOptions;
    options.ServerToken= ServerTokenValue;
    options.BaseUrl= "https://sandbox-api.uber.com/v1/";
    client :=Create(&options);
    pl:=Products{};
    pl.Latitude=result1.Coordinates.Lat;
    pl.Longitude=result1.Coordinates.Lng;
    if e := pl.get(client); e != nil {
        fmt.Println(e)
    }
    var prodid string;
    i:=0
    for _, product := range pl.Products {
        if(i == 0){
            prodid = product.ProductId
        }
    }
    var rr RideRequest;
    rr.StartLatitude=strconv.FormatFloat(result1.Coordinates.Lat, 'f', 6, 64);
    rr.StartLongitude=strconv.FormatFloat(result1.Coordinates.Lng, 'f', 6, 64);
    rr.EndLatitude=strconv.FormatFloat(result2.Coordinates.Lat, 'f', 6, 64);
    rr.EndLongitude=strconv.FormatFloat(result2.Coordinates.Lng, 'f', 6, 64);
    rr.ProductID=prodid;
    buf, _ := json.Marshal(rr)
    body := bytes.NewBuffer(buf)
    url := fmt.Sprintf(APIUrl, "requests?","access_token=AccessTokenValue")
    res, err := http.Post(url,"application/json",body)
    if err != nil {
        fmt.Println(err)
    }
    data, err := ioutil.ReadAll(res.Body)
    var rRes ReqResponse;
    err = json.Unmarshal(data, &rRes)
    ogt.UberWaitTimeEta=rRes.Eta;
    js,err := json.Marshal(ogt)
    if err != nil{
       fmt.Println("Error")
       return
    }
    ogtID=ogtID+1;
    currentPos=currentPos+1;
    rw.Header().Set("Content-Type","application/json")
    rw.Write(js)
}
func main() {
    mux := httprouter.New()
    id=0;
    tripId=0;
    currentPos=0;
    ogtID=0;
    mux.POST("/locations",postlocation)
    mux.POST("/trips",plan)
    mux.GET("/locations/:locid",getlocationdetails)
    mux.GET("/trips/:tripid",gettripdetails)
    mux.PUT("/locations/:locid",updatelocationdetails)
    mux.PUT("/trips/:tripid/request",reqtrip)
    mux.DELETE("/locations/:locid",deletelocationdetails)
    server := http.Server{
            Addr:        "0.0.0.0:8000",
            Handler: mux,
    }
    server.ListenAndServe()
}