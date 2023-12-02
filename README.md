# Logs(log search)

> Was a hackathon build. No longer maintained

Original [Project Requirement Document](TASK.md)

In this project, I build a CLI tool with two sub-commands.

- The `serve` sub-command starts a server on port `3000` to ingest data.
- The `query` sub-command queries the data from the Elasticsearch.

The database of choice is Elasticsearch. The reason for choosing Elasticsearch is mentioned in the [Design Decisions and FAQ](#design-decisions-and-faq) section. The CLI tool is written in Go.

Project Demo: [link.utkarshchourasia.in/logs-demo](https://link.utkarshchourasia.in/logs-demo)

## Requirements

- [Go](https://go.dev/dl/) 1.20 or higher
- [Docker](https://docs.docker.com/engine/install/) and Docker Compose
- [Git](https://git-scm.com/)
- [Postman](https://www.postman.com/) *optional*

## Starting Server

```bash
# Getting files
git clone https://github.com/dyte-submissions/november-2023-hiring-JammUtkarsh
cd november-2023-hiring-JammUtkarsh
git submodule init
git submodule update

# Setting up containers
cd elastic
docker compose up setup
docker compose up elasticsearch -d
cd ..

# Setting enviroment variables
export $(cat .env | xargs) # or
# export ELASTIC_USERNAME=elastic && export ELASTIC_PASSWORD=changeme

# Setting the mappings of Elasticsearch
go run elasticMapping.go # needs to be run only once after setting up the containers

# running the server | optinally, you can build and then run as well.
go run . serve 
```

Now you can ingest the data on port `3000`. Querying can be done using the CLI tool.

*You can also refer to Postman collection for further inspection and documentation for curl commands.*

I also performed a stress test on the server using Postman with the following hardware:

```text
Machine: MacBook Pro (16-inch, 2019)
OS: macOS 14.1.1 (23B81)
Processor: 2.6 GHz 6-Core Intel Core i7
RAM: 32 GB 2667 MHz DDR4

NOTE: Half of the resources were dedicated to run Docker VM.
```

The Report can be accessed from [here](https://link.utkarshchourasia.in/postman-report) or `./PostmanReports/'Ingestion Endpoint | Postman Peformance Report.pdf'`

## Build and Using CLI

Since the original tab will be occupied by the server, you can open a new tab and run the following commands.

```bash
# Assuimg you are in the same directory as before. If not, then run the first 2 commands from the previous section.
go build -o logs .
./logs --help # to see the help
```

### Sample queries

```bash
# If you want search for an error in a particular resource
./logs query --level=error --resource=server-1234

# if you want search all the logs of a particular resource in a particular time range
./logs query --resource=server-1234 --time "2021-08-01T00:00:00 2021-08-02T00:00:00" # NOTE: the range is separated by a space(' ')

# to view all the details for a particular message string
./logs query --message="token issued" --all
```

## Things that I didn't know before

Some of my major learnings are:

- How to stress test an API in Postman.
- How to perform insert and query in Elasticsearch.
- Using up Git Submodules.

## Design Decisions and FAQ

1. Why use Elasticsearch?

   - Based on my experience as s Technical Writer at SigNoz(Their whole product/business is build around telemetry collection), I was familiar with the term 'Full-Text Search' and ELK stack.
   - The Data types in *Log Data Format* is either mostly `enums`(a set of pre-defined values), `Text` and dates. So my initial instinct was to use a relational database something like Postgres(it supports these data types).
   - But the [product requirement document](https://dyte.notion.site/dyte/SDE-1-and-SDE-Intern-Assignment-6b7a7f324dc0450381b0fdb771a8ec40) clearly say that being able to store the logs and search them is an important feature.
   - Storing data was not an issue, being of search through it was. So based on the [search engine ranking](https://db-engines.com/en/ranking/search+engine) and the fact that I have never tried Elasticsearch before, I decided to use Elasticsearch.

2. Why use Elasticsearch only and not the whole ELK stack?

   - Again, the product requirement document clearly states that being able to search through the logs is an important feature. There is no mention of visualizing the data(Kibana) or collecting data from multiple sources(Logstash) is required. So I decided to use Elasticsearch only.

3. Elasticsearch has a [go-sdk](https://github.com/elastic/go-elasticsearch/), why not use that?

   - Since I had limited time, HTTP requests gave me a way to prototype the app fast. I copied the curl request from the official [REST API documentation](https://www.elastic.co/guide/en/elasticsearch/reference/current/rest-apis.html) and converted it into Go code.
   - Also, some modules are [not fully complete](https://www.elastic.co/guide/en/elasticsearch/client/go-api/current/typedapi.html) while other's didn't have Go examples in the documentation. So I decided to use HTTP requests for now.
   - **The only trade off with this approach might be [speed](https://github.com/elastic/go-elasticsearch/issues/75#issuecomment-516711935) and security practices done by the SDK.**

## Things I didn't due to time constraints

- Connecting Kafka to improve the ingestion rate.
- Containerizing the ingestion server.
- [Charmbracelet](https://github.com/charmbracelet) to better present the output.
- Make atomic commits.
- Allow user to enter username and password for Elasticsearch and store this info somewhere in $HOME
- Improved code structure.

## About Candidate

- **Name**: Utkarsh Chourasia
- **Resume**: [links.utkarshchourasia.in/resume](https://links.utkarshchourasia.in/resume)
- **LinkedIn**: [linkedin.com/in/jammutkarsh](https://www.linkedin.com/in/jammutkarsh/)
- **GitHub**: [github.com/JammUtkarsh](https://github.com/JammUtkarsh/)

Fun Fact: Built this project in about 18 engineering hours.
