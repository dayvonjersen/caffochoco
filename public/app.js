var API = (function(API){

    API.ROUTE = "/api";

    API.fetchJSON = (route, callback) => {
        return fetch(API.ROUTE+route)
            .then((r) => r.json())
            .then((data) => callback(data));
    };

    return API;
}(API || {}));
var App = ((App) => {

    /**
     *
     */
    App.bootstrap = () => {
        window.addEventListener("load", () => {
            API.fetchJSON("/", (indexData) => {
                window.dispatchEvent(new Event("hashchange"));
            });
            
            window.addEventListener("hashchange", App.router);
        });
    };

    return App;
})(App || {});
var App = ((App) => {

    App.router = () => {
        [].forEach.call(document.querySelectorAll(".page"), (el) => el.classList.remove("visible"));

        var page;
        switch(decodeURI(window.location.hash).split("/")[1]) {
        case undefined: 
        case '':             
            page = 'all-tracks';
            document.querySelector("."+page).innerHTML = `<h1>hello world</h1>`;
            break;
        case 'single-track': page = 'single-track'; break;
        default:             page = 'error'; break;
        }

        document.querySelector("."+page).classList.add("visible");
    };

    return App;
})(App || {});
App.bootstrap();
