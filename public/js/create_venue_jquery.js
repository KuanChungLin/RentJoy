$(document).ready(function() {
    // 全局變量
    let state = {
        activeTab: 'basic-information',
        spaceTypes: [],
        isLoading: true,
        hasSpaceTypes: false,
        activeTypeId: null,
        currentTypeId: null,
        radioGroup: null,
        selectedActivities: [],
        activities: [],
        venueName: '',
        venueRule: '',
        unsubscribeRule: '',
        citySelectInfo: null,
        districtSelectInfo: null,
        addressInfo: '',
        venueArea: '',
        venueCapacity: '',
        equipmentTypes: [],
        collapsedIds: [],
        equipmentGroup: [],
        photos: [],
        photosTempPath: [],
        selectedCover: '',
        managers: [],
        selectedManagerId: null,
        selectedManagerName: '',
        selectedManagerImg: '',
        selectedContact: '',
        selectedPublicPhone: '',
        selectedPrivatePhone: '',
        defaultImage: '../../images/areaManager_profile.svg',
        isDropdownVisible: false,
        isManagerDetailVisible: false,
        managerFormTitle: '',
        managerFormDescription: '',
        managerFormContact: '',
        managerFormPublicPhone: '',
        managerFormPrivatePhone: '',
        managerPhotoFile: null,
        managerPhotoTempPath: null,
        showModal: false,
        selectedSpaceType: null,
        venuePrice: '',
        venueDescription: '',
        selectedEquipment: [],
        venueManager: '',
        imageUrls: [],
        isCreatingHourPlan: true,
        isEditingHourPlan: false,
        isCreatingPeriodPlan: true,
        isEditingPeriodPlan: false,
        isSpecificHourPlan: false,
        isSpecificPeriodPlan: false,
        isEditedHourPlan: false,
        isEditedPeriodPlan: false,
        days: ['monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday'],
        dayMapping: {
            monday: '週一',
            tuesday: '週二',
            wednesday: '週三',
            thursday: '週四',
            friday: '週五',
            saturday: '週六',
            sunday: '週日'
        },
        hourPricing: {},
        periodPricing: {},
        leastRentHours: 0.5,
        timeOptions: []
    };

    

    // 初始化函數
    async function initialize() {
        try {
            await fetchSpaceTypes();
            await fetchActivities();
            await fetchEquipmentTypes();
            await fetchManagersInfo();
            initializePricing();
            state.isLoading = false;
            updateLoadingState();
        } catch (error) {
            console.error('Error in initialization:', error);
            state.isLoading = false;
            updateLoadingState();
        }
    }

    // 更新加載狀態
    function updateLoadingState() {
        if (state.isLoading) {
            $('#loading-state').show();
            $('#no-space-types').hide();
            $('#space-types-container').hide();
        } else if (!state.hasSpaceTypes) {
            $('#loading-state').hide();
            $('#no-space-types').show();
            $('#space-types-container').hide();
        } else {
            $('#loading-state').hide();
            $('#no-space-types').hide();
            $('#space-types-container').show();
        }
    }

//#region 基本資料處理

    // 獲取空間類型
    async function fetchSpaceTypes() {
        try {
            const response = await axios.get('/Create/GetSpaceTypesData');
            
            if (Array.isArray(response.data)) {
                state.spaceTypes = response.data.map(type => ({
                    ...type,
                    typeName: type.typeName || '未知類型',
                    spaceInfos: (type.spaceInfos || []).map(info => ({
                        ...info,
                        facilityName: info.facilityName || '未知設施'
                    }))
                }));
                state.hasSpaceTypes = state.spaceTypes.length > 0;
                renderSpaceTypes();
                
            } else {
                console.error('Invalid space types data format:', response.data);
                state.spaceTypes = [];
                state.hasSpaceTypes = false;
            }
        } catch (error) {
            console.error('Error fetching space types:', error);
            state.spaceTypes = [];
            state.hasSpaceTypes = false;
        }
    }

    // 渲染空間類型
    function renderSpaceTypes() {
        const container = $('.space-type-button');
        container.empty();

        state.spaceTypes.forEach(type => {
            const button = $('<button>')
                .addClass('btn btn-outline-success space')
                .text(type.typeName || '未知類型')
                .click(() => selectSpaceType(type));

            if (state.selectedSpaceType && state.selectedSpaceType.id === type.id) {
                button.addClass('active');
            }

            container.append(button);
        });
    }

    // 選擇空間類型
    function selectSpaceType(type) {
        state.selectedSpaceType = type;
        state.activeTypeId = type.id;
        renderSpaceTypes();
        renderSpaceTypeDetails();
    }

    // 渲染空間類型詳情
    function renderSpaceTypeDetails() {
        const container = $('.space-type-radio');
        const detail = $('<div>').addClass('detail');
        container.empty();
        state.spaceTypes.forEach(type => {
            if (state.activeTypeId === type.id) {
                type.spaceInfos.forEach(facility => {
                    
                    const radioGroup = $('<div>')
                        .addClass('radio-item');

                    const input = $('<input>')
                        .attr('type', 'radio')
                        .attr('id', `facility-${facility.id}`)
                        .attr('name', 'facilityId')
                        .attr('value', facility.id)
                        .change(() => {
                            state.radioGroup = facility.id;
                        });

                    const label = $('<label>')
                        .attr('for', `facility-${facility.id}`)
                        .text(facility.facilityName || '未知設施');

                    radioGroup.append(input, label);
                    detail.append(radioGroup);
                    container.append(detail);
                });
            }
        });
    }

    // 獲取活動列表
    async function fetchActivities() {
        try {
            const response = await axios.get('/Create/GetActivitiesData');
            
            if (Array.isArray(response.data)) {
                state.activities = response.data.map(activity => ({
                    ...activity,
                    name: activity.name || '未知活動'
                }));
                renderActivities();
                
            } else {
                console.error('Invalid activities data format:', response.data);
                state.activities = [];
            }
        } catch (error) {
            console.error('Error fetching activities:', error);
            state.activities = [];
        }
    }

    // 渲染活動列表
    function renderActivities() {
        const container = $('#activity-checkbox');
        container.empty();

        state.activities.forEach(activity => {
            const checkboxGroup = $('<div>')
                .addClass('checkbox-item');

            const input = $('<input>')
                .attr('type', 'checkbox')
                .attr('id', `activity-${activity.id}`)
                .attr('name', 'selectedActivities')
                .attr('value', activity.id)
                .change(function() {
                    if (this.checked) {
                        state.selectedActivities.push(activity.id);
                    } else {
                        state.selectedActivities = state.selectedActivities.filter(id => id !== activity.id);
                    }
                });

            const label = $('<label>')
                .attr('for', `activity-${activity.id}`)
                .text(activity.name || '未知活動');

            checkboxGroup.append(input, label);
            container.append(checkboxGroup);
        });
    }
//#endregion

//#region 設備處理
    // 獲取設備類型
    async function fetchEquipmentTypes() {
        try {
            const response = await axios.get('/Create/GetEquipmentTypesData');
            
            if (Array.isArray(response.data)) {
                state.equipmentTypes = response.data.map(type => ({
                    ...type,
                    type: type.typeName || '未知類型',
                    equipmentFacility: (type.equipmentInfos || []).map(facility => ({
                        ...facility,
                        id: facility.id,
                        equipmentName: facility.equipmentName || '未知設備',
                        quantity: facility.quantity || 1,
                        description: facility.description || ''
                    }))
                }));
                renderEquipmentTypes();
                
            } else {
                console.error('Invalid equipment types data format:', response.data);
                state.equipmentTypes = [];
            }
        } catch (error) {
            console.error('Error fetching equipment types:', error);
            state.equipmentTypes = [];
        }
    }

    // 渲染設備類型
    function renderEquipmentTypes() {
        const container = $('.equipment .container');
        // 保存按钮区域
        
        const buttonArea = container.find('.btn-back-next').detach();
        
        // 保留标题部分
        const header = $('<div>').append(
            $('<h3>').text('設備'),
            $('<span>').text('選填')
        );
        container.empty().append(header);

        state.equipmentTypes.forEach(equipType => {
            const equipmentSection = $('<div>')
                .addClass('equipment-classification-collapse');

            const itemArea = $('<div>')
                .addClass('classification-item-area')
                .click(() => toggleCollapse(equipType.id));

            const itemTitle = $('<div>')
                .addClass('classification-item-title')
                .append(
                    $('<div>')
                        .addClass('classification-name')
                        .text(equipType.typeName || '未知類型')
                );

            // 添加展開/收起圖標
            const isCollapsed = state.collapsedIds.includes(equipType.id);
            const expandIcon = $('<div>').addClass('expand-icon').html(`
                <svg xmlns="http://www.w3.org/2000/svg" width="30" height="30" viewBox="0 0 30 30">
                    <g fill="none" fill-rule="evenodd" stroke="#53A385">
                        <path d="M1.5 14.5c0 7.732 6.269 14 14 14 7.732 0 14-6.268 14-14s-6.268-14-14-14c-7.731 0-14 6.268-14 14z" />
                        <path d="M11.4 12.68l4.1 4.1 4.1-4.1" />
                    </g>
                </svg>
            `);

            const collapseIcon = $('<div>').addClass('collapse-icon').html(`
                <svg xmlns="http://www.w3.org/2000/svg" width="30" height="30" viewBox="0 0 30 30">
                    <g fill="none" fill-rule="evenodd">
                        <path fill="#53A385" stroke="#53A385" d="M28.5 15.5c0-7.732-6.269-14-14-14-7.732 0-14 6.268-14 14s6.268 14 14 14c7.731 0 14-6.268 14-14z" />
                        <path stroke="#FFF" d="M18.6 17.32l-4.1-4.1-4.1 4.1" />
                    </g>
                </svg>
            `);

            itemTitle.append(isCollapsed ? collapseIcon : expandIcon);
            itemArea.append(itemTitle);
            equipmentSection.append(itemArea);

            // 添加分隔線
            const hr = $('<hr>').toggle(!isCollapsed);
            equipmentSection.append(hr);

            // 展開的內容

                const expandContent = $('<div>')
                    .addClass('equipment-classification-expand')
                    .toggle(isCollapsed);
                const facilityArea = $('<div>')
                    .addClass('classification-item-area');

                equipType.equipmentFacility.forEach(facility => {
                    const facilityDetail = createFacilityDetail(facility);
                    facilityArea.append(facilityDetail);
                });

                expandContent.append(facilityArea);
                equipmentSection.append(expandContent);


            container.append(equipmentSection);
        });

        // 重新添加按钮区域
        container.append($('<div>').addClass('blank'), buttonArea);
    }

    // 創建設施詳情
    function createFacilityDetail(facility) {
        const detail = $('<div>').addClass('classification-item-detail');
        const option = $('<div>').addClass('equipment-option');

        // 設備標題
        const title = $('<div>').addClass('equipment-title');
        const checkbox = $('<input>')
            .attr('type', 'checkbox')
            .attr('id', `equipment-${facility.id}`)
            .attr('value', facility.id)
            .attr('name', `EquipmentInfos[${facility.id}].Selected`)
            .prop('checked', state.equipmentGroup.includes(facility.id))
            .change(function() {
                const facilityId = facility.id;
                const quantityNoteElement = $(this).closest('.equipment-option').find('.quantity-note');
                
                if (this.checked) {
                    if (!state.equipmentGroup.includes(facilityId)) {
                        state.equipmentGroup.push(facilityId);
                    }
                    quantityNoteElement.addClass('show');
                } else {
                    state.equipmentGroup = state.equipmentGroup.filter(id => id !== facilityId);
                    quantityNoteElement.removeClass('show');
                }
            });

        const label = $('<label>')
            .attr('for', `equipment-${facility.id}`)
            .addClass('equipment-name')
            .text(facility.equipmentName);

        title.append(checkbox, label);

        // 數量和備註
        const quantityNote = $('<div>').addClass('quantity-note');
        
        // 數量部分
        const quantity = $('<div>').addClass('equipment-quantity');
        quantity.append($('<p>').text('數量'));
        
        const calQuantity = $('<div>').addClass('cal-quantity');
        const minusBtn = $('<div>')
            .addClass('minus-button')
            .click((e) => {
                e.stopPropagation();  // 阻止事件冒泡
                decrementFacility(facility);
            });
        
        const quantityInput = $('<input>')
            .attr('type', 'number')
            .attr('name', `EquipmentInfos[${facility.id}].Quantity`)
            .attr('min', '1')
            .val(facility.quantity || 0)
            .click((e) => e.stopPropagation())  // 阻止事件冒泡
            .change(function() {
                facility.quantity = parseInt($(this).val()) || 1;
            });
        
        const plusBtn = $('<div>')
            .addClass('plus-button')
            .click((e) => {
                e.stopPropagation();  // 阻止事件冒泡
                incrementFacility(facility);
            });

        calQuantity.append(minusBtn, quantityInput, plusBtn);
        quantity.append(calQuantity);

        // 備註部分
        const note = $('<div>').addClass('equipment-note');
        note.append($('<p>').text('規格或敘述'));
        const noteInput = $('<input>')
            .attr('type', 'text')
            .attr('name', `EquipmentInfos[${facility.id}].Description`)
            .val(facility.description || '')
            .change(function() {
                facility.description = $(this).val();
            });
        note.append(noteInput);

        quantityNote.append(quantity, note);
        option.append(title, quantityNote);
        detail.append(option);

        return detail;
    }

    // 切換折疊狀態
    function toggleCollapse(typeId) {
        const index = state.collapsedIds.indexOf(typeId);
        if (index === -1) {
            state.collapsedIds.push(typeId);
        } else {
            state.collapsedIds.splice(index, 1);
        }
        // 只切換相關區域的展開/收起狀態，而不是重新渲染整個列表
    const section = $(`.equipment-classification-collapse`).filter(function() {
        return $(this).find('.classification-name').text() === 
            state.equipmentTypes.find(type => type.id === typeId).typeName;
    });

    const isCollapsed = state.collapsedIds.includes(typeId);
    section.find('hr').toggle(!isCollapsed);
    section.find('.equipment-classification-expand').toggle(isCollapsed);
    
    // 切換展開/收起圖標
    const itemTitle = section.find('.classification-item-title');
    itemTitle.find('.expand-icon, .collapse-icon').remove();
    
    if (isCollapsed) {
        itemTitle.append(`
            <div class="collapse-icon">
                <svg xmlns="http://www.w3.org/2000/svg" width="30" height="30" viewBox="0 0 30 30">
                    <g fill="none" fill-rule="evenodd">
                        <path fill="#53A385" stroke="#53A385" d="M28.5 15.5c0-7.732-6.269-14-14-14-7.732 0-14 6.268-14 14s6.268 14 14 14c7.731 0 14-6.268 14-14z" />
                        <path stroke="#FFF" d="M18.6 17.32l-4.1-4.1-4.1 4.1" />
                    </g>
                </svg>
            </div>
        `);
    } else {
        itemTitle.append(`
            <div class="expand-icon">
                <svg xmlns="http://www.w3.org/2000/svg" width="30" height="30" viewBox="0 0 30 30">
                    <g fill="none" fill-rule="evenodd" stroke="#53A385">
                        <path d="M1.5 14.5c0 7.732 6.269 14 14 14 7.732 0 14-6.268 14-14s-6.268-14-14-14c-7.731 0-14 6.268-14 14z" />
                        <path d="M11.4 12.68l4.1 4.1 4.1-4.1" />
                    </g>
                </svg>
            </div>
        `);
    }
    }

    // 增加設備數量
    function incrementFacility(facility) {
        if (!facility.quantity) facility.quantity = 0;
        facility.quantity++;
        $(`input[name="EquipmentInfos[${facility.id}].Quantity"]`).val(facility.quantity);
    }

    // 減少設備數量
    function decrementFacility(facility) {
        if (!facility.quantity) facility.quantity = 1;
        if (facility.quantity > 1) {
            facility.quantity--;
            $(`input[name="EquipmentInfos[${facility.id}].Quantity"]`).val(facility.quantity);
        }
    }
//#endregion

//#region 管理員處理
    // 獲取管理員信息
    async function fetchManagersInfo() {
        try {
            const response = await axios.get('/Create/GetManagersInfoData');

            if (Array.isArray(response.data)) {
                state.managers = response.data.map(manager => ({
                    ...manager,
                    managerName: manager.managerName || '未知管理員',
                    managerContact: manager.managerContact || '',
                    managerPublicPhone: manager.managerPublicPhone || '',
                    managerPrivatePhone: manager.managerPrivatePhone || '',
                    managerImgUrl: manager.managerImgUrl || "/static/images/areaManager_profile.svg"
                }));
                renderManagers();
                
            } else {
                console.error('Invalid managers data format:', response.data);
                state.managers = [];
            }
        } catch (error) {
            console.error('Error fetching managers:', error);
            state.managers = [];
        }
    }

    // 渲染管理員列表
    function renderManagers() {
        const managerList = $('#manager-list');
        managerList.empty();

        state.managers.forEach(manager => {
            const managerOption = $('<div>')
                .addClass('manager-option')
                .html(`
                    <label>
                        <input type="radio" name="manager" value="${manager.id}">
                        <div class="select-row">
                            <div class="border-container">
                                <img src="${manager.managerImgUrl}" alt="" class="img">
                                <p>${manager.managerName || '未知管理員'}</p>
                            </div>
                        </div>
                    </label>
                `);
            managerList.append(managerOption);
        });

        // 點擊新增管理者按鈕
        $('#create-manager-btn').click(function(e) {
            e.stopPropagation();
            $('#manager-dropdown').hide();
            $('#create-manager-form').show();
            $('#manager-detail').hide();

            // 更新顯示
            $('#selected-manager-name').text("新增管理者");
            $('#selected-manager-img').attr('src', "/static/images/areaManager_profile.svg");
            
            // 清空管理者ID
            $('#area-merchant-id').val("");
        });

        // 點擊取消按鈕
        $('#cancel-create-manager').click(function() {
            $('#create-manager-form').hide();
            if (state.selectedManagerId) {
                $('#manager-detail').show();
            }
        });

        // 點擊 select-box 時顯示/隱藏下拉選單
        $('#manager-select-box').click(function(e) {
            e.stopPropagation();
            $('#manager-dropdown').toggle();
            // 如果已經選擇過管理者，確保管理者詳細資訊是顯示的
            if (state.selectedManagerId) {
                $('#manager-detail').show();
                $('#create-manager-form').hide();
            }
        });

        // 點擊管理者選項時處理
        $('#manager-list .manager-option').click(function(e) {
            e.stopPropagation();
            const selectedId = $(this).find('input[type="radio"]').val();
            const selectedManager = state.managers.find(m => m.id === parseInt(selectedId));
            if (selectedManager) {
                state.selectedManagerId = selectedManager.id;
                state.selectedManagerName = selectedManager.managerName;
                state.selectedManagerImg = selectedManager.managerImgUrl;
                
                // 更新顯示
                $('#selected-manager-name').text(selectedManager.managerName);
                $('#selected-manager-img').attr('src', selectedManager.managerImgUrl);

                // 更新管理者詳細資訊
                $('.manager-detail-contact').text(selectedManager.managerContact);
                $('.manager-detail-mobile').text(selectedManager.managerPrivatePhone);
                $('.manager-detail-phone').text(selectedManager.managerPublicPhone);
                
                // 設定管理者ID
                $('#area-merchant-id').val(selectedManager.id);
                
                // 顯示管理者詳細資訊，隱藏新增表單
                $('#manager-detail').show();
                $('#create-manager-form').hide();
                
                // 隱藏下拉選單
                $('#manager-dropdown').hide();
            }
        });
    }
//#endregion
    // 事件處理
    $('#area-info-manager').change(function() {
        state.venueManager = $(this).val();
    });

    $('#area-info-name').on('input', function() {
        state.venueName = $(this).val();
    });

    $('#area-info-rule').on('input', function() {
        state.venueRule = $(this).val();
    });

    $('#area-info-unsubscribe-rule').on('input', function() {
        state.unsubscribeRule = $(this).val();
    });

    $('#area-info-price').on('input', function() {
        state.venuePrice = $(this).val();
    });

    $('#area-info-description').on('input', function() {
        state.venueDescription = $(this).val();
    });

//#region 圖片處理

    // 點擊新增照片按鈕
    $('#click-file-input').click(function() {
        // 觸發 input 的 click 事件
        $('#add-new-photo-file').click();
    });

    // 處理文件選擇
    $('#add-new-photo-file').change(function(e) {
        const files = e.target.files;
        if (!files || files.length === 0) return;

        const file = files[0];
        
        // 檢查文件大小
        if (file.size > 5 * 1024 * 1024) {
            swal({
                title: '圖片大小必須小於5MB',
                icon: 'error',
                button: '關閉',
                dangerMode: true
            });
            return;
        }
        
        // 檢查文件類型
        if (!['image/png', 'image/jpeg', 'image/jpg', 'image/bmp'].includes(file.type)) {
            swal({
                title: '不支持的圖片類型',
                icon: 'error',
                button: '關閉',
                dangerMode: true
            });
            return;
        }

        const reader = new FileReader();
        reader.onload = function(e) {
            const newPhotoId = Date.now();
            compressImage(e.target.result, file.type, function(compressedDataUrl) {
                const newIndex = state.photosTempPath.length;
                state.photosTempPath.push({ id: newPhotoId, url: compressedDataUrl });
                if (!state.selectedCover) {
                    state.selectedCover = newPhotoId;
                    // 如果是第一張照片，直接放在第一位
                    state.photos.unshift(file);
                } else {
                    state.photos.push(file);
                }
                
                // 只添加新的照片元素
                const area = $('.photos-area');
                const photoItem = $('<div>')
                    .addClass('photo-item img-upload')
                    .attr('data-index', newIndex);

                const coverOption = $('<div>')
                    .addClass('cover-option photo-pic')
                    .css('backgroundImage', `url(${compressedDataUrl})`);

                const coverRadio = $('<input>')
                    .attr('type', 'radio')
                    .attr('name', 'cover-photo') // 添加 name 屬性確保只能選中一個
                    .attr('id', `cover-photo-${newPhotoId}`)
                    .attr('value', newPhotoId)
                    .addClass('cover-radio black-section')
                    .prop('checked', state.selectedCover === newPhotoId)
                    .change(function() {
                        const selectedId = $(this).val();
                        state.selectedCover = selectedId;
                        
                        // 找到選中照片的索引
                        const selectedIndex = state.photosTempPath.findIndex(photo => photo.id === selectedId);
                        if (selectedIndex > 0) {
                            // 將選中的照片移到第一位
                            const selectedPhoto = state.photosTempPath.splice(selectedIndex, 1)[0];
                            state.photosTempPath.unshift(selectedPhoto);
                            
                            // 同樣移動 photos 陣列中的照片
                            const selectedFile = state.photos.splice(selectedIndex, 1)[0];
                            state.photos.unshift(selectedFile);
                            
                            // 重新渲染照片區域
                            renderPhotos();
                        }
                        
                        // 更新所有照片的封面狀態
                        $('.cover-radio').each(function() {
                            const radio = $(this);
                            const label = radio.next('label');
                            const innerCircle = label.find('.inner-circle');
                            innerCircle.toggle(radio.prop('checked'));
                        });
                    });

                const coverLabel = $('<label>')
                    .attr('for', `cover-photo-${newPhotoId}`)
                    .append(
                        $('<div>').addClass('cover-content').append(
                            $('<div>').addClass('outer-circle').append(
                                $('<div>').addClass('inner-circle').toggle(state.selectedCover === newPhotoId)
                            ),
                            $('<p>').text('設為封面')
                        )
                    );

                const deleteButton = $('<button>')
                    .addClass('cross')
                    .attr('type', 'button')
                    .click((e) => {
                        e.preventDefault();
                        e.stopPropagation();
                        removePhoto(newIndex);
                    })
                    .append(
                        $('<img>').attr('src', '/static/images/cross-orange.svg')
                    );

                coverOption.append(coverRadio, coverLabel);
                photoItem.append(coverOption, deleteButton);
                area.append(photoItem);
            });
        };
        reader.readAsDataURL(file);
    });

    // 圖片壓縮函數
    function compressImage(dataUrl, fileType, callback) {
        const img = document.createElement('img');
        img.onload = function() {
            const maxWidth = 800;
            const maxHeight = 800;
            let width = img.width;
            let height = img.height;

            // 計算縮放比例
            let ratio = 1;
            if (width > maxWidth) {
                ratio = maxWidth / width;
            }
            if (height > maxHeight) {
                ratio = Math.min(ratio, maxHeight / height);
            }

            // 如果圖片需要縮放
            if (ratio < 1) {
                width = Math.floor(width * ratio);
                height = Math.floor(height * ratio);
            }

            const canvas = document.createElement('canvas');
            const ctx = canvas.getContext('2d');
            canvas.width = width;
            canvas.height = height;

            // 使用雙線性插值算法進行縮放
            ctx.imageSmoothingEnabled = true;
            ctx.imageSmoothingQuality = 'high';
            ctx.drawImage(img, 0, 0, width, height);

            // 根據文件類型選擇壓縮質量
            let quality = 0.8;
            if (fileType === 'image/jpeg' || fileType === 'image/jpg') {
                quality = 0.85;
            } else if (fileType === 'image/png') {
                quality = 0.9;
            }

            const compressedDataUrl = canvas.toDataURL(fileType, quality);
            callback(compressedDataUrl);
        };
        img.src = dataUrl;
    }

    // 渲染照片預覽（只在初始化時使用）
    function renderPhotos() {
        const area = $('.photos-area');
        area.empty(); // 清空現有內容

        state.photosTempPath.forEach((photo, index) => {
            const photoItem = $('<div>')
                .addClass('photo-item img-upload')
                .attr('data-index', index);

            const coverOption = $('<div>')
                .addClass('cover-option photo-pic')
                .css('backgroundImage', `url(${photo.url})`);

            const coverRadio = $('<input>')
                .attr('type', 'radio')
                .attr('name', 'cover-photo') // 添加 name 屬性確保只能選中一個
                .attr('id', `cover-photo-${photo.id}`)
                .attr('value', photo.id)
                .addClass('cover-radio black-section')
                .prop('checked', state.selectedCover === photo.id)
                .change(function() {
                    const selectedId = $(this).val();
                    state.selectedCover = selectedId;
                    
                    // 找到選中照片的索引
                    const selectedIndex = state.photosTempPath.findIndex(photo => photo.id === selectedId);
                    if (selectedIndex > 0) {
                        // 將選中的照片移到第一位
                        const selectedPhoto = state.photosTempPath.splice(selectedIndex, 1)[0];
                        state.photosTempPath.unshift(selectedPhoto);
                        
                        // 同樣移動 photos 陣列中的照片
                        const selectedFile = state.photos.splice(selectedIndex, 1)[0];
                        state.photos.unshift(selectedFile);
                        
                        // 重新渲染照片區域
                        renderPhotos();
                    }
                    
                    // 更新所有照片的封面狀態
                    $('.cover-radio').each(function() {
                        const radio = $(this);
                        const label = radio.next('label');
                        const innerCircle = label.find('.inner-circle');
                        innerCircle.toggle(radio.prop('checked'));
                    });
                });

            const coverLabel = $('<label>')
                .attr('for', `cover-photo-${photo.id}`)
                .append(
                    $('<div>').addClass('cover-content').append(
                        $('<div>').addClass('outer-circle').append(
                            $('<div>').addClass('inner-circle').toggle(state.selectedCover === photo.id)
                        ),
                        $('<p>').text('設為封面')
                    )
                );

            const deleteButton = $('<button>')
                .addClass('cross')
                .attr('type', 'button')
                .click((e) => {
                    e.preventDefault();
                    e.stopPropagation();
                    removePhoto(index);
                })
                .append(
                    $('<img>').attr('src', '/static/images/cross-orange.svg')
                );

            coverOption.append(coverRadio, coverLabel);
            photoItem.append(coverOption, deleteButton);
            area.append(photoItem);
        });
    }

    // 刪除照片
    function removePhoto(photoIndex) {
        // 從 photos 陣列中移除指定索引的照片
        state.photos.splice(photoIndex, 1);

        // 檢查是否還有照片
        if (state.photosTempPath.length === 0) {
            state.selectedCover = '';
        }
        
        // 直接從 DOM 中移除對應的照片元素
        const photoElement = $(`.photo-item[data-index="${photoIndex}"]`);
        if (photoElement.length) {
            photoElement.remove();
        } 
    }
//#endregion

    // 切換標籤頁
    function setActiveTab(tabName) {
        // 隱藏所有頁籤
        $('.basic-information, .location, .spatialConfiguration, .equipment, .photo, .price, .manager').hide();
        
        // 顯示選中的頁籤
        $(`#id-${tabName}`).show();
        
        // 移除所有標題的活動狀態
        $('.header-basic-information, .header-location, .header-spatial-configuration, .header-equipment, .header-photo, .header-price, .header-manager').removeClass('active');
        
        // 為當前頁籤的標題添加活動狀態
        if (tabName === 'spatialConfiguration') {
            $(`.header-spatial-configuration`).addClass('active');
        } else {
            $(`.header-${tabName}`).addClass('active');
        }
        // 更新當前活動標籤
        state.activeTab = tabName;
    }

    // 導航欄點擊事件
    $('.nav-menu .nav-item a').on('click', function(e) {
        e.preventDefault();
        const tabName = $(this).data('tab');
        setActiveTab(tabName);
    });

    // 下一步按鈕點擊事件
    $('.next-step').on('click', function() {
        const currentTab = state.activeTab;
        let nextTab = '';
        
        switch(currentTab) {
            case 'basic-information':
                if (validateBasicInformation()) {
                    nextTab = 'location';
                }
                break;
            case 'location':
                if (validateLocation()) {
                    nextTab = 'spatialConfiguration';
                }
                break;
            case 'spatialConfiguration':
                if (validateSpatialConfiguration()) {
                    nextTab = 'equipment';
                }
                break;
            case 'equipment':
                if (validateEquipment()) {
                    nextTab = 'photo';
                }
                break;
            case 'photo':
                if (validatePhoto()) {
                    nextTab = 'price';
                }
                break;
            case 'price':
                if (validatePrice()) {
                    nextTab = 'manager';
                }
                break;
        }
        
        if (nextTab) {
            setActiveTab(nextTab);
        }
    });

    // 上一步按鈕點擊事件
    $('.pre-step').on('click', function() {
        const currentTab = state.activeTab;
        let prevTab = '';
        
        switch(currentTab) {
            case 'location':
                prevTab = 'basic-information';
                break;
            case 'spatialConfiguration':
                prevTab = 'location';
                break;
            case 'equipment':
                prevTab = 'spatialConfiguration';
                break;
            case 'photo':
                prevTab = 'equipment';
                break;
            case 'price':
                prevTab = 'photo';
                break;
            case 'manager':
                prevTab = 'price';
                break;
        }
        
        if (prevTab) {
            setActiveTab(prevTab);
        }
    });

    // 初始化顯示第一個標籤頁
    setActiveTab('basic-information');

    // 表单提交
    $('#venueForm').on('submit', function(e) {
        e.preventDefault();
        
        // 驗證基本資料
        if (!validateBasicInformation()) {
            setActiveTab('basic-information');
            return false;
        }
        
        // 驗證位置資訊
        if (!validateLocation()) {
            setActiveTab('location');
            return false;
        }
        
        // 驗證空間配置
        if (!validateSpatialConfiguration()) {
            setActiveTab('spatialConfiguration');
            return false;
        }
        
        // 驗證相片資訊
        if (!validatePhoto()) {
            setActiveTab('photo');
            return false;
        }
        
        // 驗證價格資訊
        if (!validatePrice()) {
            setActiveTab('price');
            return false;
        }
        
        // 驗證管理者資訊
        if (!validateManager()) {
            setActiveTab('manager');
            return false;
        }

        // 如果所有驗證都通過，顯示確認對話框
        swal({
            title: '確認送出審核申請？',
            text: '審核期間約7~10個工作天，將無法修改資訊\n請確認所有資料是否正確？',
            icon: 'warning',
            buttons: {
                cancel: {
                    text: '返回修改',
                    value: false,
                    visible: true,
                    className: 'btn-secondary'
                },
                confirm: {
                    text: '確認送出',
                    value: true,
                    visible: true,
                    className: 'btn-success'
                }
            }
        }).then((willSubmit) => {
            if (willSubmit) {
                // 創建 FormData 對象
                const formData = new FormData(this);
                
                // 添加設備信息
                state.equipmentGroup.forEach((equipmentId, index) => {
                    const quantity = $(`input[name="EquipmentInfos[${equipmentId}].Quantity"]`).val();
                    const description = $(`input[name="EquipmentInfos[${equipmentId}].Description"]`).val();
                    formData.append(`EquipmentInfos[${index}].Id`, equipmentId);
                    formData.append(`EquipmentInfos[${index}].Quantity`, quantity);
                    formData.append(`EquipmentInfos[${index}].Description`, description);
                });
                
                // 添加照片
                state.photos.forEach((photo, index) => {
                    formData.append('VenueImgs', photo);
                });
                
                // 添加價格信息
                const hourPricingData = flattenPricing(state.hourPricing, mapDayStringToNumber, state.leastRentHours);
                const periodPricingData = flattenPricing(state.periodPricing, mapDayStringToNumber);

                // 添加價格信息到 FormData
                formData.append('HourPricing', JSON.stringify(hourPricingData));
                formData.append('PeriodPricing', JSON.stringify(periodPricingData));
                
                
                // 發送請求
                $.ajax({
                    url: '/Create/CreateVenue',
                    type: 'POST',
                    data: formData,
                    processData: false,
                    contentType: false,
                    success: function(response) {
                        swal({
                            title: '送出成功！',
                            text: '我們將盡快審核您的場地資料',
                            icon: 'success',
                            button: '確定'
                        }).then(() => {
                            window.location.href = '/Manage/VenueManagement';
                        });
                    },
                    error: function(xhr, status, error) {
                        console.error('提交失敗:', error);
                        swal({
                            title: '送出失敗',
                            text: '請稍後再試',
                            icon: 'error',
                            button: '確定'
                        });
                    }
                });
            }
        });
        
        return false;
    });

    // 扁平化價格數據
    function flattenPricing(pricingData, dayMappingFunction, leastRentHours) {
        const flattenedData = [];
        
        Object.entries(pricingData).forEach(([day, dayData]) => {
            if (dayData.isActive) {
                dayData.timePriceSettings.forEach(setting => {
                    flattenedData.push({
                        day: dayMappingFunction(day),
                        startTime: mapTimeValueToDate(setting.startTime),
                        endTime: mapTimeValueToDate(setting.endTime),
                        price: setting.price
                    });
                });
            }
        });
        
        return {
            leastRentHours: leastRentHours,
            pricingSettings: flattenedData
        };
    }

    // 將日期字符串映射為數字
    function mapDayStringToNumber(dayString) {
        const dayMap = {
            'monday': 1,
            'tuesday': 2,
            'wednesday': 3,
            'thursday': 4,
            'friday': 5,
            'saturday': 6,
            'sunday': 7
        };
        return dayMap[dayString] || 0;
    }

    // 將時間值映射為日期對象
    function mapTimeValueToDate(timeValue) {
        const hours = Math.floor(timeValue / 2);
        const minutes = timeValue % 2 === 1 ? 30 : 0;
        const date = new Date();
        date.setUTCHours(hours, minutes, 0, 0);

        return date.toISOString();
    }

    // 驗證基本資料頁面填寫
    function validateBasicInformation() {
        // 檢查空間類型按鈕是否有被選中
        const hasActiveSpaceType = $('.space-type-button .btn').hasClass('active');
        if (!hasActiveSpaceType) {
            swal({
                title: '請選擇場地的空間類型',
                icon: 'warning',
                button: '關閉',
            });
            return false;
        }

        // 檢查空間類型詳情是否有被選中
        const hasSelectedFacility = $('.space-type-radio .detail .radio-item input:checked').length > 0;
        if (!hasSelectedFacility) {
            swal({
                title: '請選擇場地類型',
                icon: 'warning',
                button: '關閉',
            });
            return false;
        }

        // 檢查活動選項是否有被選中
        const hasSelectedActivity = $('.activity-checkbox .checkbox-item input:checked').length > 0;
        if (!hasSelectedActivity) {
            swal({
                title: '請選擇至少一個適合舉辦的活動',
                icon: 'warning',
                button: '關閉',
            });
            return false;
        }

        // 檢查場地名稱是否有被填寫
        const hasVenueName = $('#area-info-name').val() !== '';
        if (!hasVenueName) {
            swal({
                title: '請填寫場地名稱',
                icon: 'warning',
                button: '關閉',
            });
            return false;
        }

        return true;
    }

    // 驗證位置資訊頁面填寫
    function validateLocation() {
        const hasCity = $('#citySelect').val() !== null;
        if (!hasCity) {
            swal({
                title: '請選擇縣市',
                icon: 'warning',
                button: '關閉',
            });
            return false;
        }

        const hasDistrict = $('#districtSelect').val() !== '';
        if (!hasDistrict) {
            swal({
                title: '請選擇地區',
                icon: 'warning',
                button: '關閉',
            });
            return false
        }

        const hasAddressDetail = $('#addressDetail').val() !== '';
        if (!hasAddressDetail) {
            swal({
                title: '請填寫地址詳細資訊',
                icon: 'warning',
                button: '關閉',
            });
            return false;
        }

        return true;
    }

    // 驗證場地配置頁面填寫
    function validateSpatialConfiguration() {
        const hasSpaceSquare = $('#area-info-space-square').val() !== '';
        if (!hasSpaceSquare) {
            swal({
                title: '請填寫場地坪數',
                icon: 'warning',
                button: '關閉',
            });
            return false;
        }

        const hasSpacePeople = $('#area-info-space-people').val() !== '';
        if (!hasSpacePeople) {
            swal({
                title: '請填寫場地容納人數',
            });
            return false;
        }
        return true;
    }

    // 驗證設備資訊頁面填寫
    function validateEquipment() {
        // 檢查是否有選擇任何設備
        if (state.equipmentGroup.length === 0) {
            swal({
                title: '請至少選擇一個設備',
                icon: 'warning',
                button: '關閉',
            });
            return false;
        }
        
        return true;
    }

    // 驗證照片資訊頁面填寫
    function validatePhoto() {
        if (state.photos.length < 3) {
            swal({
                title: '請上傳至少3張不同角度的相片',
                icon: 'warning',
                button: '關閉',
            });
            return false;
        }
        if (!state.selectedCover) {
            swal({
                title: '請設定一張封面照片',
                icon: 'warning',
                button: '關閉',
            });
            return false;
        }
        return true;
    }
    
    // 驗證價格資訊頁面填寫
    function validatePrice() {
        // 檢查小時計費是否有任何一天被啟用
        const hasHourPricing = Object.values(state.hourPricing).some(day => day.isActive);
        // 檢查時段計費是否有任何一天被啟用
        const hasPeriodPricing = Object.values(state.periodPricing).some(day => day.isActive);

        // 檢查是否至少有一種計費方式被啟用
        if (!hasHourPricing && !hasPeriodPricing) {
            swal({
                title: '請設定至少一個價格',
                icon: 'warning',
                button: '關閉',
            });
            return false;
        }

        // 檢查每一天的時間是否有重疊
        for (const day of state.days) {
            // 檢查小時計費的時間重疊和價格
            if (state.hourPricing[day].isActive) {
                const hourSettings = state.hourPricing[day].timePriceSettings;
                // 檢查價格是否為正數
                for (let i = 0; i < hourSettings.length; i++) {
                    if (hourSettings[i].price <= 0) {
                        swal({
                            title: `${state.dayMapping[day]}的小時計費價格必須大於0`,
                            icon: 'warning',
                            button: '關閉',
                        });
                        return false;
                    }
                }
                // 檢查時間重疊
                for (let i = 0; i < hourSettings.length; i++) {
                    for (let j = i + 1; j < hourSettings.length; j++) {
                        const setting1 = hourSettings[i];
                        const setting2 = hourSettings[j];
                        
                        // 檢查時間區間是否重疊
                        if (!(setting1.endTime <= setting2.startTime || setting2.endTime <= setting1.startTime)) {
                            swal({
                                title: `${state.dayMapping[day]}的小時計費時間區間有重疊`,
                                text: '請檢查並修正時間設定',
                                icon: 'warning',
                                button: '關閉',
                            });
                            return false;
                        }
                    }
                }
            }

            // 檢查時段計費的時間重疊和價格
            if (state.periodPricing[day].isActive) {
                const periodSettings = state.periodPricing[day].timePriceSettings;
                // 檢查價格是否為正數
                for (let i = 0; i < periodSettings.length; i++) {
                    if (periodSettings[i].price <= 0) {
                        swal({
                            title: `${state.dayMapping[day]}的時段計費價格必須大於0`,
                            icon: 'warning',
                            button: '關閉',
                        });
                        return false;
                    }
                }
                // 檢查時間重疊
                for (let i = 0; i < periodSettings.length; i++) {
                    for (let j = i + 1; j < periodSettings.length; j++) {
                        const setting1 = periodSettings[i];
                        const setting2 = periodSettings[j];
                        
                        // 檢查時間區間是否重疊
                        if (!(setting1.endTime <= setting2.startTime || setting2.endTime <= setting1.startTime)) {
                            swal({
                                title: `${state.dayMapping[day]}的時段計費時間區間有重疊`,
                                text: '請檢查並修正時間設定',
                                icon: 'warning',
                                button: '關閉',
                            });
                            return false;
                        }
                    }
                }
            }
        }
    
        return true;
    }

    function validateManager() {
        const hasManager = $('#area-merchant-id').val() !== '';
        if (!hasManager) {
            swal({
                title: '請選擇管理員',
                icon: 'warning',    
                button: '關閉',
            });
            return false;
        }
        return true;
    }

// price

    // 初始化價格相關功能
    function initializePricing() {
        // 初始化小時和時段定價數據
        state.days.forEach(day => {
            state.hourPricing[day] = {
                isActive: false,
                timePriceSettings: []
            };
            state.periodPricing[day] = {
                isActive: false,
                timePriceSettings: []
            };
        });

        // 初始化時間選項
        state.timeOptions = generateTimeOptions();

        // 綁定事件處理器
        bindPricingEvents();
    }

    // 生成時間選項
    function generateTimeOptions() {
        let options = [];
        for (let i = 0; i < 24; i++) {
            let hour = i.toString().padStart(2, '0');
            options.push({ label: `${hour}:00`, value: i * 2 });
            options.push({ label: `${hour}:30`, value: i * 2 + 1 });
        }
        options.push({ label: '23:59', value: 48 });
        return options;
    }

    // 綁定價格相關事件
    function bindPricingEvents() {
        // 小時計費相關事件
        $('#create-hour-plan').click(function() {
            openHourPlan();
            if ($('#edit-period-plan').css('display') === 'block'){
                $('#edit-period-plan').css('display', 'none');
            }
        });

        $('#edit-hour-plan-btn').click(function() {
            openHourPlan();
            $('#edit-hour-plan').css('display', 'none');
            $('#edit-period-plan').css('display', 'none');
        });

        $('#back-to-hour-plan').click(function() {
            saveAndBackToPlan();
            $('#edit-hour-plan').css('display', 'block');
            if ($('#create-period-plan-wrapper').css('display') === 'none') {
                $('#edit-period-plan').css('display', 'block');
            }
        });

        // 時段計費相關事件
        $('#create-period-plan').click(function() {
            openPeriodPlan();
            if ($('#edit-hour-plan').css('display') === 'block'){
                $('#edit-hour-plan').css('display', 'none');
            }
        });

        $('#edit-period-plan-btn').click(function() {
            openPeriodPlan();
            $('#edit-period-plan').css('display', 'none');
            $('#edit-hour-plan').css('display', 'none');
        });

        $('#back-to-period-plan').click(function() {
            saveAndBackToPlan();
            $('#edit-period-plan').css('display', 'block');
            if ($('#create-hour-plan-wrapper').css('display') === 'none') {
                $('#edit-hour-plan').css('display', 'block');
            }
        });

        // 開關切換事件
        $('#hour-switch').change(function() {
            toggleHourPlan(this.checked);
        });

        $('#period-switch').change(function() {
            togglePeriodPlan(this.checked);
        });

        // 最少租借小時相關事件
        $('#increment-least-hours').click(function() {
            incrementLeastRentHours();
        });

        $('#decrement-least-hours').click(function() {
            decrementLeastRentHours();
        });
    }

    // 打開小時計費設定
    function openHourPlan() {
        state.isCreatingHourPlan = false;
        state.isEditingHourPlan = false;
        state.isCreatingPeriodPlan = false;
        state.isEditingPeriodPlan = false;
        state.isSpecificHourPlan = true;
        state.isSpecificPeriodPlan = false;
        state.isEditedHourPlan = true;
        
        updatePricingUI();
    }

    // 打開時段計費設定
    function openPeriodPlan() {
        state.isCreatingHourPlan = false;
        state.isEditingHourPlan = false;
        state.isCreatingPeriodPlan = false;
        state.isEditingPeriodPlan = false;
        state.isSpecificHourPlan = false;
        state.isSpecificPeriodPlan = true;
        state.isEditedPeriodPlan = true;
        
        updatePricingUI();
    }

    // 更新價格UI顯示
    function updatePricingUI() {
        // 更新小時計費相關顯示
        if (state.isCreatingHourPlan) {
            $('#create-hour-plan-wrapper').show();
            $('#edit-hour-plan-wrapper').hide();
        } else {
            $('#create-hour-plan-wrapper').hide();
            $('#edit-hour-plan-wrapper').show();
        }
        $('#hour-plan-setting').toggle(state.isSpecificHourPlan);
        
        // 更新時段計費相關顯示
        if (state.isCreatingPeriodPlan) {
            $('#create-period-plan-wrapper').show();
            $('#edit-period-plan-wrapper').hide();
        } else {
            $('#create-period-plan-wrapper').hide();
            $('#edit-period-plan-wrapper').show();
        }
        $('#period-plan-setting').toggle(state.isSpecificPeriodPlan);
        
        // 更新價格範圍顯示
        updatePriceRanges();
        
        // 如果顯示具體設定，則生成日期設定
        if (state.isSpecificHourPlan) {
            generateHourPlanDays();
        }
        if (state.isSpecificPeriodPlan) {
            generatePeriodPlanDays();
        }
    }

    // 生成小時計費的日期設定
    function generateHourPlanDays() {
        const h3Element = $('#hour-plan-setting h3');
        
        // 清除已存在的open-day元素
        $('.open-day').remove();
        
        // 反向遍历days数组，这样可以按正确顺序插入到h3之后
        state.days.slice().reverse().forEach(day => {
            // 創建最外層的 open-day
            const openDay = $('<div>').addClass('open-day');
            
            // 創建 checkbox-switch
            const daySwitch = $('<div>').addClass('checkbox-switch');
            
            // 創建 input
            const input = $('<input>')
                .attr('type', 'checkbox')
                .attr('id', `${day}-hour-checkbox-switch`)
                .prop('checked', state.hourPricing[day].isActive);
            
            // 創建 label
            const label = $('<label>')
                .attr('for', `${day}-hour-checkbox-switch`)
                .html(`
                    <div class="switch-background">
                        <p class="close-text">關閉</p>
                        <div class="switch"></div>
                        <p class="open-text">開放</p>
                    </div>
                    <div class="close-text">
                        <div>
                            <p class="weekday-text">${state.dayMapping[day]}</p>
                        </div>
                    </div>
                    <div class="open-text">
                        <div>
                            <p class="weekday-text">${state.dayMapping[day]}</p>
                        </div>
                    </div>
                `);

            // 將 input 和 label 添加到 daySwitch
            daySwitch.append(input, label);

            // 創建 hour-price-setting-area
            const priceSettingArea = $('<div>').addClass('hour-price-setting-area');
            
            // 創建 open-label
            const openLabel = $('<div>')
                .addClass('open-label')
                .attr('style', state.hourPricing[day].isActive ? 'display: block;' : 'display: none;');
            
            // 創建 time-price-input-wrapper
            const timePriceWrapper = $('<div>').addClass('time-price-input-wrapper');
            
            // 創建 time-price-input-group
            const timePriceGroup = $('<div>')
                .addClass('time-price-input-group')
                .attr('id', `${day}-hour-input-group`);
            
            // 添加現有的時間價格設定
            state.hourPricing[day].timePriceSettings.forEach((setting, index) => {
                const timePriceInput = createTimePriceRow(day, 'hour', index, setting);
                timePriceGroup.append(timePriceInput);
            });
            
            // 創建 add-time-price-button
            const addButton = $('<div>')
                .addClass('add-time-price-button')
                .click(() => addPriceSetting(day, 'hour'))
                .html(`
                    <img src="/static/images/add-time-price-button.svg" alt="">
                    <p>新增小時定價</p>
                `);

            // 按照正確的層級結構組裝元素
            timePriceWrapper.append(timePriceGroup, addButton);
            openLabel.append(timePriceWrapper);
            priceSettingArea.append(openLabel);
            daySwitch.append(priceSettingArea);
            
            // 將 daySwitch 添加到 openDay
            openDay.append(daySwitch);
            
            // 綁定 checkbox 事件
            input.change(function() {
                toggleHourPricing(day, this.checked);
                openLabel.css('display', this.checked ? 'block' : 'none');
                
                // 如果激活且没有现有设置，创建初始时间价格设置
                if (this.checked && state.hourPricing[day].timePriceSettings.length === 0) {
                    const initialSetting = {
                        startTime: 0,
                        endTime: 1,
                        price: 0
                    };
                    const timePriceInput = createTimePriceRow(day, 'hour', 0, initialSetting);
                    timePriceGroup.append(timePriceInput);
                    state.hourPricing[day].timePriceSettings.push(initialSetting);
                }
            });

            // 將 openDay 添加到 h3 之後
            h3Element.after(openDay);
        });
    }

    // 生成時段計費的日期設定
    function generatePeriodPlanDays() {
        const h3Element = $('#period-plan-setting h3');
        
        // 清除已存在的open-day元素
        $('.open-day').remove();
        
        // 反向遍历days数组，这样可以按正确顺序插入到h3之后
        state.days.slice().reverse().forEach(day => {
            // 創建最外層的 open-day
            const openDay = $('<div>').addClass('open-day');
            
            // 創建 checkbox-switch
            const daySwitch = $('<div>').addClass('checkbox-switch');
            
            // 創建 input
            const input = $('<input>')
                .attr('type', 'checkbox')
                .attr('id', `${day}-period-checkbox-switch`)
                .prop('checked', state.periodPricing[day].isActive);
            
            // 創建 label
            const label = $('<label>')
                .attr('for', `${day}-period-checkbox-switch`)
                .html(`
                    <div class="switch-background">
                        <p class="close-text">關閉</p>
                        <div class="switch"></div>
                        <p class="open-text">開放</p>
                    </div>
                    <div class="close-text">
                        <div>
                            <p class="weekday-text">${state.dayMapping[day]}</p>
                        </div>
                    </div>
                    <div class="open-text">
                        <div>
                            <p class="weekday-text">${state.dayMapping[day]}</p>
                        </div>
                    </div>
                `);

            // 將 input 和 label 添加到 daySwitch
            daySwitch.append(input, label);

            // 創建 period-price-setting-area
            const priceSettingArea = $('<div>').addClass('period-price-setting-area');
            
            // 創建 open-label
            const openLabel = $('<div>')
                .addClass('open-label')
                .attr('style', state.periodPricing[day].isActive ? 'display: block;' : 'display: none;');
            
            // 創建 time-price-input-wrapper
            const timePriceWrapper = $('<div>').addClass('time-price-input-wrapper');
            
            // 創建 time-price-input-group
            const timePriceGroup = $('<div>')
                .addClass('time-price-input-group')
                .attr('id', `${day}-period-input-group`);
            
            // 添加現有的時間價格設定
            state.periodPricing[day].timePriceSettings.forEach((setting, index) => {
                const timePriceInput = createTimePriceRow(day, 'period', index, setting);
                timePriceGroup.append(timePriceInput);
            });
            
            // 創建 add-time-price-button
            const addButton = $('<div>')
                .addClass('add-time-price-button')
                .click(() => addPriceSetting(day, 'period'))
                .html(`
                    <img src="/static/images/add-time-price-button.svg" alt="">
                    <p>新增時段定價</p>
                `);

            // 按照正確的層級結構組裝元素
            timePriceWrapper.append(timePriceGroup, addButton);
            openLabel.append(timePriceWrapper);
            priceSettingArea.append(openLabel);
            daySwitch.append(priceSettingArea);
            
            // 將 daySwitch openDay
            openDay.append(daySwitch);
            
            // 綁定 checkbox 事件
            input.change(function() {
                togglePeriodPricing(day, this.checked);
                openLabel.css('display', this.checked ? 'block' : 'none');
                
                // 如果觸發且没有現有設定，創建初始時間及價格
                if (this.checked && state.periodPricing[day].timePriceSettings.length === 0) {
                    const initialSetting = {
                        startTime: 0,
                        endTime: 1,
                        price: 0
                    };
                    const timePriceInput = createTimePriceRow(day, 'period', 0, initialSetting);
                    timePriceGroup.append(timePriceInput);
                    state.periodPricing[day].timePriceSettings.push(initialSetting);
                }
            });

            // 將 openDay 添加到 h3 之後
            h3Element.after(openDay);
        });
    }

    // 創建時間價格設定行
    function createTimePriceRow(day, type, index, setting) {
        // 創建最外層容器
        const timePriceInput = $('<div>')
            .addClass('time-price-input visible new-input')
            .attr('data-index', index);  // 添加data-index属性
        
        // 創建 input-row
        const row = $('<div>').addClass('input-row');
    
        // 創建時間選擇區域
        const timeDiv = $('<div>').addClass('time');
        
        // 創建開始時間選擇
        const startTimeSelect = $('<select>')
            .addClass('start-time')
            .attr('name', `${type}-start-time-${index}`)
            .attr('required', true);

        // 創建結束時間選擇
        const endTimeSelect = $('<select>')
            .addClass('end-time')
            .attr('name', `${type}-end-time-${index}`)
            .attr('required', true);

        // 獲取已存在的時間段
        const existingTimeRanges = [];
        $(`#${day}-${type}-input-group .time-price-input`).each(function() {
            const startTime = parseInt($(this).find('.start-time').val());
            const endTime = parseInt($(this).find('.end-time').val());
            existingTimeRanges.push({ start: startTime, end: endTime });
        });

        // 生成所有可能的時間選項
        const allTimeOptions = generateTimeOptions();
        
        // 過濾出可用的時間選項
        const availableTimeOptions = allTimeOptions.filter(option => {
            const timeValue = option.value;
            // 檢查是否在任何已存在的時間區間
            return !existingTimeRanges.some(range => 
                timeValue >= range.start && timeValue < range.end
            );
        });

        // 添加可用的時間選項
        availableTimeOptions.forEach(option => {
            startTimeSelect.append(
                $('<option>')
                    .val(option.value)
                    .text(option.label)
            );
            endTimeSelect.append(
                $('<option>')
                    .val(option.value)
                    .text(option.label)
            );
        });

        // 設定默認值
        if (setting) {
            startTimeSelect.val(setting.startTime || 0);
            endTimeSelect.val(setting.endTime || 1);
        } else {
            // 如果是新增的时间價格设定，設定默認值為第一個可用的時間
            if (availableTimeOptions.length > 0) {
                startTimeSelect.val(availableTimeOptions[0].value);
                endTimeSelect.val(availableTimeOptions[0].value + 1);
            }
        }

        // 添加change監聽事件
        startTimeSelect.change(function() {
            const selectedStartTime = parseInt($(this).val());
            if (!setting) {
                setting = {
                    startTime: selectedStartTime,
                    endTime: parseInt(endTimeSelect.val()),
                    price: 0
                };
            } else {
                setting.startTime = selectedStartTime;
            }
            updatePriceSetting(day, type, index, setting);
        });

        endTimeSelect.change(function() {
            const selectedEndTime = parseInt($(this).val());
            if (!setting) {
                setting = {
                    startTime: parseInt(startTimeSelect.val()),
                    endTime: selectedEndTime,
                    price: 0
                };
            } else {
                setting.endTime = selectedEndTime;
            }
            updatePriceSetting(day, type, index, setting);
        });

        // 添加破折號
        const dash = $('<p>').addClass('dash').text('-');
    
        // 組裝時間選擇區域
        timeDiv.append(startTimeSelect, dash, endTimeSelect);
    
        // 創建價格輸入區域
        const priceDiv = $('<div>').addClass('per-price');
    
        // 創建價格輸入
        const priceInput = $('<input>')
            .attr('type', 'number')
            .attr('name', `${type}-price-${index}`)
            .attr('min', '1')
            .attr('required', true)
            .val(setting ? setting.price : 1);  // 如果没有setting，默认设置为1
    
        // 添加价格change事件
        priceInput.change(function() {
            const newPrice = parseInt($(this).val());
            if (!setting) {
                setting = {
                    startTime: parseInt(startTimeSelect.val()),
                    endTime: parseInt(endTimeSelect.val()),
                    price: newPrice
                };
            } else {
                setting.price = newPrice;
            }
            updatePriceSetting(day, type, index, setting);
        });
    
        if (type === 'hour') {
            // 組裝價格區域
            priceDiv.append(
                $('<p>').text('每小時'),
                priceInput,
                $('<p>').text('元')
            );
        } else {
            // 組裝價格區域
            priceDiv.append(
                $('<p>').text('此時段'),
                priceInput,
                $('<p>').text('元')
            );
        }

        // 組裝整個行
        row.append(timeDiv, priceDiv);

        const deleteButton = $('<img>')
            .attr('src', '/static/images/cross-gray.svg')
            .attr('alt', '刪除')
            .click(function() {
                const currentIndex = $(this).closest('.time-price-input').attr('data-index');
                removePriceSetting(day, currentIndex, type);
            });
        row.append(deleteButton);

        timePriceInput.append(row);
    
        return timePriceInput;
    }

    // 更新價格設定
    function updatePriceSetting(day, type, index, setting) {
        if (type === 'hour') {
            state.hourPricing[day].timePriceSettings[index] = setting;
        } else {
            state.periodPricing[day].timePriceSettings[index] = setting;
        }
        updatePriceRanges();
    }

    // 保存並返回計劃選擇
    function saveAndBackToPlan() {
        state.isCreatingHourPlan = !state.isEditedHourPlan;
        state.isEditingHourPlan = state.isEditedHourPlan;
        state.isCreatingPeriodPlan = !state.isEditedPeriodPlan;
        state.isEditingPeriodPlan = state.isEditedPeriodPlan;
        state.isSpecificHourPlan = false;
        state.isSpecificPeriodPlan = false;
        updatePricingUI();
    }

    // 增加最少租借小時
    function incrementLeastRentHours() {
        state.leastRentHours = parseFloat((state.leastRentHours + 0.5).toFixed(1));
        $('#least-rent-hours').val(state.leastRentHours);
    }

    // 減少最少租借小時
    function decrementLeastRentHours() {
        if (state.leastRentHours > 0.5) {
            state.leastRentHours = parseFloat((state.leastRentHours - 0.5).toFixed(1));
            $('#least-rent-hours').val(state.leastRentHours);
        }
    }

    // 切換小時計費開關
    function toggleHourPlan(isActive) {
        state.isEditingHourPlan = isActive;
        updatePricingUI();
    }

    // 切換時段計費開關
    function togglePeriodPlan(isActive) {
        state.isEditingPeriodPlan = isActive;
        updatePricingUI();
    }

    // 切換小時計費日期開關
    function toggleHourPricing(day, isActive) {
        const isPeriodActive = state.periodPricing[day].isActive;
        if (!isPeriodActive) {
            state.hourPricing[day].isActive = isActive;
            updatePriceRanges();
        } else {
            swal({
                title: '已開啟時段計費。',
                icon: 'warning',
                button: '關閉',
            });
            $(`#${day}-hour-checkbox-switch`).prop('checked', false);
        }
    }

    // 切換時段計費日期開關
    function togglePeriodPricing(day, isActive) {
        const isHourActive = state.hourPricing[day].isActive;
        if (!isHourActive) {
            state.periodPricing[day].isActive = isActive;
            updatePriceRanges();
        } else {
            swal({
                title: '已開啟小時計費。',
                icon: 'warning',
                button: '關閉',
            });
            $(`#${day}-period-checkbox-switch`).prop('checked', false);
        }
    }

    // 添加價格設定
    function addPriceSetting(day, type) {
        let lastTimeSetting;
        let newStartTime;
        let newEndTime;
        let newPriceSetting;
        let pricingArray;

        if (type === 'hour') {
            pricingArray = state.hourPricing[day].timePriceSettings;
        } else {
            pricingArray = state.periodPricing[day].timePriceSettings;
        }

        // 如果沒有任何設定，創建初始設定
        if (pricingArray.length === 0) {
            newPriceSetting = {
                startTime: 0,
                endTime: 1,
                price: 0
            };
            pricingArray.push(newPriceSetting);
            
            // 創建新的時間價格設定行
            const timePriceInput = createTimePriceRow(day, type, pricingArray.length - 1, newPriceSetting);
            
            // 將新的時間價格設定行添加到對應的輸入組
            if (type === 'hour') {
                $(`#${day}-hour-input-group`).append(timePriceInput);
            } else {
                $(`#${day}-period-input-group`).append(timePriceInput);
            }
            return;
        }

        // 獲取最後一筆設定
        lastTimeSetting = pricingArray.slice(-1)[0];
        if (lastTimeSetting && lastTimeSetting.endTime) {
            // 確保新的開始時間是上一個時間設定的結束時間
            newStartTime = lastTimeSetting.endTime;
            // 防止時間超過24小時制，這裡的時間單位是半小時，因此一天有48單位
            if (newStartTime === 48 || newStartTime === '') {
                swal({
                    title: '無法添加新的時間設定，因為沒有有效的最後結束時間。',
                    icon: 'error',
                    button: '關閉',
                    dangerMode: true
                });
                return;
            }
            newEndTime = (newStartTime === 47) ? 48 : newStartTime + 1;
            newPriceSetting = {
                startTime: newStartTime,
                endTime: newEndTime,
                price: 0
            };
            pricingArray.push(newPriceSetting);
            
            // 創建新的時間價格設定行
            const timePriceInput = createTimePriceRow(day, type, pricingArray.length - 1, newPriceSetting);
            
            // 將新的時間價格設定行添加到對應的輸入組
            if (type === 'hour') {
                $(`#${day}-hour-input-group`).append(timePriceInput);
            } else {
                $(`#${day}-period-input-group`).append(timePriceInput);
            }

        } else {
            // 如果沒有前一個時間設定，或者最後一個時間設定沒有結束時間
            swal({
                title: '無法添加新的時間設定，因為沒有有效的最後結束時間。',
                icon: 'error',
                button: '關閉',
                dangerMode: true
            });
        }
    }

    // 移除價格設定
    function removePriceSetting(day, index, type) {
        const targetIndex = parseInt(index);
        if (type === 'hour') {
            // 删除对应的价格设定
            state.hourPricing[day].timePriceSettings.splice(targetIndex, 1);
            // 删除对应的DOM元素
            $(`#${day}-hour-input-group .time-price-input[data-index="${targetIndex}"]`).remove();
            // 更新剩余行的data-index
            $(`#${day}-hour-input-group .time-price-input`).each(function(newIndex) {
                const currentIndex = parseInt($(this).attr('data-index'));
                if (currentIndex > targetIndex) {
                    $(this).attr('data-index', currentIndex - 1);
                }
            });
        } else {
            // 删除对应的价格设定
            state.periodPricing[day].timePriceSettings.splice(targetIndex, 1);
            // 删除对应的DOM元素
            $(`#${day}-period-input-group .time-price-input[data-index="${targetIndex}"]`).remove();
            // 更新剩余行的data-index
            $(`#${day}-period-input-group .time-price-input`).each(function(newIndex) {
                const currentIndex = parseInt($(this).attr('data-index'));
                if (currentIndex > targetIndex) {
                    $(this).attr('data-index', currentIndex - 1);
                }
            });
        }
        updatePriceRanges();
    }

    // 更新價格範圍顯示
    function updatePriceRanges() {
        $('#hour-price-range').text(getPriceRange(state.hourPricing, '小時'));
        $('#period-price-range').text(getPriceRange(state.periodPricing, '時段'));
    }

    // 獲取價格範圍
    function getPriceRange(pricingData, unit) {
        let minPrice = Number.MAX_VALUE;
        let maxPrice = -1;

        Object.values(pricingData).forEach(dayPricing => {
            if (dayPricing.isActive) {
                dayPricing.timePriceSettings.forEach(priceSetting => {
                    if (priceSetting.price > 0) {
                        if (priceSetting.price < minPrice) {
                            minPrice = priceSetting.price;
                        }
                        if (priceSetting.price > maxPrice) {
                            maxPrice = priceSetting.price;
                        }
                    }
                });
            }
        });

        if (maxPrice === -1) return '未設定';
        if (minPrice === maxPrice) {
            return `${minPrice}元/${unit}`;
        }
        return `${minPrice}-${maxPrice}元/${unit}`;
    }

    // 返回上一頁按鈕點擊事件
    $('#return-to-previous-page').click(function(e) {
        e.preventDefault();
        swal({
            title: '確定要返回上一頁嗎？',
            text: '返回上一頁後，所有未儲存之資料將會遺失。',
            icon: 'warning',
            buttons: ['取消', '確定'],
        }).then(function(isConfirmed) {
            if (isConfirmed) {
                window.history.back();
            }
        });
    });

    // 初始化
    initialize();
}); 