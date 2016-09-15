"use strict";
var App = (function(A){
    A.bootstrap = function() {
window.addEventListener("load", () => {
    var data;

    fetch("/api/")
    .then((r) => r.json())
    .then((obj) => {
        data = obj;
        window.dispatchEvent(new Event("hashchange"));
    });

    window.addEventListener("hashchange", () => {
        [].forEach.call(document.querySelectorAll(".page"), (el) => el.classList.remove("visible"));

        var page;
        switch(decodeURI(window.location.hash).split("/")[1]) {
        case undefined: 
        case '':             
            page = 'all-tracks';
            document.querySelector("."+page).innerHTML = `<h1>${ data.hello }</h1>`;
            break;
        case 'single-track': page = 'single-track'; break;
        default:             page = 'error'; break;
        }

        document.querySelector("."+page).classList.add("visible");
    });
});

/*
[].forEach.call(document.querySelectorAll("a[data-ajax]"), (aElement) => {
    aElement.addEventListener("click", (e) => {
        e.preventDefault();
        history.pushState(
            {
                "some_data_here": 42,
                "template": aElement.dataset.ajax
            },
            "some_title_here?",
            aElement.href
        );
        render(aElement.dataset.ajax);
    });
});

window.addEventListener("popstate", (e) => render(e.state ? e.state.template : ""));
*/
};
    return A;
}(App || {}));
