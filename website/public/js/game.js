'use strict';
const socket = io();

// The vue instance
const vm = new Vue({
    el:"#main",
    data:{
        clueAndShape: {
            1:["circle", "Hint: I have no edges..."],
            2:["triangle", "Hint: I can be used to create all standard polygons"],
            3:["rectangle","Hint: I am a special case of a parallelogram"],
            4:["square","Hint: I am one side of the six side of a Minecraft block"],
            5:["pentagon","Hint: My shape is the same as a military base in the US"],
            6:["hexagon","Hint: I am the shape of a honeycomb cell"]
        },
        correctShape:"Something went wrong!",
        guesses:""

    },

    // Vue cycle created: all functions are available after the vue object is created.
    created: function() {

        this.correctShape = this.clueAndShape[Math.floor(Math.random() * 6) + 1];
        console.log("CorrectShape in game.js 24 " + this.correctShape[1]);


        socket.on("guesses", function(g){
            this.guesses = g.guesses;

        }.bind(this));

        socket.emit("correctAnswerToServer", {
            rightGuess: this.correctShape
        });

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
        },
        sendGuess: function(){
            console.log("HÄR ÄR JAG");
            socket.emit("forwardGuesses", {
                guess:this.guesses,
                correcAnswer: this.correctShape
            });
        }
    }

})
