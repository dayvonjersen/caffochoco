function isVisible(element) {
    var parentElement = element.parentElement;
    var offsetTop = element.offsetTop;
    while(parentElement) {
        offsetTop += parentElement.offsetTop;
        parentElement = parentElement.parentElement;
    }
    var elementHeight = element.getBoundingClientRect().height;
    return (window.pageYOffset >= offsetTop - elementHeight) && 
    (window.pageYOffset <= offsetTop + elementHeight);
}
