#!/usr/bin/env node

const http = require('http')
const port = 80

const requestHandler = (request, response) => {
  console.log(request.url)
  http.get("http://test/", res => {
    res.setEncoding("utf8");
    let body = "";
    res.on("data", data => {
      body += data;
    });
    res.on("end", () => {
      response.end(body)
    });
  });
}

const server = http.createServer(requestHandler)

server.listen(port, (err) => {
  if (err) {
    return console.log('something bad happened', err)
  }
  console.log(`server is listening on ${port}`)
})
