/* jslint node: true */
/* eslint-env node */
'use strict';


var pid;

// Require express, socket.io, vue, path and child_process
const express = require('express');
const app = express();
const http = require('http').Server(app);
const io = require('socket.io')(http); //FIXME: Remove ?
const path = require('path');
const crypto = require("crypto");
const parse = require('node-html-parser').parse;

const runGo  = require('child_process');

// Require filesystem
const fs = require('fs');

// Require Midleware for handling multipart/form-data "multer"
const multer = require("multer");

// Pick arbitrary port for server
const port = 3000;
app.set('port', (process.env.PORT || port));

// Serve static assets from public/
app.use(express.static(path.join(__dirname, 'public/')));

// Serve static assets from shapedImages/
app.use(express.static(path.join(__dirname, 'shapedImages/')));

// Serve vue from node_modules as /vue
app.use('/vue',
        express.static(path.join(__dirname, '/node_modules/vue/dist/')));

// Serve index.html directly as root page
app.get('/', function(req, res) {
    res.sendFile(path.join(__dirname, 'views/index.html'));
});

// Serve upload.html as /upload
app.get('/upload', function(req, res) {
    res.sendFile(path.join(__dirname, 'views/upload.html'));
})

app.get('/camera', function(req, res) {
    res.sendFile(path.join(__dirname, 'views/camshape.html'));
})

app.get('/game', function(req, res) {
    res.sendFile(path.join(__dirname, 'views/game.html'));
})


// Error handler for the post request.
const handleError = (err, res) => {
    res
        .status(500)
        .contentType("text/plain")
        .end("Oops! Something went wrong when attempting to upload the file!");
};

const upload = multer({
    dest: "images/"
    // Might also want to set some limits: https://github.com/expressjs/multer#limits
});

// Handles a post request
app.post(
    "/upload",

    //Accepts a single file which will be stored in req.file
    upload.single("file" /* name attribute of <file> element in your form */),

    (req, res) => {
        
        //Generates unique id for the image
        var id = crypto.randomBytes(12).toString("hex");
        console.log("Image id: " + id);

        const tempPath = req.file.path;
        //Saves the uploaded image with unique id name
        var targetPath = path.join(__dirname, "./public/images/" + id + ".jpg");

        // We simply check if the file is in .jpg format
        if (path.extname(req.file.originalname).toLowerCase() === ".jpg") {

            fs.rename(tempPath, targetPath, err => {
                if (err) return handleError(err, res);

                //execSync(Command) executes Command synchronously, that way exec(cat....) won't execute before the Go program
                runGo.execSync("go run ../src/shapeitup/main.go public/images/" + id + ".jpg " + "./shapedImages/shaped_" + id + ".jpg > output.txt",(error,stdout,stderr) => {  //FIXME: Remove output ?

                });
                fs.readFile('views/upload.html', 'utf8', (err,html)=>{
                    if(err){
                       throw err;
                    }
                    //Parses the upload.html file and adds the correct images (unique id name) to be displayed
                    const root = parse(html);
                    const outputContainer = root.querySelector('#outputContainer');
                    //outputContainer.appendChild('<img src="images/image.jpg" id="inputPicture">');
                    //outputContainer.appendChild('<img src="shapedimage2.jpg" id="outputPicture">');  //FIXME: appendChild would be better
                    outputContainer.set_content('<img src="images/' + id + '.jpg" id="inputPicture"><img src="shaped_' + id + '.jpg" id="outputPicture">');
                    res
                        .status(200)
                        .contentType("text/html")
                        //Sends the edited upload.html as response
                        .send(root.toString())
                  });

            });
        }
        // If it's not we return an error message
        else {
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

// Might not need this request
app.delete('/remove',function(req,res){
    console.log("app.Delete called!");
});

io.on('connection', (socket) => {
    //Deletes the images when leaving a page
    socket.on('delete image', (paths) => {  //FIXME: GÖR SNYGGARE
        fs.unlink(path.join(__dirname, "./public/" + paths.image), (err => {
            if (err) console.log(err);
            else {
                console.log("Deleted file: " + paths.image);
            }
          }));
        fs.unlink(path.join(__dirname, "./shapedImages/" + paths.shaped_image), (err => {
            if (err) console.log(err);
            else {
                console.log("Deleted file: " + paths.shaped_image);
           }
        }));
    });
    socket.on("startCamera", function(){
        pid = runGo.exec("go run ../src/shapeitup/cameradetect.go", (error,stdout,stderr) => {
        });
    })
    socket.on("stopCamera", function(){
        console.log("HELLOOOOOO");
        pid.kill();
    })
});

const server = http.listen(app.get('port'), function() {
    console.log('Server listening on port ' + app.get('port'));
});
