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
