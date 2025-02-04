const infoItemArray = document.querySelectorAll(".info-item");
const phoneNumber = document.querySelector(".phone-number");
const phoneValueDetail = phoneNumber.closest(".value-detail");



window.onload = () => {
    PhoneBtnSet();
    InfoItemEditSet();
    
    const passwordError = document.getElementById('passwordError').innerText;
    if (passwordError === "當前密碼不正確" || passwordError === "輸入的新密碼不相同") {
        swal({
            title: "當前密碼不正確，請重新操作",
            icon: "warning",
            button: true,
            dangerMode: true
        });
    }

    const passwordSuccess = document.getElementById('success').value;
    if (passwordSuccess === 'true') {
        swal({
            title: "密碼更換成功！",
            icon: "success",
            button: true,
            dangerMode: true
        });
    }
}

function PhoneBtnSet() {
    let existingButton = phoneValueDetail.querySelector(".create-phone") ||
        phoneValueDetail.querySelector(".value-edit");
    
    if (existingButton) {
        existingButton.parentNode.removeChild(existingButton);
    }


    if (phoneNumber.innerHTML === "尚未設定") {
        let newDiv = document.createElement("div");
        newDiv.className = "create-phone";

        let newP = document.createElement("p");
        newP.className = "create-text";
        newP.textContent = "新增電話";
        newDiv.appendChild(newP);

        let newImg = document.createElement("img");
        newImg.className = "add-phone-btn";
        newImg.src = "../images/add.svg";
        newImg.alt = "";
        newDiv.appendChild(newImg);

        phoneValueDetail.appendChild(newDiv);
        InfoItemEditSet()
    }
    else {
        let div = document.createElement('div');
        div.innerHTML = `<svg xmlns="http://www.w3.org/2000/svg" class="value-edit" x="0px" y="0px" width="30" height="30" viewBox="0 0 32 32"> <path d="M 23.90625 3.96875 C 22.859375 3.96875 21.8125 4.375 21 5.1875 L 5.1875 21 L 5.125 21.3125 L 4.03125 26.8125 L 3.71875 28.28125 L 5.1875 27.96875 L 10.6875 26.875 L 11 26.8125 L 26.8125 11 C 28.4375 9.375 28.4375 6.8125 26.8125 5.1875 C 26 4.375 24.953125 3.96875 23.90625 3.96875 Z M 23.90625 5.875 C 24.410156 5.875 24.917969 6.105469 25.40625 6.59375 C 26.378906 7.566406 26.378906 8.621094 25.40625 9.59375 L 24.6875 10.28125 L 21.71875 7.3125 L 22.40625 6.59375 C 22.894531 6.105469 23.402344 5.875 23.90625 5.875 Z M 20.3125 8.71875 L 23.28125 11.6875 L 11.1875 23.78125 C 10.53125 22.5 9.5 21.46875 8.21875 20.8125 Z M 6.9375 22.4375 C 8.136719 22.921875 9.078125 23.863281 9.5625 25.0625 L 6.28125 25.71875 Z" fill="#d8d8d8"></svg>`;
        let svgNode = div.firstElementChild;

        phoneValueDetail.insertAdjacentElement('beforeend', svgNode);
        InfoItemEditSet()
    }
}

function InfoItemEditSet() {
    infoItemArray.forEach(infoItem => {
        const valueEdit = infoItem.querySelector(".value-edit");
        const cancelBtn = infoItem.querySelector(".cancel");
        const submitBtn = infoItem.querySelector(".submit");
        const infoItemTitle = infoItem.querySelector(".info-item-title");
        const infoItemValue = infoItem.querySelector(".info-item-value");
        const infoItemInput = infoItem.querySelector(".info-item-input");
        const createPhone = infoItem.querySelector(".create-phone");

        const inputs = Array.from(infoItem.querySelectorAll("input"));
        const firstNameInput = inputs.find(input => input.name === "FirstName");
        const lastNameInput = inputs.find(input => input.name === "LastName");
        const emailInput = inputs.find(input => input.name === "Email");
        const phoneInput = inputs.find(input => input.name === "TelePhone");
        const currentPasswordInput = inputs.find(input => input.name === "CurrentPassword");
        const newPasswordInput = inputs.find(input => input.name === "NewPassword");
        const confirmNewPassword = inputs.find(input => input.name === "ConfirmNewPassword");


        function DisplayInput() {
            infoItemTitle.style.display = "none";
            infoItemValue.style.display = "none";
            infoItemInput.style.display = "block";
        }

        function HideInput() {
            infoItemTitle.style.display = "block";
            infoItemValue.style.display = "block";
            infoItemInput.style.display = "none";
        }


        if (valueEdit) {
            valueEdit.addEventListener("click", function () {
                DisplayInput();
            })
        }

        if (createPhone) {
            createPhone.addEventListener("click", function () {
                DisplayInput();
            })
        }

        if (cancelBtn) {
            cancelBtn.addEventListener("click", function () {
                if (firstNameInput && lastNameInput) {
                    var valueDetail = infoItem.closest(".info-item").querySelector(".value-detail p");
                    var namesSplit = valueDetail.textContent.split(" ");
                    if (namesSplit.length > 1) {
                        lastNameInput.value = namesSplit[1];
                        firstNameInput.value = namesSplit[0];
                    }
                } else if (emailInput) {
                    var valueDetail = infoItem.closest(".info-item").querySelector(".value-detail p");
                    emailInput.value = (valueDetail.innerText !== "尚未設定") ? valueDetail.innerText : "";
                } else if (phoneInput) {
                    var valueDetail = infoItem.closest(".info-item").querySelector(".value-detail p");
                    phoneInput.value = (valueDetail.innerText !== "尚未設定") ? valueDetail.innerText : "";
                }
                HideInput();
            })
        }

        if (submitBtn) {
            submitBtn.addEventListener("click", function (e) {
                var hasEmptyInput = inputs.some(function (input) {
                    return input.value.trim() === "";
                });

                if (hasEmptyInput) {
                    e.preventDefault();
                    swal({
                        title: "資料欄位不可為空白",
                        icon: "warning",
                        button: true,
                        dangerMode: true
                    });
                    return;
                }

                if (emailInput && !emailInput.value.match(/^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$/)) {
                    e.preventDefault();
                    swal({
                        title: "輸入的Email格式不正確",
                        icon: "warning",
                        button: true,
                        dangerMode: true
                    });
                    return;
                }

                if (phoneInput && !phoneInput.value.match(/^09\d{8}$/)) {
                    e.preventDefault();
                    swal({
                        title: "輸入的手機號碼須為09XXXXXXXX ",
                        icon: "warning",
                        button: true,
                        dangerMode: true
                    });
                    return;
                }

                if (currentPasswordInput.value === newPasswordInput.value) {
                    e.preventDefault();
                    swal({
                        title: "新舊密碼不可相同",
                        icon: "warning",
                        button: true,
                        dangerMode: true
                    });
                    return;
                }

                if (newPasswordInput.value !== confirmNewPassword.value) {
                    e.preventDefault();
                    swal({
                        title: "輸入的新密碼不相同",
                        icon: "warning",
                        button: true,
                        dangerMode: true
                    });
                    return;
                }
                
                if (firstNameInput && lastNameInput) {
                    var valueDetail = infoItem.closest(".info-item").querySelector(".value-detail p");
                    valueDetail.textContent = lastNameInput.value + " " + firstNameInput.value;
                } else if (singleInput) {
                    var valueDetail = infoItem.closest(".info-item").querySelector(".value-detail p");
                    valueDetail.innerText = singleInput.value;
                }
                HideInput();
                PhoneBtnSet();
            })
        }

    });
}

