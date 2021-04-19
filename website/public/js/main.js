'use strict';
const socket = io();


const vm = new Vue({
    el:"#main",
    data:{
        str:"",
    },

    created: function() {

        socket.on('runGo', function(result){
            //TODO: läs output från fil och skicka till sidan hemsidan
            console.log("här är resultatet " + result.output);
            this.str = result.output;
        }.bind(this));

    },

    methods:{
        retractSidebar: function(){
            let sidebar = document.getElementById("sideBar");
            sidebar.style.display="none";


        },

        extendSidebar: function(){
            let sidebar = document.getElementById("sideBar");
            sidebar.style.display="block";


        }
    }

})
