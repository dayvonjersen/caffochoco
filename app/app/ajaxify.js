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
