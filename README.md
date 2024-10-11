# genApiDocGo

**genApiDocGo** is an open-source tool written in Go for generating API documentation in Swagger's JSON format. It’s primarily designed for use with Express.js frameworks, allowing for automatic extraction and documentation of HTTP routes and methods.

We plan to add more languages and frameworks in the future.

## Features

- **Automatic Documentation**: Extracts API endpoint information directly from the code.
- **Supports Swagger JSON Format**: Generates Swagger-compatible JSON files for easy integration with tools like Swagger UI.
- **Express.js Compatibility**: Designed to work with Express.js-based projects, enabling quick documentation generation.

## Installation

> [!NOTE]
> You must have Go version 1.22 installed.  

At the moment to use it, you have to download and run it. 
```sh
go run ./src/ <path_project_to_doc> 
```

## Example

Given a route: 
```js
// @api_generate_doc
/* Description */
or 
/**
 * Description
 */
router.post('/', async (req, res) => {
  // Endpoint code here
});
```
**genApiDocGo** will extract and document the route, generating a Swagger JSON entry for it.

## Contributing

Contributions are welcome! Please submit a pull request or open an issue if you find any bugs or have feature requests.

## License

This project is licensed under the Apache-2.0 License. See the [LICENSE](LICENSE) file for details.
