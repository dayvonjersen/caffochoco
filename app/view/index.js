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
