# COVID Statistics API for LU
> Very unofficial, but uses the data provided by the university publicly

To see the statistics, see [here](https://portal.lancaster.ac.uk/intranet/cms/coronavirus/covid-19-statistics)

This API exists only to allow for easier manipulation of the data.

### API

- **Cases Today** [get]

  `/api/v1/today`

  Will return:
  - `204` if today's data has not yet been published (or scraped)
  - `500` if this crappy code messed up
  - `200` with a json summary of the cases today if successful

- **Total Cases Summary** [get]

  `/api/v1/summary`

  Will return:
  - `204` if data has not yet been published (or scraped)
  - `500` if this crappy code messed up
  - `200` with a json summary of the total cases if successful, e.g.
      ```json
      {"Staff Cases":5,"Student Cases":117,"Total Cases":122}
      ```
