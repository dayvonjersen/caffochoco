var View = ((View) => {

    View.index = (data) => {
        let ret = `<div class="tabbar">
            <nav>
            <ul class="papertabs">
            `;

        data.sections.forEach((section, idx) => {
            ret+=`<li><a href="#"${idx == 0 ? ' class="active"' : ''}>${section.title}
            <span class="paperripple">
                        <span class="circle"></span>
                    </span></a></li>`;
        });

            
        ret += `</ul></nav></div><a href="/test" data-ajax>test</a>`;
        setTimeout(()=>{
            let scriptElement = document.createElement("script");
            scriptElement.type = "text/javascript";
            scriptElement.src = "google-io-tabbed-nav.js";
            document.head.appendChild(scriptElement);
        }, 1000);
        return ret;
    };

    return View;
})(View || {});
