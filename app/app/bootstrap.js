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
