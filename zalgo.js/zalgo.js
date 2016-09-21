const fs = require('fs');
const path = require('path');
const process = require('process');
const htmlparser2 = require('htmlparser2');
const nodeSass = require('node-sass');

let kebabￚcase2PascalCase = (txt) => {return s="",txt.split('-').forEach((t)=>s+=t[0].toUpperCase()+t.substr(1)),s}

function Devour(htmlfile) {

var element = {
    name: "",
    tag: "",
    style: "",
    template: "",
    script: "",
};

var inElement, currentElement;
let Parser = new htmlparser2.Parser( 
    {
        onopentag: (name, attr) => {
            // console.log("opentag:",name,attr);

            switch(name) {
            case "element":
                if(! ("name" in attr))  throw new Error("Missing ``name'' attribute for <element>");
                element.name = kebabￚcase2PascalCase(attr.name);
                element.tag = attr.name;
            case "style":
            case "template":
            case "script":
                inElement = name;
                break;
            default:
                if(inElement == "template") {
                    element.template += "<"+name;
                    for(prop in attr) {
                        element.template += ` ${prop}="${attr[prop]}"`;
                    }
                    element.template += ">";
                }
            }
        },
        onclosetag: (name) => {
            // console.log("closetag:", name);

            switch(name) {
            case "element":
            case "style":
            case "template":
            case "script":
                return;
            }

            if(inElement == "template") {
                element.template += `</${name}>`;
            }
        },
        ontext: (text) => {
            // console.log("text:",text);

            if(!text.trim())  return;
            element[inElement] += text;
        }
    }, {
        decodeEntities: true,
        xmlMode: false,
    }
);

Parser.write(fs.readFileSync(htmlfile).toString());

element.template = element.template.replace(/[\\"']/g, '\\$&').replace(/\u0000/g, '\\0').replace(/\n+/g, '').replace(/\s+/g, ' ');

return {
    js: `(function(){
${element.script}
Register("${element.name}", "${element.tag}", Z, "${element.template}");
})();`,
    css: nodeSass.renderSync({data:`${element.tag} { ${element.style} }`}).css.toString()
}
};

// console.log(module.paths);
let staticFiles = [
    "webcomponents.js/CustomElements.min.js",
    "../hecomes.js",
];
let staticBuffers = [];
staticFiles.forEach((file, idx) => {
    module.paths.forEach((dir) => {
        if(staticBuffers.length > idx)  return;
        try {
            fs.statSync(dir + path.sep + file);
            staticBuffers.push(fs.readFileSync(dir + path.sep + file));
        } catch(e) {}
    });
    if(staticBuffers.length <= idx)  throw new Error(file, "not found in module.paths");
});

function Bundle(elements, jsFile = "build/build.js", cssFile = "build/build.css") {
    let js = "";
    let css = "";
    elements.forEach((element) => {
        js += element.js;
        css += element.css;
    });

    let bufs = staticBuffers.slice(0);
    bufs.push(Buffer.from(js));
    fs.writeFileSync(jsFile, Buffer.concat(bufs));

    fs.writeFileSync(cssFile, css);
}

module.exports = {Devour,Bundle};
