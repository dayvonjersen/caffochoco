const fs = require('fs');
const htmlparser2 = require('htmlparser2');
const nodeSass = require('node-sass');

let snake2PascalCase = (txt) => {return s="",txt.split('-').forEach((t)=>s+=t[0].toUpperCase()+t.substr(1)),s}

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
                element.name = snake2PascalCase(attr.name);
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

function Bundle(elements, jsFile = "build/build.js", cssFile = "build/build.css") {
    let js = "";
    let css = "";
    elements.forEach((element) => {
        js += element.js;
        css += element.css;
    });
    
    fs.writeFileSync(jsFile, 
        Buffer.concat([
            fs.readFileSync("node_modules/webcomponents.js/CustomElements.min.js"),
            fs.readFileSync("hecomes.js"),
            Buffer.from(js)
        ])
    );

    fs.writeFileSync(cssFile, css);
}

module.exports = {Devour,Bundle};
