var http = require('http');

http
  .createServer(function (req, res) {
    res.writeHead(200, {
      'Transfer-Encoding': 'chunked',
      'Content-Type': 'text/event-stream',
      'Access-Control-Allow-Origin': '*',
    });

    setInterval(function () {
      var packet = 'data: {"message":"' + new Date().getTime() + '"}\n\n';
      res.write(packet);
    }, 1000);
  })
  .listen(8080);

// var http = require('http');
// var requests = [];

// var server = http.Server(function (req, res) {
//   var clientIP = req.socket.remoteAddress;
//   var clientPort = req.socket.remotePort;

//   res.on('close', function () {
//     console.log('client ' + clientIP + ':' + clientPort + ' died');

//     for (var i = requests.length - 1; i >= 0; i--) {
//       if (requests[i].ip == clientIP && requests[i].port == clientPort) {
//         requests.splice(i, 1);
//       }
//     }
//   });

//   res.writeHead(200, {
//     'Content-Type': 'text/event-stream',
//     'Access-Control-Allow-Origin': '*',
//     'Cache-Control': 'no-cache',
//     Connection: 'keep-alive',
//   });

//   requests.push({ ip: clientIP, port: clientPort, res: res });

//   res.write(': connected.\n\n');
// });

// server.listen(8080);

// setInterval(function test() {
//   broadcast('poll', 'test message');
// }, 2000);

// function broadcast(rtype, msg) {
//   var lines = msg.split('\n');

//   for (var i = requests.length - 1; i >= 0; i--) {
//     requests[i].res.write('event: ' + rtype + '\n');
//     for (var j = 0; j < lines.length; j++) {
//       if (lines[j]) {
//         requests[i].res.write('data: ' + lines[j] + '\n');
//       }
//     }
//     requests[i].res.write('\n');
//   }
// }
