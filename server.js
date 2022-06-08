const http = require('http');

http.createServer((req, res) => {
    console.log(req.url)
    res.statusCode = 200;
    res.write('Hello!');
    res.end();
}).listen(3000, () => {
    console.log('Listening on 3000')
})
