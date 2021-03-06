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
        window.addEventListener('beforeunload',function(){
            var paths = {
                image : document.getElementById("inputPicture").getAttribute("src"),
                shaped_image : document.getElementById("outputPicture").getAttribute("src")
            };
            socket.emit("delete image", paths);
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
