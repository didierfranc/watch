import http from "http";

http
  .createServer((req, res) => {
    res.end("hello");
  })
  .listen("8080", () => console.log("listening on port 8080"));

process.on("SIGINT", function () {
  process.exit(0);
});
