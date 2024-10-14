function toggleAccordion(event) {
    const accordionItemHeader = event.currentTarget;
    const accordionItemBody = accordionItemHeader.nextElementSibling;

    if (accordionItemBody.style.maxHeight) {
        accordionItemBody.style.maxHeight = null;
    } else {
        accordionItemBody.style.maxHeight = accordionItemBody.scrollHeight + "px";
    }
}

window.onscroll = function() {
    const backToTopButton = document.getElementById('backToTop');
    if (document.body.scrollTop > 20 || document.documentElement.scrollTop > 20) {
        backToTopButton.style.display = "block";
    } else {
        backToTopButton.style.display = "none";
    }
};

document.getElementById('backToTop').onclick = function() {
    window.scrollTo({
        top: 0,
        behavior: 'smooth'
    });
};
