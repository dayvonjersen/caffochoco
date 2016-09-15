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
