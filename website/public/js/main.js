'use strict';
const socket = io();

// The vue instance
const vm = new Vue({
    el:"#main",
    data:{
        str:"",
    },

    // Vue cycle created: all functions are available after the vue object is created.
    created: function() {

        // Receive output message from server and bind it to this.str
        socket.on('runGo', function(result){
            //TODO: läs output från fil och skicka till sidan hemsidan
            console.log("här är resultatet " + result.output);
            this.str = result.output;
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


        }
    }

})
