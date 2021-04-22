/* jslint node: true */
/* eslint-env node */
'use strict';

// Require express, socket.io, vue, path and child_process
const express = require('express');
const app = express();
const http = require('http').Server(app);
const io = require('socket.io')(http); //TA BORT ?
const path = require('path');
const crypto = require("crypto");
const parse = require('node-html-parser').parse;

const runGo  = require('child_process');
const cat  = require('child_process');

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

        var id = crypto.randomBytes(12).toString("hex"); //ANTON    
        console.log(id);  //ANTON

        const tempPath = req.file.path;
        var targetPath = path.join(__dirname, "./public/images/" + id + ".jpg");  //ändrade till var ANTON

        // We simply check if the file is in .jpg format
        if (path.extname(req.file.originalname).toLowerCase() === ".jpg") {

            fs.rename(tempPath, targetPath, err => {
                if (err) return handleError(err, res);

                //execSync(Command) executes Command synchronously, that way exec(cat....) won't execute before the Go program
                runGo.execSync("go run ../src/shapeitup/main.go public/images/" + id + ".jpg " + id + "  > output.txt",(error,stdout,stderr) => {  //ANTON SKITA I OUTPUT ?

                });

                fs.readFile('views/upload.html', 'utf8', (err,html)=>{
                    if(err){
                       throw err;
                    }
                 
                    const root = parse(html);
                    const outputContainer = root.querySelector('#outputContainer');
                    //outputContainer.appendChild('<img src="images/image.jpg" id="inputPicture">');
                    //outputContainer.appendChild('<img src="shapedimage2.jpg" id="outputPicture">');  //VARFÖR FUNKAR INTE APPEND CHILD
                    outputContainer.set_content('<img src="images/' + id + '.jpg" id="inputPicture"><img src="shaped_' + id + '.jpg" id="outputPicture">');
                    res
                        .status(200)
                        .contentType("text/html")
                        //.sendFile(path.join(__dirname, 'views/upload.html'))
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


const server = http.listen(app.get('port'), function() {
    console.log('Server listening on port ' + app.get('port'));
});
