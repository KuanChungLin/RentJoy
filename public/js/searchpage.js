$(document).ready(function () {
    var isLoading = false;
    var page = 1;
    var hasMoreData = true;
    const city = $('#search-params').data('city');
    const district = $('#search-params').data('district');
    const filterData = $('#searchNavbar').serializeArray().map(item => {
        if (item.name === 'City') {
            return { ...item, value: city };
        } else if (item.name === 'District') {
            return { ...item, value: district };
        } else {
            return item;
        }
    });
    console.log(filterData)


    $(window).on('scroll', function () {
        if (hasMoreData && $(window).scrollTop() + $(window).height() >= $(document).height() - 100 && !isLoading) {
            isLoading = true;
            page++;

            var pageItem = filterData.find(function (item) { return item.name === "Page"; });

            if (pageItem) {
                pageItem.value = page;
            }

            filterData.push({ name: "Page", value: page });

            var formData = $.param(filterData);
            console.log(formData)

            if (!$('#loading').length) {
                $('#searchResult').append('<div id="loading" class="spinner"><div class="rect1"></div><div class="rect2"></div><div class="rect3"></div><div class="rect4"></div><div class="rect5"></div></div>');
            }

            $.ajax({
                url: '/SearchPageLoading',
                type: 'GET',
                data: formData,
                success: function (response) {
                    $('#loading').remove();
                    if (response.VenueInfos && response.VenueInfos.length > 0) {
                        response.VenueInfos.forEach(venue => {
                            const venueHtml = `
                                <div class="search-product rounded">
                                    <a href="/Venue/VenuePage?venueId=${venue.venueId}" class="product-imga">
                                        <img src="${venue.venueImgUrl}" class="product-img">
                                    </a>
                                    <div class="product-info-group">
                                        <a href="/Venue/VenuePage?venueId=${venue.venueId}" class="product-name">
                                            <h3>${venue.venueName}</h3>
                                        </a>
                                        <div class="product-spec-group mb-3">
                                            <div class="product-spec">
                                                <div class="mb5">
                                                    <span class="product-icon">
                                                        <img src="/static/images/house.png">
                                                        ${venue.venueOwner}
                                                    </span>
                                                </div>
                                                <div class="mb5">
                                                    <span class="product-icon">
                                                        <img src="/static/images/marker.png">
                                                        ${venue.venueCity} ${venue.venueDistrict}
                                                    </span>
                                                    <span class="product-icon">
                                                        <img src="/static/images/person.png">
                                                        ${venue.numberOfPeople} 人
                                                    </span>
                                                </div>
                                            </div>
                                            <div class="product-price">
                                                <a href="/Venue/VenuePage?venueId=${venue.venueId}">
                                                    <button type="button" class="btn btn-price">${venue.venuePrice}</button>
                                                </a>
                                            </div>
                                        </div>
                                        <div class="product-tags">
                                            <span class="product-icon">
                                                <img src="/static/images/tags.png">
                                                ${venue.activityTags ? venue.activityTags.map(tag => 
                                                    `<a href="/SearchPage?activityId=${tag.Id}">${tag.activityName}</a><b>/</b>`
                                                ).join('') : ''}
                                            </span>
                                        </div>
                                    </div>
                                </div>`;
                            $('#searchResult').append(venueHtml);
                        });
                    }
                    
                    if (response.EndOfData) {  // 使用 EndOfData 判斷
                        hasMoreData = false;
                        $(window).off('scroll');
                        return;
                    }

                    $('#searchResult').append(response);

                },
                error: function (xhr, status, error) {
                    $('#loading').remove();
                    console.error("發生錯誤: " + status);
                    hasMoreData = false;
                    $('#searchResult').append('<div class="text-center text-danger">載入資料時發生錯誤</div>');
                },
                complete: function () {
                    isLoading = false;
                }
            });
        }
    });
});








$('.rent-time').on('change', function () { 
    if ($(this).is(':checked')) {
        $('.rent-time').not(this).prop('checked', false);
    }
});

$('.rent-day').on('change', function () {
    if ($(this).is(':checked')) {
        $('.rent-day').not(this).prop('checked', false);
    }
});
