const titleBasicInformationBtn = document.querySelector('.header-basic-information')
const titleLocationBtn = document.querySelector('.header-location')
const titleSpatialConfigurationBtn = document.querySelector('.header-spatial-configuration')
const titleEquipmentBtn = document.querySelector('.header-equipment')
const titlePhotoBtn = document.querySelector('.header-photo')
const titlePriceBtn = document.querySelector('.header-price')
const titleManagerBtn = document.querySelector('.header-manager')

const basicInformationContent = document.querySelector('.basic-information');
const locationContent = document.querySelector('.location');
const spatialConfigurationContent = document.querySelector('.spatialConfiguration');
const equipmentContent = document.querySelector('.equipment');
const photoContent = document.querySelector('.photo');
const priceContent = document.querySelector('.price');
const managerContent = document.querySelector('.manager');
const allContent = document.querySelectorAll('.my-space-set > div')

//抓取space-type-radio 
const typeRadiosContent = document.querySelector('.space-type-radio');
const typeFoodSpaceBtn = document.querySelector('.food-space')
const typeStudySpaceBtn = document.querySelector('.study-space')
const typeBusinessSpaceBtn = document.querySelector('.business-space')
const typeMultipurposeSpaceBtn = document.querySelector('.multipurpose-space')
const typeSportsSpaceBtn = document.querySelector('.sports-space')
const typeMusicSpaceBtn = document.querySelector('.music-space')
const typeEveryoneSpaceBtn = document.querySelector('.everyone-space')
const typeOtherSpaceBtn = document.querySelector('.other-space')
const detailFoodRadio = document.querySelector('.detail-food')
const detailStudyRadio = document.querySelector('.detail-study')
const detailBusinessRadio = document.querySelector('.detail-business')
const detailMultipurposeRadio = document.querySelector('.detail-multipurpose')
const detailSportsRadio = document.querySelector('.detail-sports')
const detailMusicRadio = document.querySelector('.detail-music')
const detailEveryoneRadio = document.querySelector('.detail-everyone')
const detailOtherRadio = document.querySelector('.detail-other')
const detailAllRadio = document.querySelectorAll('.space-type-radio>div');

const basicInformationNextBtn = document.querySelector('.basic-information .btn-next button')
const locationBackBtn = document.querySelector('.spatialConfiguration .btn-back-next .back')
const locationNextBtn = document.querySelector('.spatialConfiguration .next')

basicInformationNextBtn.addEventListener("click", () => {
    contentViewChange(locationContent);
})
locationBackBtn.addEventListener("click", () => {
    contentViewChange(basicInformationContent);
})
locationNextBtn.addEventListener("click", () => {
    contentViewChange(spatialConfigurationContent);
})


const PagesTotalBtn = document.querySelectorAll('.nav-item')
const activityCheckboxes = document.querySelectorAll('.activity-checkbox input[type="checkbox"]');
const nextButton = document.querySelector('.btn-next');


function contentViewChange(contentToShow) {
    allContent.forEach(content => {
        if (content === contentToShow) {
            content.style.display = '';
        } else {
            content.style.display = 'none';
        }
    });
}

function typeSpaceRadioChange(detailToShow) {
    detailAllRadio.forEach(detail => {
        if (detail === detailToShow) {
            detail.style.display = '';
        } else {
            detail.style.display = 'none';
        }
    })
}

function checkActivityCheckboxes() {
    let anyCheckboxChecked = false;
    activityCheckboxes.forEach(function (checkbox) {
        if (checkbox.checked) {
            anyCheckboxChecked = true;
        }
    });

    // 根据复选框的状态显示或隐藏按钮
    if (anyCheckboxChecked) {
        nextButton.style.display = 'block';
    } else {
        nextButton.style.display = 'none';
    }
}

window.onload = () => {
    checkActivityCheckboxes();

    PagesTotalBtn.forEach(btn => {
        btn.addEventListener('click', () => {
            // 重置所有按鈕的樣式
            PagesTotalBtn.forEach(btn => {
                btn.querySelector('a').style.color = "#9A9A9A";
                btn.style.border = 'none';
            });
            // 將點擊的按鈕添加樣式
            btn.querySelector('a').style.color = '#53a385';
            btn.style.borderLeft = 'solid 1px #d8d8d8';
            btn.style.borderRight = 'solid 1px #d8d8d8';
        });
    });

    activityCheckboxes.forEach(function (checkbox) {
        checkbox.addEventListener('change', checkActivityCheckboxes);
    });


}