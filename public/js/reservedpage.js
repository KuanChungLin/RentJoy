$("#nextStepBtn").on('click', function () {
    if ($('select[name="Activity"]').val() == null) {
        swal({
            title: "資料欄位不可為空白",
            icon: "error",
            button: "關閉",
            dangerMode: true
        });
        return;
    }
    if (!/^[1-9]\d*$/.test($('input[name="UserCount"]').val())) {
        swal({
            title: "請確認使用人數",
            icon: "error",
            button: "關閉",
            dangerMode: true
        });
        return;
    }
    if (!$('input[name="agree-rule"]').is(':checked')) {
        swal({
            title: "請同意使用及取消規範",
            icon: "error",
            button: "關閉",
            dangerMode: true
        });
        return;
    }
    $(".step1").toggle();
    $(".step2").toggle();
})



$(function () {
    var debounceTimeout;
    $('#submitBtn').on('click', function (event) {
        //event.preventDefault();
        clearTimeout(debounceTimeout); // 清除之前的防抖計時器
        debounceTimeout = setTimeout(function () {
            if ($('input[name="LastName"]').val() == "") {
                swal({
                    title: "請輸入姓",
                    icon: "error",
                    button: "關閉",
                    dangerMode: true
                });
                console.log("1");
                event.preventDefault();

                return;
            }
            if ($('input[name="FirstName"]').val() == "") {
                swal({
                    title: "請輸入名",
                    icon: "error",
                    button: "關閉",
                    dangerMode: true
                });
                console.log("2");
                event.preventDefault();

                return;
            }
            if (!/^09\d{8}$/.test($('input[name="Phone"]').val())) {
                swal({
                    title: "請輸入正確的台灣手機號碼",
                    icon: "error",
                    button: "關閉",
                    dangerMode: true
                });
                console.log("3");
                event.preventDefault();

                return;
            }
            var emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
            if (!emailRegex.test($('input[name="Email"]').val())) {
                swal({
                    title: "請輸入正確的email",
                    icon: "error",
                    button: "關閉",
                    dangerMode: true
                });
                console.log("4");
                event.preventDefault();

                return;
            }
            document.getElementById('form').submit();
            swal({
                title: "正在處理您的請求...",
                customClass: {
                    title: 'custom-title'
                },
                button: false,
                closeOnClickOutside: false,
                closeOnEsc: false,
                content: {
                    element: "div",
                    attributes: {
                        innerHTML: `<p class="swal">請稍等，我們正在創建您的訂單。</p>
                                    <div class="spinner"><div class="cube1"></div><div class="cube2"></div></div>`
                    },
                }
            });
        }, 300); 
    });
});

$("#backToProduct").on('click', function () {
    window.history.back();
})

$("#backToStep1").on('click', function () {
    $(".step1").toggle();
    $(".step2").toggle();
})

