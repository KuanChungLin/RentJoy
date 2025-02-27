$(document).ready(function () {
    const searchParams = $('#search-params');
    const selectedCity = searchParams.data('city');
    const selectedDistrict = searchParams.data('district');
    console.log(searchParams)
    console.log(selectedCity)
    console.log(selectedDistrict)

    $.ajax({
        url: "/static/json/TaiwanAddress_Simple.json",
        dataType: "JSON",
    }).done(cityData => {
        // 生成縣市選項
        cityData.forEach(item => {
            const option = $(`<option value="${item.city}">${item.city}</option>`);
            $("#citySelect").append(option);
        });

        // 設置預設城市和地區
        if (selectedCity) {
            $("#citySelect").val(selectedCity);
            updateDistrictSelect(selectedCity, cityData);

            if (selectedDistrict) {
                $('#districtSelect').val(selectedDistrict);
            }
        }

        // 城市選擇後更新鄉鎮區
        $("#citySelect").on('change', function (e) {
            const citySelectVal = e.target.value;
            
            if (citySelectVal === '') {
                $("#districtSelect").html('<option value="">請先選擇縣市</option>');
                return;
            }

            updateDistrictSelect(citySelectVal, cityData);
        });

        // 攔截表單提交，確保城市和地區被正確提交
        // $('#searchNavbar').on('submit', function() {
        //     // 確保城市和地區被正確設置
        //     const cityVal = $('#citySelect').val();
        //     const districtVal = $('#districtSelect').val();

        //     // 如果城市已選但地區未選，阻止提交並提示
        //     if (cityVal && !districtVal) {
        //         alert('請選擇地區');
        //         return false;
        //     }
        // });
    });
});

function updateDistrictSelect(citySelectVal, cityData) {
    const $districtSelect = $('#districtSelect');
    $districtSelect.empty();
    $districtSelect.append('<option value="">請選擇地區</option>');

    const areaData = cityData.find(item => item.city === citySelectVal).districts;

    areaData.forEach(item => {
        const option = `<option value="${item.district}">${item.district}</option>`;
        $districtSelect.append(option);
    });
}