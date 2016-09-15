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
        App.routes = {
            "/": {
                template: "index",
                endpoint: "/",
            },
        };
        window.addEventListener("load", () => {
            API.fetchJSON("/", (data) => {
                // add routes to releases...
                window.dispatchEvent(new Event("hashchange"));
            });
            
            window.addEventListener("hashchange", App.router);
        });
    };

    return App;
})(App || {});
var App = ((App) => {

    App.router = () => {
        if(location.pathname in App.routes) {
            var r = App.routes[location.pathname];
            API.fetchJSON(r.endpoint, (data) => {
                View.render(r.template, data);
            });
        } else {
            View.render("notfound");
        }
    };

    return App;
})(App || {});
var View = ((View) => {

    View.index = (data) => {
        return `<h1>hello ${data.hello}</h1>`;
    };

    return View;
})(View || {});
var View = ((View) => {

    View.release = (data) => {
        return `<h1>release</h1>`;
    };

    return View;
})(View || {});
var View = ((View) => {

    View.render = (template, data) => {
        let html;
        if(template in View) {
            html = View[template](data);
        } else {
            html =`Missing template "${template}"`;
        }
        document.body.innerHTML = html;
    };

    return View;
})(View || {});
App.bootstrap();
