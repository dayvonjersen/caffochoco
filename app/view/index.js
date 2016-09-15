var View = ((View) => {

    View.index = (data) => {
        return `<h1>hello ${data.hello}</h1>`;
    };

    return View;
})(View || {});
