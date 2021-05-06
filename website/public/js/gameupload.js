'use strict';
const socket = io();

// The vue instance
const vm = new Vue({
    el:"#main",
    data:{

        correctShape:"Something went wrong!",
        guesses:""

    },

    // Vue cycle created: all functions are available after the vue object is created.
    created: function() {


        socket.on("guesses", function(g){
            this.guesses = g.guesses;
        }.bind(this));

        socket.on("correctAnswer", function(answer){
            this.correctShape = answer.rightGuess
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
        }
    }

})
