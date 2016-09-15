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
