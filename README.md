# COVID Statistics API for LU
> Very unofficial, but uses the data provided by the university publicly

To see the statistics, see [here](https://portal.lancaster.ac.uk/intranet/cms/coronavirus/covid-19-statistics). The API is currently hosted at `https://lucovid.willfantom.dev`.

This API exists only to allow for easier manipulation of the data.

### ⚠️ Warning

If you can't reach `portal.lancaster.ac.uk`, you can't host the API. This might be because you are using CloudFlare DNS.

### API

- **Cases Today** [get]

  `/api/v1/today`

  Will return:
  - `204` if today's data has not yet been published (or scraped)
  - `500` if this crappy code messed up
  - `200` with a json summary of the cases today if successful

- **Cases for Given Day** [get]

  `/api/v1/day`

  Required Parameters:
    - `day`: e.g. `02`
    - `month`: e.g. `Oct`
    - `year`: e.g. `2020`
    > (example builds `October 2nd 2020`)

  Will return:
  - `204` if data has not yet been published (or scraped)
  - `500` if this crappy code messed up
  - `200` with a json summary of the cases for that day if successful, e.g.
      ```json
      {"Date":"2020-10-02T00:00:00Z","Campus":4,"City":2,"Staff":0}
      ```

- **Total Cases Summary** [get]

  `/api/v1/summary`

  Will return:
  - `204` if data has not yet been published (or scraped)
  - `500` if this crappy code messed up
  - `200` with a json summary of the total cases if successful, e.g.
      ```json
      {"Staff Cases":5,"Student Cases":117,"Total Cases":122}
      ```

- **Complete Raw** [get]

  `/api/v1/raw`

  Will return:
  - `204` if data has not yet been published (or scraped)
  - `500` if this crappy code messed up
  - `200` with a json summary of the total cases as given in the table if successful, e.g.
      ```json
      [  {"ID":1,"Date":"2020-10-01T00:00:00Z","Campus":1,"City":2,"Staff":0},
         {"ID":2,"Date":"2020-10-02T00:00:00Z","Campus":4,"City":2,"Staff":0},
         {"ID":3,"Date":"2020-10-03T00:00:00Z","Campus":3,"City":1,"Staff":0},
         {"ID":4,"Date":"2020-10-04T00:00:00Z","Campus":4,"City":4,"Staff":0}  ]
      ```
