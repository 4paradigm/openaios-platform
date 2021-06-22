var http = require('http');

http
  .createServer(function (req, res) {
    res.writeHead(200, {
      'Transfer-Encoding': 'chunked',
      'Content-Type': 'text/event-stream',
      'Access-Control-Allow-Origin': '*',
    });

    setInterval(function () {
      var packet = 'event: hello_event\ndata: {"message":"' + new Date().getTime() + '"}\n\n\n\n';
      res.write(packet);
      res.flushHeaders();
    }, 1000);
  })
  .listen(9009);

  