# Documentation

## Docker
-  make run/docker 
-  


<!-- 
Notification payload
-->

http://localhost:8080/api/v1/recall/fda?product=peanut butter&category=food&date=20200101
{
  "name": "firsti",
  "store": "firsti",
  "company": "firsti",
  "date": "02-04-2019",
  "country": "usai",
  "category": "firsti",
  "phone": "1234567890",
  "url": "firsti"
}

  "id": "faf2d7cb-9067-4d89-8d45-23f2f01f1b46",
  "user_id": "57a5c3b5-a11a-4eae-a71c-607b47c112f9",
  "recall_id": "firsti",
  "fda_description": "firsti",
  "date": "02-04-2019",



<!-- `

func CheckProductMatch(ctx context.Context, productName string) ([]Recall, error) {
    base := "https://api.fda.gov/food/enforcement.json"
    query := fmt.Sprintf("%s?search=product_description:\"%s\"&limit=10", base, url.QueryEscape(productName))

    req, _ := http.NewRequestWithContext(ctx, "GET", query, nil)
    res, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    var out struct {
        Results []Recall `json:"results"`
    }
    if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
        return nil, err
    }

    return out.Results, nil
} 



terms := strings.Fields("doritos nacho")
query := "search=" + strings.Join(lo.Map(terms, func(t string, _ int) string {
    return fmt.Sprintf("product_description:\"%s\"", t)
}), "+")

search=product_description:"doritos"+product_description:"nacho"



They also have:

Drug recalls: /drug/enforcement/

Device recalls: /device/enforcement/

_, err := db.Exec(`INSERT INTO recalls (...) 
                   VALUES (...) 
                   ON CONFLICT (recall_number) DO NOTHING`)


--
import jsoniter "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary
// Replace
json.Marshal(data)
// Instead of
encoding/json.Marshal(data)

var httpClient = &http.Client{
    Timeout: 5 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:       100,
        IdleConnTimeout:    90 * time.Second,
        DisableKeepAlives:  false,
    },
}
-->
