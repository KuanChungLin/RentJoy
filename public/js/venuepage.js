const lastMonthBtn = document.querySelectorAll(".last-month-btn");
const reservedLastMonthBtn = document.querySelector("#lBtn3");
const footerReservedLastMonthBtn = document.querySelector("#lBtn4");
const nextMonthBtn = document.querySelectorAll(".next-month-btn");
const reservedNextMonthBtn = document.querySelector("#rBtn3");
const footerReservedNextMonthBtn = document.querySelector("#rBtn4");
const dateInput = document.querySelector("#selectDate");
const venueData = $('#venue-data');
const reservedDate = venueData.data('reservedDate');
const openDayOfWeek = venueData.data('openDayOfWeek');
const minRentHours = venueData.data('minRentHours');
const venue_Id = venueData.data('venueId');


const today = new Date();
let currentYear;
let currentMonth; //從1開始 1~12
let nextYear;
let nextMonth;
let startTimeData = [];

$(function () {
    init();
});
//改變視窗大小檢查
$(window).resize(function () {
    mobileShowREservedModal();
});


//lastMonthBtn事件區
lastMonthBtn.forEach(mon => {
    mon.addEventListener("click", () => {
        currentMonth--;
        if (currentMonth < 1) {
            currentYear--;
            currentMonth = 12;
        }
        nextYear = currentMonth === 12 ? currentYear + 1 : currentYear;
        nextMonth = currentMonth === 12 ? 1 : currentMonth + 1;
        showTitle(currentYear, currentMonth, ".date-title");
        renderingCalendar(currentYear, currentMonth, ".date-area1");
        showTitle(nextYear, nextMonth, ".date-title2");
        renderingCalendar(nextYear, nextMonth, ".date-area2");
        showTitle(currentYear, currentMonth, ".date-title3");
        renderingCalendar(currentYear, currentMonth, ".date-area3");
        showTitle(currentYear, currentMonth, ".date-title4");
        renderingCalendar(currentYear, currentMonth, ".date-area4");

        getVenueAvailableTime(".calender-group .curr-date-group .date-hover", ".reserved-calender",".footer-reserved-calender");
        getVenueAvailableTime(".reserved-calender .curr-date-group .date-hover", ".reserved-calender");
        getVenueAvailableTime(".next-date-group .date-hover", ".reserved-calender");
        getVenueAvailableTime(".footer-reserved-calender .date-hover", ".footer-reserved-calender");

        setDateStr(".curr-date-group .date-hover", "#selectDate");
        setNextMonDateStr(".next-date-group .date-hover", "#selectDate");
        setNextMonDateStr(".next-date-group .date-hover", "#footerSelectDate");
        setDateStr(".curr-date-group .date-hover", "#footerSelectDate");
        setDateStr(".footer-reserved-calender .date-hover", "#footerSelectDate");
        setDateStr(".footer-reserved-calender .date-hover", "#selectDate");

        clearAllSelectedAndDetails(".calender-group .curr-date-group .date-hover");
        clearAllSelectedAndDetails(".next-date-group .date-hover");
        clearAllSelectedAndDetails(".reserved-calender .curr-date-group .date-hover");
        clearAllSelectedAndDetails(".footer-reserved-calender .date-hover");
        mobileShowREservedModal()
        clearBoostrapModalJs();
    })
});

reservedLastMonthBtn.addEventListener("click", () => {
    currentMonth--;
    if (currentMonth < 1) {
        currentYear--;
        currentMonth = 12;
    }
    nextYear = currentMonth === 12 ? currentYear + 1 : currentYear;
    nextMonth = currentMonth === 12 ? 1 : currentMonth + 1;
    showTitle(currentYear, currentMonth, ".date-title");
    renderingCalendar(currentYear, currentMonth, ".date-area1");
    showTitle(nextYear, nextMonth, ".date-title2");
    renderingCalendar(nextYear, nextMonth, ".date-area2");
    showTitle(currentYear, currentMonth, ".date-title3");
    renderingCalendar(currentYear, currentMonth, ".date-area3");
    showTitle(currentYear, currentMonth, ".date-title4");
    renderingCalendar(currentYear, currentMonth, ".date-area4");

    getVenueAvailableTime(".calender-group .curr-date-group .date-hover", ".reserved-calender",".footer-reserved-calender");
    getVenueAvailableTime(".reserved-calender .curr-date-group .date-hover", ".reserved-calender");
    getVenueAvailableTime(".next-date-group .date-hover", ".reserved-calender");
    getVenueAvailableTime(".footer-reserved-calender .date-hover", ".footer-reserved-calender");

    setDateStr(".curr-date-group .date-hover", "#selectDate");
    setNextMonDateStr(".next-date-group .date-hover", "#selectDate");
    setNextMonDateStr(".next-date-group .date-hover", "#footerSelectDate");
    setDateStr(".curr-date-group .date-hover", "#footerSelectDate");
    setDateStr(".footer-reserved-calender .date-hover", "#footerSelectDate");
    setDateStr(".footer-reserved-calender .date-hover", "#selectDate");

    clearAllSelectedAndDetails(".calender-group .curr-date-group .date-hover");
    clearAllSelectedAndDetails(".next-date-group .date-hover");
    clearAllSelectedAndDetails(".reserved-calender .curr-date-group .date-hover");
    clearAllSelectedAndDetails(".footer-reserved-calender .date-hover");
    clearBoostrapModalJs();
})

footerReservedLastMonthBtn.addEventListener("click", () => {
    currentMonth--;
    if (currentMonth < 1) {
        currentYear--;
        currentMonth = 12;
    }
    nextYear = currentMonth === 12 ? currentYear + 1 : currentYear;
    nextMonth = currentMonth === 12 ? 1 : currentMonth + 1;
    showTitle(currentYear, currentMonth, ".date-title");
    renderingCalendar(currentYear, currentMonth, ".date-area1");
    showTitle(nextYear, nextMonth, ".date-title2");
    renderingCalendar(nextYear, nextMonth, ".date-area2");
    showTitle(currentYear, currentMonth, ".date-title3");
    renderingCalendar(currentYear, currentMonth, ".date-area3");
    showTitle(currentYear, currentMonth, ".date-title4");
    renderingCalendar(currentYear, currentMonth, ".date-area4");

    getVenueAvailableTime(".calender-group .curr-date-group .date-hover", ".reserved-calender",".footer-reserved-calender");
    getVenueAvailableTime(".reserved-calender .curr-date-group .date-hover", ".reserved-calender");
    getVenueAvailableTime(".next-date-group .date-hover", ".reserved-calender");
    getVenueAvailableTime(".footer-reserved-calender .date-hover", ".footer-reserved-calender");

    setDateStr(".curr-date-group .date-hover", "#selectDate");
    setNextMonDateStr(".next-date-group .date-hover", "#selectDate");
    setNextMonDateStr(".next-date-group .date-hover", "#footerSelectDate");
    setDateStr(".curr-date-group .date-hover", "#footerSelectDate");
    setDateStr(".footer-reserved-calender .date-hover", "#footerSelectDate");
    setDateStr(".footer-reserved-calender .date-hover", "#selectDate");

    clearAllSelectedAndDetails(".calender-group .curr-date-group .date-hover");
    clearAllSelectedAndDetails(".next-date-group .date-hover");
    clearAllSelectedAndDetails(".reserved-calender .curr-date-group .date-hover");
    clearAllSelectedAndDetails(".footer-reserved-calender .date-hover");
    mobileShowREservedModal()
    clearBoostrapModalJs();
})

//nextMonthBtn事件區
nextMonthBtn.forEach(mon => {
    mon.addEventListener("click", () => {
        currentMonth++;
        if (currentMonth > 12) {
            currentYear++;
            currentMonth = 1;
        }
        nextYear = currentMonth === 12 ? currentYear + 1 : currentYear;
        nextMonth = currentMonth === 12 ? 1 : currentMonth + 1;
        showTitle(currentYear, currentMonth, ".date-title");
        renderingCalendar(currentYear, currentMonth, ".date-area1");
        showTitle(nextYear, nextMonth, ".date-title2");
        renderingCalendar(nextYear, nextMonth, ".date-area2");
        showTitle(currentYear, currentMonth, ".date-title3");
        renderingCalendar(currentYear, currentMonth, ".date-area3");
        showTitle(currentYear, currentMonth, ".date-title4");
        renderingCalendar(currentYear, currentMonth, ".date-area4");

        getVenueAvailableTime(".calender-group .curr-date-group .date-hover", ".reserved-calender",".footer-reserved-calender");
        getVenueAvailableTime(".reserved-calender .curr-date-group .date-hover", ".reserved-calender");
        getVenueAvailableTime(".next-date-group .date-hover", ".reserved-calender");
        getVenueAvailableTime(".footer-reserved-calender .date-hover", ".footer-reserved-calender");

        setDateStr(".curr-date-group .date-hover", "#selectDate");
        setNextMonDateStr(".next-date-group .date-hover", "#selectDate");
        setNextMonDateStr(".next-date-group .date-hover", "#footerSelectDate");
        setDateStr(".curr-date-group .date-hover", "#footerSelectDate");
        setDateStr(".footer-reserved-calender .date-hover", "#footerSelectDate");
        setDateStr(".footer-reserved-calender .date-hover", "#selectDate");

        clearAllSelectedAndDetails(".calender-group .curr-date-group .date-hover");
        clearAllSelectedAndDetails(".next-date-group .date-hover");
        clearAllSelectedAndDetails(".reserved-calender .curr-date-group .date-hover");
        clearAllSelectedAndDetails(".footer-reserved-calender .date-hover");
        mobileShowREservedModal()
        clearBoostrapModalJs();
    });
})

reservedNextMonthBtn.addEventListener("click", () => {
    currentMonth++;
    if (currentMonth > 12) {
        currentYear++;
        currentMonth = 1;
    }
    nextYear = currentMonth === 12 ? currentYear + 1 : currentYear;
    nextMonth = currentMonth === 12 ? 1 : currentMonth + 1;
    showTitle(currentYear, currentMonth, ".date-title");
    renderingCalendar(currentYear, currentMonth, ".date-area1");
    showTitle(nextYear, nextMonth, ".date-title2");
    renderingCalendar(nextYear, nextMonth, ".date-area2");
    showTitle(currentYear, currentMonth, ".date-title3");
    renderingCalendar(currentYear, currentMonth, ".date-area3");
    showTitle(currentYear, currentMonth, ".date-title4");
    renderingCalendar(currentYear, currentMonth, ".date-area4");

    getVenueAvailableTime(".calender-group .curr-date-group .date-hover", ".reserved-calender",".footer-reserved-calender");
    getVenueAvailableTime(".reserved-calender .curr-date-group .date-hover", ".reserved-calender");
    getVenueAvailableTime(".next-date-group .date-hover", ".reserved-calender");
    getVenueAvailableTime(".footer-reserved-calender .date-hover", ".footer-reserved-calender");

    setDateStr(".curr-date-group .date-hover", "#selectDate");
    setNextMonDateStr(".next-date-group .date-hover", "#selectDate");
    setNextMonDateStr(".next-date-group .date-hover", "#footerSelectDate");
    setDateStr(".curr-date-group .date-hover", "#footerSelectDate");
    setDateStr(".footer-reserved-calender .date-hover", "#footerSelectDate");
    setDateStr(".footer-reserved-calender .date-hover", "#selectDate");

    clearAllSelectedAndDetails(".calender-group .curr-date-group .date-hover");
    clearAllSelectedAndDetails(".next-date-group .date-hover");
    clearAllSelectedAndDetails(".reserved-calender .curr-date-group .date-hover");
    clearAllSelectedAndDetails(".footer-reserved-calender .date-hover");
    clearBoostrapModalJs();
});

footerReservedNextMonthBtn.addEventListener("click", () => {
    currentMonth++;
    if (currentMonth > 12) {
        currentYear++;
        currentMonth = 1;
    }
    nextYear = currentMonth === 12 ? currentYear + 1 : currentYear;
    nextMonth = currentMonth === 12 ? 1 : currentMonth + 1;
    showTitle(currentYear, currentMonth, ".date-title");
    renderingCalendar(currentYear, currentMonth, ".date-area1");
    showTitle(nextYear, nextMonth, ".date-title2");
    renderingCalendar(nextYear, nextMonth, ".date-area2");
    showTitle(currentYear, currentMonth, ".date-title3");
    renderingCalendar(currentYear, currentMonth, ".date-area3");
    showTitle(currentYear, currentMonth, ".date-title4");
    renderingCalendar(currentYear, currentMonth, ".date-area4");

    getVenueAvailableTime(".calender-group .curr-date-group .date-hover", ".reserved-calender",".footer-reserved-calender");
    getVenueAvailableTime(".reserved-calender .curr-date-group .date-hover", ".reserved-calender");
    getVenueAvailableTime(".next-date-group .date-hover", ".reserved-calender");
    getVenueAvailableTime(".footer-reserved-calender .date-hover", ".footer-reserved-calender");

    setDateStr(".curr-date-group .date-hover", "#selectDate");
    setNextMonDateStr(".next-date-group .date-hover", "#selectDate");
    setNextMonDateStr(".next-date-group .date-hover", "#footerSelectDate");
    setDateStr(".curr-date-group .date-hover", "#footerSelectDate");
    setDateStr(".footer-reserved-calender .date-hover", "#footerSelectDate");
    setDateStr(".footer-reserved-calender .date-hover", "#selectDate");

    clearAllSelectedAndDetails(".calender-group .curr-date-group .date-hover");
    clearAllSelectedAndDetails(".next-date-group .date-hover");
    clearAllSelectedAndDetails(".reserved-calender .curr-date-group .date-hover");
    clearAllSelectedAndDetails(".footer-reserved-calender .date-hover");

    mobileShowREservedModal()
    clearBoostrapModalJs();
});

$("#deviceModalOpen").click(function () {
    $("#deviceModal").toggle();
    $("#overlay").toggle();
    $("body").addClass('modal-open');
})

$("#deviceModalClose").click(function () {
    $("#deviceModal").toggle();
    $("#overlay").toggle();
    $("body").removeClass('modal-open');
})

$(".product-img-group").click(function () {
    $("#imgCarousel").toggle();
    $("#overlay").toggle();
    $("body").addClass('modal-open');
})

$("#carouselClose").click(function () {
    $("#imgCarousel").toggle();
    $("#overlay").toggle();
    $("body").removeClass('modal-open');
})

$('#selectDate').click(function () {
    $(".reserved-calender").css("display", "flex");
})

$('#footerSelectDate').click(function () {
    $(".footer-reserved-calender").css("display", "flex");
})
$(document).on('click', function (event) {
    if (!$(event.target).is('#selectDate') && (!$(event.target).closest('.reserved-calender').length || $(event.target).is('.date-hover'))) {
        $('.reserved-calender').hide();
    }
})

$(document).on('click', function (event) {
    if (!$(event.target).is('#footerSelectDate') && (!$(event.target).closest('.footer-reserved-calender').length || $(event.target).is('.date-hover'))) {
        $('.footer-reserved-calender ').hide();
    }
})

$('body').on("click", ".product-price .select-time", function () {
    differentCalculate(this, ".product-price")
    showPerReservedDetails(this, ".product-price");
    updateDisabledLinkPer();
});

$('body').on("click", ".footer-product-price .select-time", function () {
    differentCalculate(this, ".footer-product-price");
    showPerReservedDetails(this, ".footer-product-price");
    updateDisabledLinkPer()
});
$('body').on("click", ".product-price .time .start-time li button", function () {
    setTimeout(function () {
        showHourReservedDetails.call(this, '.product-price .start-time .start', '.product-price .end-time .end');
    }.bind(this), 20);
});
$('body').on("click", ".product-price .time .end-time li button", function () {
    setTimeout(function () {
        showHourReservedDetails.call(this, '.product-price .start-time .start', '.product-price .end-time .end');
    }.bind(this), 20);
});
$('body').on("click", "#footerCalenderModal .time .start-time li button", function () {
    setTimeout(function () {
        showHourReservedDetails.call(this, '#footerCalenderModal .start-time .start', '#footerCalenderModal .end-time .end');
    }.bind(this), 20);
});
$('body').on("click", "#footerCalenderModal .time .end-time li button", function () {
    setTimeout(function () {
        showHourReservedDetails.call(this, '#footerCalenderModal .start-time .start', '#footerCalenderModal .end-time .end');
    }.bind(this), 20);
});


$(".product-price .self-set-time-group").click(function (event) {
    $(".product-price .select-time").removeClass("btn-click");
    if (!$(event.target).closest('button').length) {
        $(this).toggleClass("btn-click");
        $(".self-set-time").toggle();
        $(".self-set-time-group .arrow").toggleClass("d-flex");
        $(".self-set-time-span").toggle();
        $(".start-time .dropdown-toggle").text("開始時間");
        $(".end-time .dropdown-toggle").text("結束時間");
    }
    if ($(".product-price .self-set-start-time").text() != "開始時間" && $(".product-price .self-set-end-time").text() != "結束時間") {
        showPerReservedDetails(this, ".product-price");
    }
})

$(".footer-product-price .self-set-time-group").click(function (event) {
    $(".footer-product-price .select-time").removeClass("btn-click");
    if (!$(event.target).closest('button').length) {
        $(this).toggleClass("btn-click");
        $(".self-set-time").toggle();
        $(".self-set-time-group .arrow").toggleClass("d-flex");
        $(".self-set-time-span").toggle();
        $(".start-time .dropdown-toggle").text("開始時間");
        $(".end-time .dropdown-toggle").text("結束時間");
    }
    if ($(".footer-product-price .self-set-start-time").text() != "開始時間" && $(".footer-product-price .self-set-end-time").text() != "結束時間") {
        showPerReservedDetails(this, ".footer-product-price");
    }
})
$('body').on("click", ".start-time ul button", function (event) {
    $(".start-time .dropdown-toggle").text($(this).text());
    $(".start-time .dropdown-toggle").attr('data-hour-price', $(this).data('hourPrice'))
});

$('body').on("click", ".end-time ul button", function (event) {
    $(".end-time .dropdown-toggle").text($(this).text());
    $(".end-time .dropdown-toggle").attr('data-hour-price', $(this).data('hourPrice'))
});
$('body').on("click", ".calender-group .curr-date-group .date-hover", function (event) {
    updatePerTimeDataSelect();
    updateDisabledLinkPer()
});
$('body').on("click", ".reserved-calender .curr-date-group .date-hover", function (event) {
    updatePerTimeDataSelect();
    updateDisabledLinkPer()
});
$('body').on("click", ".next-date-group .date-hover", function (event) {
    updatePerTimeDataSelect();
    updateDisabledLinkPer()
});
$('body').on("click", ".footer-reserved-calender .date-hover", function (event) {
    updatePerTimeDataSelect();
    updateDisabledLinkPer()
});
$('body').on("click", ".product-price .end-menu", function (event) {
    setTimeout(function () {
        updateDisabledLinkHour(".product-price");
    }.bind(this), 0);
});
$('body').on("click", ".product-price .start-menu", function (event) {
    setTimeout(function () {
        updateDisabledLinkHour(".product-price");
    }.bind(this), 0);
});
$('body').on("click", "#footerCalenderModal .start-menu", function (event) {
    setTimeout(function () {
        updateDisabledLinkHour("#footerCalenderModal");
    }.bind(this), 0);
});
$('body').on("click", "#footerCalenderModal .end-menu", function (event) {
    setTimeout(function () {
        updateDisabledLinkHour("#footerCalenderModal");
    }.bind(this), 0);
});

$('.product-price .reserved .btn').click(function () {
    saveTimePriceToCookie('.product-price');
});
$('#footerCalenderModal .reserved .btn').click(function () {
    saveTimePriceToCookie('.footer-product-price');
});

$('body').on("click", ".product-price .start-hour-time", function () {
    setTimeout(function () {
        var startTime = $(this).parent().parent().siblings(".start").text().trim();
        updateEndTimeOptions(startTime,'.product-price');
    }.bind(this), 0);
});
$('body').on("click", "#footerCalenderModal .start-hour-time", function () {
    setTimeout(function () {
        var startTime = $(this).parent().parent().siblings(".start").text().trim();
        updateEndTimeOptions(startTime, '#footerCalenderModal');
    }.bind(this), 0);
});







//function
function init() {
    currentYear = today.getFullYear();
    currentMonth = today.getMonth() + 1;
    nextYear = currentMonth === 12 ? currentYear + 1 : currentYear;
    nextMonth = currentMonth === 12 ? 1 : currentMonth + 1;
    showTitle(currentYear, currentMonth, ".date-title");
    renderingCalendar(currentYear, currentMonth, ".date-area1");
    showTitle(nextYear, nextMonth, ".date-title2");
    renderingCalendar(nextYear, nextMonth, ".date-area2");
    showTitle(currentYear, currentMonth, ".date-title3");
    renderingCalendar(currentYear, currentMonth, ".date-area3");
    showTitle(currentYear, currentMonth, ".date-title4");
    renderingCalendar(currentYear, currentMonth, ".date-area4");
    setDateStr(".curr-date-group .date-hover", "#selectDate");
    setNextMonDateStr(".next-date-group .date-hover", "#selectDate");
    setNextMonDateStr(".next-date-group .date-hover", "#footerSelectDate");
    setDateStr(".footer-reserved-calender .date-hover", "#footerSelectDate");
    setDateStr(".curr-date-group .date-hover", "#footerSelectDate");
    showReservedArea(".body-calender .date-hover", ".product-price");
    showReservedArea(".footer-reserved-calender .date-hover", ".footer-product-price");
    mobileShowREservedModal();
    getVenueAvailableTime(".calender-group .curr-date-group .date-hover", ".reserved-calender",".footer-reserved-calender");
    getVenueAvailableTime(".reserved-calender .curr-date-group .date-hover", ".reserved-calender");
    getVenueAvailableTime(".next-date-group .date-hover", ".reserved-calender");
    getVenueAvailableTime(".footer-reserved-calender .date-hover", ".footer-reserved-calender");
    clearAllSelectedAndDetails(".calender-group .curr-date-group .date-hover");
    clearAllSelectedAndDetails(".next-date-group .date-hover");
    clearAllSelectedAndDetails(".reserved-calender .curr-date-group .date-hover");
    clearAllSelectedAndDetails(".footer-reserved-calender .date-hover");
    clearBoostrapModalJs();
}

function updateEndTimeOptions(startTime,father)  {
    // 清空所有結束時間選項
    $('.end-menu li').remove();
    let thirtyTimes = 1;
    // 填充結束時間選項
    const endTimeMenu = $('.end-menu');
    startTimeData.forEach((time, index) => {
        // 檢查開始時間加上每30分鐘之後的選項是否存在
        //console.log(index, startTimeData.length - 1)
        var nextThirtyMinutes = addThirtyMinutes(startTime);
        var [endtimehours, endtimeMinutes] = time.split(':').map(Number);
        var endTotalMinutes = endtimehours * 60 + endtimeMinutes + 30;//結束時間比開始時間多30分鐘
        var endTotalTimeStr = `${Math.floor(endTotalMinutes / 60).toString().padStart(2, '0')}:${(endTotalMinutes % 60).toString().padStart(2, '0')}`;
        //console.log(nextThirtyMinutes, endTotalMinutes)
        if (startTimeData.includes(nextThirtyMinutes[0]) && nextThirtyMinutes[1] + 30 * thirtyTimes == endTotalMinutes) {
            option = `<li>
                        <button class="dropdown-item end-hour-time" type="button">${endTotalTimeStr}</button>
                      </li>`;
            endTimeMenu.append(option);
            thirtyTimes++;
        }
    });
    handleHourTimeSelection(".start-menu", ".end-menu");
}
function addThirtyMinutes(time) {
    var timeSplit = time.split(':');
    var hours = parseInt(timeSplit[0]);
    var minutes = parseInt(timeSplit[1]);
    var totalMinutes = hours * 60 + minutes;
    minutes += 30;
    if (minutes >= 60) {
        hours += 1;
        minutes -= 60;
    }
    if (hours < 10) hours = '0' + hours;
    if (minutes < 10) minutes = '0' + minutes;
    return [hours + ':' + minutes, totalMinutes];
}
function updateHourTimeOptions() {
    $(".start-time .btn").change(function () {
        // 獲取當前選擇的開始時間
        let selectedStart = $(this).text().trim();

        // 將開始時間分解為小時和分鐘
        let [startHour, startMinute] = selectedStart.split(':').map(Number);

        // 將開始時間轉換成分鐘
        let startMinutes = startHour * 60 + startMinute;
    })
}
function updateDisabledLinkHour(father) {
    if ($(`${father} .start-menu`).siblings('.btn').text().trim() !== "開始時間" && $(`${father} .end-menu`).siblings('.btn').text().trim() !== "結束時間") {
        $(`${father} a.disable`).removeClass('disable');
    } else {
        $(`${father} .reserved a`).addClass('disable');
    }
}
function updateDisabledLinkPer() {
    if ($('[data-selected="true"]').length > 0) {
        if ($('a.disable').length > 0) {
            // 如果存在，删除disabled
            $('a.disable').removeClass('disable');
        }
    } else {
        $('.reserved a').addClass('disable');
    }
}
function updatePerTimeDataSelect(){
    if ($('[data-selected="true"]').length > 0) {
        $('[data-selected="true"]').removeClass('data-selected');
    }
}
function handleHourTimeSelection(startMenuClass, endMenuClass) {
        // 先清除之前的選擇狀態
        $(endMenuClass).siblings('.btn').text('結束時間');
        // 獲取當前選擇的開始時間
        let selectedStart = $(startMenuClass).siblings('.btn').text().trim().slice(0, 5);
        // 將開始時間分解為小時和分鐘
        let [startHour, startMinute] = selectedStart.split(':').map(Number);
        // 將開始時間轉換成分鐘
        let startMinutes = startHour * 60 + startMinute;
        // 最少單位小時 將小時轉換成分鐘
        let minMinutes = minRentHours * 60;
        // 遍歷結束時間的下拉選單選項
        $(endMenuClass).find('li button').each(function () {
            // 獲取選項的時間並分解為小時和分鐘
            let [endHour, endMinute] = $(this).text().trim().split(':').map(Number);
            // 將結束時間轉換成分鐘
            let endMinutes = endHour * 60 + endMinute;
            // 如果結束時間早於開始時間，則禁用該選項，否則啟用
            if (endMinutes - startMinutes < minMinutes) {
                $(this).remove();
            }
        });
}
function saveTimePriceToCookie(father) {
    Cookies.remove('TimeDetailCookie');
    class TimeDetailCookie {
        constructor(timeSlotIds, startTime, endTime, reservedDay, venueId) {
            this.TimeSlotIds = timeSlotIds;
            this.StartTime = startTime;
            this.EndTime = endTime;
            this.ReservedDay = reservedDay;
            this.VenueId = venueId;
        }
    }

    let timeSlotIds = [];
    let startTime;
    let endTime;
    let reservedDay;
    if ($(father).find('.time-period-group').length > 0) {
        //時段的資料傳遞
        let UtcDay = new Date($(`#selectDate`).text());
        UtcDay.setHours(UtcDay.getHours() + 8);
        reservedDay = UtcDay;
        $(`${father} .select-time`).each(function () {
            // 如果此元素被選中
            if ($(this).attr('data-selected') === 'true') {
                timeSlotIds.push($(this).attr('data-bill-id')); // 添加到id陣列中
            }
        });
    } else {
        //小時的資料傳遞
        let UtcDay = new Date($(`#selectDate`).text());
        UtcDay.setHours(UtcDay.getHours() + 8);
        reservedDay = UtcDay;
        let startStr = $(father).find('.start-time .dropdown-toggle').text().trim();
        let endStr = $(father).find('.end-time .dropdown-toggle').text().trim();
        let [stHours, stMinutes] = startStr.split(':').map(Number);
        let [etHours, etMinutes] = endStr.split(':').map(Number);
        let etSeconds = 0
        if (etHours == 24 && etMinutes == 0) {
            etHours = 23;
            etMinutes = 59;
            etSeconds = 59
        }
        //console.log([etHours, etMinutes])

        let startTimeData = new Date(reservedDay.getTime());
        startTimeData.setHours(+stHours + 8, +stMinutes, 0, 0);
        let endTimeData = new Date(reservedDay.getTime());
        endTimeData.setHours(+etHours + 8, +etMinutes, +etSeconds, 0);
        //console.log(startTimeData, endTimeData, reservedDay)
        startTime = startTimeData;
        endTime = endTimeData;
    }
    if (timeSlotIds.length === 0 && startTime || timeSlotIds.length > 0 && !startTime) {
        let timeDetailsDataCookie = new TimeDetailCookie(timeSlotIds, startTime, endTime, reservedDay, venue_Id);
        let DataCookie = JSON.stringify(timeDetailsDataCookie);

        Cookies.set("TimeDetailCookie", DataCookie, { expires: 60 * 60, path: '/Venues' });
    }
}
function clearAllSelectedAndDetails(targetDate) {
    $(targetDate).click(function (e) {
        // 移除所有 'data-selected' 属性
        $(`.select-time`).removeAttr('data-selected');
        // 移除所有 '.product-price-details'
        $(`.product-price-details`).empty();
        $(`.product-price-details`).hide();
    })


}
function getVenueAvailableTime(targetDate, showPrice,showPriceFooter) {
    $(targetDate).click(function (e) {
        //清除原本選項
        if (showPriceFooter) {
            $(showPriceFooter).nextAll('.time').remove();
            $(showPriceFooter).nextAll('.time-period-group').remove();
        }
        $(showPrice).nextAll('.time').remove();
        $(showPrice).nextAll('.time-period-group').remove();
        //敲api
        let dateTitle = $(this).closest(".calendar-area").siblings(".calender-header").find("h2").text().trim();
        let dateYearMonth = dateTitle.replace(/ /g, '').replace('/', '-');
        //console.log(new Date(`${dateYearMonth}-${$(this).text()}`).toISOString());
        //console.log(new Date(`${dateYearMonth}-${$(this).text().padStart(2, 0)}`).toISOString());
        $.ajax({
            url: '/Venues/GetAllAvailableTime',
            data: {
                selectDay: new Date((`${dateYearMonth}-${$(this).text().padStart(2, 0)}`)).toJSON(),
                venueId: venue_Id
            },
            dataType: 'json',
            type: 'GET',
            success: function (data) {
                //console.log(data);
                //成功就渲染畫面
                if (data.length === 0) {
                    let newPerGroup = `
                            <div class="border border-top-0 p-3 time-period-group">
                                <div class="time-period no-options d-flex">
                                    此日期已無預訂時間，請更換日期
                                </div>
                            </div>
                        `;
                    $(showPrice).after(newPerGroup);
                    if (showPriceFooter) {
                        $(showPriceFooter).after(newPerGroup);
                    }
                }
                else if (data[0].rateTypeId == "1") {
                    let newStartItem = '';
                    let newEndItem = '';
                    let i = 0;
                    startTimeData = [];
                    data.forEach(item => {
                        let minRentMintues = minRentHours * 60;//startMinutes + minRentMintues -30
                        let relateTimes = (minRentMintues - 30); //最小小時+中間選走時段與開始時間的關係(畫陣列)
                        let compareindex = i + relateTimes / 30;
                        if (compareindex >= data.length) {
                            compareindex = data.length - 1;
                        }
                        let compareStartTime = new Date(data[compareindex].startTime);
                        let compareStartTimeMintues = compareStartTime.getHours() * 60 + compareStartTime.getMinutes();
                        let start = new Date(item.startTime);
                        let end = new Date(item.endTime);
                        let startMinutes = start.getHours() * 60 + start.getMinutes();
                        let startTimeStrToData = (start.getHours() < 10 ? "0" : "") + start.getHours() + ":" + (start.getMinutes() < 10 ? "0" : "") + start.getMinutes();
                        startTimeData.push(startTimeStrToData);
                        if (startMinutes + relateTimes == compareStartTimeMintues ) {
                            // 將 item 的開始時間和結束時間格式化為 "HH:MM" 的格式
                            let startTimeStr = (start.getHours() < 10 ? "0" : "") + start.getHours() + ":" + (start.getMinutes() < 10 ? "0" : "") + start.getMinutes();
                            let totalMinutes = startMinutes + minRentMintues;
                            let endTimeStr = `${Math.floor(totalMinutes / 60).toString().padStart(2, '0')}:${(totalMinutes % 60).toString().padStart(2, '0')}`;
                            newStartItem += `<li>
                                <button class="dropdown-item start-hour-time" type="button" data-type-id=${item.Type} data-min-hour=${item.MinHour} data-hour-price=${item.price}>${startTimeStr}</button>
                            </li>`;
                        }
                        i++;
                    })

                    let newHourGroup = `<div class="time border border-top-0 p-2">
                    <div class="dropdown d-flex start-time">
                        <button class="btn btn-sm dropdown-toggle start border-0" type="button" data-bs-toggle="dropdown"
                                aria-expanded="false">
                            開始時間
                        </button>
                        <ul class="dropdown-menu start-menu">
                           ${newStartItem}
                        </ul>
                    </div>
                    <div class="arrow">→</div>
                    <div class="dropdown end-time">
                        <button class="btn btn-sm dropdown-toggle end border-0" type="button" data-bs-toggle="dropdown"
                                aria-expanded="false">
                            結束時間
                        </button>
                        <ul class="dropdown-menu end-menu">
                            
                        </ul>
                    </div>
                </div>`;
                    $(showPrice).after(newHourGroup);
                    if (showPriceFooter) {
                        $(showPriceFooter).after(newHourGroup);
                    }

                    if ($(".product-price .time .start-time .start-menu").find('li').length === 0) {
                        let noOptions = `
                            <div class="border border-top-0 p-3 time-period-group">
                                <div class="time-period no-options d-flex">
                                    此日期已無預訂時間，請更換日期
                                </div>
                            </div>
                        `;
                        $(showPrice).siblings('.time').remove();
                        $(showPriceFooter).siblings('.time').remove();
                        $(showPrice).after(noOptions);
                        $(showPriceFooter).after(noOptions);
                    }
                }
                else if (data[0].rateTypeId == "2") {
                    let newItem = '';
                    data.forEach(item => {
                        let start = new Date(item.startTime);
                        let end = new Date(item.endTime);

                        // 將 item 的開始時間和結束時間格式化為 "HH:MM" 的格式
                        let startTimeStr = ("0" + start.getHours()).slice(-2) + ":" + ("0" + start.getMinutes()).slice(-2);
                        let endTimeStr = ("0" + end.getHours()).slice(-2) + ":" + ("0" + end.getMinutes()).slice(-2);

                        newItem += `
                            <div class="btn btn-outline-secondary rounded-pill my-2 d-flex  select-time" data-bill-id=${item.billingRateId}>
                                ${startTimeStr} - ${endTimeStr}
                                <span class="border-start px-2 ms-auto select-price">$${item.price}</span>
                            </div>
                        `;
                    });

                    let newPerGroup = `
                            <div class="border border-top-0 p-3 time-period-group">
                                <div class="time-period d-flex">
                                    請選擇時段
                                </div>
                                <div class="d-flex flex-column time-period-select">
                                    ${newItem}
                                </div>
                            </div>
                        `;

                    $(showPrice).after(newPerGroup);
                    if (showPriceFooter) {
                        $(showPriceFooter).after(newPerGroup);
                    }
                }
            },
            error: function (jqXHR, textStatus, errorThrown) {
                console.error('Error: ' + textStatus, errorThrown);
            }
        });

    });
}
function renderingCalendar(year, month, dateAreaSelector) {
    const firstDateOfCurrentMonth = new Date(year, month - 1, 1);
    const lastDateOfCurrentMonth = new Date(year, month, 0)
    let start = 1 - firstDateOfCurrentMonth.getDay();
    let end = lastDateOfCurrentMonth.getDate() + (6 - lastDateOfCurrentMonth.getDay());

    const dateArea = document.querySelector(dateAreaSelector);
    dateArea.innerHTML = "";
    for (start; start <= end; start++) {
        const curr = new Date(year, month - 1, start);
        const dateDom = document.createElement("div");
        dateDom.classList.add("col");

        const dateEl = document.createElement("div");
        dateEl.classList.add("h-100", "w-100", "text-center", "date", "rounded", "date-hover");
        dateEl.textContent = curr.getDate();

        if (isDateReserved(curr, reservedDate)) {
            dateEl.classList.add("reserved-color");
        }
        if (start <= 0 || start > lastDateOfCurrentMonth.getDate()) {
            dateEl.textContent = ""
            dateEl.classList.remove("date", "date-hover", "reserved-color")
        };
        
        //日期顯示需要劃掉:  明天(包括明天)之前不可預訂, 沒提供該星期可租    
        if (curr.getTime() <= today.getTime() + (24 * 60 * 60 * 1000) || !openDayOfWeek.includes(curr.getDay()) ) {
            dateEl.classList.add("text-decoration-line-through", "opacity-50");
            dateEl.classList.remove("date-hover", "reserved-color")
        };
        

        dateDom.append(dateEl);
        dateArea.append(dateDom.cloneNode(true));
        $("#footerCalenderModal .date-hover").removeAttr('data-bs-toggle');
    }
}

function isDateReserved(currDate, reservedDateArr) {
    // 將當前日期轉換為 ISO 8601 格式的字符串，(會自動轉成UTC)
    var currDateStr = currDate.toISOString().slice(0, 10);
    
    for (var i = 0; i < reservedDateArr.length; i++) {
        // 將 reservedDate 中的日期轉換為 ISO 8601 格式的字符串
        var reservedDateStr = new Date(reservedDateArr[i]).toISOString().slice(0, 10);

        if (currDateStr === reservedDateStr) {
            return true;
        }
    }
    return false;
}


function setDateStr(targetDate, showDate) {
    $(targetDate).click(function (e) {
        // console.log("點到了");
        $(showDate).text(`${currentYear} / ${currentMonth.toString().padStart(2, 0)} / ${e.target.textContent.padStart(2, 0)}`);
    })
}
function setNextMonDateStr(targetDate, showDate) {
    $(targetDate).click(function (e) {
        // console.log("點到了");
        $(showDate).text(`${nextYear} / ${nextMonth.toString().padStart(2, 0)} / ${e.target.textContent.padStart(2, 0)}`);
    })
}
function getDateStr(date) {
    // return '2024-01-09'
    return `${date.getFullYear()}-${(date.getMonth() + 1)
        .toString()
        .padStart(2, "0")}-${date.getDate().toString().padStart(2, "0")}`;
}
function showTitle(year, month, titleSelector) {
    const title = document.querySelector(titleSelector);
    title.textContent = `${year} / ${month.toString().padStart(2, 0)}`;
}
function showReservedArea(target, father) {
    $(target).click(function () {
        $(".time-period-group").show();
        $(".time").css("display", "flex");
        $(`${father} .start`).text("開始時間");
        $(`${father} .end`).text("結束時間");
    });
}
function showPerReservedDetails(target, father) {
    if ($(target).hasClass('btn-click')) {
        // 當按鈕被點擊時，添加 'data-selected' 屬性
        $(target).attr('data-selected', 'true');
    } else {
        // 如果已經點擊的按鈕再次被點擊，刪除 'data-selected' 屬性
        $(target).removeAttr('data-selected');
    }
    // 建立一個變數來儲存總價
    let totalPrice = 0;
    // 清空明細
    $(`${father} .product-price-details`).empty();
    // 檢查每一個按鈕，如果按鈕有 'data-selected' 屬性，則在明細中增加一行
    $(`${father} .btn`).each(function () {
        if ($(this).attr('data-selected') === 'true') {
            var billId = $(this).data('bill-id');
            var newTime = $(this).clone().children().remove().end().text().trim();
            var newPrice = $(this).find('.select-price').text().replace(/\$/g, "").trim();

            // 將價格加到總價上
            totalPrice += Number(newPrice);
            var newDetails = `<div class="d-flex justify-content-between pb-1" data-bill-id=${billId}>
                                時段 ${newTime}
                                <span>$${newPrice}</span></div>`;

            $(`${father} .product-price-details`).append(newDetails);
        }
    });
    // 在明細的末尾添加總價
    var totalPriceRow = `<div class="d-flex justify-content-between pt-2 mt-3 border-top">
                            TWD
                            <span class="large-txt">$${totalPrice}</span>
                         </div>`;
    $(`${father} .product-price-details`).append(totalPriceRow);
    // 顯示明細
    if ($(`${father} .btn-click`).length === 0) {
        // 如果按鈕沒有 'btn-click' 類，則隱藏明細
        $(`${father} .product-price-details`).hide();
    } else {
        $(`${father} .product-price-details`).show();
    };
}
function showHourReservedDetails(startTime, endTime) {
    $(startTime).parent().parent().siblings('.product-price-details').hide();
    $(startTime).parent().parent().siblings('.product-price-details').empty();
    if ($(startTime).first().text().trim() !== "開始時間" && $(endTime).first().text().trim() !== "結束時間") {
        $(startTime).parent().parent().siblings('.product-price-details').show();
        // 獲取開始和結束時間的小時數和分鐘數及小時的價錢
        const [startHours, startMins] = $(startTime).text().trim().split(':').map(Number);
        const [endHours, endMins] = $(endTime).text().trim().split(':').map(Number);
        const hourPrice = $(startTime).data('hourPrice');
        // 計算此期間的總小時數，包含分鐘
        let hours = 0;
        let mins = 0;
        hours = endHours - startHours;
        mins = endMins - startMins;

        // 結果是負數，表示跨越了整點
        if (mins < 0) {
            hours -= 1;
            mins += 60;
        }
        // 剩餘的分鐘數轉換成小時
        const totalHours = hours + mins / 60;

        const totalFee = hourPrice * totalHours;

        const newHourDetail = `
        <div class="d-flex justify-content-between pb-1">
            $${hourPrice} x ${totalHours} 小時
            <span>$${totalFee}</span>
        </div>
        <div class="d-flex justify-content-between pt-3 mt-3 border-top">
            TWD
            <span class="large-txt">$${totalFee}</span>
        </div>`;

        $('.product-price-details').append(newHourDetail);
    }
}
function differentCalculate(target, father) {
    $(target).toggleClass("btn-click");
    $(`${father} .self-set-time-group`).removeClass("btn-click");
    $(".self-set-time").hide();
    $(".self-set-time-span").show();
    $(`${father} .self-set-time-group .arrow`).removeClass("d-flex");
    $(`${father} .start-time .dropdown-toggle`).text("開始時間");
    $(`${father} .end-time .dropdown-toggle`).text("結束時間");
}

function mobileShowREservedModal() {
    if ($(window).width() < 1025) {
        $('.date-hover').attr('data-bs-toggle', 'modal');
        $('.date-hover').attr('data-bs-target', '#footerCalenderModal');
    } else {
        $('.date-hover').removeAttr('data-bs-toggle');
        $('.date-hover').removeAttr('data-bs-target');
    }
}
function clearBoostrapModalJs() {
    $("#footerCalenderModal .date-hover").removeAttr('data-bs-toggle');
};