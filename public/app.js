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


    App.ajaxify = () => {
        [].forEach.call(
            document.querySelectorAll("a[data-ajax]"), 
            (anchorElement) => {
                anchorElement.addEventListener("click", (e) => {
                    e.preventDefault();
                    history.pushState(
                        {"state":"state"},
                        "title",
                        anchorElement.pathname
                    );
                    App.router(anchorElement.pathname);
                });
                anchorElement.removeAttribute("data-ajax");
            }
        );
    };

    return App;
})(App || {});
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
var View = ((View) => {

    View.index = (data) => {
        let ret = `<div class="tabbar">
            <nav>
            <ul class="papertabs">
            `;

        data.sections.forEach((section) => {
            ret+=`<li><a href="#">${section.title}
            <span class="paperripple">
                        <span class="circle"></span>
                    </span></a></li>`;
        });

            
        ret += `</ul></nav></div><a href="/test" data-ajax>test</a>`;
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
            html = `Missing template "${template}"`;
        }
        document.body.innerHTML = html;
        App.ajaxify();
    };

    return View;
})(View || {});
App.bootstrap();
