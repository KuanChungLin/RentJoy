const galleryLeftBtn = document.querySelector(".gallery-leftbtn");
const galleryRightBtn = document.querySelector(".gallery-rightbtn");
const galleryArea = document.querySelector(".gallery-area");
const galleryItems = document.querySelectorAll(".gallery-single");
const galleryClipper = document.querySelector(".gallery-clipper");



let currentMarginLeft = parseInt(window.getComputedStyle(galleryArea).marginLeft);

window.addEventListener('load', function () {
    const galleryAreaWidth = galleryItems.length * 203;
    galleryArea.style.width = galleryAreaWidth + 'px';

    fetch('/static/json/TaiwanAddress_Simple.json')
        .then(response => response.json())      
        .then(cityData => {                      
            cityData.forEach(item => {
                const option = document.createElement('option');
                option.value = item.city;
                option.textContent = item.city;

                document.querySelector("#citySelect").appendChild(option);
            });
        });

    galleryLeftBtn.addEventListener("click", handleGalleryLeftClick);
    galleryRightBtn.addEventListener("click", handleGalleryRightClick);
    window.addEventListener('resize', resetMarginOnResize);
});


function handleGalleryLeftClick() {
    if (currentMarginLeft === 0) {
        return;
    }

    currentMarginLeft += 203;
    galleryArea.style.marginLeft = currentMarginLeft + 'px';
};

function handleGalleryRightClick() {
    const breakpoints = [1260, 1023, 820, 510, 0];
    const marginValues = Array.from({ length: Math.min(galleryItems.length, 5) }, (_, i) => -(galleryItems.length - (i + 1)) * 203).sort((a, b) => b - a);
    let screenWidth = window.innerWidth;

    let targetMargin;
    for (let i = 0; i < breakpoints.length; i++) {
        if (screenWidth > breakpoints[i]) {
            targetMargin = marginValues[i];
            break;
        }
    }

    if (currentMarginLeft === targetMargin) {
        return;
    }

    currentMarginLeft -= 203;
    galleryArea.style.marginLeft = currentMarginLeft + 'px';
};

function resetMarginOnResize() {
    currentMarginLeft = 0;
    galleryArea.style.marginLeft = currentMarginLeft + 'px';
};