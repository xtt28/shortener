
# xtt28/shortener

Instantly generate shortened links online with custom aliases.


## Features

- Highly performant and optimized
- Users can create custom aliases for links
- Responsive user interface
- Supports sharing via native share sheet
- Supports HTTPS
## Screenshots

![Whole Page, Success](https://i.imgur.com/VkIrLuy.jpeg)
![Success, Zoomed In](https://i.imgur.com/R2rTp4a.png)
![Failure](https://i.imgur.com/OQ8CTF1.png)
## Badges

[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://choosealicense.com/licenses/mit/)

## Run Locally

1. Clone the project: `git clone https://github.com/xtt28/shortener`
2. Switch to the project directory: `cd shortener`
3. Create your `.env` file (see the Environment Variables section)
4. Start the server: `go run .`
5. Go to <http://localhost:8080> in your web browser

## Running Tests

To run tests, run the following command:

```bash
  go test ./...
```


## Environment Variables

To run this project, you will need to add the following environment variables to your .env file:

`ROOT` - the URL to your instance

`PORT` - the port to run your instance on

### TLS

The following environment variables concern TLS:

`TLS_ENABLED` - whether to enable TLS

`TLS_CERT` - the path to the TLS certificate file

`TLS_KEY` - the path to the TLS key file

Please refer to the `template.env` file for an example.

## API Reference

#### Create new shortened link

```http
  POST /api/create
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `destination` | `string` | **Required**. The URL to make the short link resolve to. |
| `alias` | `string` | An alias that can be used in addition to your link's ID when sharing it. |

#### Redirect to short link destination

```http
  GET /v/${id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. The ID of the shortened link. |


## Tech

**Client:** Bootstrap, htmx

**Server:** Go, Gin, GORM, SQLite


## Contributing

Contributions are always welcome. Please:
1. write tests for your code;
2. run it through `go fmt` and `go vet`.

## License

[MIT](https://choosealicense.com/licenses/mit/)

