'use strict';
const socket = io();

// The vue instance
const vm = new Vue({
    el:"#main",
    data:{

        correctShape:"",
        guesses:"",
        success:"",
        failed:"",
        listOfGuesses:[]

    },

    // Vue cycle created: all functions are available after the vue object is created.
    created: function() {


        socket.on("guesses", function(g){
            this.guesses = "Your guesses are: " + g.guesses;
            this.listOfGuesses = g.guesses.split(", ");

        }.bind(this));

        socket.on("correctAnswerToClient", function(answer){

            this.correctShape =  answer.rightGuess[0];

            console.log("listOfGuesses: "+this.listOfGuesses);
            console.log("correctshape: "+this.correctShape);
            console.log("list length is: " + this.listOfGuesses.length);
            console.log("guesses:"+ this.guesses);

            if((this.listOfGuesses.length - 1) < 4){
                this.failed="";

                let successElem = document.getElementById("success");
                let failedElem = document.getElementById("failed");
                if (this.listOfGuesses.includes(this.correctShape)) {


                    successElem.style.display="inline";
                    failedElem.style.display="none";
                    this.correctShape = "The right answer is: " + answer.rightGuess[0];
                    this.success = "You guessed the correct shape!";
                }
                else {

                    this.success="";
                    failedElem.style.display="inline";
                    successElem.style.display="none";

                    this.correctShape = "The right answer is: " + answer.rightGuess[0];
                    this.failed = "You didn't guess the correct shape...";
                }
            }
            else{

                let failedElem = document.getElementById("failed");
                this.success="";
                this.correctShape="";
                this.guesses="";
                failedElem.style.display="inline";
                this.failed = "Your image contained to many shapes, the maximum is 3";

            }

            }.bind(this));

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
