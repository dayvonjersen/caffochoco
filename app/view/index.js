var View = ((View) => {

    View.index = (data) => {
        return `<h1>hello ${data.hello}</h1><a href="/test" data-ajax>test</a>`;
    };

    return View;
})(View || {});
