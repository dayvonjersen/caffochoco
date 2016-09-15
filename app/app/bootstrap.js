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
            "/test": {
                template: "release",
                endpoint: "/",
            },
        };
        window.addEventListener("load", () => {
            API.fetchJSON("/", (data) => {
                // add routes to releases...
                App.router(location.pathname);
            });
            
            window.addEventListener("popstate", () => {
                App.router(location.pathname);
            });
        });
    };

    return App;
})(App || {});
