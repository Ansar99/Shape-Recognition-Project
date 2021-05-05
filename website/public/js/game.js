'use strict';
const socket = io();

// The vue instance
const vm = new Vue({
    el:"#main",
    data:{
        shapes: ["circle", "square", "rectangle", "triangle", "pentagon","hexagon","heptagon","octagon","nonagon"],
    },

    // Vue cycle created: all functions are available after the vue object is created.
    created: function() {
        console.log(this.shapes);
        window.onload = function () {

            // Definitions
            var canvas = document.getElementById("paint-canvas");
            var context = canvas.getContext("2d");
            var boundings = canvas.getBoundingClientRect();

            // Specifications
            var mouseX = 0;
            var mouseY = 0;
            context.strokeStyle = 'black'; // initial brush color
            context.lineWidth = 2; // initial brush width
            var isDrawing = false;



            // Mouse Down Event
            canvas.addEventListener('mousedown', function(event) {
                setMouseCoordinates(event);
                isDrawing = true;

                // Start Drawing
                context.beginPath();
                context.moveTo(mouseX, mouseY);
            });

            // Mouse Move Event
            canvas.addEventListener('mousemove', function(event) {
                setMouseCoordinates(event);

                if(isDrawing){
                    context.lineTo(mouseX, mouseY);
                    context.stroke();
                }
            });

            // Mouse Up Event
            canvas.addEventListener('mouseup', function(event) {
                setMouseCoordinates(event);
                isDrawing = false;
            });

            // Handle Mouse Coordinates
            function setMouseCoordinates(event) {
                mouseX = event.clientX - boundings.left;
                mouseY = event.clientY - boundings.top;
            }

            // Handle Clear Button
            var clearButton = document.getElementById('clear');

            clearButton.addEventListener('click', function() {
                context.clearRect(0, 0, canvas.width, canvas.height);
            });

            // Handle Save Button
            var saveButton = document.getElementById('save');

            saveButton.addEventListener('click', function() {
                var imageName = prompt('Please enter image name');
                var canvasDataURL = canvas.toDataURL('image/jpeg',1.0);
                var a = document.createElement('a');
                a.href = canvasDataURL;
                a.download = imageName || 'drawing';
                a.click();
            });
        };


    },

    methods:{
        // Retracts the sidebar
        retractSidebar: function(){
            let sidebar = document.getElementById("sideBar");
            sidebar.style.display="none";

        },
        // Extends the sidebar
        extendSidebar: function(){
            let sidebar = document.getElementById("sideBar");
            sidebar.style.display="block";
        },
        backToMain: function(){
            window:history.go(-1);
        }
    }

})
