'use strict';

const vm = new Vue({
    el:"#main",
    data:{
        sidebarD:''
    },

    created: function() {
        console.log("Demo");

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
