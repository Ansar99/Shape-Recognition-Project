/* jslint node: true */
/* eslint-env node */
'use strict';

// Require express, socket.io, and vue
const express = require('express');
const app = express();
const http = require('http').Server(app);
const io = require('socket.io')(http);
const path = require('path');

const runGo  = require('child_process');
const cat  = require('child_process');


// file system module från nodejs
const fs = require('fs');

//paket från npm för att hantera mutlipart/form-data
const multer = require("multer");

// Pick arbitrary port for server
const port = 3000;
app.set('port', (process.env.PORT || port));

// Serve static assets from public/
app.use(express.static(path.join(__dirname, 'public/')));

app.use(express.static(path.join(__dirname, 'shapedImages/')));


// Serve vue from node_modules as vue/
app.use('/vue',
        express.static(path.join(__dirname, '/node_modules/vue/dist/')));

// Serve index.html directly as root page
app.get('/', function(req, res) {
    res.sendFile(path.join(__dirname, 'views/index.html'));
});

app.get('/upload', function(req, res) {
    res.sendFile(path.join(__dirname, 'views/upload.html'));
})

const handleError = (err, res) => {
    res
        .status(500)
        .contentType("text/plain")
        .end("Oops! Something went wrong!");
};

const upload = multer({
    dest: "images/"
    // Might also want to set some limits: https://github.com/expressjs/multer#limits
});


app.post(
    "/upload",
    upload.single("file" /* name attribute of <file> element in your form */),
    (req, res) => {
        const tempPath = req.file.path;
        const targetPath = path.join(__dirname, "./public/images/image.jpg");

        if (path.extname(req.file.originalname).toLowerCase() === ".jpg") {
            fs.rename(tempPath, targetPath, err => {
                if (err) return handleError(err, res);

                res
                    .status(200)
                    .contentType("text/html")
                    .sendFile(path.join(__dirname, 'views/upload.html'))
                //exec.execSync("export GOPATH=$(../src)");
                //execSync exekverar kommandot synkront, dvs cat.exec körs efter att runGo.execsync är klar

                runGo.execSync("go run ../src/shapeitup/main.go public/images/image.jpg  > output.txt",(error,stdout,stderr) => {

                });

                io.on('connection',function(socket){
                    cat.exec("cat output.txt",(error,stdout,stderr) =>{
                        console.log(stdout);
                        socket.emit('runGo', { output: stdout} );
                    });

                });

            });
        } else {
            fs.unlink(tempPath, err => {
                if (err) return handleError(err, res);

                res
                    .status(403)
                    .contentType("text/plain")
                    .end("Only .jpg files are allowed!");
            });
        }
    }
);

app.delete('/remove',function(req,res){
    console.log("app.Delete called!");
});
const server = http.listen(app.get('port'), function() {
    console.log('Server listening on port ' + app.get('port'));
});
