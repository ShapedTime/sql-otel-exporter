# SQL Query OpenTelemetry Exporter
This application is designed to run a given SQL select query on a specified database and export the results to an OpenTelemetry collector. The data can be exported as either a gauge or a histogram.

## How to Use
1. Clone the repository to your local machine.
2. Ensure you have Docker installed and running.
3. Run `docker-compose up` in the root directory of the project. This will start the necessary services defined in the `docker-compose.yml` file.

## Configuration
The application is configured using a `config.yml` file. Here's an example of how to structure your `config.yml`:
```yaml
driver: postgres
datasource: host=localhost user=test password=test123 dbname=test sslmode=disable
queries:
  - name: my_query
    sql: SELECT COUNT(*) FROM my_table
    trackingColumn: id
    defaultTrackingValue: "0"
    intervalSeconds: 60
    metricName: my_metric
    valueColumn: count
    description: Description of my metric
    attributeColumns: [column1, column2]
    staticAttributes:
      attribute1: value1
      attribute2: value2
    dataType: gauge
    buckets: [0, 10, 100, 1000]
    ignoreError: false
```

Here's what each field means:

- `driver`: The database driver to use. Currently, only postgres is supported.
- `datasource`: The connection string for the database.
- `queries`: An array of queries to run. Each query has the following fields:
    - `name`: A unique name for the query.
    - `sql`: The SQL select query to run.
    - `trackingColumn`: The column to use for tracking the last processed row.(Optional)
    - `defaultTrackingValue`: The default value for the tracking column.(Optional)
    - `intervalSeconds`: How often to run the query, in seconds.
    - `metricName`: The name of the metric to export.
    - `valueColumn`: The column to use for the metric value.
    - `description`: A description of the metric.
    - `attributeColumns`: An array of columns to use for metric attributes.
    - `staticAttributes`: A map of static attributes to add to each metric.
    - `dataType`: The type of the metric. Can be either gauge or histogram.
    - `buckets`: An array of bucket boundaries for histogram metrics.
    - `ignoreError`: Whether to ignore errors when running the query or exporting the metric.

Please replace the values with your actual database connection details and queries.

# Contributing
Contributions are welcome! Please submit a pull request or create an issue to get started.