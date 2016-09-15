var App = ((App) => {

    App.router = (pathname) => {
        if(pathname in App.routes) {
            var r = App.routes[pathname];
            API.fetchJSON(r.endpoint, (data) => {
                View.render(r.template, data);
            });
        } else {
            View.render("notfound");
        }
    };

    return App;
})(App || {});
