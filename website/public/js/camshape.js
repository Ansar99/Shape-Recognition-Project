'use strict';
const socket = io();

// The vue instance
const vm = new Vue({
    el:"#main",
    data:{
        str:"",
    },

    // Vue cycle created: all functions are available after the vue object is created.
    methods:{
        backToMain: function(){
            window:history.go(-1);
        },
        startCamera: function(){
            socket.emit("startCamera")
        },
        stopCamera: function(){
            socket.emit("stopCamera")
        }
    }

})
