# COVID Statistics API for LU
> Very unofficial, but uses the data provided by the university publicly

To see the statistics, see [here](https://portal.lancaster.ac.uk/intranet/cms/coronavirus/covid-19-statistics). The API is currently hosted at `https://lucovid.willfantom.dev`. The data can also be found through the telegram bot `@LU_Covid19_bot`.

This API exists only to allow for easier manipulation of the data.

### ⚠️ Warning

If you can't reach `portal.lancaster.ac.uk`, you can't host the API. This might be because you are using CloudFlare DNS.

### API

- **Cases Today** [get]

  `/api/v1/today`

  Telegram bot function: `/today`

  Will return:
  - `204` if today's data has not yet been published (or scraped)
  - `500` if this crappy code messed up
  - `200` with a json summary of the cases today if successful

- **Most Recent Daily Numbers** [get]

  `/api/v1/recent`

  Telegram bot function: `/recent`

  Will return:
  - `204` if data has not yet been published (or scraped)
  - `500` if this crappy code messed up
  - `200` with a json summary of the cases for that day if successful, e.g.
      ```json
      {"Date":"2020-10-02T00:00:00Z","Campus":4,"City":2,"Staff":0}
      ```

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

- **Total Cases** [get]

  `/api/v1/totals`

  Will return:
  - `204` if data has not yet been published (or scraped)
  - `500` if this crappy code messed up
  - `200` with a json summary of the total cases if successful, e.g.
      ```json
      {
        "starting date": "Thu, 01 Oct 2020 00:00:00 UTC",
        "ending date": "Thu, 08 Oct 2020 00:00:00 UTC",
        "staff total": 5,
        "student total": 117,
        "total cases": 122
      }
      ```

- **Average Daily Cases** [get]

  `/api/v1/average`

  Optional Parameters:
    - `days`: e.g. `7` | if not provided, will average across the whole dataset
    > (example returns average over the 7 most recently provided daily data)

  Will return:
  - `204` if data has not yet been published (or scraped)
  - `500` if this crappy code messed up
  - `200` with a json summary of the total cases if successful, e.g.
      ```json
      {
        "average cases": 23,
        "ending date": "Sun, 11 Oct 2020 00:00:00 UTC",
        "starting date": "Sun, 04 Oct 2020 00:00:00 UTC"
      }
      ```

- **Complete Raw** [get]

  `/api/v1/raw`

  Will return:
  - `204` if data has not yet been published (or scraped)
  - `500` if this crappy code messed up
  - `200` with a json summary of the total cases as given in the table if successful, e.g.
      ```json
      [  {"Date":"2020-10-01T00:00:00Z","Campus":1,"City":2,"Staff":0},
         {"Date":"2020-10-02T00:00:00Z","Campus":4,"City":2,"Staff":0},
         {"Date":"2020-10-03T00:00:00Z","Campus":3,"City":1,"Staff":0},
         {"Date":"2020-10-04T00:00:00Z","Campus":4,"City":4,"Staff":0}  ]
      ```
